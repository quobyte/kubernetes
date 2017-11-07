#!/bin/bash

ansible-playbook ./kubespray/kubespray/cluster.yml -i ./kubespray/inventory/vagrant -b -v
