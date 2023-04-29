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
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/klauspost/asmfmt"
)

const buildTags = "//go:build !noasm && arm64\n"

var (
	commentLine   = regexp.MustCompile(`^\s*;.*$`)
	attributeLine = regexp.MustCompile(`^\s+\..+$`)
	nameLine      = regexp.MustCompile(`^\w+:.+$`)
	labelLine     = regexp.MustCompile(`^[A-Z0-9]+_\d+:.*$`)
	codeLine      = regexp.MustCompile(`^\s+\w+.+$`)
	symbolLine    = regexp.MustCompile(`^\w+\s+<\w+>:$`)
	dataLine      = regexp.MustCompile(`^\w+:\s+\w+\s+.+$`)
	registers     = []string{"R0", "R1", "R2", "R3"}
)

type Line struct {
	Label    string
	Assembly string
	Binary   []string
}

func (line *Line) String() string {
	var builder strings.Builder
	if len(line.Label) > 0 {
		builder.WriteString(line.Label)
		builder.WriteString(":\n")
	}

	builder.WriteString("\t")
	if strings.HasPrefix(line.Assembly, "j") {
		splits := strings.Split(line.Assembly, ".")
		op := strings.TrimSpace(splits[0])
		operand := splits[1]
		builder.WriteString(fmt.Sprintf("%s %s", strings.ToUpper(op), operand))
	} else {
		pos := 0
		for pos < len(line.Binary) {
			if pos > 0 {
				builder.WriteString("; ")
			}
			if len(line.Binary)-pos >= 8 {
				builder.WriteString(fmt.Sprintf("QUAD $0x%v%v%v%v%v%v%v%v",
					line.Binary[pos+7], line.Binary[pos+6], line.Binary[pos+5], line.Binary[pos+4],
					line.Binary[pos+3], line.Binary[pos+2], line.Binary[pos+1], line.Binary[pos]))
				pos += 8
			} else if len(line.Binary)-pos >= 4 {
				builder.WriteString(fmt.Sprintf("LONG $0x%v%v%v%v",
					line.Binary[pos+3], line.Binary[pos+2], line.Binary[pos+1], line.Binary[pos]))
				pos += 4
			} else if len(line.Binary)-pos >= 2 {
				builder.WriteString(fmt.Sprintf("WORD $0x%v%v", line.Binary[pos+1], line.Binary[pos]))
				pos += 2
			} else {
				builder.WriteString(fmt.Sprintf("BYTE $0x%v", line.Binary[pos]))
				pos += 1
			}
		}
		builder.WriteString("\t// ")
		builder.WriteString(line.Assembly)
	}
	builder.WriteString("\n")
	return builder.String()
}

func parseAssembly(path string) ([]Function, error) {
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
		functionName string
		labelName    string
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case attributeLine.MatchString(line):
			continue
		case commentLine.MatchString(line):
			continue
		case labelLine.MatchString(line):
			labelName = strings.Split(line, ":")[0]
			labelName = labelName[1:]
			current.Lines = append(current.Lines, Line{Label: labelName})
		case nameLine.MatchString(line):
			functionName = strings.Split(line, ":")[0]
			functions = append(functions, Function{
				Name:  functionName,
				Lines: make([]Line, 0),
			})
			current = &functions[len(functions)-1]
		case codeLine.MatchString(line):
			asm := strings.Split(line, ";")[0]
			asm = strings.TrimSpace(asm)
			if labelName == "" {
				current.Lines = append(current.Lines, Line{Assembly: asm})
			} else {
				current.Lines[len(current.Lines)-1].Assembly = asm
				labelName = ""
			}
		default:
			// println("unexpected line: " + line)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return functions, nil
}

func parseObjectDump(dump string, functions []Function) error {
	var (
		functionName string
		functionIdx  int
		current      *Function
		lineNumber   int
	)

	for i, line := range strings.Split(dump, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case symbolLine.MatchString(line):
			functionName = strings.Split(line, "<")[1]
			functionName = strings.Split(functionName, ">")[0]
			current = &functions[functionIdx]
			lineNumber = 0
			functionIdx++
		case dataLine.MatchString(line):
			data := strings.Split(line, ":")[1]
			data = strings.TrimSpace(data)
			splits := strings.Split(data, " ")

			var binary []string
			var assembly string
			for i, s := range splits {
				if s == "" || unicode.IsSpace(rune(s[0])) {
					assembly = strings.Join(splits[i:], " ")
					assembly = strings.TrimSpace(assembly)
					break
				}
				binary = append(binary, s)
			}

			if lineNumber >= len(current.Lines) {
				return fmt.Errorf("%d: unexpected objectdump line: %s", i, line)
			}

			current.Lines[lineNumber].Binary = binary
			lineNumber++
		}
	}
	return nil
}

func generateGoAssembly(path string, functions []Function) error {
	// generate code
	var builder strings.Builder
	builder.WriteString(buildTags)
	builder.WriteString("// AUTO-GENERATED BY GOCC -- DO NOT EDIT\n")
	for _, function := range functions {
		builder.WriteString(fmt.Sprintf("\nTEXT ·%v(SB), $0-32\n", function.Name))
		for i, param := range function.Parameters {
			builder.WriteString(fmt.Sprintf("\tMOVD %s+%d(FP), %s\n", param, i*8, registers[i]))
		}

		for _, line := range function.Lines {
			builder.WriteString(line.String())
		}
	}

	// write file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}(f)
	bytes, err := asmfmt.Format(strings.NewReader(builder.String()))
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}
