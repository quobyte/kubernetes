# Quobyte on Kubernetes and CoreOS

## Prerequisites

- Own Docker Registry or import the images on all nodes
- MountFlags=shared
- Kuberntes 1.1.+
- enabled `extensions/v1beta1`

## Testing local

See [CoreOS Kubernetes](https://github.com/coreos/coreos-kubernetes/blob/master/Documentation/kubernetes-on-vagrant.md)

### Change ../generic/worker-install.sh

#### MountFlags

```Bash
# Line 170
local TEMPLATE=/etc/systemd/system/docker.service.d/slave-mount-flags.conf
[ -f $TEMPLATE ] || {
   echo "TEMPLATE: $TEMPLATE"
   mkdir -p $(dirname $TEMPLATE)
   cat << EOF > $TEMPLATE
[Service]
MountFlags=shared
EOF
}
```

#### Quobyte images

Do this on each Node (until I found a solution to automate this or use a docker registry -> do this in worker-install)

```Bash
$ export REPO_ID=<repo-id>
$ export QUOBYTE_VERSION=1.2.6
$ curl https://support.quobyte.com/repo/2/$REPO_ID/quobyte-docker/quobyte-server-image_$QUOBYTE_VERSION.tar.bzip2 | bunzip2 -c | docker load
$ curl https://support.quobyte.com/repo/2/$REPO_ID/quobyte-docker/quobyte-client-image_$QUOBYTE_VERSION.tar.bzip2 | bunzip2 -c | docker load
```

### First Registry Node

```Bash
$ docker run --rm -ti --privileged \
 -v /mnt:/devices \
 --name bootstrap \
 quobyte-server:$QUOBYTE_VERSION qbootstrap -y -f /devices &&
sudo mkdir /mnt/dev1 &&
sudo mv -f /mnt/QUOBYTE_DEV_SETUP /mnt/dev1/ &&
ls -la /mnt/dev1/
```

### Other Registry Nodes

```Bash
$ docker run --rm -ti --privileged \
 -v /mnt:/devices \
 --name bootstrap \
 quobyte-server:$QUOBYTE_VERSION qmkdev -f -t registry /devices &&
sudo mkdir -p /mnt/dev1 &&
sudo mv -f /mnt/QUOBYTE_DEV_SETUP /mnt/dev1/ &&
ls -la /mnt/dev1/
```

## Setup DNS (local)

Run on each worker vagrant box

```bash
$ sudo sh -c 'echo -e "172.17.4.201  w1
172.17.4.202  w2
172.17.4.203  w3" > /etc/hosts'
```

## Label Nodes

```bash
$ kubectl label no --all --output=json quobyte=registry
$ kubectl get no
```

## Run the Registry

```bash
$ kubectl create -f quobyte/quobyte-registry-ds.yaml
$ kubectl create -f quobyte/quobyte-registry-svc.yaml

# Check Registry
$ kubectl describe ds quobyte-registry
$ kubectl describe svc quobyte-registry
```

## Run the web console and api server

```bash
$ kubectl create -f quobyte/quobyte-webconsole-api-rc.yaml
$ kubectl create -f quobyte/quobyte-webconsole-api-svc.yaml

# check
$ kubectl describe rc quobyte-webconsole-api
$ kubectl describe svc quobyte-webconsole-api
```

## Run the Client

```bash
$ kubectl create -f quobyte/quobyte-client-ds.yaml
```

## Todo

- [x] Registry as DaemonSet
- [x] Registry service
- [x] Node Selector -> "quobyte=registry"
- [x] Node set label -> "quobyte=registry"
- [x] Client as DaemonSet (on all nodes)
- [x] Set hostnames in vagrant
- [x] api, webconsole -> as a Pod in RC
- [x] service für api + webconsole (k8s api access?)
- [x] meta-data
- [x] data
- [ ] Kubernetes namespace "quobyte"
- [ ] add liveness probe(s)
- [ ] add readiness probe(s)
- [ ] restart api/webconsole when service is already running (stucks)
- [ ] proxy Webconsole for access

## Nice to have

- [ ] [Helm](http://helm.sh) package
- [ ] move qmboostrap/qmkdev into scope -> Run Kubernetes Job before starting DaemonSet (needs to create the /dev1 dir and mv the setup file)
- [ ] Quobyte and Rocket (test)?
- [ ] check needed capabilities (don't run as privileged)
