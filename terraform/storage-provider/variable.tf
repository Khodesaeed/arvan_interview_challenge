variable "kube_config_path" {
  type    = string
  default = "~/.kube/config"
}

variable "release_name" {
  type    = string
  default = "longhorn"
}

variable "namespace" {
  type    = string
  default = "longhorn-system"
}

variable "chart_version" {
  type    = string
  default = "1.9.1"
}

variable "create_storage_class" {
  type    = bool
  default = true
}

variable "storage_class_name" {
  type    = string
  default = "longhorn"
}

variable "reclaim_policy" {
  type    = string
  default = "Delete"
}

variable "volume_binding_mode" {
  type    = string
  default = "WaitForFirstConsumer"
}
