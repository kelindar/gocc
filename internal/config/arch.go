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
	"runtime"
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
	case "apple":
		return Apple(), nil
	case "neon":
		return Neon(), nil
	case "avx2":
		return Avx2(), nil
	case "avx512":
		return Avx512(), nil
	default:
		return nil, fmt.Errorf("unsupported architecture: %s", arch)
	}
}

// ------------------------------------- AMD64 -------------------------------------

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

// Avx2 returns a configuration for AMD64 architecture with AVX2 support
func Avx2() *Arch {
	arch := AMD64()
	arch.ClangFlags = append(arch.ClangFlags, "-masm=intel", "-mavx2", "-mfma")
	return arch
}

// Avx512 returns a configuration for AMD64 architecture with AVX512 support
func Avx512() *Arch {
	arch := AMD64()
	arch.ClangFlags = append(arch.ClangFlags, "-masm=intel", "-mavx", "-mfma", "-mavx512f", "-mavx512dq")
	return arch
}

// ------------------------------------- Linux ARM64 -------------------------------------

// ARM64 returns a configuration for ARM64 architecture
func ARM64() *Arch {
	return &Arch{
		Name:       "arm64",
		Attribute:  regexp.MustCompile(`^\s+\..+$`),
		Function:   regexp.MustCompile(`^\w+:.*$`),
		Label:      regexp.MustCompile(`^.[A-Z0-9]+_\d+:.*$`),
		Code:       regexp.MustCompile(`^\s+\w+.+$`),
		Symbol:     regexp.MustCompile(`^\w+\s+<\w+>:$`),
		Data:       regexp.MustCompile(`^\w+:\s+\w+\s+.+$`),
		Comment:    regexp.MustCompile(`^\s*//.*$`),
		Registers:  []string{"R0", "R1", "R2", "R3"},
		BuildTags:  "//go:build !noasm && !darwin && arm64\n",
		CommentCh:  "//",
		CallOp:     "MOVD",
		ClangFlags: []string{"--target=aarch64-linux-gnu"},
	}
}

// Neon returns a configuration for ARM64 architecture with NEON support
func Neon() *Arch {
	arch := ARM64()
	arch.ClangFlags = append(arch.ClangFlags, "-mfpu=neon", "-mfloat-abi=hard")
	return arch
}

func SVE() *Arch {
	arch := ARM64()
	arch.ClangFlags = append(arch.ClangFlags, "-mfpu=sve", "-mfloat-abi=hard")
	return arch
}

// ------------------------------------- Apple ARM64 -------------------------------------

// Apple returns a configuration for ARM64 architecture. On my M1 mac, supported features are:
// AESARM, ASIMD, ASIMDDP, ASIMDHP, ASIMDRDM, ATOMICS, CRC32, DCPOP, FCMA, FP, FPHP, GPA, JSCVT, LRCPC, PMULL, SHA1, SHA2, SHA3, SHA512
func Apple() *Arch {
	if runtime.GOOS != "darwin" {
		arch := ARM64()
		arch.Disassembler = []string{"llvm-objdump-15"}
		arch.BuildTags = "//go:build !noasm && darwin && arm64\n"
		arch.ClangFlags = []string{"--target=aarch64-apple-darwin", "-mfpu=neon-vfpv4", "-mfloat-abi=hard"}
		return arch
	}

	return &Arch{
		Name:      "arm64",
		Attribute: regexp.MustCompile(`^\s+\..+$`),
		Function:  regexp.MustCompile(`^\w+:.*$`),
		Label:     regexp.MustCompile(`^[A-Z0-9]+_\d+:.*$`),
		Code:      regexp.MustCompile(`^\s+\w+.+$`),
		Symbol:    regexp.MustCompile(`^\w+\s+<\w+>:$`),
		Data:      regexp.MustCompile(`^\w+:\s+\w+\s+.+$`),
		Comment:   regexp.MustCompile(`^\s*;.*$`),
		Registers: []string{"R0", "R1", "R2", "R3"},
		BuildTags: "//go:build !noasm && darwin && arm64\n",
		CommentCh: ";",
		CallOp:    "MOVD",
	}
}
