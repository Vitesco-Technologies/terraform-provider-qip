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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/rr"
)

var ErrNonUniqueRR = errors.New("non unique RR found")

func resourceV4AddressRR() *schema.Resource {
	return &schema.Resource{
		Description: "Managing additional RR for IPv4 address objects in QIP. Only supports A records as of now.",

		CreateContext: resourceV4AddressRRCreate,
		ReadContext:   resourceV4AddressRRRead,
		UpdateContext: resourceV4AddressRRUpdate,
		DeleteContext: resourceV4AddressRRDelete,

		Schema: map[string]*schema.Schema{
			"address": {
				Description: "IPv4 address to attach a RR to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Hostname for the address. (e.g. `entry-extra` or `*.entry-extra`)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"domain_name": {
				Description: "DNS Zone for the additional RR.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},

		// Importer: &schema.ResourceImporter{
		// 	StateContext: schema.ImportStatePassthroughContext,
		// },
	}
}

// getIDFromRR uses the base64 encoded JSON representation of a record as ID.
func getIDFromRR(record *rr.RR) (string, error) {
	data, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("could not encode JSON: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// getRRFromID returns a rr.RR decoded from the ID, which is a base64 encoded JSON representation of a record.
func getRRFromID(id string) (*rr.RR, error) {
	data, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, fmt.Errorf("could not decode base64: %w", err)
	}

	var record rr.RR

	err = json.Unmarshal(data, &record)
	if err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}

	return &record, nil
}

func resourceV4AddressRRCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	//nolint:forcetypeassert
	var (
		err     error
		client  = meta.(*terraformClient).QIPClient
		address = d.Get("address").(string)
		name    = d.Get("name").(string)
		domain  = d.Get("domain_name").(string)
	)

	fqdn := name + "." + domain

	record := rr.NewAForObject(fqdn, address)

	err = rr.Create(client, record)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := getIDFromRR(record)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	tflog.Trace(ctx, "Created RR for V4Address "+d.Id())

	return nil
}

func resourceV4AddressRRLoad(_ context.Context, d *schema.ResourceData, meta any) (*rr.RR, error) {
	if d.Id() == "" {
		return nil, ErrIDRequiredToLoad
	}

	idRecord, err := getRRFromID(d.Id())
	if err != nil {
		return nil, err
	}

	records, err := rr.LoadAllForObject(meta.(*terraformClient).QIPClient, idRecord.InfraAddr)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	var singleRecord *rr.RR

	// Search for a single record
	for _, record := range records {
		if record.Equal(idRecord) {
			if singleRecord == nil {
				singleRecord = record
			} else {
				return nil, ErrNonUniqueRR
			}
		}
	}

	return singleRecord, nil
}

func resourceV4AddressRRRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	record, err := resourceV4AddressRRLoad(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if record == nil {
		return diag.Errorf("could not find a record for id: %s", d.Id())
	}

	// Currently we have no state to refresh, all attributes are identifying

	return nil
}

//nolint:forcetypeassert
func resourceV4AddressRRUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	record, err := resourceV4AddressRRLoad(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if record == nil {
		return diag.Errorf("could not find a record for id: %s", d.Id())
	}

	updatedRecord := *record

	// Only allow to change the record owner (FQDN as of now)
	updatedRecord.Owner = d.Get("name").(string) + "." + d.Get("domain_name").(string)

	err = rr.Update(meta.(*terraformClient).QIPClient, record, &updatedRecord)
	if err != nil {
		return diag.FromErr(err)
	}

	// update ID after the change - because some values are identifying
	id, err := getIDFromRR(&updatedRecord)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return nil
}

//nolint:forcetypeassert
func resourceV4AddressRRDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	record, err := resourceV4AddressRRLoad(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Id() == "" || record == nil {
		// Nothing to delete
		return nil
	}

	err = rr.Delete(meta.(*terraformClient).QIPClient, record)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
