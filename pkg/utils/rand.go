/*
Copyright 2024 Vitesco Technologies Group AG

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

package utils

import "crypto/rand"

var chars = "abcdefghijklmnopqrstuvwxyz"

// ShortID returns a short alphabetic id generated randomly from all lowercase characters.
func ShortID(length int) string {
	ll := len(chars)
	buf := make([]byte, length)

	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}

	for i := 0; i < length; i++ {
		buf[i] = chars[int(buf[i])%ll]
	}

	return string(buf)
}
