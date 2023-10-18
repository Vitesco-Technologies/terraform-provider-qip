/*
Copyright 2023 Vitesco Technologies Group AG

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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/v4address"
)

func dataSourceV4Address() *schema.Resource {
	return &schema.Resource{
		Description: "IPv4 address object in QIP.",

		ReadContext: dataSourceV4AddressRead,

		Schema: schemaV4Address(true),
	}
}

func dataSourceV4AddressRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*terraformClient) //nolint:forcetypeassert

	addr, err := v4address.Load(client.QIPClient, d.Get("address").(string))
	if err != nil {
		return diag.Errorf("could not find IPv4 object: %s", err)
	}

	d.SetId(addr.ObjectAddr)

	err = d.Set("subnet", addr.SubnetAddr)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", addr.ObjectName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("description", addr.ObjectDesc)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("object_class", addr.ObjectClass)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("domain_name", addr.DomainName)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
