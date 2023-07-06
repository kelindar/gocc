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

package asm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/kelindar/gocc/internal/config"
)

// ParseAssembly parses the assembly file and returns a list of functions
func ParseAssembly(arch *config.Arch, path string) ([]Function, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}(file)

	var (
		functions    = make([]Function, 0, 8)
		current      *Function
		constant     *Const
		functionName string
		labelName    string
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {

		// Handle constant lines and attach them to the current label
		case arch.Const.MatchString(line):
			constant.Lines = append(constant.Lines, parseConst(arch, line))

		// Skip attirubtes and comment lines
		case arch.Attribute.MatchString(line):
			continue
		case arch.Comment.MatchString(line):
			continue

		// Handle assembly labels. We could potentially have multiple labels per line if
		// compiler decides to generate no-op instructions.
		case arch.Label.MatchString(line):
			labelName = strings.Split(line, ":")[0]
			labelName = labelName[1:]
			constant = &Const{Label: labelName} // reset the current constant
			switch {
			case current == nil: // No function yet
			case len(current.Lines) == 0:
				current.Lines = append(current.Lines, Line{Labels: []string{labelName}})
			case current.Lines[len(current.Lines)-1].Assembly == "": // Previous line was a label
				current.Lines[len(current.Lines)-1].Labels = append(current.Lines[len(current.Lines)-1].Labels, labelName)
			default:
				current.Lines = append(current.Lines, Line{Labels: []string{labelName}})
			}

		// Handle assembly function name
		case arch.Function.MatchString(line):
			functionName = strings.Split(line, ":")[0]
			functions = append(functions, Function{
				Name:  functionName,
				Lines: make([]Line, 0),
			})
			current = &functions[len(functions)-1]
			labelName = "" // Reset current label

			// If we have a constant, attach it to the current function
			if len(constant.Lines) > 0 {
				current.Consts = append(current.Consts, *constant)
			}

		// Handle assembly instructions
		case arch.Code.MatchString(line):
			code := strings.Split(line, arch.CommentCh)[0]
			code = strings.TrimSpace(code)
			if labelName == "" {
				current.Lines = append(current.Lines, Line{Assembly: code})
			} else {
				current.Lines[len(current.Lines)-1].Assembly = code
				labelName = ""
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return functions, nil
}

// ParseObjectDump parses the output of objdump file and returns a list of functions
func ParseObjectDump(arch *config.Arch, dump string, functions []Function) error {
	var (
		functionName string
		functionIdx  int
		current      *Function
		lineNumber   int
	)

	for i, line := range strings.Split(dump, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case arch.Symbol.MatchString(line):
			functionName = strings.Split(line, "<")[1]
			functionName = strings.Split(functionName, ">")[0]
			current = &functions[functionIdx]
			lineNumber = 0
			functionIdx++
		case arch.Data.MatchString(line):
			data := strings.Split(line, ":")[1]
			data = strings.TrimSpace(data)
			splits := strings.Split(data, " ")
			var (
				binary   []string
				assembly string
			)

			for i, s := range splits {
				if s == "" || unicode.IsSpace(rune(s[0])) {
					assembly = strings.Join(splits[i:], " ")
					assembly = strings.TrimSpace(assembly)
					break
				}

				// If the binary representation is not separated with spaces, split it
				switch {
				case len(s) > 2:
					// Iterate backwards
					for i := len(s) - 2; i >= 0; i -= 2 {
						binary = append(binary, s[i:i+2])
					}
				default:
					binary = append(binary, s)
				}
			}

			switch {
			case assembly == "":
				return fmt.Errorf("try to increase --insn-width of objdump")
			case strings.HasPrefix(assembly, "nop"):
				continue
			case assembly == "xchg   %ax,%ax":
				continue
			case strings.HasPrefix(assembly, "cs nopw"):
				continue
			case lineNumber >= len(current.Lines):
				return fmt.Errorf("%d: unexpected objectdump line: %s, please compare assembly with objdump output", i, line)
			}

			current.Lines[lineNumber].Binary = binary
			lineNumber++
		}
	}
	return nil
}
