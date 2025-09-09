# ⚙️ Ansible Playbook for Kubernetes Cluster Installation with Kubespray

This directory contains the Ansible playbooks and configurations used to provision a robust Kubernetes cluster using Kubespray. The setup is highly customized to utilize modern, high-performance components, including Cilium as the CNI and Containerd as the container runtime.

## ✨ Key Features

* **Kubernetes Version**: Deploys Kubernetes `v1.28`.
* **Container Runtime**: Uses `containerd` as the container runtime interface (CRI).
* **Network Plugin (CNI)**: Implements **Cilium** for CNI, leveraging eBPF for high-performance networking, security, and observability.
* **Kube-Proxy Replacement**: Cilium is configured to replace `kube-proxy`, which optimizes network performance by using eBPF instead of iptables.
* **Node IPAM Load Balancing**: The `Node IPAM` feature is enabled to expose `LoadBalancer` type services on the node's public IP addresses, providing a simple yet powerful way to expose applications externally.
* **Cluster Structure**: The cluster is composed of one control plane node and two worker nodes for a simple and scalable setup.

## 📁 Files and Configuration

### `inventory/mycluster/hosts.yaml`

This file defines the cluster topology, assigning roles to each virtual machine.

```yaml
all:
  children:
    etcd:
      hosts:
        node1: {}
    k8s_cluster:
      children:
        kube_control_plane: {}
        kube_node: {}
    kube_control_plane:
      hosts:
        node1: {}
    kube_node:
      hosts:
        node2: {}
        node3: {}
  hosts:
    node1:
      ansible_host: 185.226.118.149
      ansible_user: ubuntu
    node2:
      ansible_host: 185.231.180.48
      ansible_user: ubuntu
    node3:
      ansible_host: 185.231.183.180
      ansible_user: ubuntu
```

### `inventory/mycluster/group_vars/k8s_cluster/k8s-cluster.yml`

This file contains the core cluster configuration variables, including the selection of Cilium as the CNI and the activation of key features.

```yaml
# Choose CNI plugin
kube_network_plugin: cilium

cilium_kube_proxy_replacement: "true"
cilium_k8s_service_host: 185.226.118.149
cilium_k8s_service_port: 6443
# Optional: disable others
kube_network_plugin_multus: false

kube_apiserver_cert_extra_sans:
  - 185.226.118.149

# Cilium config overrides
cilium_enable_ipv4: true
cilium_enable_ipv6: false
cilium_hubble_enabled: true

# Custom Cilium config (passed into ConfigMap)
cilium_config_extra_vars:
  enable-node-ipam: "true"
  default-lb-service-ipam: nodeipam

cilium_host_network: true
```

### `kubespray.sh`

This script provides a convenient way to run the Kubespray playbook using a Docker container, ensuring all dependencies are met without affecting your local environment.

```bash
#!/bin/bash

docker run --rm -it \
  -v "$(pwd)/inventory":/kubespray/inventory \
  -v "$HOME/.ssh/arvan_interview":/root/.ssh/id_rsa \
  -v "$HOME/.ssh/known_hosts":/root/.ssh/known_hosts \
  -v "$(pwd)/kubespray_cache":/tmp/kubespray_cache \
  -e "ANSIBLE_CONFIG=/kubespray/ansible.cfg" \
  quay.io/kubespray/kubespray:v2.28.1  \
  ansible-playbook -i inventory/mycluster/hosts.yaml cluster.yml -e ansible_ssh_private_key_file=/root/.ssh/id_rsa -e host_key_checking=False --become
```

## 🚀 Getting Started

1. **Configure Inventory**: Ensure the `hosts.yaml` file in `inventory/mycluster` correctly lists the `ansible_host` and `ansible_user` for your virtual machines.
1. **Run the Playbook**: Execute the `kubespray.sh` script to start the installation. This script will handle all the necessary steps, from setting up the `containerd` runtime to configuring Cilium.

```bash 
./kubespray.sh
```

This will deploy your Kubernetes cluster with the specified Cilium CNI and kube-proxy replacement.
