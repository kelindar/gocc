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
	"os"
	"testing"

	"github.com/kelindar/gocc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestParseAssembly(t *testing.T) {
	fn, err := ParseAssembly(config.AMD64(), "../../fixtures/test_avx.s")
	assert.NoError(t, err)
	assert.Len(t, fn, 1)
	assert.Len(t, fn[0].Consts, 1)
	assert.Len(t, fn[0].Lines, 135)
	for _, line := range fn[0].Lines {
		assert.NotEmpty(t, line.Assembly)
		assert.Empty(t, line.Binary)
	}
}

func TestParseObjectDump(t *testing.T) {
	fn, err := ParseAssembly(config.AMD64(), "../../fixtures/test_avx.s")
	assert.NoError(t, err)

	dump, err := os.ReadFile("../../fixtures/test_avx.o.txt")
	assert.NoError(t, err)

	assert.NoError(t, ParseObjectDump(config.AMD64(), string(dump), fn))
	assert.Len(t, fn, 1)
	assert.Len(t, fn[0].Consts, 1)
	assert.Len(t, fn[0].Lines, 135)
	for _, line := range fn[0].Lines {
		assert.NotEmpty(t, line.Assembly)
		assert.NotEmpty(t, line.Binary)
	}
}
