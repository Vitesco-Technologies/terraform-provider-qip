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

package rr_test

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/rr"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/test"
)

func TestLoadAllForObject(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", test.QIPServer+"/api/v1/"+test.QIPOrg+"/rr.json",
		httpmock.NewStringResponder(200, `{
			"list": [
			  {
				"owner": "*.test.int.example.com",
				"classType": "IN",
				"rrType": "A",
				"data1": "192.0.2.50",
				"publishing": "ALWAYS",
				"ttl": -1,
				"infraType": "OBJECT",
				"infraAddr": "192.0.2.50",
				"tombstoned": 0,
				"isCreatingReverseZoneRR": false,
				"isDefaultRR": false
			  }
			]
		  }`))

	records, err := rr.LoadAllForObject(c, "192.0.2.50")
	require.NoError(t, err)

	if assert.Len(t, records, 1) {
		assert.Equal(t, "192.0.2.50", records[0].InfraAddr)
		assert.Equal(t, "192.0.2.50", records[0].Data1)
	}
}

func TestCreate(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("POST", test.QIPServer+"/api/v1/"+test.QIPOrg+"/rr",
		httpmock.NewStringResponder(200, `OK`))

	record := rr.NewAForObject("*.test.int.example.com", "192.0.2.50")

	err := rr.Create(c, record)
	require.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("PUT", test.QIPServer+"/api/v1/"+test.QIPOrg+"/rr",
		httpmock.NewStringResponder(200, `OK`))

	oldRecord := rr.NewAForObject("*.test.int.example.com", "192.0.2.50")
	newRecord := rr.NewAForObject("*.test2.int.example.com", "192.0.2.50")

	err := rr.Update(c, oldRecord, newRecord)
	require.NoError(t, err)
}

func TestDelete(t *testing.T) {
	c, cleanup := test.GetTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("DELETE", test.QIPServer+"/api/v1/"+test.QIPOrg+"/rr",
		httpmock.NewStringResponder(200, `OK`))

	oldRecord := rr.NewAForObject("*.test2.int.example.com", "192.0.2.50")

	err := rr.Delete(c, oldRecord)
	require.NoError(t, err)
}
