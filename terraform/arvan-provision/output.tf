locals {
  kubespray_nodes = [
    for idx, vm in arvan_abrak.built_by_terraform :
    {
      name = "node${idx + 1}"
      ip   = try(
        vm.public_ip,
        vm.public_ips[0],
        vm.ip,
        vm.address,
        vm.ipv4,
        vm.ip_address,
        vm.networks[0].ip,
        vm.networks[0].private_ip,
        vm.networks[0].ip_address,
        null
      )
    }
  ]
}

output "kubespray_inventory" {
  description = "Kubespray-style inventory with split roles"
  value = {
    all = {
      hosts = {
        for n in local.kubespray_nodes :
        n.name => {
          ansible_host = n.ip
          ansible_user = "ubuntu" # adjust to your VM image
        }
      }
      children = {
        kube_control_plane = {
          hosts = {
            node1 = {}
          }
        }
        etcd = {
          hosts = {
            node1 = {}
          }
        }
        kube_node = {
          hosts = {
            node2 = {}
            node3 = {}
          }
        }
        k8s_cluster = {
          children = {
            kube_control_plane = {}
            kube_node          = {}
          }
        }
      }
    }
  }
}
