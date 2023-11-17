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

package v4subnet_test

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/test"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/v4subnet"
)

func TestLoad(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", test.QIPServer+"/api/v1/"+test.QIPOrg+"/v4subnet/192.0.2.0.json",
		httpmock.NewStringResponder(200,
			`{"subnetAddress":"192.0.2.0","subnetMask":"255.255.255.0","subnetName":"test-subnet"}`))

	addr, err := v4subnet.Load(c, "192.0.2.0")
	require.NoError(t, err)
	assert.Equal(t, "192.0.2.0", addr.SubnetAddress)
	assert.Equal(t, "test-subnet", addr.SubnetName)
}

func TestAccLoad(t *testing.T) {
	c := test.GetIntegrationTestClient(t)

	testSubnet := os.Getenv("QIP_TEST_SUBNET")
	if testSubnet == "" {
		t.Skip("can not run without QIP_TEST_SUBNET")
	}

	subnet, err := v4subnet.Load(c, testSubnet)
	require.NoError(t, err)
	assert.Equal(t, testSubnet, subnet.SubnetAddress)
	assert.NotEmpty(t, subnet.SubnetName)
}
