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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/utils"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"qip": func() (*schema.Provider, error) { //nolint:unparam
		return New("dev")(), nil
	},
}

var (
	ipv4AddressRe = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
	// ipv4CIDRRe checks for an IPV4 CIDR with prefix length 8-32.
	ipv4PrefixLengthRe   = regexp.MustCompile(`([89]|[1-2][0-9]|3[0-2])`)
	ipv4CIDRRe           = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}/` + ipv4PrefixLengthRe.String() + `$`)
	stringNoWhitespaceRe = regexp.MustCompile(`^\S+$`)
	stringNonEmptyRe     = regexp.MustCompile(`^.+$`)
)

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	_ = getRequiredEnv(t, "QIP_SERVER")
	_ = getRequiredEnv(t, "QIP_USERNAME")
	_ = getRequiredEnv(t, "QIP_PASSWORD")
}

func getRequiredEnv(t *testing.T, name string) string {
	t.Helper()

	value := os.Getenv(name)
	if value == "" {
		t.Skip("env:" + name + " must be set for tests")
	}

	return value
}

func getRandomName(prefix string) string {
	return prefix + "-" + utils.ShortID(4)
}

func stringRe(text string) *regexp.Regexp {
	return regexp.MustCompile(`^` + regexp.QuoteMeta(text) + `$`)
}
