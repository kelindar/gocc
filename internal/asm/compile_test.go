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
	"encoding/json"
	"os"
	"testing"

	"github.com/kelindar/gocc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	b, err := os.ReadFile("../../fixtures/test_avx.json")
	assert.NoError(t, err)

	var fn []Function
	assert.NoError(t, json.Unmarshal(b, &fn))

	asm, err := Generate(config.AMD64(), fn)
	assert.NoError(t, err)

	assert.Contains(t, string(asm), "GLOBL LCPI0_0<>(SB), (8+16), $32")
	assert.Contains(t, string(asm), "TEXT ·uint8_mul(SB), $0-32")
}
