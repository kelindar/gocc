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

func TestFind(t *testing.T) {
	// Test to check if the function returns the first executable found in the system
	result, err := find([]string{"nonexistent-executable", "echo", "cd", "cmd"})
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// Test to check if the function returns an error if no executable is found
	result, err = find([]string{"nonexistent-executable"})
	assert.Error(t, err)
	assert.Empty(t, result)
}
