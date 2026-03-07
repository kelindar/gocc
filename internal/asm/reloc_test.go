package asm

import (
	"testing"

	"github.com/kelindar/gocc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestPatchConstRelocations(t *testing.T) {
	fn := Function{
		Name: "foo",
		Consts: []Const{
			{
				Label: "LCPI0_0",
				Lines: []ConstLine{
					{Size: 2, Value: 0x00ff},
					{Size: 2, Value: 0x00ff},
				},
			},
		},
		Lines: []Line{
			{
				Assembly:   "vmovdqa 0x0(%rip),%ymm0",
				Binary:     []string{"c5", "fd", "6f", "05", "00", "00", "00", "00"},
				Relocation: "LCPI0_0-0x4",
			},
			{
				Assembly: "ret",
				Binary:   []string{"c3"},
			},
		},
	}

	fn.relocate()
	assert.True(t, fn.Consts[0].Inline)

	// code size is 9 bytes, so disp32 points to first inline const at +1
	assert.Equal(t, []string{"c5", "fd", "6f", "05", "01", "00", "00", "00"}, fn.Lines[0].Binary)
}

func TestGenerateInlinesRelocatedConsts(t *testing.T) {
	fn := []Function{
		{
			Name: "foo",
			Consts: []Const{
				{
					Label: "LCPI0_0",
					Lines: []ConstLine{
						{Size: 2, Value: 0x00ff},
						{Size: 2, Value: 0x00ff},
					},
				},
			},
			Lines: []Line{
				{
					Assembly:   "vmovdqa 0x0(%rip),%ymm0",
					Binary:     []string{"c5", "fd", "6f", "05", "00", "00", "00", "00"},
					Relocation: "LCPI0_0-0x4",
				},
				{
					Assembly: "ret",
					Binary:   []string{"c3"},
				},
			},
		},
	}
	for i := range fn {
		fn[i].relocate()
	}

	asm, err := Generate(config.AMD64(), fn)
	assert.NoError(t, err)
	output := string(asm)

	assert.NotContains(t, output, "GLOBL LCPI0_0<>(SB)")
	assert.Contains(t, output, "QUAD $0x00000001056ffdc5")
	assert.Contains(t, output, "WORD $0x00ff")
}

func TestPatchConstRelocationsFallbackToAssemblySymbol(t *testing.T) {
	fn := Function{
		Name: "foo",
		Consts: []Const{
			{
				Label: "LCPI0_0",
				Lines: []ConstLine{
					{Size: 2, Value: 0x00ff},
				},
			},
		},
		Lines: []Line{
			{
				Assembly: "vpbroadcastw ymm0, word ptr [rip + .LCPI0_0]",
				Binary:   []string{"c4", "e2", "7d", "79", "05", "00", "00", "00", "00"},
			},
		},
	}

	fn.relocate()
	assert.True(t, fn.Consts[0].Inline)
	assert.Equal(t, []string{"c4", "e2", "7d", "79", "05", "00", "00", "00", "00"}, fn.Lines[0].Binary)
}
