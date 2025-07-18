package codegen_golang

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/xplosunn/tenecs/ir"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/exp/maps"
)

type GeneratedProgram struct {
	PackageCode         string
	ImportsCode         string
	UserspaceCode       string
	StandardLibraryCode string
	MainCode            string
}

func (program GeneratedProgram) String() string {
	return program.PackageCode + "\n\n" +
		program.ImportsCode + "\n\n" +
		program.UserspaceCode + "\n\n" +
		program.StandardLibraryCode + "\n\n" +
		program.MainCode
}

type Import string

func GenerateProgramNonRunnable(program *ir.Program) GeneratedProgram {
	return generate(false, program, nil, nil)
}

func GenerateProgramMain(program *ir.Program, targetMain ir.Reference) GeneratedProgram {
	return generate(false, program, &targetMain, nil)
}

func GenerateProgramTest(program *ir.Program, foundTests FoundTests) GeneratedProgram {
	return generate(true, program, nil, &foundTests)
}

func generate(testMode bool, program *ir.Program, targetMain *ir.Reference, foundTests *FoundTests) GeneratedProgram {
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
		code := GenerateStructFunction(structFunc)
		userspaceDecs += fmt.Sprintf("var %s any = %s\n", structFuncName.Name, code)
	}

	standardLibraryCode := ""

	nativeFuncNames := maps.Keys(program.NativeFunctions)
	SortNativeFunctionRefs(nativeFuncNames)
	for _, nativeFuncName := range nativeFuncNames {
		imports, dec := generateNativeFunction(nativeFuncName)
		allImports = append(allImports, imports...)
		standardLibraryCode += dec
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

	importStrings := []string{}
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
		MainCode:            main,
	}
}

func generateNativeFunction(nativeFunctionRef ir.NativeFunctionRef) ([]Import, string) {
	if nativeFunctionRef == (ir.NativeFunctionRef{
		Package: "tenecs_go",
		Name:    "Main",
	}) {
		return []Import{}, `
func tenecs_go__Main() any {
log := func(generics []string, msg any) any {
println(msg.(map[string]any)["value"].(string))
return nil
}
console := map[string]any{
"_log": log,
}
refCreator := map[string]any{
"_new": func(generics []string, value any) any {
var ref any = value
return map[string]any{
"_get": func(generics []string) any {
return ref
},
"_set": func(generics []string, value any) any {
ref = value
return nil
},
"_modify": func(generics []string, f any) any {
ref = f.(func(any) any)(ref)
return nil
},
}
},
}
runtime := map[string]any{
"_console": console,
"_ref": refCreator,
}
return func(generics []string, run any) any {
return run.(func([]string, any) any)(generics, runtime)
}
}
`
	}
	if nativeFunctionRef == (ir.NativeFunctionRef{
		Package: "tenecs_go",
		Name:    "Runtime",
	}) {
		return []Import{}, ``
	}
	panic("TODO generateNativeFunction " + nativeFunctionRef.Package + " " + nativeFunctionRef.Name)
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

func GenerateStructFunction(structFunc *types.Function) string {
	args := ""
	for i, arg := range structFunc.Arguments {
		if i > 0 {
			args += ", "
		}
		args += arg.Name + " any"
	}
	constructor := fmt.Sprintf("func (%s) any {\n", args)
	constructor += "return map[string]any{\n"
	for _, arg := range structFunc.Arguments {
		constructor += fmt.Sprintf(`"%s": %s,`+"\n", arg.Name, arg.Name)
	}
	constructor += "}\n"
	constructor += "}"

	return constructor
}

func GenerateTestRunnerMain(
	varsImplementingUnitTestSuite []ir.Reference,
	varsImplementingUnitTest []ir.Reference,
	varsImplementingGoIntegrationTest []ir.Reference,
) ([]Import, string) {
	panic("TODO GenerateTestRunnerMain")
}

func GenerateMain(varToInvoke ir.Reference) ([]Import, string) {
	main := fmt.Sprintf(`func main() {
%s()
}`, varToInvoke.Name)

	return []Import{}, main
}

func GeneratePackageDeclaration(declarationName ir.Reference, declarationExpression ir.TopLevelFunction) ([]Import, string) {
	imports := []Import{}
	result := "func " + declarationName.Name + "("
	for i, parameterName := range declarationExpression.ParameterNames {
		if i > 0 {
			result += ", "
		}
		result += parameterName + " any"
	}
	result += ") any {\n"
	for _, statement := range declarationExpression.Body {
		additionalImports, statementCode := GenerateStatement(statement)
		imports = append(imports, additionalImports...)
		result += statementCode + "\n"
	}
	result += "}"
	return imports, result
}

func GenerateFunction(function ir.TopLevelFunction) ([]Import, string) {
	imports := []Import{}
	args := "generics []string"
	for _, paramName := range function.ParameterNames {
		args += ", "
		args += paramName + " any"
	}

	body := ""
	for _, statement := range function.Body {
		additionalImports, stmtCode := GenerateStatement(statement)
		imports = append(imports, additionalImports...)
		body += stmtCode + "\n"
	}

	funcCode := fmt.Sprintf("func (%s) any {\n%s}", args, body)
	return []Import{}, funcCode
}

func GenerateStatement(statement ir.Statement) ([]Import, string) {
	switch s := statement.(type) {
	case ir.Return:
		imports, exprCode := GenerateExpression(s.ReturnExpression)
		return imports, fmt.Sprintf("return %s", exprCode)
	case ir.VariableDeclaration:
		imports, exprCode := GenerateExpression(s.Expression)
		varName := s.Name
		return imports, fmt.Sprintf("%s := %s\n_ = %s", varName, exprCode, varName)
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
		args := "[]string{"
		for _, generic := range expr.GenericsPassed {
			args += fmt.Sprintf(`"%s",`, generic)
		}
		args += "}"
		castTarget := ".(func([]string"
		for _, arg := range expr.Arguments {
			args += ", "
			castTarget += ",any"

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
		return GenerateFunction(ir.TopLevelFunction{
			ParameterNames: expr.ParameterNames,
			Body:           expr.Block,
		})
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
	case ir.EqualityComparison:
		imports, leftCode := GenerateExpression(expr.Left)
		rightImports, rightCode := GenerateExpression(expr.Right)
		imports = append(imports, rightImports...)
		return imports, fmt.Sprintf("%s == %s", leftCode, rightCode)
	case ir.If:
		imports, condCode := GenerateExpression(expr.Condition)
		thenBlock := ""
		for _, stmt := range expr.ThenBlock {
			stmtImports, stmtCode := GenerateStatement(stmt)
			imports = append(imports, stmtImports...)
			thenBlock += stmtCode + "\n"
		}
		elseBlock := ""
		for _, stmt := range expr.ElseBlock {
			stmtImports, stmtCode := GenerateStatement(stmt)
			imports = append(imports, stmtImports...)
			elseBlock += stmtCode + "\n"
		}
		return imports, fmt.Sprintf("if %s {\n%s} else {\n%s}", condCode, thenBlock, elseBlock)
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

type FoundTests struct {
	UnitTests          []ir.Reference
	UnitTestSuites     []ir.Reference
	GoIntegrationTests []ir.Reference
}
