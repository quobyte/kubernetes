# Quick Start Guide to use Quobyte from OpenShift Origin

Accessing Quobyte volumes from kubernetes is currently partly supported.
We fully support mounting PersistentVolumes, but dynamic provisioning with storage
classes is not yet supported and requires an external storage
provisioner as described
[here](https://docs.openshift.org/latest/install_config/provisioners.html).

## Client setup

In an OpenShift environment, you need to create a service
account which runs the client container.

```bash
$ oc create serviceaccount quobyteclientsrv -n quobyte
$ oc adm policy add-scc-to-user privileged -n quobyte -z quobyteclientsrv
```
Then, uncomment the line
```yaml
serviceAccountName: quobyteclientsrv
```
in client-ds.yaml or client-certificate-ds.yaml.

The remaining steps are the same as found in the [client quick setup](client_quick_setup.md).
