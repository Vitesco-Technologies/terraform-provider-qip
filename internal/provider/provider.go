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
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/qip"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"server": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Base URL of the QIP Server (e.g. https://qip.example.com). (env: `QIP_SERVER`)",
					DefaultFunc: schema.EnvDefaultFunc("QIP_SERVER", nil),
				},
				"org": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Organization name inside QIP (e.g. Example). (env: `QIP_ORG`)",
					DefaultFunc: schema.EnvDefaultFunc("QIP_ORG", nil),
				},
				"username": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Username to authenticate against the QIP REST API. (env: `QIP_USERNAME`)",
					DefaultFunc: schema.EnvDefaultFunc("QIP_USERNAME", nil),
				},
				"password": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "Password to authenticate against the QIP REST API. (env: `QIP_PASSWORD`)",
					DefaultFunc: schema.EnvDefaultFunc("QIP_PASSWORD", nil),
				},
				"request_timeout": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Timeout of HTTP requests of the provider in seconds.",
					Default:     qip.DefaultTimeout.Seconds(),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"qip_v4address": dataSourceV4Address(),
				"qip_v4subnet":  dataSourceV4Subnet(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"qip_v4address":    resourceV4Address(),
				"qip_v4address_rr": resourceV4AddressRR(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type terraformClient struct {
	QIPClient *qip.Client
}

func configure(_ string, _ *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		client := &terraformClient{}

		//nolint:forcetypeassert
		var (
			err            error
			server         = d.Get("server").(string)
			org            = d.Get("org").(string)
			username       = d.Get("username").(string)
			password       = d.Get("password").(string)
			requestTimeout = d.Get("request_timeout").(int)
		)

		if server == "" || org == "" || username == "" || password == "" {
			return nil, diag.Errorf("Unable to create QIP client: server, org, username and password must be set")
		}

		client.QIPClient, err = qip.NewClient(server, org)
		if err != nil {
			return nil, diag.Errorf("could not setup QIP Client: %s", err)
		}

		client.QIPClient.Client.Timeout = time.Duration(requestTimeout) * time.Second

		err = client.QIPClient.Login(username, password)
		if err != nil {
			return nil, diag.Errorf("could not authenticate against QIP API: %s", err)
		}

		return client, nil
	}
}
