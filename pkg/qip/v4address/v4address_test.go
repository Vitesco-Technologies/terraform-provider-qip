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

package v4address_test

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/test"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/v4address"
)

func TestLoad(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", test.QIPServer+"/api/v1/"+test.QIPOrg+"/v4address/192.0.2.50.json",
		httpmock.NewStringResponder(200, `{"objectAddr":"192.0.2.50","subnetAddr":"192.0.2.0","objectName":"test-host"}`))

	addr, err := v4address.Load(c, "192.0.2.50")
	assert.NoError(t, err)
	assert.Equal(t, "192.0.2.50", addr.ObjectAddr)
}

func TestCreate(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	addr := &v4address.V4Address{
		SubnetAddr: "192.0.2.0",
		ObjectAddr: "192.0.2.50",
		ObjectName: "test-host",
	}

	httpmock.RegisterResponder("POST", test.QIPServer+"/api/v1/"+test.QIPOrg+"/v4address",
		httpmock.NewStringResponder(200, ""))

	err := v4address.Create(c, addr)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	addr := &v4address.V4Address{
		SubnetAddr: "192.0.2.0",
		ObjectAddr: "192.0.2.50",
		ObjectName: "test-host",
	}

	httpmock.RegisterResponder("PUT", test.QIPServer+"/api/v1/"+test.QIPOrg+"/v4address",
		httpmock.NewStringResponder(200, ""))

	err := v4address.Update(c, addr)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("DELETE", test.QIPServer+"/api/v1/"+test.QIPOrg+"/v4address/192.0.2.55/",
		httpmock.NewStringResponder(200, ""))

	err := v4address.Delete(c, "192.0.2.55")
	assert.NoError(t, err)
}

func TestE2E(t *testing.T) {
	c := test.GetIntegrationTestClient(t)
	testSubnet, _ := getTestSubnet(t)

	addr := os.Getenv("QIP_TEST_SUBNET_IP")
	if addr == "" {
		t.Skip("can not test without QIP_TEST_SUBNET_IP")
	}

	addrObj := &v4address.V4Address{
		ObjectAddr:  addr,
		SubnetAddr:  testSubnet,
		ObjectName:  "terraform-provider-qip",
		ObjectDesc:  "Terraform integration testing",
		ObjectClass: "Virtualized Server",
	}

	err := v4address.Create(c, addrObj)
	assert.NoError(t, err)

	err = v4address.Delete(c, addr)
	assert.NoError(t, err)
}

func TestE2E_WithSelect(t *testing.T) {
	c := test.GetIntegrationTestClient(t)
	testSubnet, testAddrRange := getTestSubnet(t)

	addr, err := v4address.CreateSelected(c, testSubnet, testAddrRange)
	assert.NoError(t, err)
	assert.NotEmpty(t, addr)

	// Not tested here - we update the selected address
	// err = qip.DeleteSelectedV4Address(c, addr)
	// assert.NoError(t, err)

	addrObj := &v4address.V4Address{
		ObjectAddr:  addr,
		SubnetAddr:  testSubnet,
		ObjectName:  "terraform-provider-qip",
		ObjectDesc:  "Terraform integration testing",
		ObjectClass: "Virtualized Server",
	}

	err = v4address.Update(c, addrObj)
	assert.NoError(t, err)

	updatedObj, err := v4address.Load(c, addr)
	assert.NoError(t, err)
	assert.Equal(t, "Virtualized Server", updatedObj.ObjectClass)
	assert.NotEmpty(t, updatedObj.ObjectDesc)
	assert.NotEqual(t, "None", updatedObj.DomainName)

	err = v4address.Delete(c, addr)
	assert.NoError(t, err)
}

func getTestSubnet(t *testing.T) (string, *v4address.SelectedAddrRange) {
	t.Helper()

	var (
		testSubnet           = os.Getenv("QIP_TEST_SUBNET")
		addrRange            *v4address.SelectedAddrRange
		testSubnetRangeStart = os.Getenv("QIP_TEST_SUBNET_RANGE_START")
		testSubnetRangeEnd   = os.Getenv("QIP_TEST_SUBNET_RANGE_END")
	)

	if testSubnet == "" {
		t.Skip("can not run without QIP_TEST_SUBNET")
	}

	if testSubnetRangeStart != "" && testSubnetRangeEnd != "" {
		addrRange = &v4address.SelectedAddrRange{
			StartAddress: testSubnetRangeStart,
			EndAddress:   testSubnetRangeEnd,
		}
	}

	return testSubnet, addrRange
}
