// Copyright 2026 Roman Atachiants
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

func TestGenerateABI_AMD64MixedArgs(t *testing.T) {
	fn := []Function{
		{
			Name: "f32_axpy",
			Params: []Param{
				{Name: "x", IsPointer: true},
				{Name: "y", IsPointer: true},
				{Name: "size", Type: "uint64_t"},
				{Name: "alpha", Type: "float"},
			},
		},
	}

	asm, err := Generate(config.AMD64(), fn)
	assert.NoError(t, err)
	output := string(asm)

	// SysV ABI for amd64 requires float args in XMM registers, independent from GP args.
	assert.Regexp(t, `MOVQ\s+x\+0\(FP\), DI`, output)
	assert.Regexp(t, `MOVQ\s+y\+8\(FP\), SI`, output)
	assert.Regexp(t, `MOVQ\s+size\+16\(FP\), DX`, output)
	assert.Contains(t, output, "MOVSS alpha+24(FP), X0")
	assert.NotContains(t, output, "MOVQ alpha+24(FP), CX")
}

func TestGenerateABI_ARM64MixedArgs(t *testing.T) {
	fn := []Function{
		{
			Name: "f32_axpy",
			Params: []Param{
				{Name: "x", IsPointer: true},
				{Name: "y", IsPointer: true},
				{Name: "size", Type: "uint64_t"},
				{Name: "alpha", Type: "float"},
			},
		},
	}

	asm, err := Generate(config.ARM64(), fn)
	assert.NoError(t, err)
	output := string(asm)

	// AAPCS64 requires float args in FP/SIMD registers, independent from GP args.
	assert.Regexp(t, `MOVD\s+x\+0\(FP\), R0`, output)
	assert.Regexp(t, `MOVD\s+y\+8\(FP\), R1`, output)
	assert.Regexp(t, `MOVD\s+size\+16\(FP\), R2`, output)
	assert.Contains(t, output, "FMOVS alpha+24(FP), F0")
	assert.NotContains(t, output, "MOVD alpha+24(FP), R3")
}

func TestGenerateABI_ARM64FramePointerSafety(t *testing.T) {
	fn := []Function{
		{
			Name: "unsafe_fp_chain",
			Lines: []Line{
				{
					Assembly: "stp x29, x30, [sp, #-16]!",
					Binary:   []string{"fd", "7b", "bf", "a9"},
				},
				{
					Assembly: "mov x29, sp",
					Binary:   []string{"fd", "03", "00", "91"},
				},
				{
					Assembly: "ldp x29, x30, [sp], #16",
					Binary:   []string{"fd", "7b", "c1", "a8"},
				},
				{
					Assembly: "ret",
					Binary:   []string{"c0", "03", "5f", "d6"},
				},
			},
		},
	}

	asm, err := Generate(config.ARM64(), fn)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "unsafe arm64 frame-pointer")
	assert.Nil(t, asm)
}
