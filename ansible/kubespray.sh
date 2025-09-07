#!/bin/bash


docker run --rm -it \
  -v "$(pwd)/inventory":/kubespray/inventory \
  -v "$HOME/.ssh/arvan_interview":/root/.ssh/id_rsa \
  -v "$HOME/.ssh/known_hosts":/root/.ssh/known_hosts \
  -v "$(pwd)/kubespray_cache":/tmp/kubespray_cache \
  -e "ANSIBLE_CONFIG=/kubespray/ansible.cfg" \
  quay.io/kubespray/kubespray:v2.28.1  \
  ansible-playbook -i inventory/mycluster/hosts.yaml cluster.yml -e ansible_ssh_private_key_file=/root/.ssh/id_rsa -e host_key_checking=False --become --tags cilium 
