// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package include

import (
	"github.com/elastic/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("springboothystrixbeat", "fields.yml", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}