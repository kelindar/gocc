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

package cc

import (
	"fmt"
	"os/exec"
)

// Disassembler resolves disassembler to use.
func Disassembler() (string, error) {
	return find([]string{
		"objdump", "llvm-objdump", "llvm-objdump-17", "llvm-objdump-16",
		"llvm-objdump-15", "llvm-objdump-14", "llvm-objdump-13",
	})
}

// find looks for a particular executable in the system
func find(versions []string) (string, error) {
	for _, v := range versions {
		if _, err := exec.LookPath(v); err == nil {
			return v, nil
		}
	}

	return "", fmt.Errorf("gocc: '%s' executable not found)", versions[0])
}
