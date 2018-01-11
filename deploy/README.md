# Getting started with Quobyte Services on Kubernetes

This setup guide assumes that you have 4 servers available.
Quobyte stores its data on dedicated disk drives, so please make sure that
the servers have some spare, unformatted SSDs or HDDs available.

## Deploying Quobyte on a Kubernetes Cluster

The following steps will deploy Quobyte as a cluster-wide /var/lib/kubelet/plugins/kubernetes.io~quobyte mount on a Kubernetes (>=1.1) cluster. If you use Kubernetes (>=1.4) you can use the official Quobyte Kubernetes volume plugin.

Quobyte services use dedicated disk drives to store data.
A typical Quobyte cluster consists of 3 registries, 3 metadata servers and at least 3 data servers.

Since we know the servers and their connected disk drives in advance, we can layout a plan for the
Quobyte deployment. To run a Quobyte service on a node, we apply labels for 'metadata', 'data', 'registry', and 'client'.
Kubernetes will make sure that the services are scheduled on the nodes.

## Prerequisites
### Format Quobyte Devices
Quobyte is designed to run on dedicated disk drives which are formatted with ext4 or xfs and
are initialized as a Quobyte device.

Log in the host machines and prepare the devices. In this example, we expect
`/dev/sd{b,c,d}` to be available.

```bash
$ ssh $host

# prepare the target mount dirs
$ for i in `seq 1 1 3`; do mkdir -p /mnt/quobyte/data_$i; done

# partition and format disk drives
$ for dev in sdb sdc sdd; do \
    parted /dev/${dev} mklabel gpt; \
    parted /dev/${dev} mkpart primary 2048s 100%; \
    mkfs.xfs -isize=1024 /dev/${dev}1; \
  done

# persist to fstab.

# retrieve the PARTUUIDs for the newly created devices
$ for dev in sdb sdc sdd; do \
  blkid -s PARTUUID -o value /dev/${dev}1; \
  done

# for each disk, add a line to fstab
PARTUUID=<PARTUUID from blkid> /mnt/quobyte/data_$i xfs relatime,logbufs=8,logbsize=256k,swalloc,allocsize=131072k

# for SSDs we recommend mount options:
PARTUUID=<PARTUUID from blkid> /mnt/quobyte/data_$i relatime,nodiscard
```

When all fstab entries are created, mount the fresh disks with `mount -a`.

Next, we need to bootstrap a registry and setup the devices to be recognized by Quobyte.
Then, prepare the other registries, metadata, and data devices as they fit your cluster setup.
Please note that every device can host all three types in parallel, but for initial setup,
it can only be initialized for a single purpose. Our target cluster will each have /mnt/quobyte/data_1 as registry
and metadata enabled, and data_2/3 as data disks.

To achieve this, we will setup the first machine like:
```bash
$ ssh $host
$ curl -O https://raw.githubusercontent.com/quobyte/kubernetes/master/tools/qbootstrap
$ chmod +x qbootstrap

# create a bootstrapped registry on first disk
$ sudo ./qbootstrap /mnt/quobyte/data_1

# and two data devices
$ curl -O https://raw.githubusercontent.com/quobyte/kubernetes/master/tools/qmkdev
$ chmod +x qmkdev
$ sudo ./qmkdev -f -s $(uuidgen) -t DATA /mnt/quobyte/data_2
$ sudo ./qmkdev -f -s $(uuidgen) -t DATA /mnt/quobyte/data_3
```

On the other 3 machines:

```bash
$ ssh $host
$ curl -O https://raw.githubusercontent.com/quobyte/kubernetes/master/tools/qmkdev
$ chmod +x qmkdev

# create registry - will be used for METADATA as well
$ sudo ./qmkdev -f -s $(uuidgen) -t REGISTRY /mnt/quobyte/data_1
$ sudo ./qmkdev -f -s $(uuidgen) -t DATA /mnt/quobyte/data_2
$ sudo ./qmkdev -f -s $(uuidgen) -t DATA /mnt/quobyte/data_3
```

Please note that the target cluster will use shared REGISTRY / METADATA devices.
Hence, we first set up all these devices as REGISTRY, and later upgrade them
to also host METADATA.

#### MountFlags

Ensure MountFlags are shared or not set:
```bash
$ systemctl cat docker.service | grep MountFlags=shared
```

If not you can set the MountFlags as shared with the following commands:
```bash
$ cat << EOF > /etc/systemd/system/docker.service.d/slave-mount-flags.conf
[Service]
MountFlags=shared
EOF

$ systemctl daemon-reload && systemctl restart docker
```

#### NTP

Ensure that ntp is running on all of your nodes hosting any Quobyte service otherwise this can lead to a non working cluster.


## Starting the Quobyte Components

### Create "quobyte" Namespace

```bash
$ cd deploy
$ kubectl create -f quobyte-ns.yaml
```

**Additional**

Set the default namespace of `kubectl` to the quobyte namespace:

```bash
$ export CONTEXT=$(kubectl config view | awk '/current-context/ {print $2}')

$ kubectl config set-context $CONTEXT --namespace=quobyte

$ kubectl config view | grep namespace:
namespace: quobyte
```

If you like to reset the to the `default` namespace just repeat the steps above and replace quobyte with `default`.

### Create Configuration

Create the memory configuration for the Quobyte resources (if you change this file don't forget to adjust the resource limits for the components):

```bash
$ kubectl create -f config.yaml
```

### Create Services

```bash
$ kubectl create -f quobyte-services.yaml
```

### Deploy Registry

Start with labeling one node as `registry`, we use this node as bootstrap node.
If you have prepared physical disks with qbootstrap, make sure that you label that host as `registry`.
Also, only label the bootstrap node. All other nodes will be added in a second step.

```bash
$ kubectl label nodes <bootstrap-node> quobyte_registry="true"
```

Now start the deployment for the registry.
The DaemonSet will automatically deploy a `registry` on the nodes:

```bash
$ kubectl create -f registry-ds.yaml
$ kubectl -n quobyte get po --watch
```

### Deploy the API and the Web UI

```bash
$ kubectl create -f webconsole-deployment.yaml
```

We have created a Quobyte management container which holds the relevant API clients to
control and monitor the cluster.
As soon as the webconsole/api pod is running:

```bash
$ kubectl create -f qmgmt-pod.yaml
$ kubectl -n quobyte exec -it qmgmt-pod -- qmgmt -u api device list
Id  Host                                      Mode                   Disk Used  Disk Avail  Services    LED Mode
 1  registry-39v1s                            ONLINE                      4 GB       40 GB  REGISTRY    OFF

$ kubectl -n quobyte exec -it qmgmt-pod -- qmgmt -u api registry list
Id  Host                                         Mode
 1  registry-39v1s                               ONLINE
```

Also, you can now connect to the Webconsole.

```bash
$ kubectl port-forward <webconsole-pod> 8080 > /dev/null &
```
and point your browser to [http://localhost:8080](http://localhost:8080) to log in with default credentials "admin" and password "quobyte".

### Deploy Remaining Registries and Upgrade Devices

As soon as the first registry is up and running, mark the remaining nodes which contain REGISTRY devices.
Wait until all registry devices show up online.
```bash
$ kubectl label nodes <node-1> <node-n> quobyte_registry="true"
$ watch "kubectl -n quobyte exec -it qmgmt-pod -- qmgmt -u api device list | grep REGISTRY"
```

You can skip this step, if you have prepared all METADATA devices with qmkdev already.
```bash

$ kubectl -n quobyte exec -it qmgmt-pod -- qmgmt -u api device list | grep REGISTRY
  1  registry-26tj1                            ONLINE                     34 MB       21 GB  REGISTRY    OFF       
  4  registry-716mq                            ONLINE                    168 MB       21 GB  REGISTRY    OFF       
 10  registry-fnr2x                            ONLINE                     34 MB       21 GB  REGISTRY    OFF       
  7  registry-fxwc2                            ONLINE                    168 MB       21 GB  REGISTRY    OFF

# for all relevant device ids update the devices to also host METADATA:
kubectl -n quobyte exec -it qmgmt-pod -- qmgmt -u api device update add-type <device-id> METADATA
```

### Deploy the Metadata and Data Service on each Node

Add a label to all nodes that should run the Quobyte `metadata` and/or `data` service:

```bash
$ kubectl label nodes <node-1> <node-n> quobyte_metadata="true"
$ kubectl label nodes <node-1> <node-n> quobyte_data="true"
```

Now you can start the DaemonSets to run the pods on these nodes.

```bash
$ kubectl create -f data-ds.yaml
$ kubectl create -f metadata-ds.yaml
```

Now all Quobyte service are up and running. Check that all prepared devices are actually online:

```bash
$ kubectl -n quobyte exec -it qmgmt-pod -- qmgmt -u api device list
Id  Host                                      Mode                   Disk Used  Disk Avail  Services    LED Mode
 11  data-0rhrs                                ONLINE                     38 MB       21 GB  DATA        OFF       
 12  data-0rhrs                                ONLINE                     38 MB       21 GB  DATA        OFF       
  2  data-1kv8v                                ONLINE                     38 MB       21 GB  DATA        OFF       
  3  data-1kv8v                                ONLINE                     38 MB       21 GB  DATA        OFF       
  8  data-8lsd5                                ONLINE                     38 MB       21 GB  DATA        OFF       
  9  data-8lsd5                                ONLINE                     38 MB       21 GB  DATA        OFF       
  5  data-tjkqs                                ONLINE                     38 MB       21 GB  DATA        OFF       
  6  data-tjkqs                                ONLINE                     38 MB       21 GB  DATA        OFF       
  1  metadata-0wtnc                            ONLINE                    168 MB       21 GB  METADATA REGISTRY  OFF       
 10  metadata-jdtzk                            ONLINE                     34 MB       21 GB  METADATA REGISTRY  OFF       
  7  metadata-md78r                            ONLINE                    168 MB       21 GB  METADATA REGISTRY  OFF       
  4  metadata-qkjhx                            ONLINE                    168 MB       21 GB  METADATA REGISTRY  OFF    
```





## Clients and Volumes

To actually access data from Quobyte, we require to have a Quobyte Client running on every host.
First, create a Quobyte volume, which hosts some test data.

```bash
$ kubectl -n quobyte exec -it qmgmt-pod -- qmgmt -u api volume create testVolume root root BASE 0777
Success. Created new volume with volume uuid ba843ca0-40e2-4f05-ae41-2813d6262201
```

Then you can mount all volumes on each node at `/var/lib/kubelet/plugins/kubernetes.io~quobyte` (default plugin directory) using the client daemonset:

```bash
$ kubectl create -f client-ds.yaml
$ kubectl label nodes <node-1> <node-n> quobyte_client="true"
```

Log into one of the hosts with ssh and check that the volume is mounted:

```bash
$ grep "quobyte" /proc/mounts
quobyte@10.244.4.3:7866|10.244.3.3:7866|10.244.2.4:7866/cluster on /var/lib/kubelet/plugins/kubernetes.io~quobyte type fuse (rw,nosuid,nodev,noatime,user_id=0,group_id=0,default_permissions)
```

## Deploy a Pod using the Cluster Storage

We prepared a simple demo pod which uses the previously prepared testVolume to demonstrate shared file system access.
It uses the [Kubernetes Quobyte Plugin](https://kubernetes.io/docs/concepts/storage/volumes/#quobyte).

```bash
$ kubectl create -f example-pod.yaml

$ kubectl logs example -f
Starting with a fresh state
Tue May 10 10:01:39 UTC 2016
Tue May 10 10:01:44 UTC 2016
Tue May 10 10:01:49 UTC 2016
```

Now delete and recreate the pod:

```bash
$ kubectl delete pod example --grace-period 0
$ kubectl create -f example-pod.yaml
```
Independently from the node the pod comes up on, the pod will find its old state:

```bash
$ kubectl logs example -f
Found old state starting at Tue May 10 10:01:39 UTC 2016
Tue May 10 10:07:32 UTC 2016
Tue May 10 10:07:37 UTC 2016
Tue May 10 10:07:42 UTC 2016
```

## Access Webconsole and API

### Local port-forward

```bash
$ kubectl port-forward <webconsole-pod> 8080 > /dev/null &
```

Access the Webconsole at <http://localhost:8080>

### Ingress

```bash
$ kubectl create -f ingress.yaml
```

For more information about `ingress` have a look at the [contrib repo](https://github.com/kubernetes/contrib/tree/master/ingress/controllers) and the [docs](http://kubernetes.io/docs/user-guide/ingress)

### NodePort

For local development you can also use a [NodePort](http://kubernetes.io/docs/user-guide/services/#type-nodeport):

```
$ cat webconsole-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: webconsole
spec:
  type: NodePort
  ports:
    - name: web80
      targetPort: 8080
      port: 80
      protocol: TCP
    - name: web
      port: 8080
      protocol: TCP
  selector:
    role: webconsole

$ cat api-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: api
spec:
  type: NodePort
  ports:
    - name: api80
      targetPort: 7860
      port: 80
      protocol: TCP
    - name: api
      port: 7860
      protocol: TCP
  selector:
    role: webconsole
```

# Acknowledgement

- [sttts](https://github.com/sttts)
- [johscheuer](https://github.com/johscheuer)
