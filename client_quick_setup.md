# Quick Start Guide to use Quobyte from Kubernetes

## Prerequisites
If you use Docker 1.12 or older, you need to  add the MountFlags parameter.
For newer versions of Docker, this is not required anymore.

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

## Client setup
This guide assumes that you have a dedicated Quobyte instance running and you
want to provide access to Quobyte volumes to pods running in Kubernetes.

To access a Quobyte volume a pod has to run on a Kubernetes node which has a
Quobyte client running. The client runs inside of a Pod and makes the Quobyte
storage accessible to other pods.

Quobyte clients run in the quobyte namespace.
```bash
$ cd deploy
$ kubectl create -f quobyte-ns.yaml
```

To connect to Quobyte, the client needs to resolve the address of the registry.
It is configured in the client-ds.yaml DaemonSet definition:
```yaml
env:
  - name: QUOBYTE_REGISTRY
    value: registry.quobyte
```
If you have a QNS (Quobyte Name Service), the value for QUOBYTE_REGISTRY might
look like `<qnsid>.myquobyte.net`

If you have a certificate for the client, it is stored as a Secret and
mounted into the client Pod as client.cfg.

First create a file that contains only the certificate information
(<ca>, <cert>, and <key> blocks) and store it as a secret.
```bash
$ kubectl -n quobyte create secret generic client-config --from-file /tmp/client.cfg
```

```bash
$ kubectl -n quobyte create -f client-ds.yaml
or
$ kubectl -n quobyte create -f client-certificate-ds.yaml
```

The deployed DaemonSet starts client pods on all nodes marked as `quobyte_client`.
This can either be done manually, or by using the Quobyte operator.

```bash
$ kubectl label nodes <node-1> <node-n> quobyte_client="true"
```

When the client pod is up and running, you should see a mount point on the Kubernetes node
at `/var/lib/kubelet/plugins/kubernetes.io~quobyte`.

##Benchmarking

For easy testing and benchmarking we provide a fio-container which uses
Quobyte volumes. By default, it will start writing to volume `fio-test`

```bash
$ kubectl create -f fio-benchmark-ds.yaml
```
This will start a single Pod on a node which is marked as quobyte_client.
The container is designed to put load on the volume, so you can scale it:

```bash
$ kubectl scale --replicas=100 deployment fio-benchmark
```

## Volume Access

Quobyte volumes can be accessed from Pods in multiple ways. Either, via
the QuobyteVolumeSource like in `deploy/example-pod.yaml`

```yaml
volumes:
- name: quobytevolume
  quobyte:
    registry: LBIP:7861
    volume: testVolume
    readOnly: false
    user: username
    group: groupname
```

or as a PersistentVolumeClaim, which is defined in the same namespace, where
the accessing pod is running. The claim (volumes/claim.json)

```json
{
  "kind": "PersistentVolumeClaim",
  "apiVersion": "v1",
  "metadata": {
    "name": "test"
  },
  "spec": {
    "accessModes": [
      "ReadWriteOnce"
    ],
    "resources": {
      "requests": {
        "storage": "3Gi"
      }
    },
    "storageClassName": "base"
  }
}
```

```bash
$ kubectl -n quobyte create -f volumes/claim.json
```

is mounted in the Pod like (see volumes/example-pod.yaml):

```yaml
volumes:
  - name: quobytepvc
    persistentVolumeClaim:
      claimName: test
```

The PersistentVolumeClaim is bound to a PersistentVolume which is created by an
administrator for existing Quobyte Volumes or dynamically by the StorageClass
(volumes/storageclass.yaml).

If a Quobyte Volume already exists, the administrator can make it available
as a PersistentVolume by creating a resource (volumes/pv.yaml).
Please note that Quobyte currently ignores the required capacity.storage field,
since its using internal quota mechanisms.

```yaml
kind: PersistentVolume
apiVersion: v1
metadata:
  name: test
  labels:
    type: quobyte
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteOnce
  storageClassName: "base"
  quobyte:
    registry: LBIP:7861
    volume: test
    readOnly: false
    user: username
    group: groupname
```

```bash
$ kubectl -n quobyte create -f volumes/pv.yaml
```

For dynamic provisioning, a StorageClass is created, which manages the
lifecycle of Quobyte volumes.

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
   name: base
provisioner: kubernetes.io/quobyte
parameters:
    quobyteAPIServer: "http://api.quobyte:7860"
    registry: "registry.quobyte:7861"
    adminSecretName: "quobyte-admin-secret"
    adminSecretNamespace: "kube-system"
    user: "username"
    group: "groupname"
    quobyteConfig: "BASE"
    quobyteTenant: "DEFAULT"
    createQuota: "False"
```

To enable it, first create the Quobyte admin secret according to your API
credentials in volumes/quobyte-admin-secret.yaml (defaults to admin:quobyte):
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: quobyte-admin-secret
type: "kubernetes.io/quobyte"
data:
  password: cXVvYnl0ZQ==
  user: YWRtaW4=
type: kubernetes.io/quobyte
```

The password and user strings are base64 encoded and can be created like
`echo "quobyte" | base64`.

```bash
$ kubectl -n kube-system create -f volumes/quobyte-admin-secret.yaml
$ kubectl -n quobyte create -f volumes/storageclass.yaml
```

## Quobyte Tenants and Kubernetes Namespaces

Quobyte supports multiple tenants and provides a secure mapping of containers to
 users known in Quobyte.
For a longer read, please see the article on the Quobyte blog:
[The State of Secure Storage Access in Container Infrastructures](https://www.quobyte.com/blog/2017/03/17/the-state-of-secure-storage-access-in-container-infrastructures/)

TODO: document this


## Further Readings

- https://github.com/kubernetes/examples/tree/master/staging/volumes/quobyte

- https://kubernetes.io/docs/tasks/configure-pod-container/configure-persistent-volume-storage/#create-a-persistentvolume

- https://kubernetes.io/docs/concepts/storage/storage-classes/#quobyte
