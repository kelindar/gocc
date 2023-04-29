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
	"fmt"
	"regexp"
)

// Arch represents a context for a specific architecture
type Arch struct {
	Name         string         // Architecture name
	Attribute    *regexp.Regexp // Parses assembly attributes
	Function     *regexp.Regexp // Parses assembly function names
	Label        *regexp.Regexp // Parses assembly labels
	Code         *regexp.Regexp // Parses assembly code
	Symbol       *regexp.Regexp // Parses assembly symbols
	Data         *regexp.Regexp // Parses assembly data
	Comment      *regexp.Regexp // Parses assembly comments
	Registers    []string       // Registers to use
	BuildTags    string         // Golang build tags
	CommentCh    string         // Assembly comment character
	CallOp       string         // Call instruction to use to move the params onto the stack
	Disassembler []string       // Disassembler to use and flags
	ClangFlags   []string       // Flags for clang
}

// For returns a configuration for a given architecture
func For(arch string) (*Arch, error) {
	switch arch {
	case "amd64":
		return AMD64(), nil
	case "arm64":
		return ARM64(), nil
	default:
		return nil, fmt.Errorf("unsupported architecture: %s", arch)
	}
}

// AMD64 returns a configuration for AMD64 architecture
func AMD64() *Arch {
	return &Arch{
		Name:         "amd64",
		Attribute:    regexp.MustCompile(`^\s+\..+$`),
		Function:     regexp.MustCompile(`^\w+:.*$`),
		Label:        regexp.MustCompile(`^\.[A-Z0-9]+_\d+:.*$`),
		Code:         regexp.MustCompile(`^\s+\w+.+$`),
		Symbol:       regexp.MustCompile(`^\w+\s+<\w+>:$`),
		Data:         regexp.MustCompile(`^\w+:\s+\w+\s+.+$`),
		Comment:      regexp.MustCompile(`^\s*#.*$`),
		Registers:    []string{"DI", "SI", "DX", "CX"},
		BuildTags:    "//go:build !noasm && amd64\n",
		CommentCh:    "#",
		CallOp:       "MOVQ",
		Disassembler: []string{"objdump", "--insn-width", "16"},
	}
}

// ARM64 returns a configuration for ARM64 architecture
func ARM64() *Arch {
	return &Arch{
		Name:         "arm64",
		Attribute:    regexp.MustCompile(`^\s+\..+$`),
		Function:     regexp.MustCompile(`^\w+:.*$`),
		Label:        regexp.MustCompile(`^[A-Z0-9]+_\d+:.*$`),
		Code:         regexp.MustCompile(`^\s+\w+.+$`),
		Symbol:       regexp.MustCompile(`^\w+\s+<\w+>:$`),
		Data:         regexp.MustCompile(`^\w+:\s+\w+\s+.+$`),
		Comment:      regexp.MustCompile(`^\s*@.*$`),
		Registers:    []string{"R0", "R1", "R2", "R3"},
		BuildTags:    "//go:build !noasm && arm64\n",
		CommentCh:    "@",
		CallOp:       "MOVD",
		Disassembler: []string{"aarch64-linux-gnu-objdump"},
		ClangFlags:   []string{"--target=arm-linux-gnueabihf", "-march=armv7-a", "-mfpu=neon-vfpv4", "-mfloat-abi=hard"},
	}
}
