package codegen_golang

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/ir"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"golang.org/x/exp/maps"
)

type GeneratedProgram struct {
	PackageCode         string
	ImportsCode         string
	UserspaceCode       string
	StandardLibraryCode string
	HelpersCode         string
	MainCode            string
}

func (program GeneratedProgram) String() string {
	return program.PackageCode + "\n\n" +
		program.ImportsCode + "\n\n" +
		program.UserspaceCode + "\n\n" +
		program.StandardLibraryCode + "\n\n" +
		program.HelpersCode + "\n\n" +
		program.MainCode
}

type Import string

func GenerateProgramNonRunnable(program *ir.Program) GeneratedProgram {
	return generate(false, program, nil, nil)
}

func GenerateProgramMain(program *ir.Program, targetMain ir.Reference) GeneratedProgram {
	return generate(false, program, &targetMain, nil)
}

func GenerateProgramTest(program *ir.Program, foundTests codegen.FoundTests) GeneratedProgram {
	return generate(true, program, nil, &foundTests)
}

//go:embed embed/helpers.go
var helpers string

func generate(testMode bool, program *ir.Program, targetMain *ir.Reference, foundTests *codegen.FoundTests) GeneratedProgram {
	programDeclarationNames := []ir.Reference{}
	for declarationName, _ := range program.Declarations {
		programDeclarationNames = append(programDeclarationNames, declarationName)
	}
	SortReferences(programDeclarationNames)

	userspaceDecs := ""
	allImports := []Import{}
	for _, declarationName := range programDeclarationNames {
		for decName, decExp := range program.Declarations {
			if decName != declarationName {
				continue
			}
			imports, dec := GeneratePackageDeclaration(decName, decExp)
			userspaceDecs += dec + "\n"
			allImports = append(allImports, imports...)
		}
	}

	structNames := maps.Keys(program.StructFunctions)
	SortReferences(structNames)
	for _, structFuncName := range structNames {
		structFunc := program.StructFunctions[structFuncName]
		imports, code := GenerateStructFunction(structFunc)
		allImports = append(allImports, imports...)
		userspaceDecs += fmt.Sprintf("var %s any = %s\n", structFuncName.Name, code)
	}

	standardLibraryCode := ""

	nativeFuncNames := append([]ir.NativeFunctionRef{}, program.NativeFunctions...)
	SortNativeFunctionRefs(nativeFuncNames)
	for _, nativeFuncName := range nativeFuncNames {
		imports, dec := generateNativeFunction(nativeFuncName)
		allImports = append(allImports, imports...)
		standardLibraryCode += dec
	}

	stdLibStructNames := maps.Keys(program.StdLibStructFunctions)
	SortReferences(stdLibStructNames)
	for _, structFuncName := range stdLibStructNames {
		structFunc := program.StdLibStructFunctions[structFuncName]
		imports, code := GenerateStructFunction(structFunc)
		allImports = append(allImports, imports...)
		standardLibraryCode += fmt.Sprintf(`func %s() any {
return %s
}
`, structFuncName.Name, code)
	}

	main := ""

	if !testMode {
		if targetMain != nil {
			imports, mainCode := GenerateMain(*targetMain)
			main = mainCode
			allImports = append(allImports, imports...)
		}
	} else {
		imports, mainCode := GenerateTestRunnerMain(foundTests.UnitTestSuites, foundTests.UnitTests, foundTests.GoIntegrationTests)
		main = mainCode
		allImports = append(allImports, imports...)
	}

	importsFromHelpers := []Import{"slices", "strings"}

	importStrings := []string{}
	for _, importPkg := range importsFromHelpers {
		importStrings = append(importStrings, string(importPkg))
	}
	for _, importPkg := range allImports {
		importStrings = append(importStrings, string(importPkg))
	}
	sort.Strings(importStrings)
	importStrings = removeDuplicates(importStrings)

	imports := "import (\n"
	for _, importPkg := range importStrings {
		imports += fmt.Sprintf(`	"%s"`, importPkg) + "\n"
	}
	imports += ")\n"

	return GeneratedProgram{
		PackageCode:         "package main",
		ImportsCode:         imports,
		UserspaceCode:       userspaceDecs,
		StandardLibraryCode: standardLibraryCode,
		HelpersCode:         dropBeforeEmbedStart(helpers),
		MainCode:            main,
	}
}

func dropBeforeEmbedStart(s string) string {
	splitBeforeaAndAfterStart := strings.Split(s, "// embed start")
	if len(splitBeforeaAndAfterStart) != 2 {
		panic("// embed start : " + s)
	}
	lines := strings.Split(splitBeforeaAndAfterStart[1], "\n")
	if len(lines) > 1 {
		return strings.Join(lines[1:], "\n")
	}
	panic("expected multiple lines")
}

func removeDuplicates(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GenerateNonFunctionTypeName(nonFunctionTypeName ir.NonFunctionTypeName) ([]Import, string) {
	imports := []Import{}
	middleAdd := ""
	for _, concat := range nonFunctionTypeName.CommaSpaceSeparated {
		exprImports, exprCode := GenerateExpression(concat)
		imports = append(imports, exprImports...)
		middleAdd += " + " + exprCode
	}
	suffixAdd := ""
	if nonFunctionTypeName.Suffix != "" {
		suffixAdd += fmt.Sprintf(` + "%s"`, nonFunctionTypeName.Suffix)
	}
	return imports, fmt.Sprintf(`"%s"%s%s`, nonFunctionTypeName.Prefix, middleAdd, suffixAdd)
}

func GenerateStructFunction(structFunc *ir.StructFunction) ([]Import, string) {
	args := structFunc.GenericsArgumentName + " any"
	for _, argName := range structFunc.OtherArgumentNames {
		args += ", " + argName + " any"
	}
	constructor := fmt.Sprintf("func (%s) any {\n", args)
	constructor += "return map[string]any{\n"

	imports, returnTypeExpression := GenerateNonFunctionTypeName(structFunc.ReturnType)

	constructor += fmt.Sprintf(`"$type": %s,`+"\n", returnTypeExpression)
	for _, argName := range structFunc.OtherArgumentNames {
		constructor += fmt.Sprintf(`"%s": %s,`+"\n", argName, argName)
	}
	constructor += "}\n"
	constructor += "}"

	return imports, constructor
}

func GenerateTestRunnerMain(
	varsImplementingUnitTestSuite []ast.Ref,
	varsImplementingUnitTest []ast.Ref,
	varsImplementingGoIntegrationTest []ast.Ref,
) ([]Import, string) {
	testRunnerUnitTestSuiteArgs := ""
	for i, v := range varsImplementingUnitTestSuite {
		if i > 0 {
			testRunnerUnitTestSuiteArgs += ", "
		}
		_, exprCode := GenerateExpression(ir.InvocationOverTopLevelFunction{
			Over: ir.RefToIR(v),
		})
		testRunnerUnitTestSuiteArgs += exprCode
	}
	testRunnerUnitTestArgs := ""
	for i, v := range varsImplementingUnitTest {
		if i > 0 {
			testRunnerUnitTestArgs += ", "
		}
		_, exprCode := GenerateExpression(ir.InvocationOverTopLevelFunction{
			Over: ir.RefToIR(v),
		})
		testRunnerUnitTestArgs += exprCode
	}
	testRunnerGoIntegrationTestArgs := ""
	for i, v := range varsImplementingGoIntegrationTest {
		if i > 0 {
			testRunnerGoIntegrationTestArgs += ", "
		}
		_, exprCode := GenerateExpression(ir.InvocationOverTopLevelFunction{
			Over: ir.RefToIR(v),
		})
		testRunnerGoIntegrationTestArgs += exprCode
	}
	imports, runner := generateTestRunner()
	return imports, fmt.Sprintf(`func main() {
runTests([]any{%s}, []any{%s}, []any{%s})
}

%s
`, testRunnerUnitTestSuiteArgs, testRunnerUnitTestArgs, testRunnerGoIntegrationTestArgs, runner)
}

//go:embed embed/runTests.go
var runTests string

func generateTestRunner() ([]Import, string) {
	imports := []Import{"fmt", "reflect", "encoding/json"}

	result := dropBeforeEmbedStart(runTests)

	return imports, result
}

//go:embed embed/createMainRuntime.go
var createMainRuntime string

func GenerateMain(varToInvoke ir.Reference) ([]Import, string) {
	main := fmt.Sprintf(`func main() {
runtime := createMainRuntime()
%s().(map[string]any)["_main"].(func(any, any) any)([]any{}, runtime)
}
`, varToInvoke.Name)

	// Remove the first line (package declaration) from the embedded string
	runtime := dropBeforeEmbedStart(createMainRuntime)

	return []Import{}, main + "\n" + runtime
}

func GeneratePackageDeclaration(declarationName ir.Reference, declarationExpression ir.TopLevelFunction) ([]Import, string) {
	imports := []Import{}
	result := "func " + declarationName.Name + "("
	result += ") any {\n"

	statement := declarationExpression.Body
	additionalImports, statementCode := GenerateStatement(statement)
	imports = append(imports, additionalImports...)
	result += statementCode + "\n"

	result += "}"
	return imports, result
}

func GenerateFunction(function ir.LocalFunction) ([]Import, string) {
	imports := []Import{}
	args := ""
	for i, paramName := range function.ParameterNames {
		if i > 0 {
			args += ", "
		}
		args += paramName + " any"
	}

	body := ""
	for _, statement := range function.Block {
		additionalImports, stmtCode := GenerateStatement(statement)
		imports = append(imports, additionalImports...)
		body += stmtCode + "\n"
	}

	funcCode := fmt.Sprintf("func (%s) any {\n%s}", args, body)
	return imports, funcCode
}

func GenerateStatement(statement ir.Statement) ([]Import, string) {
	asExpression, ok := statement.(ir.Expression)
	if ok {
		imports, exrCode := GenerateExpression(asExpression)
		return imports, "_ = " + exrCode
	}
	switch s := statement.(type) {
	case ir.Return:
		imports, exprCode := GenerateExpression(s.ReturnExpression)
		return imports, fmt.Sprintf("return %s", exprCode)
	case ir.VariableDeclaration:
		imports, exprCode := GenerateExpression(s.Expression)
		varName := s.Name
		return imports, fmt.Sprintf("var %s any\n%s = %s\n_ = %s", varName, varName, exprCode, varName)
	case ir.InvocationOverTopLevelFunction:
		imports, exprCode := GenerateExpression(s)
		return imports, exprCode
	case ir.If:
		imports, condCode := GenerateExpression(s.Condition)
		thenBlock := ""
		for _, stmt := range s.ThenBlock {
			stmtImports, stmtCode := GenerateStatement(stmt)
			imports = append(imports, stmtImports...)
			thenBlock += stmtCode + "\n"
		}
		elseBlock := ""
		for _, stmt := range s.ElseBlock {
			stmtImports, stmtCode := GenerateStatement(stmt)
			imports = append(imports, stmtImports...)
			elseBlock += stmtCode + "\n"
		}
		return imports, fmt.Sprintf("if any(%s).(bool) {\n%s} else {\n%s}", condCode, thenBlock, elseBlock)
	case ir.IfForTypeMatching:
		imports, condCode := func() ([]Import, string) {
			expr := s.Condition
			imports, leftCode := GenerateExpression(expr.Left)
			rightImports, rightCode := GenerateExpression(expr.Right)
			imports = append(imports, rightImports...)
			imports = append(imports, "slices")
			return imports, fmt.Sprintf("_, ok := any(%s).([]any); (ok && slices.Contains(any(%s).([]any), %s) || (!ok && (%s == %s)))", leftCode, leftCode, rightCode, leftCode, rightCode)
		}()
		thenBlock := ""
		for _, stmt := range s.ThenBlock {
			stmtImports, stmtCode := GenerateStatement(stmt)
			imports = append(imports, stmtImports...)
			thenBlock += stmtCode + "\n"
		}
		elseBlock := ""
		for _, stmt := range s.ElseBlock {
			stmtImports, stmtCode := GenerateStatement(stmt)
			imports = append(imports, stmtImports...)
			elseBlock += stmtCode + "\n"
		}
		return imports, fmt.Sprintf("if %s {\n%s} else {\n%s}", condCode, thenBlock, elseBlock)
	case ir.Invocation:
		imports, exprCode := GenerateExpression(s)
		return imports, exprCode
	default:
		panic(fmt.Sprintf("unsupported statement type: %T", statement))
	}
}

func GenerateExpression(expression ir.Expression) ([]Import, string) {
	switch expr := expression.(type) {
	case ir.Literal:
		return []Import{}, GenerateLiteral(expr.Value)
	case ir.Reference:
		return []Import{}, expr.Name
	case ir.FieldAccess:
		imports, overCode := GenerateExpression(expr.Over)
		return imports, fmt.Sprintf(`%s.(map[string]any)["%s"]`, overCode, expr.FieldName)
	case ir.InvocationOverTopLevelFunction:
		imports, overCode := GenerateExpression(expr.Over)
		return imports, fmt.Sprintf("%s()", overCode)
	case ir.Invocation:
		imports, overCode := GenerateExpression(expr.Over)
		args := ""
		if expr.CodePointAsFirstArgument != nil {
			args += fmt.Sprintf(`"%s:%d"`, expr.CodePointAsFirstArgument.FileName, expr.CodePointAsFirstArgument.Line)
		}
		castTarget := ".(func("
		if expr.CodePointAsFirstArgument != nil {
			castTarget += "string"
		}
		for i, arg := range expr.Arguments {
			if i > 0 || expr.CodePointAsFirstArgument != nil {
				args += ", "
				castTarget += ", "
			}
			castTarget += "any"

			argImports, argCode := GenerateExpression(arg)
			imports = append(imports, argImports...)
			args += argCode
		}
		castTarget += ") any)"

		switch expr.Over.(type) {
		case ir.LocalFunction:
			castTarget = ""
		}

		return imports, fmt.Sprintf("%s%s(%s)", overCode, castTarget, args)
	case ir.LocalFunction:
		return GenerateFunction(expr)
	case ir.ObjectInstantiation:
		imports := []Import{}
		fields := ""
		sortedFieldNames := maps.Keys(expr.Fields)
		sort.Strings(sortedFieldNames)
		for _, fieldName := range sortedFieldNames {
			fieldExpr := expr.Fields[fieldName]
			if fields != "" {
				fields += ", "
			}
			additionalImports, fieldCode := GenerateExpression(fieldExpr)
			imports = append(imports, additionalImports...)
			fields += fmt.Sprintf(`"%s": %s`, fieldName, fieldCode)
		}
		return imports, fmt.Sprintf("map[string]any{%s}", fields)
	case ir.ArrayInstantiation:
		imports := []Import{}
		code := "[]any{"
		for i, arg := range expr.Arguments {
			if i > 0 {
				code += ", "
			}
			argImports, argCode := GenerateExpression(arg)
			imports = append(imports, argImports...)
			code += argCode
		}
		code += "}"
		return imports, code
	case ir.NonFunctionTypeName:
		return GenerateNonFunctionTypeName(expr)
	case ir.TypeNamesOrConcat:
		resultImports := []Import{}
		resultCode := "typeNameCollapseRecursive([]any{"
		for i, genericExpression := range expr.Expressions {
			if i > 0 {
				resultCode += `, `
			}
			exprImports, exprCode := GenerateExpression(genericExpression)
			resultImports = append(resultImports, exprImports...)
			resultCode += exprCode
		}
		resultCode += "})"
		return resultImports, resultCode
	case ir.ArrayIndexAccess:
		imports, overCode := GenerateExpression(expr.Over)
		return imports, fmt.Sprintf("%s.([]any)[%d]", overCode, expr.Index)
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
		return "nil"
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
