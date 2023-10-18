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

//go:generate go run github.com/Vitesco-Technologies/terraform-provider-qip/pkg/utils/qip_type -type V4Address -package v4address
package v4address

import (
	"errors"
	"fmt"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/rest"
)

var (
	ErrBothAddrRequired   = errors.New("ObjectAddr and SubnetAddr is required")
	ErrObjectNameRequired = errors.New("ObjectName is required")
)

func Load(client *qip.Client, address string) (*V4Address, error) {
	request, err := rest.NewRequest("GET", client.APITenantURL("v4address", address+".json"), nil)
	if err != nil {
		return nil, fmt.Errorf("could not build get request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not load V4Address: %w", err)
	}

	var o V4Address

	err = rest.UnmarshalResponse(response, &o)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON result: %w", err)
	}

	return &o, nil
}

// Create a V4Address from input values.
//
// Required fields:
//   - ObjectAddr or SubnetAddr
//   - ObjectName
func Create(client *qip.Client, addr *V4Address) error {
	if addr.ObjectAddr == "" || addr.SubnetAddr == "" {
		return ErrBothAddrRequired
	} else if addr.ObjectName == "" {
		return ErrObjectNameRequired
	}

	request, err := rest.NewRequest("POST", client.APITenantURL("v4address"), addr)
	if err != nil {
		return fmt.Errorf("could not build create request: %w", err)
	}

	_, err = client.Do(request)
	if err != nil {
		return fmt.Errorf("could not create V4Address: %w", err)
	}

	return nil
}

// Update an existing object or converts a select address into an object.
//
// Warning, all fields should be set during an update, leaving out e.g. domain will unset a domain association.
//
// Recommended uses:
//   - LoadV4Address -> Update
//   - SelectV4Address -> Update
func Update(client *qip.Client, addr *V4Address) error {
	if addr.ObjectAddr == "" || addr.SubnetAddr == "" {
		return ErrBothAddrRequired
	} else if addr.ObjectName == "" {
		return ErrObjectNameRequired
	}

	request, err := rest.NewRequest("PUT", client.APITenantURL("v4address"), addr)
	if err != nil {
		return fmt.Errorf("could not build update request: %w", err)
	}

	_, err = client.Do(request)
	if err != nil {
		return fmt.Errorf("could not update V4Address: %w", err)
	}

	return nil
}

// Delete an object and frees its address in the subnet.
func Delete(client *qip.Client, addr string) error {
	request, err := rest.NewRequest("DELETE", client.APITenantURL("v4address", addr, "/"), addr)
	if err != nil {
		return fmt.Errorf("could not build delete request: %w", err)
	}

	_, err = client.Do(request)
	if err != nil {
		return fmt.Errorf("could not delete V4Address: %w", err)
	}

	return nil
}
