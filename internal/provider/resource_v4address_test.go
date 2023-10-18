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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceV4Address(t *testing.T) {
	address := getRequiredEnv(t, "QIP_TEST_ACC_RESOURCE_IP")
	subnet := getRequiredEnv(t, "QIP_TEST_ACC_RESOURCE_SUBNET")

	name := getRandomName("terraform-qip")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "qip_v4address" "test" {
						address = "` + address + `"
						subnet  = "` + subnet + `"
						name    = "` + name + `"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("qip_v4address.test", "address", stringRe(address)),
					resource.TestMatchResourceAttr("qip_v4address.test", "subnet", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address.test", "name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address.test", "domain_name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address.test", "object_class", regexp.MustCompile(`^Virtualized Server$`)),
				),
			},
			{
				Config: `
					resource "qip_v4address" "test" {
						address = "` + address + `"
						subnet  = "` + subnet + `"
						name    = "` + name + `"

						description = "Added description"
						object_class = "Server"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("qip_v4address.test", "object_class", regexp.MustCompile(`^Server$`)),
					resource.TestMatchResourceAttr("qip_v4address.test", "description", regexp.MustCompile(`^Added`)),
				),
			},
		},
	})
}

func TestAccResourceV4Address_WithSelect(t *testing.T) {
	subnet := getRequiredEnv(t, "QIP_TEST_ACC_RESOURCE_SUBNET")

	testSrc := `
	resource "qip_v4address" "test" {
		subnet = "` + subnet + `"
		name   = "` + getRandomName("terraform-qip") + `"
	}
	`

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testSrc,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("qip_v4address.test", "address", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address.test", "subnet", stringRe(subnet)),
					resource.TestMatchResourceAttr("qip_v4address.test", "name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address.test", "domain_name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address.test", "object_class", stringNonEmptyRe),
				),
			},
		},
	})
}

func TestAccResourceV4Address_WithSelectRange(t *testing.T) {
	var (
		subnet     = getRequiredEnv(t, "QIP_TEST_ACC_RESOURCE_SUBNET")
		rangeStart = getRequiredEnv(t, "QIP_TEST_ACC_RESOURCE_SUBNET_START")
		rangeEnd   = getRequiredEnv(t, "QIP_TEST_ACC_RESOURCE_SUBNET_END")
	)

	testSrc := `
	resource "qip_v4address" "test" {
		subnet = "` + subnet + `"
		name   = "` + getRandomName("terraform-qip") + `"

		subnet_range_start = "` + rangeStart + `"
		subnet_range_end = "` + rangeEnd + `"
	}
	`

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testSrc,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("qip_v4address.test", "address", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address.test", "subnet", regexp.MustCompile("^"+regexp.QuoteMeta(subnet)+"$")),
					resource.TestMatchResourceAttr("qip_v4address.test", "name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address.test", "domain_name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("qip_v4address.test", "object_class", regexp.MustCompile(`^.+$`)),
				),
			},
		},
	})
}
