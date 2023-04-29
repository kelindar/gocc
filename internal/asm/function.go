// Copyright 2022 gorse Project Authors
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
	"fmt"
	"strings"
)

type Function struct {
	Name       string
	Position   int
	Parameters []string
	Lines      []Line
}

// Line represents a line of assembly code
type Line struct {
	Labels   []string
	Binary   []string
	Assembly string
}

// String returns the string representation of a line in PLAN9 assembly
func (line *Line) String() string {
	var builder strings.Builder
	for _, label := range line.Labels {
		builder.WriteString(label)
		builder.WriteString(":\n")
	}

	builder.WriteString("\t")
	if strings.HasPrefix(line.Assembly, "j") {
		splits := strings.Split(line.Assembly, ".")
		op := strings.TrimSpace(splits[0])
		operand := splits[1]
		builder.WriteString(fmt.Sprintf("%s %s", strings.ToUpper(op), operand))
	} else {
		pos := 0
		for pos < len(line.Binary) {
			if pos > 0 {
				builder.WriteString("; ")
			}
			if len(line.Binary)-pos >= 8 {
				builder.WriteString(fmt.Sprintf("QUAD $0x%v%v%v%v%v%v%v%v",
					line.Binary[pos+7], line.Binary[pos+6], line.Binary[pos+5], line.Binary[pos+4],
					line.Binary[pos+3], line.Binary[pos+2], line.Binary[pos+1], line.Binary[pos]))
				pos += 8
			} else if len(line.Binary)-pos >= 4 {
				builder.WriteString(fmt.Sprintf("LONG $0x%v%v%v%v",
					line.Binary[pos+3], line.Binary[pos+2], line.Binary[pos+1], line.Binary[pos]))
				pos += 4
			} else if len(line.Binary)-pos >= 2 {
				builder.WriteString(fmt.Sprintf("WORD $0x%v%v", line.Binary[pos+1], line.Binary[pos]))
				pos += 2
			} else {
				builder.WriteString(fmt.Sprintf("BYTE $0x%v", line.Binary[pos]))
				pos += 1
			}
		}
		builder.WriteString("\t// ")
		builder.WriteString(line.Assembly)
	}
	builder.WriteString("\n")
	return builder.String()
}
