/*
Copyright 2024 Vitesco Technologies Group AG

SPDX-License-Identifier: Apache-2.0

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/v4subnet"
)

func dataSourceV4Subnet() *schema.Resource {
	return &schema.Resource{
		Description: "IPv4 address object in QIP.",

		ReadContext: dataSourceV4SubnetRead,

		Schema: map[string]*schema.Schema{
			"address": {
				Description:      "IPv4 subnet address.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateIPV4Address,
			},
			"address_cidr": {
				Description: "IPv4 subnet address in CIDR notation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"mask": {
				Description: "IPv4 subnet network mask.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"prefix_length": {
				Description: "IPv4 CIDR prefix length for the subnet.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"network": {
				Description: "IPv4 network address this subnet belongs to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of the V4 subnet.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the V4 subnet.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domains": {
				Description: "List of domains for the V4 subnet.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dns_servers": {
				Description: "List of DNS servers preferred for this subnet.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ntp_servers": {
				Description: "List of NTP time servers preferred for this subnet.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_routers": {
				Description: "List of default routers for this subnet.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceV4SubnetRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*terraformClient) //nolint:forcetypeassert

	subnet, err := v4subnet.Load(client.QIPClient, d.Get("address").(string))
	if err != nil {
		return diag.Errorf("could not find IPv4 object: %s", err)
	}

	d.SetId(subnet.SubnetAddress)

	// Calculate prefix length from netmask
	mask := net.IPMask(net.ParseIP(subnet.SubnetMask).To4())
	prefixLength, _ := mask.Size()

	values := map[string]any{
		"address":         subnet.SubnetAddress,
		"address_cidr":    fmt.Sprintf("%s/%d", subnet.SubnetAddress, prefixLength),
		"mask":            subnet.SubnetMask,
		"prefix_length":   prefixLength,
		"network":         subnet.NetworkAddress,
		"name":            subnet.SubnetName,
		"description":     subnet.SubnetDescription,
		"domains":         subnet.Domains.Name,
		"dns_servers":     subnet.PreferredDNSServers.Name,
		"ntp_servers":     subnet.PreferredTimeServers.Name,
		"default_routers": subnet.DefaultRouters.Name,
	}

	for k, v := range values {
		err = d.Set(k, v)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
