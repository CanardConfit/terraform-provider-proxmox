/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestDiskInstantiation tests whether the Disk instance can be instantiated.
func TestDiskInstantiation(t *testing.T) {
	t.Parallel()

	s := Disk()
	if s == nil {
		t.Fatalf("Cannot instantiate Disk")
	}
}

// TestDiskSchema tests the Disk schema.
func TestDiskSchema(t *testing.T) {
	t.Parallel()

	s := Disk().Schema

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentDiskSuffix,
		mkResourceVirtualEnvironmentDiskNode,
		mkResourceVirtualEnvironmentDiskSize,
		mkResourceVirtualEnvironmentDiskStorage,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentDiskVmId,
		mkResourceVirtualEnvironmentDiskFormat,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentDiskName,
		mkResourceVirtualEnvironmentDiskId,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentDiskSuffix:  schema.TypeString,
		mkResourceVirtualEnvironmentDiskNode:    schema.TypeString,
		mkResourceVirtualEnvironmentDiskSize:    schema.TypeString,
		mkResourceVirtualEnvironmentDiskStorage: schema.TypeString,
		mkResourceVirtualEnvironmentDiskVmId:    schema.TypeInt,
		mkResourceVirtualEnvironmentDiskFormat:  schema.TypeString,
	})
}
