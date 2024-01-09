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

package v4address

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/rest"
)

type SelectedAddrRange struct {
	StartAddress string `json:"startAddress"`
	EndAddress   string `json:"endAddress"`
}

// selectMutex for ensuring that no other selectedv4address operation happens at the same time.
var selectMutex sync.Mutex //nolint:gochecknoglobals

var ErrNoSelection = errors.New("no object address was returned")

// CreateSelected reserves a new address within the range and returns the objectAddr.
//
// If you don't want to create the IP, you need to free it, not sure if it will expire.
func CreateSelected(client *qip.Client, subnet string, addrs *SelectedAddrRange) (string, error) {
	var body any

	if addrs != nil {
		body = struct {
			AddrRange []*SelectedAddrRange `json:"addrRange"`
		}{
			AddrRange: []*SelectedAddrRange{addrs},
		}
	}

	request, err := rest.NewRequest("PUT", client.APITenantURL("selectedv4address", subnet+".json"), body)
	if err != nil {
		return "", fmt.Errorf("could not build select request: %w", err)
	}

	selectMutex.Lock()

	response, err := client.Do(request)
	if err != nil {
		selectMutex.Unlock()

		return "", fmt.Errorf("could not create SelectedV4Address: %w", err)
	}

	selectMutex.Unlock()

	var addr V4Address

	err = rest.UnmarshalResponse(response, &addr)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal v4address: %w", err)
	} else if addr.ObjectAddr == "" {
		return "", ErrNoSelection
	}

	return addr.ObjectAddr, nil
}

// DeleteSelected clears the reservation in the API for an address.
func DeleteSelected(client *qip.Client, addr string) error {
	request, err := rest.NewRequest("DELETE", client.APITenantURL("selectedv4address", addr, "/"), nil)
	if err != nil {
		return fmt.Errorf("could not build delete request: %w", err)
	}

	_, err = client.Do(request)
	if err != nil {
		return fmt.Errorf("could not delete SelectedV4Address: %w", err)
	}

	// Notes for error handling:
	// - Unknown address should return "Internal Server Error - IP address [address] does not have an object associated with it"
	//   but it fails with a NullPointerException

	return nil
}
