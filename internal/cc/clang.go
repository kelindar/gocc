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

package cc

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kelindar/gocc/internal/config"
)

// Compiler represents a C/C++ compiler.
type Compiler struct {
	arch  *config.Arch
	clang string
}

// NewCompiler creates a new compiler.
func NewCompiler(arch *config.Arch) (*Compiler, error) {
	clang, err := config.FindClang()
	if err != nil {
		return nil, err
	}

	return &Compiler{
		arch:  arch,
		clang: clang,
	}, nil
}

// compile compiles the C source file to assembly and then to object.
func (c *Compiler) Compile(source, assembly, object string, args ...string) error {
	args = append(args, "-mno-red-zone", "-mstackrealign", "-mllvm", "-inline-threshold=1000",
		"-fno-asynchronous-unwind-tables", "-fno-exceptions", "-fno-rtti", "-ffast-math")
	args = append(args, c.arch.ClangFlags...)

	// Compile to assembly first
	if _, err := runCommand(c.clang, append([]string{"-S", "-c", source, "-o", assembly}, args...)...); err != nil {
		return err
	}

	// Use clang to compile to object
	_, err := runCommand(c.clang, append([]string{"-c", assembly, "-o", object}, args...)...)
	return err
}

// runCommand runs a command and extract its output.
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	fmt.Printf("Running %s %s\n", name, strings.Join(args, " "))

	output, err := cmd.CombinedOutput()
	if err == nil {
		return string(output), nil
	}

	switch {
	case output != nil:
		return "", errors.New(string(output))
	default:
		return "", err
	}
}
