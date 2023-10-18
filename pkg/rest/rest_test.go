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

package rest_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/rest"
)

type testStruct struct {
	Something string `json:"something"`
}

const testStructJSON = `{"something":"test"}`

var testStructData = testStruct{"test"}

func TestNewRESTRequest(t *testing.T) {
	request, err := rest.NewRequest("GET", "http://localhost/resource", nil)
	assert.NoError(t, err)
	assert.Nil(t, request.Body)

	request, err = rest.NewRequest("POST", "http://localhost/login", testStructData)
	assert.NoError(t, err)
	assert.NotNil(t, request.Body)

	data, err := io.ReadAll(request.Body)
	assert.NoError(t, err)
	assert.Equal(t, testStructJSON, string(data))
}

func TestUnmarshalRESTResponse(t *testing.T) {
	response := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(testStructJSON)),
	}

	var o testStruct
	err := rest.UnmarshalResponse(response, &o)
	assert.NoError(t, err)
	assert.Equal(t, "test", o.Something)
}
