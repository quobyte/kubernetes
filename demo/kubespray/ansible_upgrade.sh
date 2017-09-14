!/bin/bash

ansible-playbook kubespray/upgrade-cluster.yml -i inventory/vagrant -b -v
