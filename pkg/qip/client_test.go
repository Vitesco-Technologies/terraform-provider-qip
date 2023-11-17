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

package qip_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip"
	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip/test"
)

func TestClient_Login(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	c, err := qip.NewClient(test.QIPServer, test.QIPOrg)
	require.NoError(t, err)

	httpmock.RegisterResponder("POST", test.QIPServer+"/api/login",
		func(req *http.Request) (*http.Response, error) {
			body := make(map[string]interface{})
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, ""), nil //nolint:nilerr
			}

			if user, ok := body["username"].(string); ok && user != "unknown-user" {
				resp := httpmock.NewStringResponse(200, "")
				resp.Header.Set("authentication", "THIS_WOULD_BE_A_BASE64_TOKEN")

				return resp, nil
			}

			return httpmock.NewStringResponse(401, ""), nil
		})

	err = c.Login("unknown-user", "dummy-password")
	require.NoError(t, err)

	var targetErr *qip.HTTPUnauthorizedError

	require.ErrorAs(t, err, &targetErr)

	err = c.Login("admin", "password123")
	require.NoError(t, err)
	assert.Equal(t, "THIS_WOULD_BE_A_BASE64_TOKEN", c.AuthToken)
}
