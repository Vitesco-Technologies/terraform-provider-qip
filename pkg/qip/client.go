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

package qip

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Vitesco-Technologies/terraform-provider-qip/pkg/rest"
)

type Client struct {
	BaseURL   string
	OrgName   string
	AuthToken string
	Client    *http.Client
}

var ErrNoAuthToken = errors.New("no authentication token was returned in header")

const DefaultTimeout = 20 * time.Second

func NewClient(baseURL, orgName string) (*Client, error) {
	// validate URL by parsing it
	if _, err := url.Parse(baseURL); err != nil {
		return nil, fmt.Errorf("base URL is not valid: %w", err)
	}

	return &Client{
		BaseURL: baseURL,
		OrgName: orgName,
		Client: &http.Client{
			Timeout: DefaultTimeout,
		},
	}, nil
}

func (c *Client) Login(username, password string) error {
	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Expires  uint16 `json:"expires"`
	}{
		username,
		password,
		10 * 60, // token will be valid for 10 minutes
	}

	request, err := rest.NewRequest("POST", c.apiURL("login"), body)
	if err != nil {
		return fmt.Errorf("could not build login request: %w", err)
	}

	// Clear the token now
	c.AuthToken = ""

	response, err := c.Do(request)
	if err != nil {
		return err
	}

	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}

	c.AuthToken = response.Header.Get("authentication")
	if c.AuthToken == "" {
		return ErrNoAuthToken
	}

	return nil
}

// Do executes and returns the http.Response.
//
// For this implementation, status codes are checked and error is returned accordingly.
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	if c.AuthToken != "" {
		// Pass auth token to request if set
		request.Header.Set("Authentication", "Token "+c.AuthToken)
	}

	response, err := c.Client.Do(request)

	// Read all of body and store in buffer
	var rawBody []byte
	if response != nil && response.Body != nil {
		rawBody, err = io.ReadAll(response.Body)
		if err == nil {
			response.Body.Close()

			// re-insert the body as buffer to the response
			response.Body = io.NopCloser(bytes.NewBuffer(rawBody))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	switch {
	case response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices:
		return response, nil
	case response.StatusCode >= http.StatusMultipleChoices && response.StatusCode < http.StatusBadRequest:
		return response, &HTTPUnexpectedRedirectError{response}
	case response.StatusCode == http.StatusUnauthorized:
		return response, &HTTPUnauthorizedError{response}
	case response.StatusCode == http.StatusNotFound:
		return response, &HTTPNotFoundError{response}
	}

	// try to parse error from
	errorBody := struct {
		Error string `json:"error"`
	}{}
	_ = json.Unmarshal(rawBody, &errorBody)

	if response.StatusCode >= 400 && response.StatusCode < 500 {
		return response, &HTTPClientError{errorBody.Error, response}
	}

	// response.StatusCode >= 500
	return response, &HTTPServerError{errorBody.Error, response}
}

// apiURL builds a full URL from base and specified parts.
func (c *Client) apiURL(path ...string) string {
	path = append([]string{"api"}, path...)

	fullURL, err := url.JoinPath(c.BaseURL, path...)
	if err != nil {
		// Errors here should not happen, we validate the URL earlier
		panic(fmt.Errorf("could not join URL: %w", err))
	}

	return fullURL
}

func (c *Client) APITenantURL(path ...string) string {
	path = append([]string{"v1", c.OrgName}, path...)

	return c.apiURL(path...)
}

// func (c *Client) ApiGlobalURL(path ...string) string {
// 	path = append([]string{"global", "v1"}, path...)
// 	return c.apiUrl(path...)
// }
