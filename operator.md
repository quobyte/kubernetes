# Prerequisites
- Kubernetes 1.9 with mountPropagation feature enabled.
- Vanilla Kubernetes 1.10
- Unformatted storage devices (use wipefs to clean devices)

# Deploy Operator
The Quobyte operator helps to run a full Quobyte cluster in kubernetes, or to
operate Quobyte clients to provide persisted volumes to pods.

In any case, you first need to create the quobyte namespace and install the operator.
```
kubectl create -f quobyte-ns.yaml
kubectl -n quobyte create -f operator.yaml
kubectl -n quobyte create -f operator-config.yaml
```

# Configure Quobyte Services
To run the Quobyte services in kubernetes, first edit the services-config.yaml file,
and determine which node should run which services.

A new Quobyte instance needs to be bootstrapped. This can be achieved by either formatting
a device and using the qbootstrap tool, or by marking a node as `registry.bootstrap_node`
in the services-config.yaml. A registry pod which is started on that node will
create an ephemeral bootstrap registry, which can be used to start, format,
and persist the new cluster. When persisted registry devices are added, it is safe
and recommended to delete the ephemeral bootstrap registry node.
```
kubectl -n quobyte create -f services.yaml
kubectl -n quobyte create -f services-config.yaml
```

When you update the container or add or remove nodes, reapply the config and
let the operator work.

```
kubectl -n quobyte apply -f services-config.yaml
```

Now you should see all pods running. The non-bootsrap registries will crash-loop
until they find a quobyte registry device. So let's create some quobyte devices.

The webconsole pod is running and you can port forward it like:
```
kubectl -n quobyte port-forward "$(kubectl get po -owide -n quobyte | grep webconsole | awk '{print $1}')" 8080:8080
```

Then point your browser to http://localhost:8080 and follow the installation wizard.


# Configure Quobyte Clients
If you want to manage Quobyte clients only, edit the client-config.yaml,
choose the latest image version and add the names of the nodes, where a
Quobyte client should run, to the `nodes` list.

```
kubectl -n quobyte create -f client.yaml
```

The client config represents the desired state of your cluster. So adjust the
version and the hosts, where the operator should deploy clients.

```
kubectl -n quobyte create -f client-config.yaml
```

Once the client-config is created, you should see pods being started on the
desired hosts.

If you add or remove clients, edit the config and update it like:

```
kubectl -n quobyte apply -f client-config.yaml
```

# Rolling Updates
The operator supports rolling updates. When you change the container version
in the client-config.yaml or services-config.yaml, and rolling updates are enabled, the operator will upgrade one node after the other.

Since restarting a client pod would terminate currently opened files, the operator will taint the node and drain all pods to other nodes. When the
Quobyte client is updated, the node is untainted again.

Quobyte service containers on the other hand, don't require to drain other pods,
but a careful timing between pod restarts, to always ensure availability of
the Quobyte services.

# Uninstall Operator
If you want to remove all services or clients, remove the config files, before
you delete the deployments or the operator. This will terminate the scheduled pods and remove the all labels, which the operator applied to any nodes.
```
kubectl -n quobyte delete -f services-config.yaml
or
kubectl -n quobyte delete -f client-config.yaml
```

# Build Operator from Source

## Requirements
1. golang 1.8+
2. glide for package management
3. docker

## Build
1. Clone the repository.
```
git clone git@github.com:quobyte/kubernetes.git quobyte-kubernetes
```
2. Compile and build binary from source.
```
cd quobyte-kubernetes/operator
export GOPATH=$(pwd)
cd src/operator
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o operator .
```
If you're building for the first time after clone run ``glide install --strip-vendor`` to get the dependencies.

3. To run operator outside cluster (skip to 4 to run operator inside cluster)
```
./operator --kubeconfig <kuberenetes-admin-conf>
```
  Follow [Deploy clients](#deploy-clients), and you can skip step 3 of deploy clients.

4. Build Docker container
```
sudo docker build -t operator -f Dockerfile.scratch .
sudo docker run -it operator
```
5. Get the ``CONTAINER ID`` of the operator.
```
sudo docker ps -l
```
6. Commit the container.
```
sudo docker commit <CONTAINER ID>  <Docker-repository-url>
```
7. Push to the container repository.
```
sudo docker push <Docker-repository-url>
```
8. Edit ``operator.yaml`` and point ``quobyte-operator`` container image to the docker image.
9.
