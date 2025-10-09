package codegen_javascript

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"

	"golang.org/x/exp/maps"

	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/ir"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
)

type GeneratedProgram struct {
	UserspaceCode       string
	StandardLibraryCode string
	HelpersCode         string
	MainCode            string
}

func (program GeneratedProgram) String() string {
	return program.UserspaceCode + "\n\n" +
		program.StandardLibraryCode + "\n\n" +
		program.HelpersCode + "\n\n" +
		program.MainCode
}

func GenerateProgramNonRunnable(program *ir.Program) GeneratedProgram {
	return generateProgram(false, program, nil)
}

func GenerateProgramTest(program *ir.Program, foundTests codegen.FoundTests) GeneratedProgram {
	return generateProgram(true, program, &foundTests)
}

func GenerateHtmlPageForWebApp(program *ir.Program, targetWebApp ir.Reference, cssFiles []string) string {
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

func generateJsOfWebApp(program *ir.Program, targetWebApp ir.Reference) string {
	result := GenerateProgramNonRunnable(program).String() + "\n"
	result += generateWebAppJsMain(targetWebApp)
	return result
}

func generateWebAppJsMain(webAppVarName ir.Reference) string {
	return fmt.Sprintf(`const webApp = %s()

let webAppState = webApp._init([])

function renderCurrentWebAppState() {
  const element = document.getElementById("toplevel_tenecs_webapp_container");
  element.innerHTML = render(webApp._view([], webAppState))
}

function updateState(event) {
  webAppState = webApp._update([], webAppState, event)
  renderCurrentWebAppState()
}

function render(htmlElement) {
  if (htmlElement["$type"] == "String") {
    return htmlElement._children.value
  }
  let result = "<" + htmlElement._name.value
  for (const property of htmlElement._properties.value) {
    result += " " + property._name.value + "="
    if (property._value["$type"] == "String") {
      result += "\"" + property._value.value + "\""
    } else {
      result += "\"updateState((" + property._value + ")())\""
    }
  }
  result += ">"
  if (htmlElement._children["$type"] == "String") {
    result += htmlElement._children.value
  } else {
    for(const child of htmlElement._children.value) {
      result += render(child)
    }
  }
  result += "</" + htmlElement._name.value + ">"
  return result
}

renderCurrentWebAppState()
`, webAppVarName.Name)
}

//go:embed embed/helpers.js
var helpers string

func generateProgram(testMode bool, program *ir.Program, foundTests *codegen.FoundTests) GeneratedProgram {
	programDeclarationNames := []ir.Reference{}
	for declarationName, _ := range program.Declarations {
		programDeclarationNames = append(programDeclarationNames, declarationName)
	}
	SortReferences(programDeclarationNames)

	userspaceDecs := ""
	for _, declarationName := range programDeclarationNames {
		for decName, decExp := range program.Declarations {
			if decName != declarationName {
				continue
			}
			userspaceDecs += GeneratePackageDeclaration(decName, decExp) + "\n"
		}
	}

	structNames := maps.Keys(program.StructFunctions)
	SortReferences(structNames)
	for _, structFuncName := range structNames {
		structFunc := program.StructFunctions[structFuncName]
		code := GenerateStructFunction(structFunc)
		userspaceDecs += fmt.Sprintf("let %s = %s\n", structFuncName.Name, code)
	}

	standardLibraryCode := ""

	nativeFuncNames := append([]ir.NativeFunctionRef{}, program.NativeFunctions...)
	SortNativeFunctionRefs(nativeFuncNames)
	for _, nativeFuncName := range nativeFuncNames {
		standardLibraryCode += generateNativeFunction(nativeFuncName)
	}

	stdLibStructNames := maps.Keys(program.StdLibStructFunctions)
	SortReferences(stdLibStructNames)
	for _, structFuncName := range stdLibStructNames {
		structFunc := program.StdLibStructFunctions[structFuncName]
		code := GenerateStructFunction(structFunc)
		standardLibraryCode += fmt.Sprintf(`function %s() {
return %s
}
`, structFuncName.Name, code)
	}

	main := ""

	if testMode {
		main = GenerateTestRunnerMain(foundTests.UnitTestSuites, foundTests.UnitTests, foundTests.GoIntegrationTests)
	}

	return GeneratedProgram{
		UserspaceCode:       userspaceDecs,
		StandardLibraryCode: standardLibraryCode,
		HelpersCode:         helpers,
		MainCode:            main,
	}
}

func GenerateTestRunnerMain(
	varsImplementingUnitTestSuite []ast.Ref,
	varsImplementingUnitTest []ast.Ref,
	varsImplementingGoIntegrationTest []ast.Ref,
) string {
	testRunnerUnitTestSuiteArgs := ""
	for i, v := range varsImplementingUnitTestSuite {
		if i > 0 {
			testRunnerUnitTestSuiteArgs += ", "
		}
		exprCode := GenerateExpression(ir.InvocationOverTopLevelFunction{
			Over: ir.RefToIR(v),
		})
		testRunnerUnitTestSuiteArgs += exprCode
	}
	testRunnerUnitTestArgs := ""
	for i, v := range varsImplementingUnitTest {
		if i > 0 {
			testRunnerUnitTestArgs += ", "
		}
		exprCode := GenerateExpression(ir.InvocationOverTopLevelFunction{
			Over: ir.RefToIR(v),
		})
		testRunnerUnitTestArgs += exprCode
	}
	if len(varsImplementingGoIntegrationTest) > 0 {
		panic("ja test runner can't run go integration tests")
	}
	runner := runTests
	return fmt.Sprintf(`function main() {
runTests([%s], [%s])
}

%s

main()
`, testRunnerUnitTestSuiteArgs, testRunnerUnitTestArgs, runner)
}

//go:embed embed/runTests.js
var runTests string

func GeneratePackageDeclaration(declarationName ir.Reference, declarationExpression ir.TopLevelFunction) string {
	result := "function " + declarationName.Name + "("
	result += ") {\n"

	statement := declarationExpression.Body
	result += GenerateStatement(statement) + "\n"

	result += "}"
	return result
}

func GenerateFunction(function ir.LocalFunction) string {
	args := ""
	for i, paramName := range function.ParameterNames {
		if i > 0 {
			args += ", "
		}
		args += paramName
	}

	body := ""
	for _, statement := range function.Block {
		stmtCode := GenerateStatement(statement)
		body += stmtCode + "\n"
	}

	return fmt.Sprintf("function (%s) {\n%s}", args, body)
}

func GenerateNonFunctionTypeName(nonFunctionTypeName ir.NonFunctionTypeName) string {
	middleAdd := ""
	for _, concat := range nonFunctionTypeName.CommaSpaceSeparated {
		middleAdd += " + " + GenerateExpression(concat)
	}
	suffixAdd := ""
	if nonFunctionTypeName.Suffix != "" {
		suffixAdd += fmt.Sprintf(` + "%s"`, nonFunctionTypeName.Suffix)
	}
	return fmt.Sprintf(`"%s"%s%s`, nonFunctionTypeName.Prefix, middleAdd, suffixAdd)
}

func GenerateStructFunction(structFunc *ir.StructFunction) string {
	args := structFunc.GenericsArgumentName
	for _, argName := range structFunc.OtherArgumentNames {
		args += ", " + argName
	}
	constructor := fmt.Sprintf("function (%s) {\n", args)
	constructor += "return ({\n"

	returnTypeExpression := GenerateNonFunctionTypeName(structFunc.ReturnType)

	constructor += fmt.Sprintf(`"$type": %s,`+"\n", returnTypeExpression)
	for _, argName := range structFunc.OtherArgumentNames {
		constructor += fmt.Sprintf(`"%s": %s,`+"\n", argName, argName)
	}
	constructor += "})\n"
	constructor += "}"

	return constructor
}

func GenerateStatement(statement ir.Statement) string {
	asExpression, ok := statement.(ir.Expression)
	if ok {
		return GenerateExpression(asExpression)
	}
	switch s := statement.(type) {
	case ir.Return:
		return fmt.Sprintf("return %s", GenerateExpression(s.ReturnExpression))
	case ir.VariableDeclaration:
		return fmt.Sprintf("let %s = %s", s.Name, GenerateExpression(s.Expression))
	case ir.InvocationOverTopLevelFunction:
		return GenerateExpression(s)
	case ir.If:
		condCode := GenerateExpression(s.Condition)
		thenBlock := ""
		for _, stmt := range s.ThenBlock {
			thenBlock += GenerateStatement(stmt) + "\n"
		}
		elseBlock := ""
		for _, stmt := range s.ElseBlock {
			elseBlock += GenerateStatement(stmt) + "\n"
		}
		return fmt.Sprintf("if (%s) {\n%s} else {\n%s}", condCode, thenBlock, elseBlock)
	case ir.IfForTypeMatching:
		condCode := func() string {
			expr := s.Condition
			//TODO FIXME the expressions here are evaluated multiple times. They might not be pure.
			leftCode := GenerateExpression(expr.Left)
			rightCode := GenerateExpression(expr.Right)
			return fmt.Sprintf("Array.isArray(%s) ? %s.includes(%s) : %s == %s", leftCode, leftCode, rightCode, leftCode, rightCode)
		}()
		thenBlock := ""
		for _, stmt := range s.ThenBlock {
			thenBlock += GenerateStatement(stmt) + "\n"
		}
		elseBlock := ""
		for _, stmt := range s.ElseBlock {
			elseBlock += GenerateStatement(stmt) + "\n"
		}
		return fmt.Sprintf("if (%s) {\n%s} else {\n%s}", condCode, thenBlock, elseBlock)
	case ir.Invocation:
		return GenerateExpression(s)
	default:
		panic(fmt.Sprintf("unsupported statement type: %T", statement))
	}
}

func GenerateExpression(expression ir.Expression) string {
	switch expr := expression.(type) {
	case ir.Literal:
		return GenerateLiteral(expr.Value)
	case ir.Reference:
		return expr.Name
	case ir.FieldAccess:
		overCode := GenerateExpression(expr.Over)
		return fmt.Sprintf(`%s["%s"]`, overCode, expr.FieldName)
	case ir.InvocationOverTopLevelFunction:
		overCode := GenerateExpression(expr.Over)
		return fmt.Sprintf("%s()", overCode)
	case ir.Invocation:
		overCode := GenerateExpression(expr.Over)
		args := ""
		if expr.CodePointAsFirstArgument != nil {
			args += fmt.Sprintf(`"%s:%d"`, expr.CodePointAsFirstArgument.FileName, expr.CodePointAsFirstArgument.Line)
		}
		for i, arg := range expr.Arguments {
			if i > 0 || expr.CodePointAsFirstArgument != nil {
				args += ", "
			}

			argCode := GenerateExpression(arg)
			args += argCode
		}

		return fmt.Sprintf("%s(%s)", overCode, args)
	case ir.LocalFunction:
		return GenerateFunction(expr)
	case ir.ObjectInstantiation:
		fields := ""
		sortedFieldNames := maps.Keys(expr.Fields)
		sort.Strings(sortedFieldNames)
		for _, fieldName := range sortedFieldNames {
			fieldExpr := expr.Fields[fieldName]
			if fields != "" {
				fields += ", "
			}
			fieldCode := GenerateExpression(fieldExpr)
			fields += fmt.Sprintf(`"%s": %s`, fieldName, fieldCode)
		}
		return fmt.Sprintf("{%s}", fields)
	case ir.ArrayInstantiation:
		code := "["
		for i, arg := range expr.Arguments {
			if i > 0 {
				code += ", "
			}
			code += GenerateExpression(arg)
		}
		code += "]"
		return code
	case ir.NonFunctionTypeName:
		return GenerateNonFunctionTypeName(expr)
	case ir.TypeNamesOrConcat:
		resultCode := "typeNameCollapseRecursive(["
		for i, genericExpression := range expr.Expressions {
			if i > 0 {
				resultCode += `, `
			}
			resultCode += GenerateExpression(genericExpression)
		}
		resultCode += "])"
		return resultCode
	case ir.ArrayIndexAccess:
		overCode := GenerateExpression(expr.Over)
		return fmt.Sprintf("%s[%d]", overCode, expr.Index)
	default:
		panic(fmt.Sprintf("unsupported expression type: %T", expression))
	}
}

func GenerateLiteral(literal parser.Literal) string {
	switch l := literal.(type) {
	case parser.LiteralString:
		return l.Value
	case parser.LiteralInt:
		value := l.Value
		if l.Negative {
			value = 0 - value
		}
		return strconv.Itoa(value)
	case parser.LiteralFloat:
		return strconv.FormatFloat(l.Value, 'f', -1, 64)
	case parser.LiteralBool:
		return strconv.FormatBool(l.Value)
	case parser.LiteralNull:
		return "null"
	default:
		panic(fmt.Sprintf("unsupported literal type: %T", literal))
	}
}

func SortReferences(refs []ir.Reference) {
	sort.Slice(refs, func(i, j int) bool {
		return refs[i].Name < refs[j].Name
	})
}

func SortNativeFunctionRefs(refs []ir.NativeFunctionRef) {
	sort.Slice(refs, func(i, j int) bool {
		if refs[i].Package != refs[j].Package {
			return refs[i].Package < refs[j].Package
		}
		return refs[i].Name < refs[j].Name
	})
}
