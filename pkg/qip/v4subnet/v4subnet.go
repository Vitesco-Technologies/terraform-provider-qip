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

package v4subnet

import (
	"fmt"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/rest"
)

// Generate the type from a JSON statement (from QIP rest-api documentation)
//go:generate go run github.com/Vitesco-Technologies/terraform-provider-qip/pkg/utils/qip_type -type V4Subnet -package v4subnet

// Load returns V4Subnet data from the API.
func Load(client *qip.Client, address string) (*V4Subnet, error) {
	request, err := rest.NewRequest("GET", client.APITenantURL("v4subnet", address+".json"), nil)
	if err != nil {
		return nil, fmt.Errorf("could not build get request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not load V4Subnet: %w", err)
	}

	var o V4Subnet

	err = rest.UnmarshalResponse(response, &o)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON result: %w", err)
	}

	return &o, nil
}
