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
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/v4address"
)

var ErrIDRequiredToLoad = errors.New("can not load object with empty id")

func resourceV4Address() *schema.Resource {
	return &schema.Resource{
		Description: "Managing an IPv4 address object in QIP.",

		CreateContext: resourceV4AddressCreate,
		ReadContext:   resourceV4AddressRead,
		UpdateContext: resourceV4AddressUpdate,
		DeleteContext: resourceV4AddressDelete,

		Schema: schemaV4Address(false),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceV4AddressCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*terraformClient) //nolint:forcetypeassert

	//nolint:forcetypeassert
	var (
		err         error
		address     = d.Get("address").(string)
		subnet      = d.Get("subnet").(string)
		name        = d.Get("name").(string)
		description = d.Get("description").(string)
		class       = d.Get("object_class").(string)
		domain      = d.Get("domain_name").(string)
		rangeStart  = d.Get("subnet_range_start").(string)
		rangeEnd    = d.Get("subnet_range_end").(string)
	)

	if subnet == "" {
		return diag.Errorf("subnet must be set")
	}

	var addressIsSelected bool

	if address == "" {
		var addressRange *v4address.SelectedAddrRange

		if rangeStart != "" && rangeEnd != "" {
			addressRange = &v4address.SelectedAddrRange{
				StartAddress: rangeStart,
				EndAddress:   rangeEnd,
			}
		}

		selectedAddress, err := v4address.CreateSelected(client.QIPClient, subnet, addressRange)
		if err != nil {
			return diag.FromErr(err)
		}

		tflog.Trace(ctx, "Got V4Address selected: "+selectedAddress)

		address = selectedAddress
		addressIsSelected = true
	}

	addr := &v4address.V4Address{
		ObjectAddr:  address,
		SubnetAddr:  subnet,
		ObjectName:  name,
		ObjectClass: class,
		ObjectDesc:  description,
		DomainName:  domain,
	}

	if addressIsSelected {
		err = v4address.Update(client.QIPClient, addr)
	} else {
		err = v4address.Create(client.QIPClient, addr)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(addr.ObjectAddr)

	diags := resourceV4AddressRead(ctx, d, meta)
	if diags != nil {
		return diags
	}

	tflog.Trace(ctx, "Created V4Address "+d.Id())

	return nil
}

func resourceV4AddressLoad(_ context.Context, d *schema.ResourceData, meta any) (*v4address.V4Address, error) {
	if d.Id() == "" {
		return nil, ErrIDRequiredToLoad
	}

	addr, err := v4address.Load(meta.(*terraformClient).QIPClient, d.Id())
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return addr, nil
}

func resourceV4AddressRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	addr, err := resourceV4AddressLoad(ctx, d, meta)
	if err != nil {
		var notFoundErr *qip.HTTPNotFoundError

		if errors.As(err, &notFoundErr) {
			// Object is not found, so reset Id and return no error
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	// Update state from object
	err = d.Set("address", addr.ObjectAddr)
	if err != nil {
		return diag.FromErr(err)
	}

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

//nolint:forcetypeassert
func resourceV4AddressUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	addr, err := resourceV4AddressLoad(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if address changed
	if d.Get("address").(string) != d.Id() {
		return diag.Errorf("address can not be changed after object was created")
	}

	addr.ObjectName = d.Get("name").(string)
	addr.ObjectDesc = d.Get("description").(string)
	addr.ObjectClass = d.Get("object_class").(string)
	addr.DomainName = d.Get("domain_name").(string)

	err = v4address.Update(meta.(*terraformClient).QIPClient, addr)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceV4AddressDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*terraformClient) //nolint:forcetypeassert

	if d.Id() == "" {
		return diag.Errorf("can not delete V4Address with empty id")
	}

	err := v4address.Delete(client.QIPClient, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
