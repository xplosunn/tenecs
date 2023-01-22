package javascript

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"strconv"
)

func Codegen(program ast.Program) (string, error) {
	result := ""
	for _, module := range program.Modules {
		result += "\n" + codegenModule(*module) + "\n"
	}

	//TODO stdlib
	//TODO find main and invoke it

	return result, nil
}

func codegenModule(module ast.Module) string {
	commaSeparatedConstructorArgNames := ""
	for i, constructorArgument := range module.ConstructorArguments {
		commaSeparatedConstructorArgNames += "_" + constructorArgument.Name
		if i > 0 {
			commaSeparatedConstructorArgNames += ", "
		}
	}
	result := fmt.Sprintf("const _%s = (%s) => {\n", module.Name, commaSeparatedConstructorArgNames)
	for variableName, variableExpression := range module.Variables {
		expressionCode := codegenExpression(variableExpression)
		result += fmt.Sprintf("%s: %s", variableName, expressionCode)
	}
	result += "}\n"
	return result
}

func codegenExpression(expression ast.Expression) string {
	caseLiteral, caseReferenceOrInvocation, caseFunction, caseDeclaration, caseIf := expression.Cases()
	if caseLiteral != nil {
		return codegenLiteral(*caseLiteral)
	} else if caseReferenceOrInvocation != nil {
		return codegenReferenceOrInvocation(*caseReferenceOrInvocation)
	} else if caseFunction != nil {
		return codegenFunction(*caseFunction)
	} else if caseDeclaration != nil {
		return codegenDeclaration(*caseDeclaration)
	} else if caseIf != nil {
		return codegenIf(*caseIf)
	} else {
		panic(fmt.Errorf("code on %v", expression))
	}
}

func codegenLiteral(expression ast.Literal) string {
	return parser.LiteralFold(
		expression.Literal,
		func(arg float64) string { return fmt.Sprintf("%f", arg) },
		func(arg int) string { return fmt.Sprintf("%d", arg) },
		func(arg string) string { return arg },
		func(arg bool) string { return strconv.FormatBool(arg) },
	)
}

func codegenReferenceOrInvocation(expression ast.WithAccessAndMaybeInvocation) string {
	result := ""
	for i, varName := range expression.DotSeparatedVars {
		if i > 0 {
			result += "."
		}
		result += "_" + varName
	}
	if expression.Arguments != nil {
		result += "("
		for i, arg := range expression.Arguments.Arguments {
			if i > 0 {
				result += ", "
			}
			result += codegenExpression(arg)
		}
		result += ")"
	}
	return result
}

func codegenFunction(expression ast.Function) string {
	result := "("
	for i, argument := range expression.VariableType.Arguments {
		if i > 0 {
			result += ", "
		}
		result += argument.Name
	}
	result += ") => {\n"
	for i, exp := range expression.Block {
		if i == len(expression.Block)-1 {
			result += "return "
		}
		result += codegenExpression(exp) + "\n"
	}
	result += "}\n"
	return result
}

func codegenDeclaration(expression ast.Declaration) string {
	result := fmt.Sprintf("const _%s = ", expression.Name)
	result += codegenExpression(expression.Expression)
	result += "\n"
	return result
}

func codegenIf(expression ast.If) string {
	panic("todo")
}
