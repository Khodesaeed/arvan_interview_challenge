output "longhorn_release_name" {
  value = helm_release.longhorn.name
}

output "longhorn_namespace" {
  value = helm_release.longhorn.namespace
}

output "longhorn_storage_class" {
  value = kubernetes_storage_class.longhorn[0].metadata[0].name
  description = "The Longhorn storage class name (if created)"
}
