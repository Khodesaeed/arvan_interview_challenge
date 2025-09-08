variable "kube_config_path" {
  type    = string
  default = "~/.kube/config"
}

variable "argo_password" {
  type        = string
  description = "Bcrypt hash of the ArgoCD admin password"
  # Example generated with:
  # htpasswd -nbBC 10 "" password | tr -d ':\n' | sed 's/$2y/$2a/'
}
