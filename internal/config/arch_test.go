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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestARM64(t *testing.T) {
	cfg, err := For("arm64")
	assert.NoError(t, err)
	assert.Contains(t, cfg.BuildTags, "arm64")
}

func TestAMD64(t *testing.T) {
	cfg, err := For("amd64")
	assert.NoError(t, err)
	assert.Contains(t, cfg.BuildTags, "amd64")
}

func TestApple(t *testing.T) {
	cfg, err := For("apple")
	assert.NoError(t, err)
	assert.Contains(t, cfg.BuildTags, "arm64")
}

func TestNeon(t *testing.T) {
	cfg, err := For("neon")
	assert.NoError(t, err)
	assert.Contains(t, cfg.BuildTags, "arm64")
}

func TestAvx2(t *testing.T) {
	cfg, err := For("avx2")
	assert.NoError(t, err)
	assert.Contains(t, cfg.BuildTags, "amd64")
}

func TestAvx512(t *testing.T) {
	cfg, err := For("avx512")
	assert.NoError(t, err)
	assert.Contains(t, cfg.BuildTags, "amd64")
}
