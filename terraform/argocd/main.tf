resource "kubernetes_namespace" "argocd" {
  metadata {
    name = "argocd"
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  namespace  = kubernetes_namespace.argocd.metadata[0].name
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "8.3.5"

  set = [
    {
      name  = "configs.secret.argocdServerAdminPassword"
      value = var.argo_password
    }
  ]
}

resource "kubernetes_manifest" "argocd_root_app" {
  manifest = {
    "apiVersion" = "argoproj.io/v1alpha1"
    "kind"       = "Application"
    "metadata" = {
      "name"      = "root-app"
      "namespace" = "argocd"
      "finalizers" = ["resources-finalizer.argocd.argoproj.io"]
    }
    "spec" = {
      "project" = "default"
      "source" = {
        "repoURL"        = "https://github.com/Khodesaeed/arvan_interview_challenge.git"
        "targetRevision" = "main"
        "path"           = "gitops/argo-apps" # this folder contains other Application manifests
      }
      "destination" = {
        "server"    = "https://kubernetes.default.svc"
        "namespace" = "argocd"
      }
      "syncPolicy" = {
        "automated" = {
          "prune"    = true
          "selfHeal" = true
        }
      }
    }
  }
  depends_on = [helm_release.argocd]
}
