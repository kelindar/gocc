package cc

import (
	"github.com/kelindar/gocc/internal/asm"
	"github.com/kelindar/gocc/internal/config"
)

type Disassembler struct {
	arch    *config.Arch
	objdump string
}

func NewDisassembler(arch *config.Arch) (*Disassembler, error) {
	objdump, err := config.FindObjdump()
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
