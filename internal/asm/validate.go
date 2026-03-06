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
	"errors"
	"strings"
)

var (
	// ErrUnsafeARM64FrameChain indicates known-unsafe arm64 frame-pointer chain
	// prologue/epilogue instructions that may break Go stack unwinding.
	ErrUnsafeARM64FrameChain = errors.New("gocc: unsafe arm64 frame-pointer prologue/epilogue detected")
)

// ValidatorFn validates a line and returns an error when it matches a rule violation.
type ValidatorFn func(line *Line) error

var validatorRules = [...]ValidatorFn{
	validateUnsafeARM64FrameChain,
}

// Validate validates a disassembled instruction line before emitting Plan 9 assembly.
func (line *Line) Validate() error {
	if line == nil {
		return nil
	}

	for _, validate := range validatorRules {
		if err := validate(line); err != nil {
			return err
		}
	}
	return nil
}

// validateUnsafeARM64FrameChain rejects known-unsafe frame-pointer chain updates
// (`stp/ldp x29, x30` with pre/post-indexed SP) by checking both assembly text
// prefixes and the corresponding encoded 32-bit instruction words from objdump.
func validateUnsafeARM64FrameChain(line *Line) error {
	asm := strings.TrimSpace(line.Assembly)
	if len(asm) >= len("stp x29, x30, [sp, #-16]!") && strings.EqualFold(asm[:len("stp x29, x30, [sp, #-16]!")], "stp x29, x30, [sp, #-16]!") ||
		len(asm) >= len("ldp x29, x30, [sp], #16") && strings.EqualFold(asm[:len("ldp x29, x30, [sp], #16")], "ldp x29, x30, [sp], #16") {
		return ErrUnsafeARM64FrameChain
	}
	if len(line.Binary) != 4 {
		return nil
	}

	// Encoded little-endian words for:
	// a9bf7bfd -> stp x29, x30, [sp, #-16]!
	// a8c17bfd -> ldp x29, x30, [sp], #16
	if hexEq(line.Binary, "a9bf7bfd") || hexEq(line.Binary, "a8c17bfd") {
		return ErrUnsafeARM64FrameChain
	}
	return nil
}

func hexEq(v []string, x string) bool {
	if len(v) != 4 || len(x) != 8 {
		return false
	}

	return hexByteEq(v[3], x[0], x[1]) &&
		hexByteEq(v[2], x[2], x[3]) &&
		hexByteEq(v[1], x[4], x[5]) &&
		hexByteEq(v[0], x[6], x[7])
}

func hexByteEq(v string, a, b byte) bool {
	return len(v) == 2 &&
		(v[0]|0x20) == (a|0x20) &&
		(v[1]|0x20) == (b|0x20)
}
