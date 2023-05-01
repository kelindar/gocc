package cc

import (
	"github.com/kelindar/gocc/internal/asm"
	"github.com/kelindar/gocc/internal/config"
)

// findObjdump resolves disassembler to use.
func findObjdump() (string, error) {
	return find([]string{
		"llvm-objdump", "llvm-objdump-17", "llvm-objdump-16",
		"llvm-objdump-15", "llvm-objdump-14", "llvm-objdump-13",
		"objdump",
	})
}

type Disassembler struct {
	arch    *config.Arch
	objdump string
}

func NewDisassembler(arch *config.Arch) (*Disassembler, error) {
	objdump, err := findObjdump()
	if err != nil {
		return nil, err
	}

	return &Disassembler{
		arch:    arch,
		objdump: objdump,
	}, nil
}

// Disassemble disassembles the object file
func (d *Disassembler) Disassemble(assemblyPath, objectPath string) ([]asm.Function, error) {

	// Parse the assembly file
	assembly, err := asm.ParseAssembly(d.arch, assemblyPath)
	if err != nil {
		return nil, err
	}

	disassembler := []string{d.objdump}
	if d.arch.Disassembler != nil {
		disassembler = d.arch.Disassembler
	}
	disassembler = append(disassembler, "-d", objectPath)

	// Run the disassembler
	dump, err := runCommand(disassembler[0], disassembler[1:]...)
	if err != nil {
		return nil, err
	}

	// Parse the object dump and map machine code to assembly
	if err := asm.ParseObjectDump(d.arch, dump, assembly); err != nil {
		return nil, err
	}

	return assembly, nil
}
