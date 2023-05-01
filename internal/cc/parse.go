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
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/kelindar/gocc/internal/asm"
	"modernc.org/cc/v3"
)

var supportedTypes = mapset.NewSet("int64_t", "uint64_t", "float", "unsignedlonglong", "longlong", "long", "unsignedlong", "int", "unsignedint")

// Parse parse C source file and extracts functions declarations.
func Parse(path string) ([]asm.Function, error) {
	source, err := redactSource(path)
	if err != nil {
		return nil, err
	}

	ast, err := cc.Parse(&cc.Config{}, nil, nil,
		[]cc.Source{{Name: path, Value: source}})
	if err != nil {
		return nil, err
	}

	var functions []asm.Function
	for _, nodes := range ast.Scope {
		if len(nodes) != 1 || nodes[0].Position().Filename != path {
			continue
		}

		node := nodes[0]
		if declarator, ok := node.(*cc.Declarator); ok {
			funcIdent := declarator.DirectDeclarator
			if funcIdent.Case != cc.DirectDeclaratorFuncParam {
				continue
			}

			if function, err := convertFunction(funcIdent); err != nil {
				return nil, err
			} else {
				functions = append(functions, function)
			}
		}
	}
	sort.Slice(functions, func(i, j int) bool {
		return functions[i].Position < functions[j].Position
	})
	return functions, nil
}

// redactSource removes code from the source and only leaves function declarations.
// This is done to avoid parsing errors when the source is not compatible with the compiler.
func redactSource(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var src strings.Builder
	src.WriteString("#define __STDC_HOSTED__ 1\n")
	src.WriteString("#define uint64_t unsigned long long\n")
	src.WriteString("#define uint32_t unsigned int\n")
	src.WriteString("#define int64_t long long\n")
	src.WriteString("#define int32_t int\n")

	var clauseCount int
	for _, line := range strings.Split(string(bytes), "\n") {
		switch {
		case strings.HasPrefix(line, "#include"):
			continue
		case strings.HasPrefix(line, "//"):
			continue
		case strings.Contains(line, "{"):
			if clauseCount == 0 {
				src.WriteString(line[:strings.Index(line, "{")+1])
				src.WriteString("\n // removed for compatibility\n")
			}
			clauseCount++
		case strings.Contains(line, "}"):
			clauseCount--
			if clauseCount == 0 {
				src.WriteString(line[strings.Index(line, "}"):])
				src.WriteRune('\n')
			}
		default:
			continue
		}
	}

	return src.String(), nil
}

// convertFunction extracts the function definition from cc.DirectDeclarator.
func convertFunction(declarator *cc.DirectDeclarator) (asm.Function, error) {
	params, err := convertFunctionParameters(declarator.ParameterTypeList.ParameterList)
	if err != nil {
		return asm.Function{}, err
	}

	return asm.Function{
		Name:       declarator.DirectDeclarator.Token.String(),
		Position:   declarator.Position().Line,
		Parameters: params,
	}, nil
}

// convertFunctionParameters extracts function parameters from cc.ParameterList.
func convertFunctionParameters(params *cc.ParameterList) ([]asm.Param, error) {
	declaration := params.ParameterDeclaration
	isPointer := declaration.Declarator.Pointer != nil
	paramName := declaration.Declarator.DirectDeclarator.Token.Value
	paramType := typeOf(declaration.DeclarationSpecifiers)

	if !isPointer && !supportedTypes.Contains(paramType) {
		position := declaration.Position()
		return nil, fmt.Errorf("gocc: [%v] unsupported type: %v\n",
			position.Filename, paramType)
	}

	paramNames := []asm.Param{{
		Name:      paramName.String(),
		Type:      paramType,
		IsPointer: isPointer,
	}}

	if params.ParameterList != nil {
		if nextParamNames, err := convertFunctionParameters(params.ParameterList); err != nil {
			return nil, err
		} else {
			paramNames = append(paramNames, nextParamNames...)
		}
	}
	return paramNames, nil
}

// typeOf returns the type of the given value, recursively.
func typeOf(v any) string {
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Ptr && rv.IsNil() {
		return ""
	}

	switch s := v.(type) {
	case *cc.TypeQualifier:
		return s.Token.String()
	case *cc.TypeSpecifier:
		return s.Token.String()
	case *cc.DeclarationSpecifiers:
		var result string
		switch s.Case {
		case cc.DeclarationSpecifiersTypeQual:
			result += typeOf(s.TypeSpecifier)
			result += typeOf(s.DeclarationSpecifiers)
		case cc.DeclarationSpecifiersTypeSpec:
			result += typeOf(s.TypeSpecifier)
			result += typeOf(s.DeclarationSpecifiers)
		default:
			panic(fmt.Sprintf("gocc: unexpected specifiers case: %v", s.Case))
		}
		return result
	default:
		panic(fmt.Sprintf("gocc: unexpected specifier type: %T", v))
	}
}
