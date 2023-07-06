// Copyright 2023 Roman Atachiants
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package asm

import (
	"testing"

	"github.com/kelindar/gocc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	fn := []Function{{
		Name:     "uint8_mul",
		Position: 10,
		Params: []Param{
			{Name: "input1", Type: "unsignedchar", IsPointer: true},
			{Name: "input2", Type: "unsignedchar", IsPointer: true},
			{Name: "output", Type: "unsignedchar", IsPointer: true},
			{Name: "size", Type: "unsignedlonglong"},
		},
	}}

	asm, err := Generate(config.AMD64(), fn)
	assert.NoError(t, err)
	assert.Equal(t, "", string(asm))
}
