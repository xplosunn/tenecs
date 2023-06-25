package codegen

import (
	"fmt"
	"github.com/xplosunn/tenecs/codegen/standard_library"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"sort"
	"strconv"
	"strings"
)

type Import string

type TrackedDeclaration struct {
	Is      IsTrackedDeclaration
	VarName string
}

type IsTrackedDeclaration string

const (
	IsTrackedDeclarationNone     IsTrackedDeclaration = ""
	IsTrackedDeclarationMain     IsTrackedDeclaration = "main"
	IsTrackedDeclarationUnitTest IsTrackedDeclaration = "unit_test"
)

func Generate(testMode bool, program *ast.Program) string {
	trackedDeclarationType := IsTrackedDeclarationMain
	if testMode {
		trackedDeclarationType = IsTrackedDeclarationUnitTest
	}

	trackedDeclarations := []string{}

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
			trackedDeclaration, imports, dec := GenerateDeclaration(declaration)
			decs += dec + "\n"
			allImports = append(allImports, imports...)
			if trackedDeclaration != nil && trackedDeclaration.Is == trackedDeclarationType {
				trackedDeclarations = append(trackedDeclarations, trackedDeclaration.VarName)
			}
		}
	}

	for structFuncName, structFunc := range program.StructFunctions {
		code := GenerateStructFunction(structFunc)
		decs += fmt.Sprintf("var %s any = %s\n", VariableName(structFuncName), code)
	}

	for nativeFuncName, nativeFuncPkg := range program.NativeFunctionPackages {
		f := standard_library.Functions[nativeFuncPkg]
		for _, impt := range f.Imports {
			allImports = append(allImports, Import(impt))
		}
		decs += fmt.Sprintf("var %s any = %s\n", VariableName(nativeFuncName), f.Code)
	}

	main := ""

	if !testMode {
		if len(trackedDeclarations) > 1 {
			panic("TODO Generate multiple mains")
		} else if len(trackedDeclarations) == 1 {
			imports, mainCode := GenerateMain(trackedDeclarations[0])
			main = mainCode
			allImports = append(allImports, imports...)
		}
	} else {
		imports, mainCode := GenerateUnitTestRunnerMain(trackedDeclarations)
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

func GenerateStructFunction(structFunc *types.Function) string {
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
%s
}
}`, args, resultMapElements)
}

func GenerateUnitTestRunnerMain(varsImplementingUnitTests []string) ([]Import, string) {
	testRunnerTestNameArgs := ""
	for i, varName := range varsImplementingUnitTests {
		if i > 0 {
			testRunnerTestNameArgs += ", "
		}
		testRunnerTestNameArgs += fmt.Sprintf(`"%s"`, strings.TrimPrefix(varName, "P"))
	}
	testRunnerTestArgs := strings.Join(varsImplementingUnitTests, ", ")
	imports, runner := GenerateTestRunner()
	return imports, fmt.Sprintf(`func main() {
runTests([]string{%s}, []any{%s})
}

%s
`, testRunnerTestNameArgs, testRunnerTestArgs, runner)

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

func GenerateDeclaration(declaration *ast.Declaration) (*TrackedDeclaration, []Import, string) {
	isTrackedDeclaration, imports, exp := GenerateExpression(&declaration.Name, declaration.Expression)
	varName := VariableName(declaration.Name)
	result := "var " + varName + " any = " + exp + "\n"

	var trackedDeclaration *TrackedDeclaration
	if isTrackedDeclaration != IsTrackedDeclarationNone {
		trackedDeclaration = &TrackedDeclaration{
			Is:      isTrackedDeclaration,
			VarName: varName,
		}
	}
	return trackedDeclaration, imports, result
}

func GenerateExpression(variableName *string, expression ast.Expression) (IsTrackedDeclaration, []Import, string) {
	caseModule, caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseArray, caseWhen := expression.ExpressionCases()
	if caseModule != nil {
		return GenerateModule(variableName, *caseModule)
	} else if caseLiteral != nil {
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
		_, imports, result := GenerateDeclaration(caseDeclaration)
		return IsTrackedDeclarationNone, imports, result
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

func GenerateReference(reference ast.Reference) ([]Import, string) {
	allImports := []Import{}
	result := VariableName(reference.Name)

	return allImports, result
}

func GenerateAccess(access ast.Access) ([]Import, string) {
	allImports := []Import{}
	result := ""

	_, imports, over := GenerateExpression(nil, access.Over)
	allImports = append(allImports, imports...)
	result += fmt.Sprintf("%s.(map[string]any)[\"%s\"]", over, access.Access)

	return allImports, result
}

func GenerateInvocation(invocation ast.Invocation) ([]Import, string) {
	allImports := []Import{}

	_, imports, over := GenerateExpression(nil, invocation.Over)
	allImports = append(allImports, imports...)

	funcArgList := ""
	argsCode := ""
	for i, argument := range invocation.Arguments {
		if i > 0 {
			funcArgList += ","
			argsCode += ", "
		}
		funcArgList += "any"

		_, imports, arg := GenerateExpression(nil, argument)
		allImports = append(allImports, imports...)
		argsCode += arg
	}

	result := fmt.Sprintf(`%s.(func(%s)any)(%s)`, over, funcArgList, argsCode)

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

	for i, expression := range function.Block {
		_, imports, exp := GenerateExpression(nil, expression)
		if i == len(function.Block)-1 {
			result += "return "
		}
		result += exp + "\n"
		allImports = append(allImports, imports...)
	}

	if len(function.Block) == 0 {
		result += "return nil\n"
	}

	result += "}"
	return allImports, result
}

func GenerateModule(variableName *string, module ast.Module) (IsTrackedDeclaration, []Import, string) {
	isTrackedDeclaration := IsTrackedDeclarationNone
	if module.Implements.Package == "tenecs.os" && module.Implements.Name == "Main" {
		isTrackedDeclaration = IsTrackedDeclarationMain
	} else if module.Implements.Package == "tenecs.test" && module.Implements.Name == "UnitTests" {
		isTrackedDeclaration = IsTrackedDeclarationUnitTest
	}

	varName := "m"
	if variableName != nil {
		varName = *variableName
	}

	allImports := []Import{}
	result := "func() any {\n"
	result += fmt.Sprintf("var %s any = map[string]any{}\n", VariableName(varName))
	moduleVariables := []string{}
	for varName, _ := range module.Variables {
		moduleVariables = append(moduleVariables, varName)
	}
	sort.Strings(moduleVariables)
	for _, variableName := range moduleVariables {
		result += fmt.Sprintf("var %s any\n", VariableName(variableName))
	}
	for _, variableName := range moduleVariables {
		exp := module.Variables[variableName]
		_, imports, expStr := GenerateExpression(&variableName, exp)
		result += fmt.Sprintf("%s = %s\n", VariableName(variableName), expStr)
		result += fmt.Sprintf("%s.(map[string]any)[\"%s\"] = %s\n", VariableName(varName), variableName, VariableName(variableName))
		allImports = append(allImports, imports...)
	}
	result += fmt.Sprintf("return P%s\n", varName)
	result += "}()"

	return isTrackedDeclaration, allImports, result
}
