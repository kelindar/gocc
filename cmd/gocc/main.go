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

	"github.com/kelindar/gocc"
	"github.com/kelindar/gocc/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	command.PersistentFlags().StringP("output", "o", "", "output directory of generated files")
	command.PersistentFlags().StringSliceP("machine-option", "m", nil, "machine option for clang")
	command.PersistentFlags().IntP("optimize-level", "O", 0, "optimization level for clang")
	command.PersistentFlags().StringP("arch", "a", "amd64", "target architecture to use")
	command.PersistentFlags().StringP("package", "p", "", "go package name to use for the stubs")
	command.PersistentFlags().BoolP("local", "l", false, "use local machine for compilation")
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
				exit(err)
			}
		}

		var options []string
		machineOptions, _ := cmd.PersistentFlags().GetStringSlice("machine-option")
		for _, m := range machineOptions {
			options = append(options, "-m"+m)
		}

		// Load the architecture
		target, _ := cmd.PersistentFlags().GetString("arch")
		optimizeLevel, _ := cmd.PersistentFlags().GetInt("optimize-level")
		packageName, _ := cmd.PersistentFlags().GetString("package")
		options = append(options, fmt.Sprintf("-O%d", optimizeLevel))

		// Compile locally or remotely
		local, _ := cmd.PersistentFlags().GetBool("local")
		switch local {
		case true:
			if err := compileLocally(target, args[0], output, packageName, options...); err != nil {
				exit(err)
			}
		default:
			if err := compileRemotely(target, args[0], output, packageName, options...); err != nil {
				exit(err)
			}
		}
	},
}

func compileRemotely(target, source, outputDir, packageName string, options ...string) error {
	remote, err := gocc.NewRemote(target, source, outputDir, packageName, options...)
	if err != nil {
		return err
	}

	return remote.Translate()
}

func compileLocally(target, source, outputDir, packageName string, options ...string) error {
	arch, err := config.For(target)
	if err != nil {
		exit(err)
	}

	local, err := gocc.NewLocal(arch, source, outputDir, packageName, options...)
	if err != nil {
		return err
	}

	return local.Translate()
}

func exit(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
