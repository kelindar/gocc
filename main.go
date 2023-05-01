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
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelindar/gocc/internal/asm"
	"github.com/kelindar/gocc/internal/cc"
	"github.com/kelindar/gocc/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	command.PersistentFlags().StringP("output", "o", "", "output directory of generated files")
	command.PersistentFlags().StringSliceP("machine-option", "m", nil, "machine option for clang")
	command.PersistentFlags().IntP("optimize-level", "O", 0, "optimization level for clang")
	command.PersistentFlags().StringP("arch", "a", "amd64", "target architecture to use")
}

func main() {
	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var command = &cobra.Command{
	Use:  "gocc source [-o output_directory]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.PersistentFlags().GetString("output")
		if output == "" {
			var err error
			if output, err = os.Getwd(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		var options []string
		machineOptions, _ := cmd.PersistentFlags().GetStringSlice("machine-option")
		for _, m := range machineOptions {
			options = append(options, "-m"+m)
		}

		// Load the architecture
		target, _ := cmd.PersistentFlags().GetString("arch")
		arch, err := config.For(target)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		optimizeLevel, _ := cmd.PersistentFlags().GetInt("optimize-level")
		options = append(options, fmt.Sprintf("-O%d", optimizeLevel))
		file, err := NewTranslator(arch, args[0], output, options...)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := file.Translate(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

// Translator translates a C file to Go assembly
type Translator struct {
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

func NewTranslator(arch *config.Arch, source string, outputDir string, options ...string) (*Translator, error) {
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

	return &Translator{
		Arch:       arch,
		Clang:      clang,
		ObjDump:    objdump,
		Source:     source,
		Assembly:   fmt.Sprintf("%s.s", noExtSourcePath),
		Object:     fmt.Sprintf("%s.o", noExtSourcePath),
		GoAssembly: filepath.Join(outputDir, fmt.Sprintf("%s.s", noExtSourceBase)),
		Go:         filepath.Join(outputDir, fmt.Sprintf("%s.go", noExtSourceBase)),
		Package:    filepath.Base(outputDir),
		Options:    options,
	}, nil
}

// Translate translates the source file to Go assembly
func (t *Translator) Translate() error {
	defer t.Cleanup()
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
	return asm.Generate(t.Arch, t.GoAssembly, functions)
}

// Cleanup cleans up the temporary files
func (t *Translator) Cleanup() {
	_ = os.Remove(t.Assembly)
	_ = os.Remove(t.Object)
}
