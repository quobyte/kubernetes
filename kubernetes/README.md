# Deploying Quobyte on a Kubernetes Cluster

The following steps will deploy Quobyte as a cluster-wide `/var/lib/kubelet/plugins/kubernetes.io~quobyte` mount on a Kubernetes (>=1.1) cluster. If you use Kubernetes (>=1.4) you can use the official Quobyte Kubernetes volume plugin.

## Prerequisites

Ensure MountFlags are shared or not set:

```bash
$ systemctl cat docker.service | grep MountFlags=shared
```

If not you can set the MountFlags as `shared` with the following commands:

```bash
$ cat << EOF > /etc/systemd/system/docker.service.d/slave-mount-flags.conf
[Service]
MountFlags=shared
EOF

$ systemctl daemon-reload && systemctl restart docker
```

## Starting the Quobyte Components

### Create "quobyte" Namespace

```bash
$ kubectl create -f quobyte-ns.yaml
```

Set the default namespace of `kubectl` to the quobyte namespace:

```bash
$ export CONTEXT=$(kubectl config view | awk '/current-context/ {print $2}')

$ kubectl config set-context $CONTEXT --namespace=quobyte

$ kubectl config view | grep namespace:
namespace: quobyte
```

If you like to reset the to the `default` namespace just repeat the steps above and replace quobyte with `default`.

### Create Configuration

Create the memory configuration for the quobyte resources (if you change this file don't forget to adjust the resource limits for the components):

```bash
$ kubectl create -f config.yaml
```

### Create Services

```bash
$ kubectl create -f quobyte-services.yaml
```

### Deploy Registry

Start with labeling one node as `registry`, we use this node as bootstrap node:

```bash
$ kubectl label nodes <bootstrap-node> quobyte_registry="true"
```

Now start the deployment for the registry:

```bash
$ kubectl create -f registry-ds.yaml
```

Wait until the pod is up now you can label other nodes as `registry` to add additional Quobyty registries.

```bash
$ kubectl get po --watch
```

The DaemonSet will automatically deploy a `registry` on these nodes:

```bash
$ kubectl label nodes <node-name-1> <node-name-n> quobyte_registry="true"
```

### Deploy the API and the Web UI

```bash
$ kubectl create -f webconsole-deployment.yaml
```

Wait that all devices are up, e.g. by

```bash
$ kubectl create -f qmgmt-pod.yaml

$ kubectl exec -it qmgmt-pod -- qmgmt -u api:7860 device list
Id  Host                                      Mode                   Disk Used  Disk Avail  Services
 1  registry-21ndb                            ONLINE                      4 GB       17 GB  REGISTRY
 2  registry-5gtgz                            ONLINE                      4 GB       17 GB  REGISTRY
 3  registry-irvss                            ONLINE                      4 GB       17 GB  REGISTRY

$ kubectl exec -it qmgmt-pod -- qmgmt -u api:7860 registry list
Id  Host                                         Mode
*   1  registry-21ndb                            ONLINE
 2  registry-5gtgz                               ONLINE
 3  registry-irvss                               ONLINE
```

### Deploy the Metadata and Data Service on each Node

Add a label to all nodes that should run the Quobyte `meta-data` and/or `data` service:

```bash
$ kubectl label nodes <node-name-1> <node-name-n> quobyte_metadata="true"
$ kubectl label nodes <node-name-1> <node-name-n> quobyte_data="true"
```

Now you can start the DaemonSets to run the pods on these nodes.

```bash
$ kubectl create -f data-ds.yaml
$ kubectl create -f metadata-ds.yaml
```

Wait that all devices are up with the upper `kubectl exec` command and check all services:

```bash
$ kubectl exec -it qmgmt-pod -- qmgmt -u api:7860 service list
Name                                      Type             UUID
data-6bdhe                                Data (D)         47638429-f0ed-4ea7-8f12-2452b17cde72
data-94uc7                                Data (D)         c307c3fc-c7fd-45e8-a3ef-2e24ecfb23d8
data-ogq7j                                Data (D)         2eb58546-a507-4a20-81b7-eccc86ff2695
metadata-czs40                            Metadata (M)     c336de21-c669-4da9-90ca-dd2d9ec7977a
metadata-q3nx7                            Metadata (M)     f9fd3200-2c46-4dd1-8926-43be69545764
metadata-zjz6c                            Metadata (M)     1c453894-aa27-4fc1-9a80-d5c54c3fde4d
registry-21ndb                            Registry (R)     87c78331-97d2-4878-9746-a3ff1c07c333
registry-5gtgz                            Registry (R)     c5c0ef49-4656-407a-845e-dc393ddcbaf4
registry-irvss                            Registry (R)     6314af41-7dac-4f75-8ee9-efcebe78db26
webconsole-3666637658-xwblz               API Proxy (A)    49edfb90-cc83-47d2-9d80-1ddc45b1b9cb
webconsole-3666637658-xwblz               Web Console      b5e35998-0ed4-4c31-a2e0-3ff2f973851e
```

## Create a Volume

Create a volume `cluster`:

```bash
$ kubectl exec -it qmgmt-pod -- qmgmt -u api:7860 volume create testVolume root root BASE 0777
```

Then you can mount all volumes on each node at `/var/lib/kubelet/plugins/kubernetes.io~quobyte` (default plugin directory) using the client daemonset:

```bash
$ kubectl create -f client-ds.yaml
```

Log into one of the nodes with ssh and check that the volume is mounted:

```bash
$ grep "quobyte" /proc/mounts
quobyte@10.244.4.3:7866|10.244.3.3:7866|10.244.2.4:7866/cluster on /var/lib/kubelet/plugins/kubernetes.io~quobyte type fuse (rw,nosuid,nodev,noatime,user_id=0,group_id=0,default_permissions)
```

## Deploy a Pod using the Cluster Storage

### Pre 1.4 kubernetes

```bash
$ kubectl create -f example-pod-pre.yaml

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

### Kubernetes 1.4+

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
