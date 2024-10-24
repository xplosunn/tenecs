package codegen

import (
	"fmt"
	"github.com/xplosunn/tenecs/codegen/standard_library"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/exp/maps"
	"sort"
	"strconv"
	"strings"
)

type Import string

type TrackedDeclaration struct {
	Is        IsTrackedDeclaration
	VarName   string
	TestSuite bool
}

type IsTrackedDeclaration string

const (
	IsTrackedDeclarationNone     IsTrackedDeclaration = ""
	IsTrackedDeclarationMain     IsTrackedDeclaration = "main"
	IsTrackedDeclarationUnitTest IsTrackedDeclaration = "unit_test"
)

func GenerateProgramMain(program *ast.Program, targetMain *string) string {
	return generate(false, program, targetMain)
}

func GenerateProgramTest(program *ast.Program) string {
	return generate(true, program, nil)
}

func generate(testMode bool, program *ast.Program, targetMain *string) string {
	trackedDeclarationMains := []string{}
	trackedDeclarationUnitTestSuites := []string{}
	trackedDeclarationUnitTests := []string{}

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
			trackedDeclaration, imports, dec := GenerateDeclaration(&program.Package, declaration, true)
			decs += dec + "\n"
			allImports = append(allImports, imports...)
			if trackedDeclaration != nil {
				if trackedDeclaration.Is == IsTrackedDeclarationMain {
					trackedDeclarationMains = append(trackedDeclarationMains, trackedDeclaration.VarName)
				} else if trackedDeclaration.Is == IsTrackedDeclarationUnitTest {
					if trackedDeclaration.TestSuite {
						trackedDeclarationUnitTestSuites = append(trackedDeclarationUnitTestSuites, trackedDeclaration.VarName)
					} else {
						trackedDeclarationUnitTests = append(trackedDeclarationUnitTests, trackedDeclaration.VarName)
					}
				}
			}
		}
	}

	structNames := maps.Keys(program.StructFunctions)
	sort.Strings(structNames)
	for _, structFuncName := range structNames {
		structFunc := program.StructFunctions[structFuncName]
		code := GenerateStructFunction(structFuncName, structFunc)
		decs += fmt.Sprintf("var %s any = %s\n", VariableName(&program.Package, structFuncName), code)
		var _ = decs
	}

	nativeFuncNames := maps.Keys(program.NativeFunctionPackages)
	sort.Strings(nativeFuncNames)
	for _, nativeFuncName := range nativeFuncNames {
		nativeFuncPkg, ok := program.NativeFunctionPackages[nativeFuncName]
		if !ok {
			panic(fmt.Sprintf("native function pkg for %s not found", nativeFuncName))
		}
		f := standard_library.Functions[nativeFuncPkg+"_"+nativeFuncName]
		if len(f.Code) == 0 {
			panic("failed to find function")
		}
		for _, impt := range f.Imports {
			allImports = append(allImports, Import(impt))
		}
		decs += fmt.Sprintf("var %s any = %s\n", VariableName(&nativeFuncPkg, nativeFuncName), f.Code)
	}

	main := ""

	if !testMode {
		var mainVar *string

		if targetMain != nil {
			for _, trackedDeclaration := range trackedDeclarationMains {
				if strings.HasSuffix(trackedDeclaration, *targetMain) {
					mainVar = &trackedDeclaration
					break
				}
			}
			if mainVar == nil {
				panic("Target main not found: " + *targetMain)
			}
		} else {
			if len(trackedDeclarationMains) > 1 {
				panic("Multiple mains without a target")
			} else if len(trackedDeclarationMains) == 1 {
				mainVar = &trackedDeclarationMains[0]
			}
		}

		if mainVar != nil {
			imports, mainCode := GenerateMain(*mainVar)
			main = mainCode
			allImports = append(allImports, imports...)
		}
	} else {
		imports, mainCode := GenerateUnitTestRunnerMain(trackedDeclarationUnitTestSuites, trackedDeclarationUnitTests)
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

	result := "package main\n\n" + imports + "\n" + decs + "\n" + main

	return result
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

func GenerateStructFunction(structName string, structFunc *types.Function) string {
	args := ""
	resultMapElements := ""
	for i, arg := range structFunc.Arguments {
		if i > 0 {
			args += ", "
		}
		args += arg.Name + " any"
		resultMapElements += fmt.Sprintf(`"%s": %s,`, arg.Name, arg.Name)
	}

	return fmt.Sprintf(`func (%s) any {
return map[string]any{
"$type": "%s",
%s
}
}`, args, structName, resultMapElements)
}

func GenerateUnitTestRunnerMain(varsImplementingUnitTestSuite []string, varsImplementingUnitTest []string) ([]Import, string) {
	testRunnerTestSuiteArgs := strings.Join(varsImplementingUnitTestSuite, ", ")
	testRunnerTestArgs := strings.Join(varsImplementingUnitTest, ", ")
	imports, runner := GenerateTestRunner()
	return imports, fmt.Sprintf(`func main() {
runUnitTests([]any{%s}, []any{%s})
}

%s
`, testRunnerTestSuiteArgs, testRunnerTestArgs, runner)

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

func VariableName(pkgName *string, name string) string {
	pkgPrefix := ""
	if pkgName != nil {
		pkgPrefix = "__" + strings.ReplaceAll(*pkgName, ".", "_") + "__"
	}
	return "P" + pkgPrefix + name
}

func GenerateDeclaration(pkgName *string, declaration *ast.Declaration, topLevel bool) (*TrackedDeclaration, []Import, string) {
	_, imports, exp := GenerateExpression(declaration.Expression)
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

	var trackedDeclaration *TrackedDeclaration
	varType := ast.VariableTypeOfExpression(declaration.Expression)
	_, caseKnownType, _, _ := varType.VariableTypeCases()
	if topLevel && caseKnownType != nil {
		if caseKnownType.Name == "UnitTestSuite" && caseKnownType.Package == "tenecs.test" {
			trackedDeclaration = &TrackedDeclaration{
				Is:        IsTrackedDeclarationUnitTest,
				VarName:   varName,
				TestSuite: true,
			}
		} else if caseKnownType.Name == "UnitTest" && caseKnownType.Package == "tenecs.test" {
			trackedDeclaration = &TrackedDeclaration{
				Is:      IsTrackedDeclarationUnitTest,
				VarName: varName,
			}
		} else if caseKnownType.Name == "Main" && caseKnownType.Package == "tenecs.go" {
			trackedDeclaration = &TrackedDeclaration{
				Is:      IsTrackedDeclarationMain,
				VarName: varName,
			}
		}
	}
	return trackedDeclaration, imports, result
}

func GenerateExpression(expression ast.Expression) (IsTrackedDeclaration, []Import, string) {
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
		return IsTrackedDeclarationNone, []Import{}, GenerateLiteral(*caseLiteral)
	} else if caseReference != nil {
		imports, result := GenerateReference(*caseReference)
		return IsTrackedDeclarationNone, imports, result
	} else if caseAccess != nil {
		imports, result := GenerateAccess(*caseAccess)
		return IsTrackedDeclarationNone, imports, result
	} else if caseInvocation != nil {
		imports, result := GenerateInvocation(*caseInvocation)
		return IsTrackedDeclarationNone, imports, result
	} else if caseFunction != nil {
		imports, result := GenerateFunction(*caseFunction)
		return IsTrackedDeclarationNone, imports, result
	} else if caseDeclaration != nil {
		_, imports, result := GenerateDeclaration(nil, caseDeclaration, false)
		return IsTrackedDeclarationNone, imports, result
	} else if caseIf != nil {
		imports, result := GenerateIf(*caseIf)
		return IsTrackedDeclarationNone, imports, result
	} else if caseList != nil {
		imports, result := GenerateList(*caseList)
		return IsTrackedDeclarationNone, imports, result
	} else if caseWhen != nil {
		imports, result := GenerateWhen(*caseWhen)
		return IsTrackedDeclarationNone, imports, result
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func GenerateList(list ast.List) ([]Import, string) {
	allImports := []Import{}
	result := "[]any{\n"
	for _, argument := range list.Arguments {
		_, imports, arg := GenerateExpression(argument)
		allImports = append(allImports, imports...)
		result += arg + ",\n"
	}
	result += "}"
	return allImports, result
}

func GenerateWhen(when ast.When) ([]Import, string) {
	allImports := []Import{}

	result := "func() any {\n"

	_, imports, over := GenerateExpression(when.Over)
	allImports = append(allImports, imports...)
	result += "var over any = " + over + "\n"

	type WhenCase struct {
		name    *string
		varType types.VariableType
		block   []ast.Expression
	}
	sortedCases := []WhenCase{}
	for variableType, block := range when.Cases {
		sortedCases = append(sortedCases, WhenCase{
			name:    when.CaseNames[variableType],
			varType: variableType,
			block:   block,
		})
	}
	sort.Slice(sortedCases, func(i, j int) bool {
		return fmt.Sprintf("%+v", sortedCases[i].varType) < fmt.Sprintf("%+v", sortedCases[j].varType)
	})

	for _, whenCase := range sortedCases {
		variableType := whenCase.varType
		block := whenCase.block

		result += fmt.Sprintf("if %s {", whenClause(variableType, false))
		if whenCase.name != nil {
			result += fmt.Sprintf("%s := over\n", VariableName(nil, *whenCase.name))
		}
		for i, expression := range block {
			_, imports, exp := GenerateExpression(expression)
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
			_, imports, exp := GenerateExpression(expression)
			if i == len(when.OtherCase)-1 {
				result += "return "
			}
			result += exp + "\n"
			allImports = append(allImports, imports...)
		}
	}

	result += "return nil\n"
	result += "}()"

	return allImports, result
}

func whenClause(variableType types.VariableType, nested bool) string {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO GenerateWhen caseTypeArgument")
	} else if caseKnownType != nil {
		if caseKnownType.Package == "" {
			return whenKnownTypeIfClause(caseKnownType, nested)
		} else {
			if nested {
				return fmt.Sprintf(`func() bool {
value, okObj := over.(map[string]any)
return value["$type"] == "%s"
}()`, caseKnownType.Name)
			} else {
				return fmt.Sprintf("value, okObj := over.(map[string]any); okObj && value[\"$type\"] == \"%s\"", caseKnownType.Name)
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

func whenKnownTypeIfClause(caseKnownType *types.KnownType, nested bool) string {
	if caseKnownType.Name == "List" {
		ofKnownType, ok := caseKnownType.Generics[0].(*types.KnownType)
		if ok {
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
}()`, whenKnownTypeIfClause(ofKnownType, true))
		} else {
			panic("TODO GenerateWhen List")
		}
	} else if caseKnownType.Name == "Void" {
		if nested {
			return `func() bool {
return over == nil
}()`
		} else {
			return "over == nil"
		}
	} else {
		if !nested {
			return fmt.Sprintf(`_, ok := over.(%s); ok`, whenKnownTypeGoType(caseKnownType))
		} else {
			return fmt.Sprintf(`func() bool {
_, ok := over.(%s)
return ok
}()`, whenKnownTypeGoType(caseKnownType))
		}
	}
}

func whenKnownTypeGoType(caseKnownType *types.KnownType) string {
	if caseKnownType.Name == "String" {
		return "string"
	} else if caseKnownType.Name == "Int" {
		return "int"
	} else if caseKnownType.Name == "Boolean" {
		return "bool"
	} else {
		panic("TODO GenerateWhen caseBasicType " + caseKnownType.Name)
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

	_, imports, over := GenerateExpression(access.Over)
	allImports = append(allImports, imports...)
	result += fmt.Sprintf("%s.(map[string]any)[\"%s\"]", over, access.Access)

	return allImports, result
}

func GenerateInvocation(invocation ast.Invocation) ([]Import, string) {
	allImports := []Import{}

	_, imports, over := GenerateExpression(invocation.Over)
	allImports = append(allImports, imports...)

	funcArgList := ""
	argsCode := ""
	for i, argument := range invocation.Arguments {
		if i > 0 {
			funcArgList += ","
			argsCode += ", "
		}
		funcArgList += "any"

		_, imports, arg := GenerateExpression(argument)
		allImports = append(allImports, imports...)
		argsCode += arg
	}

	result := fmt.Sprintf(`%s.(func(%s)any)(%s)`, over, funcArgList, argsCode)

	return allImports, result
}

func GenerateIf(caseIf ast.If) ([]Import, string) {
	allImports := []Import{}

	result := "func() any {\n"

	_, imports, conditionCode := GenerateExpression(caseIf.Condition)
	allImports = append(allImports, imports...)
	result += "if func() any { return " + conditionCode + " }().(bool) {\n"

	for i, expression := range caseIf.ThenBlock {
		_, imports, exp := GenerateExpression(expression)
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
			_, imports, exp := GenerateExpression(expression)
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
			_, imports, exp := GenerateExpression(expression)
			result += exp + "\n"
			allImports = append(allImports, imports...)
		}
	}

	if len(function.Block) == 0 {
		result += "return nil\n"
	}

	result += "}"
	return allImports, result
}

func generateLastExpressionOfBlock(expression ast.Expression) ([]Import, string) {
	_, imports, exp := GenerateExpression(expression)
	expLiteral, _, _, _, _, _, _, _, _ := expression.ExpressionCases()
	isVoid := types.VariableTypeEq(ast.VariableTypeOfExpression(expression), types.Void())
	result := ""
	if isVoid && expLiteral != nil {
		result += "return " + exp + "\n"
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
