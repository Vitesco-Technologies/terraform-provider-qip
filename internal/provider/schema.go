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
	"net"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MaxObjectDescriptionLength = 32

func schemaV4Address(forData bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"address": {
			Description: "IPv4 address.",
			Type:        schema.TypeString,
			Required:    forData,
			Optional:    !forData,
			Computed:    !forData,
		},
		"subnet": {
			Description: "Subnet of the IPv4 address.",
			Type:        schema.TypeString,
			Required:    !forData,
			Computed:    forData,
		},
		"name": {
			Description: "Hostname for the address.",
			Type:        schema.TypeString,
			Required:    !forData,
			Computed:    forData,
		},
		"description": {
			Description: "Description for the address.",
			Type:        schema.TypeString,
			Optional:    !forData,
			Computed:    forData,
		},
		"object_class": {
			Description: "Object class for the address. Must be known by the QIP server.",
			Type:        schema.TypeString,
			Optional:    !forData,
			Computed:    forData,
			Default:     ifSet(!forData, "Virtualized Server"),
		},
		"domain_name": {
			Description: "DNS Zone of the address.",
			Type:        schema.TypeString,
			Optional:    !forData,
			Computed:    true,
		},
	}

	if !forData {
		// Add schema entries only for the resource
		s["address"].ValidateDiagFunc = validateIPV4Address
		s["subnet"].ValidateDiagFunc = validateIPV4Address

		s["description"].DiffSuppressFunc = func(_, oldValue, newValue string, _ *schema.ResourceData) bool {
			// Do not change a value when the description length is larger than MaxObjectDescriptionLength
			// and the non exceeding characters are equal.
			if len(newValue) > MaxObjectDescriptionLength {
				shortenedValue := newValue[0:MaxObjectDescriptionLength]
				if oldValue == shortenedValue {
					return true
				}
			}

			return false
		}

		s["subnet_range_start"] = &schema.Schema{
			Description:      "Starting address of a range to select a free IPv4 address from. Will be passed to QIP.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validateIPV4Address,
		}

		s["subnet_range_end"] = &schema.Schema{
			Description:      "Ending address of a range to select a free IPv4 address from. Will be passed to QIP.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validateIPV4Address,
		}
	}

	return s
}

func validateIPV4Address(value interface{}, _ cty.Path) diag.Diagnostics {
	address, ok := value.(string)
	if !ok {
		return diag.Errorf("value is not a string")
	}

	ip := net.ParseIP(address)
	if ip == nil || ip.To4() == nil {
		return diag.Errorf("value is not an IPv4 address")
	}

	return nil
}

func ifSet(condition bool, value any) any {
	if condition {
		return value
	}

	return nil
}
