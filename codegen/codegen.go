package codegen

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"strconv"
)

type Import string
type MainDeclaration string
type IsMain bool

func Generate(program *ast.Program) string {
	mainDeclarations := []MainDeclaration{}

	decs := ""
	allImports := []Import{}
	for _, declaration := range program.Declarations {
		mainDeclaration, imports, dec := GenerateDeclaration(declaration)
		decs += dec + "\n"
		allImports = append(allImports, imports...)
		if mainDeclarations != nil {
			mainDeclarations = append(mainDeclarations, *mainDeclaration)
		}
	}

	main := ""

	if len(mainDeclarations) > 1 {
		panic("TODO Generate multiple mains")
	} else if len(mainDeclarations) == 1 {
		imports, mainCode := GenerateMain(string(mainDeclarations[0]))
		main = mainCode
		allImports = append(allImports, imports...)
	}

	imports := "import (\n"
	for _, importPkg := range allImports {
		imports += fmt.Sprintf(`	"%s"`, importPkg) + "\n"
	}
	imports += ")\n"

	result := "package main\n\n" + imports + "\n" + decs + "\n" + main

	return result
}

func GenerateMain(varToInvoke string) ([]Import, string) {
	imports, runtime := GenerateRuntime()
	return imports, fmt.Sprintf(`func main() {
r := runtime()
%s.(map[string]any)["main"].(func(any)any)(r)
}

func runtime() map[string]any {
return %s
}
`, varToInvoke, runtime)
}

func VariableName(name string) string {
	return "P" + name
}

func GenerateDeclaration(declaration *ast.Declaration) (*MainDeclaration, []Import, string) {
	isMain, imports, exp := GenerateExpression(declaration.Expression)
	varName := VariableName(declaration.Name)
	result := "var " + varName + " any = " + exp + "\n"
	var mainDeclaration *MainDeclaration
	if isMain {
		m := MainDeclaration(varName)
		mainDeclaration = &m
	}
	return mainDeclaration, imports, result
}

func GenerateExpression(expression ast.Expression) (IsMain, []Import, string) {
	caseModule, caseLiteral, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf, caseArray, caseWhen := expression.ExpressionCases()
	if caseModule != nil {
		return GenerateModule(*caseModule)
	} else if caseLiteral != nil {
		return false, []Import{}, GenerateLiteral(*caseLiteral)
	} else if caseReferenceAndMaybeInvocation != nil {
		imports, result := GenerateReferenceAndMaybeInvocation(*caseReferenceAndMaybeInvocation)
		return false, imports, result
	} else if caseWithAccessAndMaybeInvocation != nil {
		imports, result := GenerateWithAccessAndMaybeInvocation(*caseWithAccessAndMaybeInvocation)
		return false, imports, result
	} else if caseFunction != nil {
		imports, result := GenerateFunction(*caseFunction)
		return false, imports, result
	} else if caseDeclaration != nil {
		panic("TODO GenerateExpression caseDeclaration")
	} else if caseIf != nil {
		panic("TODO GenerateExpression caseIf")
	} else if caseArray != nil {
		panic("TODO GenerateExpression caseArray")
	} else if caseWhen != nil {
		panic("TODO GenerateExpression caseWhen")
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func GenerateLiteral(literal ast.Literal) string {
	result := ""
	parser.LiteralExhaustiveSwitch(
		literal.Literal,
		func(literal float64) { result = fmt.Sprintf("%f", literal) },
		func(literal int) { result = fmt.Sprintf("%d", literal) },
		func(literal string) { result = literal },
		func(literal bool) { result = strconv.FormatBool(literal) },
	)
	return result
}

func GenerateReferenceAndMaybeInvocation(referenceAndMaybeInvocation ast.ReferenceAndMaybeInvocation) ([]Import, string) {
	allImports := []Import{}
	result := VariableName(referenceAndMaybeInvocation.Name)

	if referenceAndMaybeInvocation.ArgumentsList != nil {
		funcArgList := ""
		argsCode := ""
		for i, argument := range referenceAndMaybeInvocation.ArgumentsList.Arguments {
			if i > 0 {
				funcArgList += ","
				argsCode += ", "
			}
			funcArgList += "any"

			_, imports, arg := GenerateExpression(argument)
			allImports = append(allImports, imports...)
			argsCode += arg
		}
		result += fmt.Sprintf(`.(func(%s)any)(%s)`, funcArgList, argsCode)
	}

	return allImports, result
}

func GenerateWithAccessAndMaybeInvocation(withAccessAndMaybeInvocation ast.WithAccessAndMaybeInvocation) ([]Import, string) {
	allImports := []Import{}
	result := ""

	_, imports, over := GenerateExpression(withAccessAndMaybeInvocation.Over)
	allImports = append(allImports, imports...)
	result += over

	for _, accessAndMaybeInvocation := range withAccessAndMaybeInvocation.AccessChain {
		result += fmt.Sprintf(`.(map[string]any)["%s"]`, accessAndMaybeInvocation.Access)
		if accessAndMaybeInvocation.ArgumentsList != nil {
			funcArgList := ""
			argsCode := ""
			for i, argument := range accessAndMaybeInvocation.ArgumentsList.Arguments {
				if i > 0 {
					funcArgList += ","
					argsCode += ", "
				}
				funcArgList += "any"

				_, imports, arg := GenerateExpression(argument)
				allImports = append(allImports, imports...)
				argsCode += arg
			}
			result += fmt.Sprintf(`.(func(%s)any)(%s)`, funcArgList, argsCode)
		}
	}

	return allImports, result
}
func GenerateFunction(function ast.Function) ([]Import, string) {
	allImports := []Import{}
	args := ""
	for i, argument := range function.VariableType.Arguments {
		args += VariableName(argument.Name) + " any"
		if i < len(function.VariableType.Arguments)-1 {
			args += ", "
		}
	}
	result := fmt.Sprintf("func (%s) any {\n", args)

	for _, expression := range function.Block {
		_, imports, exp := GenerateExpression(expression)
		result += exp + "\n"
		allImports = append(allImports, imports...)
	}

	result += "return nil\n"

	result += "}"
	return allImports, result
}

func GenerateModule(module ast.Module) (IsMain, []Import, string) {
	isMain := module.Implements.Package == "tenecs.os" && module.Implements.Name == "Main"
	allImports := []Import{}
	result := "map[string]any{\n"
	for variableName, exp := range module.Variables {
		_, imports, expStr := GenerateExpression(exp)
		result += fmt.Sprintf(`"%s": %s,`+"\n", variableName, expStr)
		allImports = append(allImports, imports...)
	}
	result += "}"
	return IsMain(isMain), allImports, result
}
