package asm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineValidate_UnsafeByAssembly(t *testing.T) {
	line := &Line{Assembly: "stp x29, x30, [sp, #-16]!"}
	err := line.Validate()
	assert.ErrorIs(t, err, ErrUnsafeARM64FrameChain)
}

func TestLineValidate_UnsafeByBinary(t *testing.T) {
	line := &Line{Binary: []string{"fd", "7b", "bf", "a9"}}
	err := line.Validate()
	assert.ErrorIs(t, err, ErrUnsafeARM64FrameChain)
}

func TestLineValidate_Safe(t *testing.T) {
	line := &Line{
		Assembly: "mov x29, sp",
		Binary:   []string{"fd", "03", "00", "91"},
	}
	err := line.Validate()
	assert.NoError(t, err)
}

func TestLineValidate_NilLine(t *testing.T) {
	var line *Line
	err := line.Validate()
	assert.NoError(t, err)
}
