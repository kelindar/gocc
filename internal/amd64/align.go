package amd64

import (
	"strconv"
	"strings"

	"github.com/kelindar/gocc/internal/asm"
	"github.com/kelindar/gocc/internal/config"
)

type alignRewrite struct {
	mnemonic string
	patch    patchKind
}

type patchKind int

const (
	patchPPDqa patchKind = iota
	patchOpcodeAps
)

var alignRewrites = map[string]alignRewrite{
	"movdqa":    {mnemonic: "movdqu", patch: patchPPDqa},
	"vmovdqa":   {mnemonic: "vmovdqu", patch: patchPPDqa},
	"vmovdqa32": {mnemonic: "vmovdqu32", patch: patchPPDqa},
	"vmovdqa64": {mnemonic: "vmovdqu64", patch: patchPPDqa},
	"movaps":    {mnemonic: "movups", patch: patchOpcodeAps},
	"vmovaps":   {mnemonic: "vmovups", patch: patchOpcodeAps},
	"movapd":    {mnemonic: "movupd", patch: patchOpcodeAps},
	"vmovapd":   {mnemonic: "vmovupd", patch: patchOpcodeAps},
}

// Rewrite rewrites alignment-required const-pool accesses to their
// unaligned-safe forms and patches corresponding instruction bytes.
func Rewrite(arch *config.Arch, functions []asm.Function) {
	if arch.Name != "amd64" {
		return
	}

	for i := range functions {
		rewriteFunction(&functions[i])
	}
}

func rewriteFunction(fn *asm.Function) {
	consts := make(map[string]struct{}, len(fn.Consts))
	for _, c := range fn.Consts {
		consts[c.Label] = struct{}{}
	}

	for i := range fn.Lines {
		line := &fn.Lines[i]
		label := line.ConstLabel()
		if label == "" {
			continue
		}
		if _, ok := consts[label]; !ok {
			continue
		}

		mnemonic, suffix := splitMnemonic(line.Assembly)
		if mnemonic == "" {
			continue
		}

		rewrite, ok := alignRewrites[strings.ToLower(mnemonic)]
		if !ok {
			continue
		}

		encoded, ok := patchEncoding(line.Binary, rewrite.patch)
		if !ok {
			continue
		}

		line.Binary = encoded
		line.Assembly = rewrite.mnemonic + suffix
	}
}

func splitMnemonic(assembly string) (string, string) {
	trimmed := strings.TrimLeft(assembly, " \t")
	if trimmed == "" {
		return "", ""
	}

	idx := strings.IndexAny(trimmed, " \t")
	if idx < 0 {
		return trimmed, ""
	}

	return trimmed[:idx], trimmed[idx:]
}

func patchEncoding(binary []string, kind patchKind) ([]string, bool) {
	encoded := decodeHex(binary)
	if len(encoded) == 0 {
		return nil, false
	}

	layout, ok := parseEncodingLayout(encoded)
	if !ok {
		return nil, false
	}

	switch kind {
	case patchPPDqa:
		ok = patchDqaToDqu(encoded, layout)
	case patchOpcodeAps:
		ok = patchApsToUps(encoded, layout)
	default:
		return nil, false
	}

	if !ok {
		return nil, false
	}

	return encodeHex(encoded), true
}

type encodingLayout struct {
	ppIndex       int
	opcodeIndex   int
	prefix66Index int
}

func patchDqaToDqu(encoded []byte, layout encodingLayout) bool {
	if !isOneOf(encoded[layout.opcodeIndex], 0x6f, 0x7f) {
		return false
	}

	if layout.ppIndex >= 0 {
		encoded[layout.ppIndex] = (encoded[layout.ppIndex] &^ 0x03) | 0x02
		return true
	}

	if layout.prefix66Index < 0 {
		return false
	}

	encoded[layout.prefix66Index] = 0xF3
	return true
}

func patchApsToUps(encoded []byte, layout encodingLayout) bool {
	switch encoded[layout.opcodeIndex] {
	case 0x28:
		encoded[layout.opcodeIndex] = 0x10
		return true
	case 0x29:
		encoded[layout.opcodeIndex] = 0x11
		return true
	default:
		return false
	}
}

func parseEncodingLayout(encoded []byte) (encodingLayout, bool) {
	switch encoded[0] {
	case 0xC5:
		if len(encoded) < 3 {
			return encodingLayout{}, false
		}
		return encodingLayout{
			ppIndex:       1,
			opcodeIndex:   2,
			prefix66Index: -1,
		}, true
	case 0xC4:
		if len(encoded) < 4 {
			return encodingLayout{}, false
		}
		return encodingLayout{
			ppIndex:       2,
			opcodeIndex:   3,
			prefix66Index: -1,
		}, true
	case 0x62:
		if len(encoded) < 5 {
			return encodingLayout{}, false
		}
		return encodingLayout{
			ppIndex:       2,
			opcodeIndex:   4,
			prefix66Index: -1,
		}, true
	default:
		return parseLegacyLayout(encoded)
	}
}

func parseLegacyLayout(encoded []byte) (encodingLayout, bool) {
	for i := 0; i+1 < len(encoded); i++ {
		if encoded[i] != 0x0F {
			continue
		}

		prefix66 := -1
		for j := 0; j < i; j++ {
			if encoded[j] == 0x66 {
				prefix66 = j
			}
		}

		return encodingLayout{
			ppIndex:       -1,
			opcodeIndex:   i + 1,
			prefix66Index: prefix66,
		}, true
	}

	return encodingLayout{}, false
}

func isOneOf(value byte, options ...byte) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}

func decodeHex(binary []string) []byte {
	encoded := make([]byte, 0, len(binary))
	for _, value := range binary {
		v, err := strconv.ParseUint(value, 16, 8)
		if err != nil {
			return nil
		}
		encoded = append(encoded, byte(v))
	}
	return encoded
}

func encodeHex(binary []byte) []string {
	encoded := make([]string, len(binary))
	for i := range binary {
		encoded[i] = strings.ToLower(strconv.FormatUint(uint64(binary[i]), 16))
		if len(encoded[i]) == 1 {
			encoded[i] = "0" + encoded[i]
		}
	}
	return encoded
}
