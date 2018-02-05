# Quobyte operator
Deploys quobyte clients on to the nodes configured in Custom Resuource QuobyteClient.

* [Deploy clients](#deploy-clients)
* [Build from source](#build-source)

Deploy clients:
-------------
1. Deploy quobyte namespace.  
```
cd ../deploy
kubectl create -f quobyte-ns.yaml
```
2. Deploy client daemonset in quobyte namespace(daemonset name must be client).
```
 kubectl create -f client-ds.yaml
```
3. Create required RBAC for quobyte operator.
```
cd ../operator
kubectl create -f operator-rbac.yaml
```
4. Make sure quobyte-operator pod is deployed and running.
```
kubectl -n quobyte get po
```
5. Deploy QuobyteClient (see the examples/quobyte-client-example.yaml).
```
kubectl create -f examples/quobyte-client-example.yaml
```

Build from source:
-----------------
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
8. Edit ``operator-rbac.yaml`` and point ``quobyte-operator`` container image to the docker image.
9. Follow [Deploy clients](#deploy-clients)
