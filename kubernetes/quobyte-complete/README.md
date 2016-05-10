# Deploying Quobyte on a Kubernetes Cluster

The following steps will deploy Quobyte as a cluster-wide `/storage/cluster` mount on a Kubernetes (>=1.1) cluster.

## Starting the Quobyte Components

### Create "quobyte" Namespace

```bash
$ kubectl create -f quobyte-ns.yaml
```

### Create Services

```bash
$ kubectl create -f api-service.yaml
$ kubectl create -f registry-service.yaml
$ kubectl create -f webconsole-service.yaml
```

### Deploy Registry

Bootstrap the replicated registry:

```bash
$ kubectl create -f  registry-bootstrap-pod.yaml
```

Wait until the registry pod is up. Then start the replication controller:

```bash
$ kubectl create -f registry-rc.yaml
```

The replication controller's selector will also match the bootstrapping pod. So there is no need to shut that down.

### Deploy the API and the Web UI

```bash
$ kubectl create -f webconsole-rc.yaml
```

Wait that all devices are up, e.g. by

```bash
$ kubectl exec -it --namespace=quobyte registry-bootstrap -- qmgmt -u api:7860 device list
```

### Deploy the Metadata and Data Service on each Node

```bash
$ kubectl create -f data-ds.yaml
$ kubectl create -f metadata-ds.yaml
```

Wait that all devices are up with the upper `kubectl exec` command.

## Create a Volume

Create a volume `cluster`:

```bash
$ kubectl exec -it --namespace=quobyte registry-bootstrap -- qmgmt -u api:7860 volume create cluster root root BASE 0777
```

Then you can mount it on each node at `/storage/cluster` using the client daemonset:

```bash
$ kubectl create -f client-ds.yaml
```

Log into one of the nodes with ssh and check that the volume is mounted:

```bash
$ mount | grep quobyte
quobyte@10.244.4.3:7866|10.244.3.3:7866|10.244.2.4:7866/cluster on /storage/cluster type fuse (rw,nosuid,nodev,noatime,user_id=0,group_id=0,default_permissions)
```

## Deploy a Pod using the Cluster Storage

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
