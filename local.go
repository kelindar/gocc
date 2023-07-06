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
package gocc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelindar/gocc/internal/asm"
	"github.com/kelindar/gocc/internal/cc"
	"github.com/kelindar/gocc/internal/config"
)

// Local translates a C file to Go assembly
type Local struct {
	Arch       *config.Arch
	Clang      *cc.Compiler
	ObjDump    *cc.Disassembler
	Source     string
	Assembly   string
	Object     string
	GoAssembly string
	Go         string
	Package    string
	Options    []string
}

// NewLocal creates a new translator that uses locally installed toolchain
func NewLocal(arch *config.Arch, source, outputDir, packageName string, options ...string) (*Local, error) {
	sourceExt := filepath.Ext(source)
	noExtSourcePath := source[:len(source)-len(sourceExt)]
	noExtSourceBase := filepath.Base(noExtSourcePath)
	clang, err := cc.NewCompiler(arch)
	if err != nil {
		return nil, err
	}

	objdump, err := cc.NewDisassembler(arch)
	if err != nil {
		return nil, err
	}

	// If package name is not provided, use the directory name of the output
	if packageName == "" {
		filepath.Base(outputDir)
	}

	return &Local{
		Arch:       arch,
		Clang:      clang,
		ObjDump:    objdump,
		Source:     source,
		Assembly:   fmt.Sprintf("%s.s", noExtSourcePath),
		Object:     fmt.Sprintf("%s.o", noExtSourcePath),
		GoAssembly: filepath.Join(outputDir, fmt.Sprintf("%s.s", noExtSourceBase)),
		Go:         filepath.Join(outputDir, fmt.Sprintf("%s.go", noExtSourceBase)),
		Package:    packageName,
		Options:    options,
	}, nil
}

// Translate translates the source file to Go assembly
func (t *Local) Translate() error {
	functions, err := cc.Parse(t.Source)
	if err != nil {
		return err
	}

	// Generate the Go stubs for the functions
	if err := asm.GenerateGoStubs(t.Arch, t.Package, t.Go, functions); err != nil {
		return err
	}

	// Compile the source file to assembly
	if err := t.Clang.Compile(t.Source, t.Assembly, t.Object, t.Options...); err != nil {
		return err
	}

	// Disassemble the object file and extract machine code
	assembly, err := t.ObjDump.Disassemble(t.Assembly, t.Object)
	if err != nil {
		return err
	}

	// Map the machine code to the assembly one
	for i, v := range assembly {
		functions[i].Lines = v.Lines
	}

	_ = t.Close()
	return asm.GenerateFile(t.Arch, t.GoAssembly, functions)
}

// Output returns the output files as a web result
func (t *Local) Output() (*WebResult, error) {
	goFile, err := os.ReadFile(t.Go)
	if err != nil {
		return nil, err
	}

	asmFile, err := os.ReadFile(t.GoAssembly)
	if err != nil {
		return nil, err
	}

	return &WebResult{
		Asm: File{
			Name: t.GoAssembly,
			Body: asmFile,
		},
		Go: File{
			Name: t.Go,
			Body: goFile,
		},
	}, nil
}

// Cleanup cleans up the temporary files
func (t *Local) Close() error {
	return errors.Join(
		os.Remove(t.Assembly),
		os.Remove(t.Object),
	)
}
