/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkResourceVirtualEnvironmentDiskId        = "id"
	mkResourceVirtualEnvironmentDiskSuffix    = "suffix"
	mkResourceVirtualEnvironmentDiskName      = "name"
	mkResourceVirtualEnvironmentDiskNode      = "node_id"
	mkResourceVirtualEnvironmentDiskStorage   = "storage_id"
	mkResourceVirtualEnvironmentDiskSize      = "size"
	mkResourceVirtualEnvironmentDiskVmId      = "vm_id"
	mkResourceVirtualEnvironmentDiskFormat    = "format"
	mkResourceVirtualEnvironmentDiskPath      = "path"
	mkResourceVirtualEnvironmentDiskSpaceUsed = "space_used"
	mkResourceVirtualEnvironmentDiskSizeGb    = "size_gb"
	mkResourceVirtualEnvironmentDiskSizeMb    = "size_mb"
	mkResourceVirtualEnvironmentDiskSizeBytes = "size_bytes"
)

// Disk returns a resource that manages disks.
func Disk() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentDiskId: {
				Type:        schema.TypeString,
				Description: "ID of the disk in the format <node>:<datastore>:vm-<vmid>-<name>.",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentDiskSuffix: {
				Type:        schema.TypeString,
				Description: "The name of the file to create.",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentDiskName: {
				Type:        schema.TypeString,
				Description: "Generated disk name.",
				Optional:    true,
				Computed:    true,
			},
			mkResourceVirtualEnvironmentDiskNode: {
				Type:        schema.TypeString,
				Description: "The cluster node name.",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentDiskStorage: {
				Type:        schema.TypeString,
				Description: "The storage identifier.",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentDiskSize: {
				Type:        schema.TypeString,
				Description: "Size in kilobyte (1024 bytes). Optional suffixes 'M' (megabyte, 1024K) and 'G' (gigabyte, 1024M).",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentDiskSizeGb: {
				Type:        schema.TypeInt,
				Description: "Disk size in gigabytes",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentDiskSizeMb: {
				Type:        schema.TypeInt,
				Description: "Disk size in megabytes",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentDiskSizeBytes: {
				Type:        schema.TypeInt,
				Description: "Disk size in bytes",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentDiskVmId: {
				Type:        schema.TypeInt,
				Description: "Specify owner VM.",
				Optional:    true,
				Default:     999,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentDiskFormat: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "raw",
			},
			mkResourceVirtualEnvironmentDiskPath: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			mkResourceVirtualEnvironmentDiskSpaceUsed: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
		CreateContext: diskCreate,
		ReadContext:   diskRead,
		UpdateContext: diskUpdate,
		DeleteContext: diskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				config := m.(proxmoxtf.ProviderConfiguration)

				client, err := config.GetClient()
				if err != nil {
					return nil, err
				}

				diskID := d.Id()
				parts := strings.SplitN(diskID, ":", 3)

				d.Set(mkResourceVirtualEnvironmentDiskId, diskID)
				d.Set(mkResourceVirtualEnvironmentDiskNode, parts[0])
				d.Set(mkResourceVirtualEnvironmentDiskStorage, parts[1])
				d.Set(mkResourceVirtualEnvironmentDiskName, parts[2])

				disk, err := client.Node(parts[0]).Storage(parts[1]).GetDatastoreFile(ctx, parts[2])
				if err != nil {
					return nil, err
				}

				parts = strings.SplitN(parts[2], "-", 3)
				vmid, _ := strconv.Atoi(parts[1])
				d.Set(mkResourceVirtualEnvironmentDiskVmId, vmid)
				d.Set(mkResourceVirtualEnvironmentDiskSuffix, parts[2])

				sizeM := (*disk.FileSize) / 1024 / 1024
				sizeG := sizeM / 1024

				var size string
				if sizeG*1024*1024*1024 == *disk.FileSize {
					size = fmt.Sprintf("%dG", sizeG)
				} else if sizeM*1024*1024 == *disk.FileSize {
					size = fmt.Sprintf("%dM", sizeG)
				} else {
					size = strconv.Itoa(int(*disk.FileSize))
				}

				d.Set(mkResourceVirtualEnvironmentDiskSizeBytes, *disk.FileSize)
				d.Set(mkResourceVirtualEnvironmentDiskSizeGb, sizeG)
				d.Set(mkResourceVirtualEnvironmentDiskSizeMb, sizeM)
				d.Set(mkResourceVirtualEnvironmentDiskFormat, disk.FileFormat)
				d.Set(mkResourceVirtualEnvironmentDiskSize, size)
				d.Set(mkResourceVirtualEnvironmentDiskPath, disk.Path)
				d.Set(mkResourceVirtualEnvironmentDiskSpaceUsed, disk.SpaceUsed)

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func diskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	suffix := d.Get(mkResourceVirtualEnvironmentDiskSuffix).(string)
	nodeId := d.Get(mkResourceVirtualEnvironmentDiskNode).(string)
	storageId := d.Get(mkResourceVirtualEnvironmentDiskStorage).(string)
	size := d.Get(mkResourceVirtualEnvironmentDiskSize).(string)
	vmId := d.Get(mkResourceVirtualEnvironmentDiskVmId).(int)
	tmpFormat, isFormatSet := d.GetOk(mkResourceVirtualEnvironmentDiskFormat)

	var format *string
	format = nil
	if isFormatSet {
		t := tmpFormat.(string)
		format = &t
	}

	name := fmt.Sprintf("vm-%d-%s", vmId, suffix)
	d.Set(mkResourceVirtualEnvironmentDiskName, name)

	body := &storage.DatastoreFileCreateRequest{
		Filename:   name,
		NodeID:     nodeId,
		StorageID:  storageId,
		FileSize:   size,
		VMID:       vmId,
		FileFormat: format,
	}

	vid, err := client.Node(nodeId).Storage(storageId).CreateDatastoreFile(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%s:%s", nodeId, *vid)
	d.SetId(id)
	return diskRead(ctx, d, m)
}

func diskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeId := d.Get(mkResourceVirtualEnvironmentDiskNode).(string)
	storageId := d.Get(mkResourceVirtualEnvironmentDiskStorage).(string)
	volumeId := d.Get(mkResourceVirtualEnvironmentDiskName).(string)

	disk, err := client.Node(nodeId).Storage(storageId).GetDatastoreFile(ctx, volumeId)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	d.Set(mkResourceVirtualEnvironmentDiskSizeBytes, *disk.FileSize)
	d.Set(mkResourceVirtualEnvironmentDiskSizeMb, *disk.FileSize/1024/1024)
	d.Set(mkResourceVirtualEnvironmentDiskSizeGb, *disk.FileSize/1024/1024/1024)

	d.Set(mkResourceVirtualEnvironmentDiskPath, disk.Path)
	d.Set(mkResourceVirtualEnvironmentDiskSpaceUsed, disk.SpaceUsed)
	return diag.FromErr(err)
}

func diskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Errorf("Cannot update a disk in-place")
}

func diskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeId := d.Get(mkResourceVirtualEnvironmentDiskNode).(string)
	storageId := d.Get(mkResourceVirtualEnvironmentDiskStorage).(string)
	volumeId := d.Get(mkResourceVirtualEnvironmentDiskName).(string)

	err = client.Node(nodeId).Storage(storageId).DeleteDatastoreFile(ctx, volumeId)
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
