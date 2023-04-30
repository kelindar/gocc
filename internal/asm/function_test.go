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

package asm

import (
	"testing"

	"github.com/kelindar/gocc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLineWord(t *testing.T) {
	line := Line{
		Assembly: "vaddps 0x40(%rdx,%rax,4),%ymm3,%ymm3",
		Binary:   []string{"c5", "e4", "58", "5c", "82", "40"},
	}
	assert.Equal(t, "\tLONG $0x5c58e4c5; WORD $0x4082\t// vaddps 0x40(%rdx,%rax,4),%ymm3,%ymm3\n",
		line.Compile(nil))
}

func TestLineByte(t *testing.T) {
	line := Line{
		Assembly: "ret",
		Binary:   []string{"c3"},
	}
	assert.Equal(t, "\tBYTE $0xc3\t// ret\n", line.Compile(nil))
}

func TestLineLabel(t *testing.T) {
	line := Line{
		Labels:   []string{"label"},
		Assembly: "ret",
		Binary:   []string{"c3"},
	}
	assert.Equal(t, "label:\n\tBYTE $0xc3\t// ret\n", line.Compile(nil))
}

func TestLineARM(t *testing.T) {
	line := Line{
		Assembly: "mov x29, sp",
		Binary:   []string{"fd", "03", "00", "91"},
	}
	assert.Equal(t, "\tWORD $0x910003fd\t// mov x29, sp\n",
		line.Compile(config.ARM64()))
}
