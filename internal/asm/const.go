package asm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/kelindar/gocc/internal/config"
)

var constDirective = regexp.MustCompile(`^\s*\.([a-z]+)\s+([^#;]+)`)

type Const struct {
	Label  string      `json:"label"`            // Label of the constant
	Lines  []ConstLine `json:"lines"`            // LInes of the constant
	Inline bool        `json:"inline,omitempty"` // Emit this constant inline in TEXT
}

type ConstLine struct {
	Size  int   `json:"size"`  // Size of the constant
	Value int64 `json:"value"` // Value of the constant
}

// Compile returns the string representation of a constant.
func (c *Const) Compile(arch *config.Arch) string {
	switch {
	case arch.Name != "amd64":
		panic("gocc: only amd64 is supported for constants")
	case c.Inline:
		return c.inline()
	default:
		return c.global()
	}
}

// inline emits constant bytes directly inside TEXT
func (c Const) inline() string {
	var out strings.Builder
	for _, d := range c.Lines {
		switch d.Size {
		case 1:
			out.WriteString(fmt.Sprintf("\tBYTE $0x%02x\n", uint8(d.Value)))
		case 2:
			out.WriteString(fmt.Sprintf("\tWORD $0x%04x\n", uint16(d.Value)))
		case 4:
			out.WriteString(fmt.Sprintf("\tLONG $0x%08x\n", uint32(d.Value)))
		case 8:
			out.WriteString(fmt.Sprintf("\tQUAD $0x%016x\n", uint64(d.Value)))
		}
	}
	return out.String()
}

// global emits constant bytes in the DATA section
func (c Const) global() string {
	var output strings.Builder
	var totalSize int
	for _, d := range c.Lines {
		instruction := fmt.Sprintf("DATA %s<>+%#04x(SB)/%d, $%#04x\n", c.Label, totalSize, d.Size, d.Value)
		output.WriteString(instruction)
		totalSize += d.Size
	}

	output.WriteString(fmt.Sprintf("GLOBL %s<>(SB), (8+16), $%d\n", c.Label, totalSize))
	return output.String()
}

// Size returns the total byte-size of the constant payload.
func (c Const) Size() int {
	size := 0
	for _, line := range c.Lines {
		size += line.Size
	}
	return size
}

// parseConst parses a line in the constant section
func parseConst(arch *config.Arch, line string) []ConstLine {
	if arch.Name != "amd64" {
		panic("gocc: only amd64 is supported for constants")
	}

	match := constDirective.FindStringSubmatch(line)
	if len(match) != 3 {
		panic(fmt.Sprintf("gocc: invalid constant line: %q", line))
	}

	dir := match[1]
	valuePart := strings.TrimSpace(match[2])

	switch dir {
	case "byte", "short", "long", "int", "quad":
		value, err := parseConstInt(valuePart)
		if err != nil {
			panic(fmt.Sprintf("gocc: invalid constant value in data: %v", err))
		}
		return []ConstLine{{
			Size:  scalarConstSize(dir),
			Value: value,
		}}
	case "zero":
		parts := splitComma(valuePart)
		if len(parts) == 0 {
			panic(fmt.Sprintf("gocc: invalid .zero constant: %q", line))
		}

		count, err := strconv.Atoi(parts[0])
		if err != nil || count < 0 {
			panic(fmt.Sprintf("gocc: invalid .zero size: %q", line))
		}

		fill := int64(0)
		if len(parts) > 1 {
			fill, err = parseConstInt(parts[1])
			if err != nil {
				panic(fmt.Sprintf("gocc: invalid .zero fill value: %q", line))
			}
		}

		lines := make([]ConstLine, count)
		for i := 0; i < count; i++ {
			lines[i] = ConstLine{Size: 1, Value: fill}
		}
		return lines
	case "fill":
		parts := splitComma(valuePart)
		if len(parts) < 3 {
			panic(fmt.Sprintf("gocc: invalid .fill constant: %q", line))
		}

		repeat, err := strconv.Atoi(parts[0])
		if err != nil || repeat < 0 {
			panic(fmt.Sprintf("gocc: invalid .fill repeat: %q", line))
		}
		size, err := strconv.Atoi(parts[1])
		if err != nil || size < 0 {
			panic(fmt.Sprintf("gocc: invalid .fill size: %q", line))
		}
		fill, err := parseConstInt(parts[2])
		if err != nil {
			panic(fmt.Sprintf("gocc: invalid .fill value: %q", line))
		}

		lines := make([]ConstLine, 0, repeat)
		for i := 0; i < repeat; i++ {
			lines = append(lines, ConstLine{
				Size:  size,
				Value: fill,
			})
		}
		return lines
	default:
		panic(fmt.Sprintf("gocc: unsupported const directive .%s", dir))
	}
}

func scalarConstSize(dir string) int {
	switch dir {
	case "byte":
		return 1
	case "short":
		return 2
	case "long", "int":
		return 4
	case "quad":
		return 8
	default:
		panic(fmt.Sprintf("gocc: unsupported scalar const directive .%s", dir))
	}
}

func assemblyConstLabel(assembly string) string {
	idx := strings.Index(assembly, ".LCPI")
	if idx < 0 {
		return ""
	}

	rest := assembly[idx+1:] // drop leading dot
	end := 0
	for end < len(rest) {
		r := rune(rest[end])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			end++
			continue
		}
		break
	}

	if end == 0 {
		return ""
	}
	return rest[:end]
}

func relocationLabel(relocation string) string {
	if relocation == "" || (!strings.Contains(relocation, ".LCPI") && !strings.Contains(relocation, "LCPI")) {
		return ""
	}

	end := len(relocation)
	for i, c := range relocation {
		if c == '+' || c == '-' || c == ' ' || c == '\t' {
			end = i
			break
		}
	}

	return relocation[:end]
}

func normalizeConstLabel(label string) string {
	return strings.TrimPrefix(label, ".")
}

func splitComma(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

func parseConstInt(raw string) (int64, error) {
	value := strings.TrimSpace(raw)
	base := 10
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "-0x") {
		base = 0
	}

	return strconv.ParseInt(value, base, 64)
}
