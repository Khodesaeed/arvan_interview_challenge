# Arvan Interview Challenge

This repository presents a comprehensive and professional solution for the Arvan interview challenge. The project demonstrates a modern **Infrastructure as Code (IaC)**, **configuration management**, and **GitOps** approach to application deployment, ensuring a reliable, automated, and reproducible **CI/CD pipeline**.

## 🚀 Key Features

* **IaC**: Provision cloud resources on **Arvan Cloud** using **Terraform**.
* **Configuration Management**: Automate server setup and application deployment with **Ansible**.
* **GitOps**: Utilize a **Git-driven** continuous deployment pipeline to automate the deployment process.
* **Containerization**: Leverage **Docker** and **Kubernetes** for application containerization and orchestration.
* **Technology Stack**: The application is built with **Go** (Golang).

## 📁 Repository Structure

The project is logically divided into the following components:

* **`app`**: Contains the source code for the Go-based application.
* **`terraform`**: Includes Terraform configurations for provisioning the required infrastructure on Arvan Cloud.
* **`ansible`**: Holds Ansible playbooks for server configuration and Kubernetes cluster setup.
* **`gitops`**: Manages the continuous deployment pipeline using **ArgoCD** and **GitHub Actions** as the single source of truth.

For a detailed breakdown of each component, please refer to their respective README files:

* [Application README](./app/README.md)
* [Terraform README](./terraform/README.md)
* [Ansible README](./ansible/README.md)
* [GitOps README](./gitops/README.md)

## ⚙️ Getting Started

To deploy this project, follow the sequential steps outlined below. Each step has detailed instructions within its respective directory's README file.

### Step 1: Provision Infrastructure with Terraform

Start by provisioning the foundational cloud infrastructure. Navigate to the `terraform/arvan-provision` directory and apply the Terraform configuration. This will provision three virtual machines within a private network on Arvan Cloud.

```bash
cd terraform/arvan-provision
terraform apply
```

### Step 2: Install Kubernetes Cluster with Ansible

Next, configure the provisioned virtual machines and set up the Kubernetes cluster using Ansible. This process leverages the Kubespray project for a robust and standardized cluster installation.

```bash
cd ansible
./kubespray.sh
```

### Step 3: Implement Storage Provider with Terraform

Once the cluster is up, install a persistent storage solution. Apply the Terraform configuration in the `terraform/storage-provider` directory to deploy Longhorn as the storage provider and create the necessary storage classes.

```bash
cd terraform/storage-provider
terraform apply
```

### Step 4: Deploy the Application with GitOps

Finally, deploy the application layer using ArgoCD and a GitOps approach. Apply the Terraform configuration in the `terraform/argocd` directory to install the ArgoCD Helm chart. This setup utilizes an App-of-Apps pattern for managing and deploying the application layer automatically.

```bash
cd terraform/argocd
terraform apply
```
