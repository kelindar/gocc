package asm

import (
	"fmt"
	"strconv"
)

// RelocateConsts rewrites RIP-relative LCPI displacements to point to
// constants embedded in the same TEXT symbol and marks those constants inline.
func (fn *Function) RelocateConsts() {
	labels := make(map[string]Const, len(fn.Consts))
	for _, c := range fn.Consts {
		labels[normalizeConstLabel(c.Label)] = c
	}

	// Determine which constants are referenced by relocation.
	ordered := make([]string, 0, len(fn.Consts))
	referenced := make(map[string]Const, len(fn.Consts))
	for _, line := range fn.Lines {
		label := line.ConstLabel()
		if label == "" {
			continue
		}

		c, ok := labels[label]
		if !ok {
			continue
		}

		if _, exists := referenced[label]; !exists {
			ordered = append(ordered, label)
		}
		referenced[label] = c
	}

	if len(referenced) == 0 {
		for i := range fn.Consts {
			fn.Consts[i].Inline = false
		}
		return
	}

	for i := range fn.Consts {
		_, inline := referenced[normalizeConstLabel(fn.Consts[i].Label)]
		fn.Consts[i].Inline = inline
	}

	// Compute code offsets.
	lineOffsets := make([]int, len(fn.Lines))
	codeSize := 0
	for i := range fn.Lines {
		lineOffsets[i] = codeSize
		codeSize += len(fn.Lines[i].Binary)
	}

	// Compute inline constant offsets.
	constOffsets := make(map[string]int, len(referenced))
	cursor := codeSize
	for _, label := range ordered {
		c := referenced[label]
		constOffsets[label] = cursor
		cursor += c.Size()
	}

	// Remove duplicated assembly labels for inlined constants. These labels come
	// from clang assembly parsing and would collide with our post-body constants.
	for i := range fn.Lines {
		if len(fn.Lines[i].Labels) == 0 {
			continue
		}

		filtered := fn.Lines[i].Labels[:0]
		for _, label := range fn.Lines[i].Labels {
			if _, inline := referenced[normalizeConstLabel(label)]; inline {
				continue
			}
			filtered = append(filtered, label)
		}
		fn.Lines[i].Labels = filtered
	}

	// Patch disp32 on each relocated instruction.
	for i := range fn.Lines {
		label := fn.Lines[i].ConstLabel()
		target, ok := constOffsets[label]
		if !ok {
			continue
		}

		encoded := decodeHex(fn.Lines[i].Binary)
		if len(encoded) < 4 {
			continue
		}

		lineOffset := lineOffsets[i]
		nextInstruction := lineOffset + len(encoded)
		disp32 := target - nextInstruction
		patchDisp32(encoded, disp32)
		fn.Lines[i].Binary = encodeHex(encoded)
	}

}

func patchDisp32(encoded []byte, disp int) {
	encoded[len(encoded)-4] = byte(disp)
	encoded[len(encoded)-3] = byte(disp >> 8)
	encoded[len(encoded)-2] = byte(disp >> 16)
	encoded[len(encoded)-1] = byte(disp >> 24)
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
		encoded[i] = fmt.Sprintf("%02x", binary[i])
	}
	return encoded
}
