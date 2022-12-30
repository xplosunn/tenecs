package javascript

import (
	"errors"
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"strconv"
)

func Codegen(parsed parser.FileTopLevel) (string, error) {
	pkg, imports, topLevelDeclarations := parser.FileTopLevelFields(parsed)
	_ = pkg
	_ = imports

	result := ""

	var moduleNameWithPublicMain *string

	for _, module := range topLevelDeclarations {
		implementing, moduleName, constructorArgs, declarations := parser.ModuleFields(*module.(*parser.Module))
		_ = implementing
		_ = constructorArgs
		result += "const _" + moduleName + " = {\n"

		for _, declaration := range declarations {
			public, name, lambda := parser.ModuleDeclarationFields(declaration)

			if public && name == "main" {
				if moduleNameWithPublicMain != nil {
					return "", errors.New("found multiple modules with public main")
				}
				moduleNameWithPublicMain = &moduleName
			}

			result += "_" + name + ": " + codegenLambda(lambda.(parser.Lambda))
		}

		result += "}\n"
	}

	if moduleNameWithPublicMain == nil {
		return "", errors.New("no module with public main found")
	}

	result += codegenStdLib()

	result += "_" + *moduleNameWithPublicMain + "._main($runtime)\n"

	return result, nil
}

func codegenLambda(lambda parser.Lambda) string {
	parameters, returnType, block := parser.LambdaFields(lambda)
	_ = returnType

	result := "("

	for i, parameter := range parameters {
		paramName, paramType := parser.ParameterFields(parameter)
		_ = paramType
		if i > 0 {
			result += ", "
		}
		result += "_" + paramName
	}

	result += ") => {\n"

	for _, invocation := range block {
		dotSeparatedVars, arguments := parser.ReferenceOrInvocationFields(invocation)
		for i, varName := range dotSeparatedVars {
			if i > 0 {
				result += "."
			}
			result += "_" + varName
		}
		result += "("
		for i, argument := range *arguments {
			if i > 0 {
				result += ", "
			}

			result += parser.LiteralFold(
				argument.(parser.LiteralExpression).Literal,
				func(arg float64) string {
					return fmt.Sprintf("%f", arg)
				},
				func(arg int) string {
					return fmt.Sprintf("%d", arg)
				},
				func(arg string) string {
					return arg
				},
				func(arg bool) string {
					return strconv.FormatBool(arg)
				},
			)
		}
		result += ")\n"
	}

	result += "}\n"

	return result
}

func codegenStdLib() string {
	return `const $runtime = {
  _console: {
    _log: (str) => console.log(str)
  }
}
`
}
