#!/bin/bash

ansible-playbook kubespray/cluster.yml -i inventory/vagrant -b -v
