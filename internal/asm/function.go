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
	"strconv"
	"strings"

	"github.com/kelindar/gocc/internal/config"
)

// ------------------------------------- Function -------------------------------------

type Function struct {
	Name     string  `json:"name"`
	Position int     `json:"position"`
	Params   []Param `json:"params"`
	Consts   []Const `json:"consts,omitempty"`
	Lines    []Line  `json:"lines"`
}

// String returns the function signature for a Go stub
func (f *Function) String() string {
	var builder strings.Builder
	builder.WriteString("\n//go:noescape,nosplit\n")
	builder.WriteString(fmt.Sprintf("func %s(", f.Name))
	for i, param := range f.Params {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(param.String())
	}
	builder.WriteString(")\n")
	return builder.String()
}

// ------------------------------------- Code -------------------------------------

// Line represents a line of assembly code
type Line struct {
	Labels   []string `json:"labels,omitempty"` // Labels for the line
	Binary   []string `json:"binary"`           // Binary representation of the line
	Assembly string   `json:"assembly"`         // Assembly representation of the line
}

// Compile returns the string representation of a line in PLAN9 assembly
func (line *Line) Compile(arch *config.Arch) string {
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
		builder.WriteString("\n")
		return builder.String()
	}

	// Special case for arm64, since it's a RISC architecture
	if arch != nil && arch.Name == "arm64" && len(line.Binary) == 4 {
		builder.WriteString(fmt.Sprintf("WORD $0x%v%v%v%v",
			line.Binary[3], line.Binary[2], line.Binary[1], line.Binary[0]))
		builder.WriteString("\t// ")
		builder.WriteString(line.Assembly)
		builder.WriteString("\n")
		return builder.String()
	}

	// Dynamic length, assuming WORD = 32-bit
	for pos := 0; pos < len(line.Binary); {
		if pos > 0 {
			builder.WriteString("; ")
		}

		switch {
		case len(line.Binary)-pos >= 8:
			builder.WriteString(fmt.Sprintf("QUAD $0x%v%v%v%v%v%v%v%v",
				line.Binary[pos+7], line.Binary[pos+6], line.Binary[pos+5], line.Binary[pos+4],
				line.Binary[pos+3], line.Binary[pos+2], line.Binary[pos+1], line.Binary[pos]))
			pos += 8
		case len(line.Binary)-pos >= 4:
			builder.WriteString(fmt.Sprintf("LONG $0x%v%v%v%v",
				line.Binary[pos+3], line.Binary[pos+2], line.Binary[pos+1], line.Binary[pos]))
			pos += 4
		case len(line.Binary)-pos >= 2:
			builder.WriteString(fmt.Sprintf("WORD $0x%v%v", line.Binary[pos+1], line.Binary[pos]))
			pos += 2
		case len(line.Binary)-pos >= 1:
			builder.WriteString(fmt.Sprintf("BYTE $0x%v", line.Binary[pos]))
			pos += 1
		}
	}

	builder.WriteString("\t// ")
	builder.WriteString(line.Assembly)
	builder.WriteString("\n")
	return builder.String()
}

// Param represents a function parameter
type Param struct {
	Type      string `json:"type"`                // Type of the parameter (C type)
	Name      string `json:"name"`                // Name of the parameter
	IsPointer bool   `json:"isPointer,omitempty"` // Whether the parameter is a pointer
}

// String returns the Go string representation of a parameter
func (p *Param) String() string {
	if p.IsPointer {
		return fmt.Sprintf("%s unsafe.Pointer", p.Name)
	}

	switch p.Type {
	case "int16_t":
		return fmt.Sprintf("%s int16", p.Name)
	case "int32_t":
		return fmt.Sprintf("%s int32", p.Name)
	case "int64_t":
		return fmt.Sprintf("%s int64", p.Name)
	case "uint16_t":
		return fmt.Sprintf("%s uint16", p.Name)
	case "uint32_t":
		return fmt.Sprintf("%s uint32", p.Name)
	case "uint64_t":
		return fmt.Sprintf("%s uint64", p.Name)
	case "float":
		return fmt.Sprintf("%s float32", p.Name)
	case "double":
		return fmt.Sprintf("%s float64", p.Name)
	case "unsignedlonglong":
		return fmt.Sprintf("%s uint64", p.Name)
	case "unsignedint":
		return fmt.Sprintf("%s uint32", p.Name)
	case "longlong":
		return fmt.Sprintf("%s int64", p.Name)
	case "int":
		return fmt.Sprintf("%s int32", p.Name)
	default:
		panic(fmt.Sprintf("gocc: unknown type %s", p.Type))
	}
}

// ------------------------------------- Constants -------------------------------------

type Const struct {
	Label string      `json:"label"` // Label of the constant
	Lines []ConstLine `json:"lines"` // LInes of the constant
}

type ConstLine struct {
	Size  int   `json:"size"`  // Size of the constant
	Value int64 `json:"value"` // Value of the constant
}

// Compile returns the string representation of a line in PLAN9 assembly
func (c *Const) Compile(arch *config.Arch) string {
	if arch.Name != "amd64" {
		panic("gocc: only amd64 is supported for constants")
	}

	var output strings.Builder
	var totalSize int
	for _, d := range c.Lines {

		// Write the DATA instruction.
		instruction := fmt.Sprintf("DATA %s<>+%#04x(SB)/%d, $%#04x\n", c.Label, totalSize, d.Size, d.Value)
		output.WriteString(instruction)
		totalSize += d.Size
	}

	// Write the GLOBL instruction.
	output.WriteString(fmt.Sprintf("GLOBL %s<>(SB), NOPTR, $%d\n", c.Label, totalSize))
	return output.String()
}

// parseConst parses a line in the constant section
func parseConst(arch *config.Arch, line string) ConstLine {
	if arch.Name != "amd64" {
		panic("gocc: only amd64 is supported for constants")
	}

	sizes := map[string]int{
		"byte":  1,
		"short": 2,
		"long":  4,
		"int":   4,
		"quad":  8,
	}

	match := arch.Const.FindStringSubmatch(line)
	typeName := match[1]
	value, err := strconv.ParseInt(match[2], 10, 64)
	if err != nil {
		panic(fmt.Sprintf("gocc: invalid constant value in data: %v", err))
	}

	return ConstLine{
		Size:  sizes[typeName],
		Value: value,
	}
}
