terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.0.2"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.38.0"
    }
  }
}

provider "kubernetes" {
  config_path = var.kube_config_path
}

provider "helm" {
  kubernetes = {
    config_path = var.kube_config_path
  }
}

resource "helm_release" "longhorn" {
  name             = var.release_name
  namespace        = var.namespace
  create_namespace = true

  repository = "https://charts.longhorn.io"
  chart      = "longhorn"
  version    = var.chart_version

  values = [
    file("${path.module}/longhorn-values.yaml")
  ]
}

# Optional StorageClass for Longhorn
resource "kubernetes_storage_class" "longhorn" {
  count = var.create_storage_class ? 1 : 0

  metadata {
    name = var.storage_class_name
  }

  storage_provisioner = "driver.longhorn.io"
  reclaim_policy = var.reclaim_policy
  volume_binding_mode = var.volume_binding_mode
}
