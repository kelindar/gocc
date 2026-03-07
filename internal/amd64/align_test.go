package amd64

import (
	"os"
	"strings"
	"testing"

	"github.com/kelindar/gocc/internal/asm"
	"github.com/kelindar/gocc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestApplyAlignmentSafety(t *testing.T) {
	fn, err := asm.ParseAssembly(config.AMD64(), "../../fixtures/test_avx.s")
	assert.NoError(t, err)

	dump, err := os.ReadFile("../../fixtures/test_avx.o.txt")
	assert.NoError(t, err)
	assert.NoError(t, asm.ParseObjectDump(config.AMD64(), string(dump), fn))
	for i := range fn[0].Lines {
		if strings.Contains(fn[0].Lines[i].Assembly, ".LCPI0_0") {
			fn[0].Lines[i].Relocation = ".LCPI0_0-0x4"
		}
	}

	Rewrite(config.AMD64(), fn)

	rewritten := 0
	for _, line := range fn[0].Lines {
		if !strings.Contains(line.Assembly, ".LCPI0_0") {
			continue
		}

		assert.NotContains(t, line.Assembly, "vmovdqa")
		assert.Contains(t, line.Assembly, "vmovdqu")
		assert.GreaterOrEqual(t, len(line.Binary), 2)
		assert.Equal(t, "c5", strings.ToLower(line.Binary[0]))
		assert.Equal(t, "fe", strings.ToLower(line.Binary[1]))
		rewritten++
	}

	assert.Equal(t, 2, rewritten)
}

func TestApplyAlignmentSafetyFallbackToAssemblyLabel(t *testing.T) {
	fn, err := asm.ParseAssembly(config.AMD64(), "../../fixtures/test_avx.s")
	assert.NoError(t, err)

	dump, err := os.ReadFile("../../fixtures/test_avx.o.txt")
	assert.NoError(t, err)
	assert.NoError(t, asm.ParseObjectDump(config.AMD64(), string(dump), fn))

	Rewrite(config.AMD64(), fn)

	rewritten := 0
	for _, line := range fn[0].Lines {
		if !strings.Contains(line.Assembly, ".LCPI0_0") {
			continue
		}

		assert.Contains(t, line.Assembly, "vmovdqu")
		rewritten++
	}

	assert.Equal(t, 2, rewritten)
}
