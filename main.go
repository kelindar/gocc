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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kelindar/gocc/internal/config"
	"github.com/spf13/cobra"
)

// listIncludePaths lists include paths used by clang.
func listIncludePaths() ([]string, error) {
	out, err := runCommand("bash", "-c", "echo | gcc -xc -E -v -")
	if err != nil {
		return nil, err
	}

	var start bool
	var paths []string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "#include <...> search starts here:") {
			start = true
		} else if strings.HasPrefix(line, "End of search list.") {
			start = false
		} else if start {
			path := strings.TrimSpace(line)
			paths = append(paths, path)
		}
	}

	return paths, nil
}

// runCommand runs a command and extract its output.
func runCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			return "", errors.New(string(output))
		} else {
			return "", err
		}
	}
	return string(output), nil
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
		file := NewTranslateUnit(arch, args[0], output, options...)
		if err := file.Translate(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

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
