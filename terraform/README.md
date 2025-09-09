### **📄 Terraform Project Documentation: Automated Kubernetes Infrastructure on ArvanCloud**

This document provides a comprehensive overview of the multi-part Terraform project designed to provision and configure a complete Kubernetes environment. The project is structured into three distinct modules, each responsible for a specific layer of the stack, from base infrastructure to the GitOps application management layer.

---

### **Module 1: Arvan Cloud Infrastructure Provisioning (`arvan-provision`) ☁️**

This foundational module is responsible for provisioning the necessary cloud infrastructure on ArvanCloud to serve as the base for the Kubernetes cluster.

**Objective:** To automate the creation of virtual machines and networking resources required for a Kubernetes cluster installation.

**Key Components & Resources:**

* **Networking 🌐:** A dedicated private network is established to ensure secure communication between cluster nodes.
    * It creates a resource named `tf_private_network` [cite: 32] with a CIDR of `10.255.255.0/24`[cite: 32].
* **Virtual Machines 🖥️:** The module provisions three virtual machines (`arvan_abrak`) which will serve as the cluster nodes[cite: 33].
    * The configuration specifies resource timeouts, including a read timeout of 10 minutes[cite: 33].
    * The VM image is dynamically selected based on the chosen distribution (`ubuntu`) and release (`22.04`)[cite: 3].
* **Configuration & Variables ⚙️:**
    * **Provider:** Utilizes the `arvancloud/iaas` provider, version `0.8.1`[cite: 8].
    * **Required Inputs:** Users must provide an `api_key` [cite: 3] and `ssh_key_name`[cite: 3]. An example `.tfvars` file is provided for these values[cite: 9].
    * **Customization:** The region (`region`) and instance plan (`chosen_plan_id`) can be customized as needed[cite: 3, 4].
* **Outputs 📋:**
    * The module generates a `kubespray_inventory` output[cite: 22]. This output is specifically formatted for use with Kubespray (an Ansible-based Kubernetes installer).
    * It assigns specific roles to the nodes, designating `node1` for the `kube_control_plane` [cite: 22] and `etcd` roles [cite: 23], and `node2` and `node3` for the `kube_node` role[cite: 24].
    * The inventory correctly maps node names to their IP addresses [cite: 20] and sets the `ansible_user` to `ubuntu`[cite: 22].

---

### **Module 2: Kubernetes Storage Provisioning (`storage-provider`) 💾**

Once the Kubernetes cluster is operational, this module deploys a persistent storage solution, ensuring stateful applications can run reliably.

**Objective:** To deploy and configure the Longhorn distributed block storage system within the Kubernetes cluster.

**Implementation Details:**

* **Helm-based Deployment  CHART:** Longhorn is deployed using its official Helm chart, managed by the `helm_release` resource[cite: 17].
    * The release is installed in the `longhorn-system` namespace by default[cite: 5].
* **Kubernetes StorageClass 🏅:** The module can automatically create a `StorageClass` resource named `longhorn`[cite: 18].
    * This enables dynamic volume provisioning for applications.
    * The `volume_binding_mode` is set to `WaitForFirstConsumer`[cite: 5, 6], which delays volume binding until a pod is scheduled, allowing for more intelligent placement decisions.
* **Configuration & Variables ⚙️:**
    * **Providers:** This module depends on the `hashicorp/helm` (version `~> 3.0.2`) and `hashicorp/kubernetes` (version `~> 2.38.0`) providers[cite: 16].
    * **Customization:** Key variables include `release_name` [cite: 5], `namespace` [cite: 5], `chart_version` [cite: 5], and a boolean `create_storage_class` [cite: 5] to control the creation of the StorageClass resource.
* **Outputs 📋:**
    * The module outputs the Helm release name, the namespace, and the name of the created StorageClass for easy reference[cite: 19].

---

### **Module 3: GitOps Bootstrapping (`argocd`) 🔄**

This module installs and configures Argo CD, implementing the "App of Apps" pattern to bootstrap the GitOps-managed environment.

**Objective:** To deploy Argo CD and establish a single root application that manages all other applications within the cluster.

**Implementation Details:**

* **Argo CD Installation 🚀:**
    * A dedicated `argocd` namespace is created[cite: 10].
    * The `argo-cd` Helm chart (version `8.3.5`) is deployed using the `helm_release` resource[cite: 10].
    * The admin password for Argo CD must be supplied as a Bcrypt hash via the `argo_password` variable[cite: 10, 14].
* **The App of Apps Pattern 🧩:**
    * A single Argo CD `Application` custom resource, named `root-app`, is created via a `kubernetes_manifest` resource[cite: 10, 11].
    * This `root-app` points to a path (`gitops/argo-apps`) within a Git repository[cite: 12]. This specific path is expected to contain the manifests for all other Argo CD applications.
* **Automated Sync Policy 🩹:**
    * The root application is configured with an automated synchronization policy[cite: 12].
    * `prune = true`: Resources that are removed from the Git repository will be automatically removed from the cluster[cite: 12].
    * `selfHeal = true`: Argo CD will automatically correct any detected drift between the live cluster state and the desired state in Git[cite: 13].
* **Configuration & Variables ⚙️:**
    * **Providers:** This module also relies on the `hashicorp/helm` (version `~> 3.0.2`) and `hashicorp/kubernetes` (version `~> 2.38.0`) providers[cite: 7].
    * **Inputs:** Requires the `kube_config_path` and the Bcrypt hash for the `argo_password`[cite: 14].

---

### **Project-wide Provider Dependencies & Versioning 🔗**

The project maintains strict versioning for its Terraform providers to ensure predictable and repeatable deployments. The versions are locked across all modules.

* **`terraform.arvancloud.ir/arvancloud/iaas`**
    * Constraint: `0.8.1` [cite: 8]
    * Version Used: `0.8.1` [cite: 2]
* **`registry.terraform.io/hashicorp/helm`**
    * Constraint: `~> 3.0.2` [cite: 7, 16]
    * Version Used: `3.0.2` [cite: 26, 29]
* **`registry.terraform.io/hashicorp/kubernetes`**
    * Constraint: `~> 2.38.0` [cite: 7, 16]
    * Version Used: `2.38.0` [cite: 26, 29]

*Note: The `.terraform.lock.hcl` files ensure that the exact same provider versions are used on every execution.* [cite: 1, 25, 28]