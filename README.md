Scripts and tools for deploying Quobyte on Kubernetes
=====================================================

**DEPRECATED in favour of [Quobyte CSI](https://github.com/quobyte/quobyte-csi)**

Currently contains:
 * **operator** The recommended way to run Quobyte clients and servers in your k8s cluster.
 * **deploy** Deployment files for Quobyte on Kubernetes
 * **volumes** Examples how to use PersistentVolumes, PersistentVolumeClaims, and StorageClassses with Quobyte
 * **tools** device initialization tools: **qbootstrap** and **qmkdev**

Please see the quick start guides for instructions:

 * Manual deployment of Quobyte [clients](client_quick_setup.md) and [services](server_quick_setup.md).
 * [Operator based setup for Quobyte clients and services](operator/README.md)

These blog posts will give you some more insights and ideas how to benefit from these integrations:
* [“Quobernetes” or How to Enable Simple and Efficient Storage Operations for Kubernetes](https://www.quobyte.com/blog/2018/08/27/quobernetes-is-kubernetes-with-quobyte-software-storage/)
* [The State of Secure Storage Access in Container Infrastructures](https://www.quobyte.com/blog/2017/03/17/the-state-of-secure-storage-access-in-container-infrastructures/)
