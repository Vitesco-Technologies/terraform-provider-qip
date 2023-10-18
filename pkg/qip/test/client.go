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

package test

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip"
)

const (
	QIPServer = "https://qip.example.com"
	QIPOrg    = "Example"
)

func GetTestClient(t *testing.T) (*qip.Client, func()) {
	t.Helper()

	httpmock.Activate()

	c, err := qip.NewClient(QIPServer, QIPOrg)
	if err != nil {
		t.Error(err)
	}

	c.AuthToken = "TEST_TOKEN"

	// err = c.Login("dummy-username", "dummy-password")
	// if err != nil {
	// 	t.Error(err)
	// }

	return c, func() {
		httpmock.DeactivateAndReset()
	}
}

func GetIntegrationTestClient(t *testing.T) *qip.Client {
	t.Helper()

	var (
		testServer   = os.Getenv("QIP_SERVER")
		testOrg      = os.Getenv("QIP_ORG")
		testUsername = os.Getenv("QIP_USERNAME")
		testPassword = os.Getenv("QIP_PASSWORD")
	)

	if testServer == "" && testOrg == "" && testUsername == "" && testPassword == "" {
		t.Skip("can not test without real QIP credentials")
	}

	c, err := qip.NewClient(testServer, testOrg)
	if err != nil {
		t.Error(err)
	}

	err = c.Login(testUsername, testPassword)
	if err != nil {
		t.Error(err)
	}

	return c
}

// SkipIfNotEnabled skips a testing.T if the named environment variable is not set.
func SkipIfNotEnabled(t *testing.T, name string) {
	t.Helper()

	enabled := os.Getenv(name)
	if enabled == "" || enabled == "0" || enabled == "false" {
		t.Skip("Test must be enabled using " + name)
	}
}
