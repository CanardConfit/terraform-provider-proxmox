---
layout: page
title: proxmox_virtual_environment_disk
parent: Resources
subcategory: Virtual Environment
---

# Resource: proxmox_virtual_environment_disk

Manages a disk.

## Example Usage

```hcl
resource "proxmox_virtual_environment_disk" "my_disk" {
  node_id    = "first-node"
  storage_id = "local-lvm"
  suffix     = "terraform"
  size       = "2G"
}
```

## Argument Reference

- `suffix` - (Required) The name of the file to create.
- `node_id` - (Required) The cluster node name.
- `storage_id` - (Required) The storage identifier.
- `size` - (Required) Size in kilobyte (1024 bytes). Optional suffixes 'M' (megabyte, 1024K) and 'G' (gigabyte, 1024M). E.g. `1G`, `4096`, `300M`.
- `vm_id` - (Optional) Specify owner VM. (defaults to `999`)
- `format` - (Optional) Format of the disk. One of `raw`, `qcow2`, or `subvol`. (defaults to `raw`)

## Attribute Reference

- `path` - The path of the disk within the datastore.
- `space_used` - Space used on the disk.

## Import

Instances can be imported using the `node_id`, `storage_id`, and the `volume_id`, e.g.,

```bash
terraform import proxmox_virtual_environment_disk.my_disk my-node:my-storage:vm-999-myvolume
```
