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

package v4address_test

import (
	"os"
	"sync"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/test"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/v4address"
)

func TestCreateSelected(t *testing.T) {
	c, cleanup := test.GetTestClient(t)

	defer cleanup()

	httpmock.RegisterResponder("PUT", test.QIPServer+"/api/v1/"+test.QIPOrg+"/selectedv4address/192.0.2.0.json",
		httpmock.NewStringResponder(200, `{"objectAddr":"192.0.2.2"}`))

	addr, err := v4address.CreateSelected(c, "192.0.2.0", nil)
	require.NoError(t, err)
	assert.Equal(t, "192.0.2.2", addr)

	httpmock.RegisterResponder("PUT", test.QIPServer+"/api/v1/"+test.QIPOrg+"/selectedv4address/192.0.2.0.json",
		httpmock.NewStringResponder(200, `{"objectAddr":"192.0.2.25"}`))

	addressRange := &v4address.SelectedAddrRange{
		StartAddress: "192.0.2.25",
		EndAddress:   "192.0.2.30",
	}
	addr, err = v4address.CreateSelected(c, "192.0.2.0", addressRange)
	require.NoError(t, err)
	assert.Equal(t, "192.0.2.25", addr)
}

func TestDeleteSelected(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("DELETE", test.QIPServer+"/api/v1/"+test.QIPOrg+"/selectedv4address/192.0.2.25/",
		httpmock.NewStringResponder(200, ``))

	err := v4address.DeleteSelected(c, "192.0.2.25")
	require.NoError(t, err)
}

// TestAccCreateBulkSelected will test if multiple selects against the QIP API fail.
//
// This is a race condition bug in the QIP API, that needs a workaround on client side.
func TestAccCreateBulkSelected(t *testing.T) {
	c := test.GetIntegrationTestClient(t)

	test.SkipIfNotEnabled(t, "QIP_TEST_ACC_BULK_SELECT_ENABLED")

	subnet, addressRange := getAccSubnet(t)

	var (
		wait  sync.WaitGroup
		count = 4
		addrs = make([]*string, count)
	)

	wait.Add(count)

	for num := 0; num < count; num++ {
		go func(num int) {
			defer wait.Done()

			addr, err := v4address.CreateSelected(c, subnet, addressRange)
			require.NoError(t, err)
			assert.NotEmpty(t, addr)

			t.Logf("Selected address #%d: %s", num, addr)

			addrs[num] = &addr
		}(num)
	}

	wait.Wait()

	// One could check here if all addresses are unique

	for num := 0; num < count; num++ {
		if addrs[num] == nil || *addrs[num] == "" {
			continue
		}

		require.NoError(t, v4address.DeleteSelected(c, *addrs[num]))
	}
}

func getAccSubnet(t *testing.T) (string, *v4address.SelectedAddrRange) {
	t.Helper()

	subnet := os.Getenv("QIP_TEST_ACC_SUBNET")
	if subnet == "" {
		t.Skip("QIP_TEST_ACC_SUBNET required")
	}

	var (
		subnetStart = os.Getenv("QIP_TEST_ACC_RESOURCE_SUBNET_START")
		subnetEnd   = os.Getenv("QIP_TEST_ACC_RESOURCE_SUBNET_END")
	)

	var addressRange *v4address.SelectedAddrRange

	if subnetStart != "" && subnetEnd != "" {
		addressRange = &v4address.SelectedAddrRange{
			StartAddress: subnetStart,
			EndAddress:   subnetEnd,
		}
	}

	return subnet, addressRange
}
