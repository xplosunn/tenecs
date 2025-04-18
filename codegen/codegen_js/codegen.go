package codegen_js

import (
	"fmt"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_js/standard_library"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/exp/maps"
	"strconv"
	"strings"
)

func GenerateProgramNonRunnable(program *ast.Program) string {
	return generateProgram(program)
}

func GenerateHtmlPageForWebApp(program *ast.Program, targetWebApp ast.Ref, cssFiles []string) string {
	cssChildren := ""
	for _, cssFile := range cssFiles {
		cssChildren += fmt.Sprintf(`<link rel="stylesheet" type="text/css" href="%s">`, cssFile)
	}
	return generateTagWithoutAttributes(
		"html",
		generateTagWithoutAttributes(
			"head",
			cssChildren,
		)+generateTagWithoutAttributes(
			"body",
			`<div id="toplevel_tenecs_webapp_container"></div>`+generateTagWithoutAttributes(
				"script",
				generateJsOfWebApp(program, targetWebApp),
			),
		),
	)
}

func generateTagWithoutAttributes(tagName string, children string) string {
	return fmt.Sprintf("<%s>%s</%s>", tagName, children, tagName)
}

func generateJsOfWebApp(program *ast.Program, targetWebApp ast.Ref) string {
	result := GenerateProgramNonRunnable(program) + "\n"
	result += generateWebAppJsMain(targetWebApp.Package, targetWebApp.Name)
	return result
}

func generateWebAppJsMain(pkgName string, targetWebApp string) string {
	webAppVarName := variableName(&pkgName, targetWebApp)
	return fmt.Sprintf(`const webApp = %s

let webAppState = webApp.init()

function renderCurrentWebAppState() {
  const element = document.getElementById("toplevel_tenecs_webapp_container");
  element.innerHTML = render(webApp.view(webAppState))
}

function updateState(event) {
  webAppState = webApp.update(webAppState, event)
  renderCurrentWebAppState()
}

function render(htmlElement) {
  let result = "<" + htmlElement.name
  for (const property of htmlElement.properties) {
    result += " " + property.name + "="
    if (typeof property.value == "string") {
      result += "\"" + property.value + "\""
    } else {
      result += "\"updateState((" + property.value + ")())\""
    }
  }
  result += ">"
  if (typeof htmlElement.children == "string") {
    result += htmlElement.children
  } else {
    for(const child of htmlElement.children) {
      result += render(child)
    }
  }
  result += "</" + htmlElement.name + ">"
  return result
}

renderCurrentWebAppState()
`, webAppVarName)
}

func GenerateProgramTest(program *ast.Program, foundTests codegen.FoundTests) string {
	result := GenerateProgramNonRunnable(program)
	result += "\n"

	testRunnerTestSuiteArgs := ""
	for i, v := range foundTests.UnitTestSuites {
		if i > 0 {
			testRunnerTestSuiteArgs += ", "
		}
		testRunnerTestSuiteArgs += variableName(&v.Package, v.Name)
	}
	testRunnerTestArgs := ""
	for i, v := range foundTests.UnitTests {
		if i > 0 {
			testRunnerTestArgs += ", "
		}
		testRunnerTestArgs += variableName(&v.Package, v.Name)
	}

	result += generateNodeTestRunner()

	result += fmt.Sprintf(`
runUnitTests([%s], [%s])
`, testRunnerTestSuiteArgs, testRunnerTestArgs)
	return result
}

func generateProgram(program *ast.Program) string {
	programDeclarationNames := []ast.Ref{}
	for decName, _ := range program.Declarations {
		programDeclarationNames = append(programDeclarationNames, decName)
	}
	ast.SortRefs(programDeclarationNames)

	decs := ""
	for _, declarationName := range programDeclarationNames {
		for decName, decExp := range program.Declarations {
			if decName != declarationName {
				continue
			}
			dec := generateDeclaration(&decName.Package, &ast.Declaration{
				Name:       decName.Name,
				Expression: decExp,
			})
			decs += dec + "\n"
		}
	}

	structNames := maps.Keys(program.StructFunctions)
	ast.SortRefs(structNames)
	for _, structFuncName := range structNames {
		structFunc := program.StructFunctions[structFuncName]
		decs += generateStructFunction(&structFuncName.Package, structFuncName.Name, structFunc) + "\n"
	}

	nativeFuncNames := maps.Keys(program.NativeFunctions)
	ast.SortRefs(nativeFuncNames)
	for _, nativeFuncName := range nativeFuncNames {
		f := standard_library.Functions[nativeFuncName.Package+"_"+nativeFuncName.Name]
		if len(f.Code) == 0 {
			panic("failed to find function " + nativeFuncName.Package + "_" + nativeFuncName.Name)
		}
		decs += fmt.Sprintf("function %s%s", variableName(&nativeFuncName.Package, nativeFuncName.Name), f.Code) + "\n"
	}

	main := ""

	result := decs + "\n" + main

	return result
}

func generateStructFunction(pkgName *string, name string, structFunc *types.Function) string {
	result := "function " + variableName(pkgName, name)
	result += "("
	for i, argument := range structFunc.Arguments {
		if i > 0 {
			result += ", "
		}
		result += argument.Name
	}
	result += ") {\n"
	result += "return ({"
	result += `  "$type": "` + name + "\""
	for _, argument := range structFunc.Arguments {
		result += ","
		result += "\n" + argument.Name + ": " + argument.Name
	}
	result += "})\n"
	result += "}"
	return result
}

func generateDeclaration(pkgName *string, declaration *ast.Declaration) string {
	_, _, _, _, caseFunction, _, _, _, _ := declaration.Expression.ExpressionCases()
	if caseFunction != nil {
		result := "function " + variableName(pkgName, declaration.Name)
		result += generateFunction(pkgName, *caseFunction, false)
		return result
	} else {
		result := "let " + variableName(pkgName, declaration.Name) + " = "
		result += generateExpression(pkgName, declaration.Expression)
		return result
	}
}

func generateExpression(pkgName *string, expression ast.Expression) string {
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
		return generateLiteral(*caseLiteral)
	} else if caseReference != nil {
		return generateReference(pkgName, *caseReference)
	} else if caseAccess != nil {
		return generateAccess(pkgName, *caseAccess)
	} else if caseInvocation != nil {
		return generateInvocation(pkgName, *caseInvocation)
	} else if caseFunction != nil {
		return generateFunction(pkgName, *caseFunction, true)
	} else if caseDeclaration != nil {
		return generateDeclaration(pkgName, caseDeclaration)
	} else if caseIf != nil {
		return generateIf(pkgName, *caseIf)
	} else if caseList != nil {
		return generateList(pkgName, *caseList)
	} else if caseWhen != nil {
		return generateWhen(pkgName, *caseWhen)
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func generateFunction(pkgName *string, function ast.Function, includeArrow bool) string {
	result := "("
	for i, argument := range function.VariableType.Arguments {
		if i > 0 {
			result += ", "
		}
		result += variableName(pkgName, argument.Name)
	}
	if includeArrow {
		result += ") => "
	} else {
		result += ") "
	}
	result += generateBlock(pkgName, function.Block)
	return result
}

func generateBlock(pkgName *string, block []ast.Expression) string {
	result := "{"
	result += generateExpressionsWithinBlock(pkgName, block)
	result += "\n}"
	return result
}

func generateExpressionsWithinBlock(pkgName *string, block []ast.Expression) string {
	result := ""
	if len(block) == 0 {
		result += "\n return null"
	} else {
		for i, expression := range block {
			if i < len(block)-1 {
				result += "\n" + generateExpression(pkgName, expression)
			} else {
				result += "\n" + "return " + generateExpression(pkgName, expression)
			}

		}
	}
	return result
}

func generateInvocation(pkgName *string, invocation ast.Invocation) string {
	result := generateExpression(pkgName, invocation.Over)
	result += "("
	for i, argument := range invocation.Arguments {
		if i > 0 {
			result += ", "
		}
		result += generateExpression(pkgName, argument)
	}
	result += ")"
	return result
}

func generateAccess(pkgName *string, access ast.Access) string {
	result := generateExpression(pkgName, access.Over)
	result += "." + access.Access
	return result
}

func generateReference(pkgName *string, reference ast.Reference) string {
	if reference.PackageName != nil {
		pkgName = reference.PackageName
	}
	return variableName(pkgName, reference.Name)
}

func generateWhen(pkgName *string, when ast.When) string {
	result := "(() => {\n"
	result += "let __over = " + generateExpression(pkgName, when.Over) + "\n"
	for _, whenCase := range when.Cases {
		result += "if (" + generateWhenClause(whenCase.VariableType, "__over") + ") {\n"
		varName := whenCase.Name
		if varName != nil {
			result += "let " + variableName(pkgName, *varName) + " = __over\n"
		}
		result += generateExpressionsWithinBlock(pkgName, whenCase.Block) + "\n"
		result += "}\n"
	}
	result += "})()"
	return result
}

func generateWhenClause(variableType types.VariableType, varName string) string {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO generateWhenClause caseTypeArgument")
	} else if caseList != nil {
		panic("TODO generateWhenClause caseList")
	} else if caseKnownType != nil {
		return generateWhenClauseKnownType(*caseKnownType, varName)
	} else if caseFunction != nil {
		panic("TODO generateWhenClause caseFunction")
	} else if caseOr != nil {
		panic("TODO generateWhenClause caseOr")
	} else {
		panic("cases on variableType")
	}
}

func generateWhenClauseKnownType(knownType types.KnownType, varName string) string {
	if knownType.Package == "" {
		if knownType.Name == "String" {
			return "typeof " + varName + `=== "string"`
		} else if knownType.Name == "Boolean" {
			return "typeof " + varName + `=== "boolean"`
		} else {
			panic("TODO generateWhenClauseKnownType " + knownType.Name)
		}
	} else {
		return fmt.Sprintf(`typeof %s === "object" && %s["$type"] === "%s"`, varName, varName, knownType.Name)
	}
}

func generateList(pkgName *string, list ast.List) string {
	result := "["
	for i, expression := range list.Arguments {
		if i > 0 {
			result += ", "
		}
		result += generateExpression(pkgName, expression)
	}
	result += "]"
	return result
}

func generateIf(pkgName *string, astIf ast.If) string {
	condition := generateExpression(pkgName, astIf.Condition)
	thenBlock := generateBlock(pkgName, astIf.ThenBlock)
	elseBlock := generateBlock(pkgName, astIf.ElseBlock)
	return "(" + condition + ") ? (() => " + thenBlock + ")() : (() => " + elseBlock + ")()"
}

func generateLiteral(literal ast.Literal) string {
	result := ""
	parser.LiteralExhaustiveSwitch(
		literal.Literal,
		func(literal float64) { result = fmt.Sprintf("%f", literal) },
		func(literal int) { result = fmt.Sprintf("%d", literal) },
		func(literal string) { result = literal },
		func(literal bool) { result = strconv.FormatBool(literal) },
		func() { result = "null" },
	)
	return result
}

func variableName(pkgName *string, name string) string {
	prefix := ""
	if pkgName != nil {
		prefix = strings.ReplaceAll(*pkgName, ".", "_") + "__"
	}
	return prefix + name
}
