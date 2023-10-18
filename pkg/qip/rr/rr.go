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

package rr

import (
	"fmt"
	"net/url"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/rest"
)

//go:generate go run github.com/Vitesco-Technologies/terraform-provider-qip/pkg/utils/qip_type -type RR -package rr

// Equal checks if two RR are equal by checking the identifying attributes.
func (record *RR) Equal(otherRecord *RR) bool {
	return record.Owner == otherRecord.Owner &&
		record.ClassType == otherRecord.ClassType &&
		record.RRType == otherRecord.RRType &&
		record.InfraType == otherRecord.InfraType &&
		record.InfraAddr == otherRecord.InfraAddr &&
		record.InfraFQDN == otherRecord.InfraFQDN
}

// DeleteInfo is a subset of RR to delete a RR (singleDelete is added for deletion).
type DeleteInfo struct {
	Owner        string `json:"owner,omitempty"`
	RRType       string `json:"rrType,omitempty"`
	InfraType    string `json:"infraType,omitempty"`
	InfraFQDN    string `json:"infraFQDN,omitempty"` //nolint:tagliatelle
	InfraAddr    string `json:"infraAddr,omitempty"`
	SingleDelete bool   `json:"singleDelete"`
}

const (
	PublishingAlways = "ALWAYS"

	InfraTypeObject = "OBJECT"

	/* Unused other values for InfraType.
	InfraTypeV6Address     = "V6ADDRESS"
	InfraTypeZone          = "ZONE"
	InfraTypeV4ReverseZone = "V4REVERSEZONE"
	InfraTypeV6ReverseZone = "V6REVERSEZONE"
	InfraTypeNode          = "NODE"
	InfraTypeAll           = "ALL"
	*/
)

type loadResult struct {
	List []*RR
}

// NewAForObject returns a RR for a A record belonging to an object in QIP.
//
// Owner is the respective FQDN for the DNS entry.
func NewAForObject(owner, address string) *RR {
	return &RR{
		Owner:                   owner,
		ClassType:               "IN",
		RRType:                  "A",
		Data1:                   address,
		InfraType:               InfraTypeObject,
		InfraAddr:               address,
		Publishing:              PublishingAlways,
		TTL:                     -1,
		IsCreatingReverseZoneRR: false,
		IsDefaultRR:             false,
	}
}

func LoadAllForObject(client *qip.Client, address string) ([]*RR, error) {
	query := url.Values{}
	query.Set("address", address)
	query.Set("type", InfraTypeObject)
	query.Set("getDefaultRRs", "false")

	request, err := rest.NewRequest("GET", client.APITenantURL("rr.json")+"?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not build get request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not load RR: %w", err)
	}

	var results loadResult

	err = rest.UnmarshalResponse(response, &results)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON result: %w", err)
	}

	return results.List, nil
}

func Create(client *qip.Client, rr *RR) error {
	request, err := rest.NewRequest("POST", client.APITenantURL("rr"), rr)
	if err != nil {
		return fmt.Errorf("could not build create request: %w", err)
	}

	_, err = client.Do(request)
	if err != nil {
		return fmt.Errorf("could not create RR: %w", err)
	}

	return nil
}

func Update(client *qip.Client, oldRR, newRR *RR) error {
	data := map[string]*RR{
		"oldRRRec":     oldRR,
		"updatedRRRec": newRR,
	}

	request, err := rest.NewRequest("PUT", client.APITenantURL("rr"), data)
	if err != nil {
		return fmt.Errorf("could not build update request: %w", err)
	}

	_, err = client.Do(request)
	if err != nil {
		return fmt.Errorf("could not update RR: %w", err)
	}

	return nil
}

// Delete will remove a RR from QIP in connection to the belonging object.
//
// Note: this copies values from an RR instance to DeleteInfo, so the API understands the deletion request.
// Sending a simple RR objects yields a NullPointerException within the API.
// This is not really well documented, you will notice the "singleDelete" attribute in the model, but not the example.
func Delete(client *qip.Client, rr *RR) error {
	deleteInfo := &DeleteInfo{
		Owner:        rr.Owner,
		RRType:       rr.RRType,
		InfraType:    rr.InfraType,
		InfraFQDN:    rr.InfraFQDN,
		InfraAddr:    rr.InfraAddr,
		SingleDelete: true,
	}

	request, err := rest.NewRequest("DELETE", client.APITenantURL("rr"), deleteInfo)
	if err != nil {
		return fmt.Errorf("could not build delete request: %w", err)
	}

	_, err = client.Do(request)
	if err != nil {
		return fmt.Errorf("could not delete RR: %w", err)
	}

	return nil
}
