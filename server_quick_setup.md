# Quick Start Guide to run Quobyte in Kubernetes

## Prerequisites

Quobyte is designed to run on dedicated disk drives which are formatted with ext4 or xfs and
are initialized as a Quobyte device.
Quobyte 2.0 features the Device Inspector, which can detect unformatted devices
on a node and format and mount them for use in Quobyte. The Device Inspector
requires the [feature gate](https://kubernetes.io/docs/reference/feature-gates/) "MountPropagation=true" to be set.
This feature is in alpha State on Kubernetes 1.8 and requires manual interaction to enable.
If it's not set, the Device Inspector is able to detect and format devices, but they will not be mounted.

If you plan to create the cluster with the help of the Device Inspector, you can
start a directory based ephemeral registry, or prepare and mount a registry device.

Instead of creating any devices, just schedule an ephemeral bootstrap registry
first, start registries, data, and metadata services on the target pods
and let the device inspector create the devices.
Afterwards, it is safe to delete the bootstrap registry again.

First, create the quobyte namespace, and set up a config and services.

```bash
kubectl create -f quobyte-ns.yaml
kubectl -n quobyte create -f config.yaml
kubectl -n quobyte create -f quobyte-services.yaml

kubectl -n quobyte create -f registry-bootstrap-ds.yaml
kubectl -n quobyte create -f registry-ds.yaml
kubectl -n quobyte create -f data-ds.yaml
kubectl -n quobyte create -f metadata-ds.yaml
```

Choose any node to be the initial bootstrap registry.

```bash
kubectl label nodes qb1 quobyte_bootstrap="true"
```
As soon as the registry is up, start the webconsole, api, and qmgmt pods:

```bash
kubectl create -f webconsole-deployment.yaml
kubectl create -f qmgmt-pod.yaml
```

When all pods are up, you should be able to log in to your initial cluster with your preferred browser:
```bash
kubectl port-forward <webconsole-pod> 8080:8080
```

To schedule other registries, data, and metadata servers, label the nodes accordingly.
The services will show up in the webconsole and the Device Inspector will help
to set up the devices.

```bash
kubectl label nodes <node> quobyte_metadata="true"
kubectl label nodes <node> quobyte_registry="true"
kubectl label nodes <node> quobyte_metadata="true"
```

Once you have three other registries set up on physical devices, it's safe
to remove the `quobyte_bootstrap` label and delete the registry pod.
