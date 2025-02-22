locals {
  datastore_id = element(data.proxmox_virtual_environment_datastores.example.datastore_ids, index(data.proxmox_virtual_environment_datastores.example.datastore_ids, "local-lvm"))
}

resource "proxmox_virtual_environment_disk" "my_disk" {
  node_id    = data.proxmox_virtual_environment_nodes.example.names[0]
  storage_id = local.datastore_id
  suffix     = "terraform"
  size       = "2G"
}
