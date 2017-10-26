## Setting up Test environment with Vagrant

For a fast demo setup, we use a Vagrant based 4-machine cluster, where each server has additional 3 disk drives attached.

```bash
$ cd demo/vagrant
$ vagrant up
$ vagrant ssh-config
```

We use kubespray to bootstrap and setup the Kubernetes cluster.
We provide an inventory file for the newly created cluster `demo/kubespray/inventory/vagrant`.
Please make sure that the *ansible_port* and *ansible_ssh_private_key_file* match.


If the 4 machines are running and you are able to connect to them like:
```bash
$ cd demo/vagrant
$ vagrant ssh qb1
```
we're good to apply some kubespray.

```bash
$ cd demo/kubespray
$ ./clone_kubespray
$ ./ansible_cluster.sh
```

Make sure that `kubectl` [is installed](https://kubernetes.io/docs/tasks/tools/install-kubectl/ "Install and Set Up kubectl") on your machine.

To configure and use your newly created cluster, you can run:

```bash
$ mkdir -p $HOME/.kube/certs/qb
$ cd demo/vagrant/
$ vagrant ssh qb1 -- -t sudo cat /etc/kubernetes/ssl/admin-qb1.pem > $HOME/.kube/certs/qb/qb-admin.pem
$ vagrant ssh qb1 -- -t sudo cat /etc/kubernetes/ssl/admin-qb1-key.pem > $HOME/.kube/certs/qb/qb-admin-key.pem
$ vagrant ssh qb1 -- -t sudo cat /etc/kubernetes/ssl/ca.pem > $HOME/.kube/certs/qb/qb-ca.pem

$ kubectl config set-credentials qb-admin \
  --certificate-authority=$HOME/.kube/certs/qb/qb-ca.pem \
  --client-key=$HOME/.kube/certs/qb/qb-admin-key.pem \
  --client-certificate=$HOME/.kube/certs/qb/qb-admin.pem
$ kubectl config set-cluster qb --server=https://127.0.0.1:6443 \
  --certificate-authority=$HOME/.kube/certs/qb/qb-ca.pem

$ kubectl config set-context qb --cluster=qb --user=qb-admin
$ kubectl config use-context qb
```

Your cluster should be available now:

```bash
$ kubectl get nodes
NAME      STATUS    AGE       VERSION
qb1       Ready     5m        v1.7.3+coreos.0
qb2       Ready     5m        v1.7.3+coreos.0
qb3       Ready     5m        v1.7.3+coreos.0
qb4       Ready     5m        v1.7.3+coreos.0
```
