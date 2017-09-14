#!/bin/bash

ansible-playbook kubespray/reset.yml -i inventory/vagrant -b -v
