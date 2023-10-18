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
	"net/http"
)

// HTTPUnexpectedRedirectError - status 3XX represents a redirect to another resource.
type HTTPUnexpectedRedirectError struct {
	Response *http.Response
}

func (e *HTTPUnexpectedRedirectError) Error() string {
	return "HTTP 3XX unexpected redirect"
}

// HTTPUnauthorizedError - status 401 represents not authorized.
type HTTPUnauthorizedError struct {
	Response *http.Response
}

func (e *HTTPUnauthorizedError) Error() string {
	return "HTTP 401 Authentication failed"
}

// HTTPNotFoundError - status 404 represents resource not found.
type HTTPNotFoundError struct {
	Response *http.Response
}

func (e *HTTPNotFoundError) Error() string {
	return "HTTP 404 Not Found"
}

// HTTPClientError - status 4XX represents a client error.
type HTTPClientError struct {
	Message  string
	Response *http.Response
}

func (e *HTTPClientError) Error() string {
	s := "HTTP 4XX other client error"
	if e.Message != "" {
		s += ": " + e.Message
	}

	return s
}

// HTTPServerError - status 5XX represents various server errors.
type HTTPServerError struct {
	Message  string
	Response *http.Response
}

func (e *HTTPServerError) Error() string {
	s := "HTTP 500 Server Error"
	if e.Message != "" {
		s += ": " + e.Message
	}

	return s
}
