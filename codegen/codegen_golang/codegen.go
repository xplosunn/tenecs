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

func GenerateProgramMain(program *ast.Program, targetMain ast.Ref) string {
	return generate(false, program, &targetMain, nil)
}

func GenerateProgramTest(program *ast.Program, foundTests codegen.FoundTests) string {
	return generate(true, program, nil, &foundTests)
}

func generate(testMode bool, program *ast.Program, targetMain *ast.Ref, foundTests *codegen.FoundTests) string {
	programDeclarationNames := []ast.Ref{}
	for declarationName, _ := range program.Declarations {
		programDeclarationNames = append(programDeclarationNames, declarationName)
	}
	ast.SortRefs(programDeclarationNames)

	decs := ""
	allImports := []Import{}
	for _, declarationName := range programDeclarationNames {
		for decName, decExp := range program.Declarations {
			if decName != declarationName {
				continue
			}
			imports, dec := GeneratePackageDeclaration(decName.Package, decName.Name, decExp, program.StructTypeArgumentMatchFields)
			decs += dec + "\n"
			allImports = append(allImports, imports...)
		}
	}

	structNames := maps.Keys(program.StructFunctions)
	ast.SortRefs(structNames)
	for _, structFuncName := range structNames {
		structFunc := program.StructFunctions[structFuncName]
		code := GenerateStructFunction(structFunc)
		structName := strings.ReplaceAll(structFuncName.Package, ".", "_") + "_" + structFuncName.Name
		decs += GenerateStructDefinition(structName, structFunc) + "\n"
		decs += fmt.Sprintf("var %s any = %s\n", VariableName(&structFuncName.Package, structFuncName.Name), code)
	}

	nativeFuncNames := maps.Keys(program.NativeFunctions)
	ast.SortRefs(nativeFuncNames)
	for _, nativeFuncName := range nativeFuncNames {
		f := standard_library.Functions[nativeFuncName.Package+"_"+nativeFuncName.Name]
		caseNativeFunction, caseStructFunction := f.FunctionCases()
		if caseNativeFunction != nil {
			f := caseNativeFunction
			for _, impt := range f.Imports {
				allImports = append(allImports, Import(impt))
			}
			decs += fmt.Sprintf("var %s any = %s\n", VariableName(&nativeFuncName.Package, nativeFuncName.Name), f.Code)
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
			decs += fmt.Sprintf("var %s any = %s\n", VariableName(&nativeFuncName.Package, caseStructFunction.Struct.Name), code)
		} else {
			panic("failed to find function")
		}
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
		structDefinition += fmt.Sprintf("%s any\n", VariableName(nil, arg.Name))
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
		args += VariableName(nil, arg.Name) + " any"
	}
	constructor := fmt.Sprintf("func (%s) any {\n", args)
	constructor += fmt.Sprintf("return %s {\n", generateTypeName(structFunc.ReturnType))
	for _, arg := range structFunc.Arguments {
		constructor += fmt.Sprintf("%s,\n", VariableName(nil, arg.Name))
	}
	constructor += "}\n"
	constructor += "}"

	return constructor
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
		testRunnerUnitTestSuiteArgs += VariableName(&v.Package, v.Name)
	}
	testRunnerUnitTestArgs := ""
	for i, v := range varsImplementingUnitTest {
		if i > 0 {
			testRunnerUnitTestArgs += ", "
		}
		testRunnerUnitTestArgs += VariableName(&v.Package, v.Name)
	}
	testRunnerGoIntegrationTestArgs := ""
	for i, v := range varsImplementingGoIntegrationTest {
		if i > 0 {
			testRunnerGoIntegrationTestArgs += ", "
		}
		testRunnerGoIntegrationTestArgs += VariableName(&v.Package, v.Name)
	}
	imports, runner := GenerateTestRunner()
	return imports, fmt.Sprintf(`func main() {
runTests([]any{%s}, []any{%s}, []any{%s})
}

%s
`, testRunnerUnitTestSuiteArgs, testRunnerUnitTestArgs, testRunnerGoIntegrationTestArgs, runner)

}

func GenerateMain(varToInvoke ast.Ref) ([]Import, string) {
	imports, runtime := GenerateRuntime()
	return imports, fmt.Sprintf(`func main() {
r := runtime()
%s.(tenecs_go_Main)._main.(func(any)any)(r)
}

func runtime() tenecs_go_Runtime{
return %s
}
`, VariableName(&varToInvoke.Package, varToInvoke.Name), runtime)
}

func VariableName(pkgName *string, name string) string {
	prefix := "_"
	if pkgName != nil {
		prefix = strings.ReplaceAll(*pkgName, ".", "_") + "__"
	}
	return prefix + name
}

func GeneratePackageDeclaration(declarationPackage string, declarationName string, declarationExpression ast.Expression, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
	imports, exp := GenerateExpression(declarationExpression, structTypeArgumentMatchFields)
	varName := VariableName(&declarationPackage, declarationName)
	result := fmt.Sprintf(`var %s any
var _ = func() any {
%s = %s
return nil
}()
`, varName, varName, exp)
	return imports, result
}

func GenerateDeclaration(declaration *ast.Declaration, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
	imports, exp := GenerateExpression(declaration.Expression, structTypeArgumentMatchFields)
	varName := VariableName(nil, declaration.Name)
	result := fmt.Sprintf(`var %s any
var _ = func() any {
%s = %s
return nil
}()
`, varName, varName, exp)
	result += "_ = " + varName + "\n"
	return imports, result
}

func GenerateExpression(expression ast.Expression, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
		return []Import{}, GenerateLiteral(*caseLiteral)
	} else if caseReference != nil {
		imports, result := GenerateReference(*caseReference)
		return imports, result
	} else if caseAccess != nil {
		imports, result := GenerateAccess(*caseAccess, structTypeArgumentMatchFields)
		return imports, result
	} else if caseInvocation != nil {
		imports, result := GenerateInvocation(*caseInvocation, structTypeArgumentMatchFields)
		return imports, result
	} else if caseFunction != nil {
		imports, result := GenerateFunction(*caseFunction, structTypeArgumentMatchFields)
		return imports, result
	} else if caseDeclaration != nil {
		imports, result := GenerateDeclaration(caseDeclaration, structTypeArgumentMatchFields)
		return imports, result
	} else if caseIf != nil {
		imports, result := GenerateIf(*caseIf, structTypeArgumentMatchFields)
		return imports, result
	} else if caseList != nil {
		imports, result := GenerateList(*caseList, structTypeArgumentMatchFields)
		return imports, result
	} else if caseWhen != nil {
		imports, result := GenerateWhen(*caseWhen, structTypeArgumentMatchFields)
		return imports, result
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func GenerateList(list ast.List, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
	allImports := []Import{}
	result := "[]any{\n"
	for _, argument := range list.Arguments {
		imports, arg := GenerateExpression(argument, structTypeArgumentMatchFields)
		allImports = append(allImports, imports...)
		result += arg + ",\n"
	}
	result += "}"
	return allImports, result
}

func GenerateWhen(when ast.When, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
	allImports := []Import{}

	result := "func() any {\n"

	imports, over := GenerateExpression(when.Over, structTypeArgumentMatchFields)
	allImports = append(allImports, imports...)
	//TODO FIXME this "over" name might clash
	result += "var over any = " + over + "\n"

	type WhenCase struct {
		name    *string
		varType types.VariableType
		block   []ast.Expression
	}

	for _, whenCase := range when.Cases {
		variableType := whenCase.VariableType
		block := whenCase.Block

		result += fmt.Sprintf("if %s {", whenClause("over", variableType, false, structTypeArgumentMatchFields))
		if whenCase.Name != nil {
			result += fmt.Sprintf("%s := over\n", VariableName(nil, *whenCase.Name))
			result += fmt.Sprintf("_ = %s\n", VariableName(nil, *whenCase.Name))
		}
		for i, expression := range block {
			imports, exp := GenerateExpression(expression, structTypeArgumentMatchFields)
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
			imports, exp := GenerateExpression(expression, structTypeArgumentMatchFields)
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

func whenClause(varName string, variableType types.VariableType, nested bool, structTypeArgumentMatchFields map[ast.Ref][]string) string {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO GenerateWhen caseTypeArgument")
	} else if caseList != nil {
		return whenListIfClause(caseList, nested, structTypeArgumentMatchFields)
	} else if caseKnownType != nil {
		if caseKnownType.Package == "" {
			return whenKnownTypeIfClause(varName, caseKnownType, nested)
		} else {
			typeArgumentMatchFields := structTypeArgumentMatchFields[ast.Ref{
				Package: caseKnownType.Package,
				Name:    caseKnownType.Name,
			}]
			if len(typeArgumentMatchFields) != len(caseKnownType.Generics) {
				panic(fmt.Sprintf("len(typeArgumentMatchFields) != len(caseKnownType.Generics), %d != %d", len(typeArgumentMatchFields), len(caseKnownType.Generics)))
			}

			additionalClauses := ""
			for i, generic := range caseKnownType.Generics {
				matchFieldName := typeArgumentMatchFields[i]
				nestedVarName := fmt.Sprintf("%s__%d", varName, i)
				additionalClauses += fmt.Sprintf(` && (func() bool {
  %s := %s.(%s).%s
  return %s
}())`, nestedVarName, varName, generateTypeName(caseKnownType), VariableName(nil, matchFieldName), whenClause(nestedVarName, generic, true, structTypeArgumentMatchFields))
			}

			if nested {
				return fmt.Sprintf(`func() bool {
_, okObj := %s.(%s)
return okObj%s
}()`, varName, generateTypeName(caseKnownType), additionalClauses)
			} else {
				return fmt.Sprintf("_, okObj := %s.(%s); okObj%s", varName, generateTypeName(caseKnownType), additionalClauses)
			}

		}
	} else if caseFunction != nil {
		panic("TODO GenerateWhen caseFunction")
	} else if caseOr != nil {
		result := ""
		for i, elem := range caseOr.Elements {
			result += whenClause(varName, elem, true, structTypeArgumentMatchFields)
			if i < len(caseOr.Elements)-1 {
				result += " || "
			}
		}
		return result
	} else {
		panic(fmt.Errorf("cases on %v", variableType))
	}
}

func whenListIfClause(caseList *types.List, nested bool, structTypeArgumentMatchFields map[ast.Ref][]string) string {
	nestedClause := whenClause("over", caseList.Generic, true, structTypeArgumentMatchFields)

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

func whenKnownTypeIfClause(varName string, caseKnownType *types.KnownType, nested bool) string {
	if caseKnownType.Name == "Void" {
		if nested {
			return fmt.Sprintf(`func() bool {
return %s == nil
}()`, varName)
		} else {
			return fmt.Sprintf("%s == nil", varName)
		}
	} else {
		if !nested {
			return fmt.Sprintf(`_, ok := %s.(%s); ok`, varName, generateTypeName(caseKnownType))
		} else {
			return fmt.Sprintf(`func() bool {
_, ok := %s.(%s)
return ok
}()`, varName, generateTypeName(caseKnownType))
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

func GenerateAccess(access ast.Access, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
	allImports := []Import{}
	result := ""

	imports, over := GenerateExpression(access.Over, structTypeArgumentMatchFields)
	allImports = append(allImports, imports...)
	typeName := generateTypeName(ast.VariableTypeOfExpression(access.Over))
	result += fmt.Sprintf("%s.(%s).%s", over, typeName, VariableName(nil, access.Access))

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

func GenerateInvocation(invocation ast.Invocation, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
	allImports := []Import{}

	imports, over := GenerateExpression(invocation.Over, structTypeArgumentMatchFields)
	allImports = append(allImports, imports...)

	funcArgList := ""
	argsCode := ""
	for i, argument := range invocation.Arguments {
		if i > 0 {
			funcArgList += ","
			argsCode += ", "
		}
		funcArgList += "any"

		imports, arg := GenerateExpression(argument, structTypeArgumentMatchFields)
		allImports = append(allImports, imports...)
		argsCode += arg
	}

	_, _, _, overFunction, _ := ast.VariableTypeOfExpression(invocation.Over).VariableTypeCases()
	if overFunction == nil {
		panic("expected function for invocation")
	}

	if overFunction.CodePointAsFirstArgument {
		funcArgList = "any," + funcArgList
		argsCode = fmt.Sprintf("\"%s:%d\", ", invocation.CodePoint.FileName, invocation.CodePoint.Line) + argsCode
	}

	result := fmt.Sprintf(`%s.(func(%s)any)(%s)`, over, funcArgList, argsCode)

	return allImports, result
}

func GenerateIf(caseIf ast.If, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
	allImports := []Import{}

	result := "func() any {\n"

	imports, conditionCode := GenerateExpression(caseIf.Condition, structTypeArgumentMatchFields)
	allImports = append(allImports, imports...)
	result += "if func() any { return " + conditionCode + " }().(bool) {\n"

	for i, expression := range caseIf.ThenBlock {
		imports, exp := GenerateExpression(expression, structTypeArgumentMatchFields)
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
			imports, exp := GenerateExpression(expression, structTypeArgumentMatchFields)
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

func GenerateFunction(function ast.Function, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
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
			imports, exp := generateLastExpressionOfBlock(expression, structTypeArgumentMatchFields)
			result += exp
			allImports = append(allImports, imports...)
		} else {
			imports, exp := GenerateExpression(expression, structTypeArgumentMatchFields)
			result += exp + "\n"
			allImports = append(allImports, imports...)
		}
	}

	result += "}"
	return allImports, result
}

func generateLastExpressionOfBlock(expression ast.Expression, structTypeArgumentMatchFields map[ast.Ref][]string) ([]Import, string) {
	imports, exp := GenerateExpression(expression, structTypeArgumentMatchFields)
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
