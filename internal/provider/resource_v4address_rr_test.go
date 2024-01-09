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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceV4AddressRR(t *testing.T) {
	subnet := getRequiredEnv(t, "QIP_TEST_ACC_RESOURCE_SUBNET")
	name := getRandomName("terraform-qip-rr")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "qip_v4address" "test" {
						subnet  = "` + subnet + `"
						name    = "` + name + `"
					}

					resource "qip_v4address_rr" "test" {
						name        = "` + name + `-extra"
						address     = qip_v4address.test.address
						domain_name = qip_v4address.test.domain_name
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("qip_v4address.test", "name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address_rr.test", "address", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address_rr.test", "name", stringRe(name+"-extra")),
				),
			},
			{
				Config: `
					resource "qip_v4address" "test" {
						subnet  = "` + subnet + `"
						name    = "` + name + `"
					}

					resource "qip_v4address_rr" "test" {
						name        = "*.` + name + `"
						address     = qip_v4address.test.address
						domain_name = qip_v4address.test.domain_name
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("qip_v4address.test", "name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address_rr.test", "address", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address_rr.test", "name", stringRe(`*.`+name)),
				),
			},
		},
	})
}
