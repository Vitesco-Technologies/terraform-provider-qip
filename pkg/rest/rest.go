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

package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewRequest(method, url string, body any) (*http.Request, error) {
	var buf io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("could not marshal to JSON: %w", err)
		}

		buf = bytes.NewBuffer(data)
	}

	request, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, fmt.Errorf("building request failed: %w", err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

func UnmarshalResponse(response *http.Response, v any) error {
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("could not read the full response body: %w", err)
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("could not unmarshal JSON: %w", err)
	}

	return nil
}
