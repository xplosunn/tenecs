package codegen_golang

import (
	"fmt"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_golang/standard_library"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"sort"
	"strconv"
	"strings"
)

type Import string

func GenerateProgramNonRunnable(program *ast.Program) string {
	return generate(false, program, nil, nil)
}

func GenerateProgramMain(program *ast.Program, targetMain string) string {
	return generate(false, program, &targetMain, nil)
}

func GenerateProgramTest(program *ast.Program, foundTests codegen.FoundTests) string {
	return generate(true, program, nil, &foundTests)
}

func generate(testMode bool, program *ast.Program, targetMain *string, foundTests *codegen.FoundTests) string {
	programDeclarationNames := []string{}
	for _, declaration := range program.Declarations {
		programDeclarationNames = append(programDeclarationNames, declaration.Name)
	}
	sort.Strings(programDeclarationNames)

	decs := ""
	allImports := []Import{}
	for _, declarationName := range programDeclarationNames {
		for _, declaration := range program.Declarations {
			if declaration.Name != declarationName {
				continue
			}
			imports, dec := GenerateDeclaration(&program.Package, declaration, true)
			decs += dec + "\n"
			allImports = append(allImports, imports...)
		}
	}

	structNames := maps.Keys(program.StructFunctions)
	sort.Strings(structNames)
	for _, structFuncName := range structNames {
		structFunc := program.StructFunctions[structFuncName]
		code := GenerateStructFunction(structFunc)
		structName := strings.ReplaceAll(program.Package, ".", "_") + "_" + structFuncName
		decs += GenerateStructDefinition(structName, structFunc) + "\n"
		decs += fmt.Sprintf("var %s any = %s\n", VariableName(&program.Package, structFuncName), code)
	}

	nativeFuncNames := maps.Keys(program.NativeFunctionPackages)
	sort.Strings(nativeFuncNames)
	for _, nativeFuncName := range nativeFuncNames {
		nativeFuncPkg, ok := program.NativeFunctionPackages[nativeFuncName]
		if !ok {
			panic(fmt.Sprintf("native function pkg for %s not found", nativeFuncName))
		}
		f := standard_library.Functions[nativeFuncPkg+"_"+nativeFuncName]
		caseNativeFunction, caseStructFunction := f.FunctionCases()
		if caseNativeFunction != nil {
			f := caseNativeFunction
			for _, impt := range f.Imports {
				allImports = append(allImports, Import(impt))
			}
			decs += fmt.Sprintf("var %s any = %s\n", VariableName(&nativeFuncPkg, nativeFuncName), f.Code)
		} else if caseStructFunction != nil {
			constructorArguments := []types.FunctionArgument{}
			for _, field := range caseStructFunction.FieldNamesSorted {
				constructorArguments = append(constructorArguments, types.FunctionArgument{
					Name:         field,
					VariableType: caseStructFunction.Fields[field],
				})
			}
			constructor := &types.Function{
				Generics:   caseStructFunction.Struct.DeclaredGenerics,
				Arguments:  constructorArguments,
				ReturnType: caseStructFunction.Struct,
			}
			code := GenerateStructFunction(constructor)
			decs += fmt.Sprintf("var %s any = %s\n", VariableName(&nativeFuncPkg, caseStructFunction.Struct.Name), code)
		} else {
			panic("failed to find function")
		}
	}

	main := ""

	if !testMode {
		if targetMain != nil {
			imports, mainCode := GenerateMain(&program.Package, *targetMain)
			main = mainCode
			allImports = append(allImports, imports...)
		}
	} else {
		imports, mainCode := GenerateUnitTestRunnerMain(&program.Package, foundTests.UnitTestSuites, foundTests.UnitTests)
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

	stdLibStructs := GenerateStdLibStructs()

	result := "package main\n\n" + imports + "\n" + decs + "\n" + stdLibStructs + "\n" + main

	return result
}

func GenerateStdLibStructs() string {
	stdLibStructNames := maps.Keys(standard_library.Functions)
	slices.Sort(stdLibStructNames)
	stdLibStructs := ""
	for _, name := range stdLibStructNames {
		function := standard_library.Functions[name]
		_, caseStructFunction := function.FunctionCases()
		if caseStructFunction == nil {
			continue
		}
		constructorArguments := []types.FunctionArgument{}
		for _, field := range caseStructFunction.FieldNamesSorted {
			constructorArguments = append(constructorArguments, types.FunctionArgument{
				Name:         field,
				VariableType: caseStructFunction.Fields[field],
			})
		}
		constructor := &types.Function{
			Generics:   caseStructFunction.Struct.DeclaredGenerics,
			Arguments:  constructorArguments,
			ReturnType: caseStructFunction.Struct,
		}
		stdLibStructs += GenerateStructDefinition(name, constructor) + "\n"
	}
	return stdLibStructs
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

func GenerateStructDefinition(structName string, structFunc *types.Function) string {
	structDefinition := "type " + structName + " struct {\n"
	for _, arg := range structFunc.Arguments {
		structDefinition += fmt.Sprintf("%s any\n", arg.Name)
	}
	structDefinition += "}"
	return structDefinition
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
	constructor += fmt.Sprintf("return %s {\n", generateTypeName(structFunc.ReturnType))
	for _, arg := range structFunc.Arguments {
		constructor += fmt.Sprintf("%s,\n", arg.Name)
	}
	constructor += "}\n"
	constructor += "}"

	return constructor
}

func GenerateUnitTestRunnerMain(pkgName *string, varsImplementingUnitTestSuite []string, varsImplementingUnitTest []string) ([]Import, string) {
	testRunnerTestSuiteArgs := ""
	for i, v := range varsImplementingUnitTestSuite {
		if i > 0 {
			testRunnerTestSuiteArgs += ", "
		}
		testRunnerTestSuiteArgs += VariableName(pkgName, v)
	}
	testRunnerTestArgs := ""
	for i, v := range varsImplementingUnitTest {
		if i > 0 {
			testRunnerTestArgs += ", "
		}
		testRunnerTestArgs += VariableName(pkgName, v)
	}
	imports, runner := GenerateTestRunner()
	return imports, fmt.Sprintf(`func main() {
runUnitTests([]any{%s}, []any{%s})
}

%s
`, testRunnerTestSuiteArgs, testRunnerTestArgs, runner)

}

func GenerateMain(pkgName *string, varToInvoke string) ([]Import, string) {
	imports, runtime := GenerateRuntime()
	return imports, fmt.Sprintf(`func main() {
r := runtime()
%s.(tenecs_go_Main).main.(func(any)any)(r)
}

func runtime() tenecs_go_Runtime{
return %s
}
`, VariableName(pkgName, varToInvoke), runtime)
}

func VariableName(pkgName *string, name string) string {
	prefix := "_"
	if pkgName != nil {
		prefix = strings.ReplaceAll(*pkgName, ".", "_") + "__"
	}
	return prefix + name
}

func GenerateDeclaration(pkgName *string, declaration *ast.Declaration, topLevel bool) ([]Import, string) {
	imports, exp := GenerateExpression(declaration.Expression)
	varName := VariableName(pkgName, declaration.Name)
	result := fmt.Sprintf(`var %s any
var _ = func() any {
%s = %s
return nil
}()
`, varName, varName, exp)
	if !topLevel {
		result += "_ = " + varName + "\n"
	}
	return imports, result
}

func GenerateExpression(expression ast.Expression) ([]Import, string) {
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
		return []Import{}, GenerateLiteral(*caseLiteral)
	} else if caseReference != nil {
		imports, result := GenerateReference(*caseReference)
		return imports, result
	} else if caseAccess != nil {
		imports, result := GenerateAccess(*caseAccess)
		return imports, result
	} else if caseInvocation != nil {
		imports, result := GenerateInvocation(*caseInvocation)
		return imports, result
	} else if caseFunction != nil {
		imports, result := GenerateFunction(*caseFunction)
		return imports, result
	} else if caseDeclaration != nil {
		imports, result := GenerateDeclaration(nil, caseDeclaration, false)
		return imports, result
	} else if caseIf != nil {
		imports, result := GenerateIf(*caseIf)
		return imports, result
	} else if caseList != nil {
		imports, result := GenerateList(*caseList)
		return imports, result
	} else if caseWhen != nil {
		imports, result := GenerateWhen(*caseWhen)
		return imports, result
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func GenerateList(list ast.List) ([]Import, string) {
	allImports := []Import{}
	result := "[]any{\n"
	for _, argument := range list.Arguments {
		imports, arg := GenerateExpression(argument)
		allImports = append(allImports, imports...)
		result += arg + ",\n"
	}
	result += "}"
	return allImports, result
}

func GenerateWhen(when ast.When) ([]Import, string) {
	allImports := []Import{}

	result := "func() any {\n"

	imports, over := GenerateExpression(when.Over)
	allImports = append(allImports, imports...)
	result += "var over any = " + over + "\n"

	type WhenCase struct {
		name    *string
		varType types.VariableType
		block   []ast.Expression
	}

	for _, whenCase := range when.Cases {
		variableType := whenCase.VariableType
		block := whenCase.Block

		result += fmt.Sprintf("if %s {", whenClause(variableType, false))
		if whenCase.Name != nil {
			result += fmt.Sprintf("%s := over\n", VariableName(nil, *whenCase.Name))
		}
		for i, expression := range block {
			imports, exp := GenerateExpression(expression)
			if i == len(block)-1 {
				result += "return "
			}
			result += exp + "\n"
			allImports = append(allImports, imports...)
		}
		result += "}\n"
	}
	if when.OtherCase != nil {
		if when.OtherCaseName != nil {
			result += VariableName(nil, *when.OtherCaseName) + " := over\n"
		}
		for i, expression := range when.OtherCase {
			imports, exp := GenerateExpression(expression)
			if i == len(when.OtherCase)-1 {
				result += "return "
			}
			result += exp + "\n"
			allImports = append(allImports, imports...)
		}
	} else {
		result += "return nil\n"
	}
	result += "}()"

	return allImports, result
}

func whenClause(variableType types.VariableType, nested bool) string {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO GenerateWhen caseTypeArgument")
	} else if caseList != nil {
		return whenListIfClause(caseList, nested)
	} else if caseKnownType != nil {
		if caseKnownType.Package == "" {
			return whenKnownTypeIfClause(caseKnownType, nested)
		} else {
			if nested {
				return fmt.Sprintf(`func() bool {
_, okObj := over.(%s)
return okObj
}()`, generateTypeName(caseKnownType))
			} else {
				return fmt.Sprintf("_, okObj := over.(%s); okObj", generateTypeName(caseKnownType))
			}

		}
	} else if caseFunction != nil {
		panic("TODO GenerateWhen caseFunction")
	} else if caseOr != nil {
		result := ""
		for i, elem := range caseOr.Elements {
			result += whenClause(elem, true)
			if i < len(caseOr.Elements)-1 {
				result += " || "
			}
		}
		return result
	} else {
		panic(fmt.Errorf("cases on %v", variableType))
	}
}

func whenListIfClause(caseList *types.List, nested bool) string {
	nestedClause := ""
	if ofKnownType, ok := caseList.Generic.(*types.KnownType); ok {
		nestedClause = whenKnownTypeIfClause(ofKnownType, true)
	} else if ofList, ok := caseList.Generic.(*types.List); ok {
		nestedClause = whenListIfClause(ofList, true)
	} else {
		panic("TODO GenerateWhen List")
	}

	return fmt.Sprintf(`func() bool {
arr, ok := over.([]any)
if !ok {
	return false
}
if len(arr) == 0 {
	return true
}
for _, over := range arr {
	ok := %s
	if !ok {
		return false
	}
}
return true
}()`, nestedClause)
}

func whenKnownTypeIfClause(caseKnownType *types.KnownType, nested bool) string {
	if caseKnownType.Name == "Void" {
		if nested {
			return `func() bool {
return over == nil
}()`
		} else {
			return "over == nil"
		}
	} else {
		if !nested {
			return fmt.Sprintf(`_, ok := over.(%s); ok`, generateTypeName(caseKnownType))
		} else {
			return fmt.Sprintf(`func() bool {
_, ok := over.(%s)
return ok
}()`, generateTypeName(caseKnownType))
		}
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
		func() { result = "nil" },
	)
	return result
}

func GenerateReference(reference ast.Reference) ([]Import, string) {
	allImports := []Import{}
	result := VariableName(reference.PackageName, reference.Name)

	return allImports, result
}

func GenerateAccess(access ast.Access) ([]Import, string) {
	allImports := []Import{}
	result := ""

	imports, over := GenerateExpression(access.Over)
	allImports = append(allImports, imports...)
	typeName := generateTypeName(ast.VariableTypeOfExpression(access.Over))
	result += fmt.Sprintf("%s.(%s).%s", over, typeName, access.Access)

	return allImports, result
}

func generateTypeName(varType types.VariableType) string {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		return "any"
	} else if caseList != nil {
		panic("TODO generateTypeName caseList")
	} else if caseKnownType != nil {
		if caseKnownType.Package == "" {
			if caseKnownType.Name == "String" {
				return "string"
			} else if caseKnownType.Name == "Int" {
				return "int"
			} else if caseKnownType.Name == "Boolean" {
				return "bool"
			} else {
				panic("TODO generateTypeName caseBasicType " + caseKnownType.Name)
			}
		} else {
			return strings.ReplaceAll(caseKnownType.Package, ".", "_") + "_" + caseKnownType.Name
		}
	} else if caseFunction != nil {
		funcArgList := ""
		for i, _ := range caseFunction.Arguments {
			if i > 0 {
				funcArgList += ","
			}
			funcArgList += "any"
		}
		return fmt.Sprintf(`func(%s)any`, funcArgList)
	} else if caseOr != nil {
		return "any"
	} else {
		panic("cases on variableType")
	}
}

func GenerateInvocation(invocation ast.Invocation) ([]Import, string) {
	allImports := []Import{}

	imports, over := GenerateExpression(invocation.Over)
	allImports = append(allImports, imports...)

	funcArgList := ""
	argsCode := ""
	for i, argument := range invocation.Arguments {
		if i > 0 {
			funcArgList += ","
			argsCode += ", "
		}
		funcArgList += "any"

		imports, arg := GenerateExpression(argument)
		allImports = append(allImports, imports...)
		argsCode += arg
	}

	result := fmt.Sprintf(`%s.(func(%s)any)(%s)`, over, funcArgList, argsCode)

	return allImports, result
}

func GenerateIf(caseIf ast.If) ([]Import, string) {
	allImports := []Import{}

	result := "func() any {\n"

	imports, conditionCode := GenerateExpression(caseIf.Condition)
	allImports = append(allImports, imports...)
	result += "if func() any { return " + conditionCode + " }().(bool) {\n"

	for i, expression := range caseIf.ThenBlock {
		imports, exp := GenerateExpression(expression)
		if i == len(caseIf.ThenBlock)-1 {
			result += "return "
		}
		result += exp + "\n"
		allImports = append(allImports, imports...)
	}

	if len(caseIf.ElseBlock) == 0 {
		result += "}\n"
		result += "return nil\n"
	} else {
		result += "} else {\n"
		for i, expression := range caseIf.ElseBlock {
			imports, exp := GenerateExpression(expression)
			if i == len(caseIf.ElseBlock)-1 {
				result += "return "
			}
			result += exp + "\n"
			allImports = append(allImports, imports...)
		}
		result += "}\n"
	}

	result += "}()"

	return allImports, result
}

func GenerateFunction(function ast.Function) ([]Import, string) {
	allImports := []Import{}
	args := ""
	for i, argument := range function.VariableType.Arguments {
		args += VariableName(nil, argument.Name) + " any"
		if i < len(function.VariableType.Arguments)-1 {
			args += ", "
		}
	}
	result := fmt.Sprintf("func (%s) any {\n", args)

	for i, expression := range function.Block {
		if i == len(function.Block)-1 {
			imports, exp := generateLastExpressionOfBlock(expression)
			result += exp
			allImports = append(allImports, imports...)
		} else {
			imports, exp := GenerateExpression(expression)
			result += exp + "\n"
			allImports = append(allImports, imports...)
		}
	}

	result += "}"
	return allImports, result
}

func generateLastExpressionOfBlock(expression ast.Expression) ([]Import, string) {
	imports, exp := GenerateExpression(expression)
	expLiteral, _, _, _, _, _, _, _, _ := expression.ExpressionCases()
	isVoid := types.VariableTypeEq(ast.VariableTypeOfExpression(expression), types.Void())
	result := ""
	if isVoid && expLiteral != nil {
		result += "return nil\n"
	} else {
		if !isVoid {
			result += "return "
		}
		result += exp + "\n"
		if isVoid {
			result += "return nil\n"
		}
	}
	return imports, result
}
