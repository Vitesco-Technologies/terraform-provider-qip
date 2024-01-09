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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceV4Address(t *testing.T) {
	address := os.Getenv("QIP_TEST_ACC_DATA_IP")
	if address == "" {
		t.Skip("must set QIP_TEST_ACC_DATA_IP for this test")
	}

	testSrc := `
	data "qip_v4address" "test" {
		address = "` + address + `"
	  }
	`

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testSrc,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.qip_v4address.test", "id", stringRe(address)),
					resource.TestMatchResourceAttr("data.qip_v4address.test", "subnet", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("data.qip_v4address.test", "name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("data.qip_v4address.test", "domain_name", stringNoWhitespaceRe),
					resource.TestMatchResourceAttr("data.qip_v4address.test", "object_class", stringNoWhitespaceRe),
				),
			},
		},
	})
}
