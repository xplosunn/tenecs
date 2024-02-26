package godsl

import (
	"fmt"
	"strings"
)

func Print(godsl GoDSL) string {
	imports, code := PrintImportsAndCode(godsl)
	for i, imp := range imports {
		imports[i] = `"` + imp + `"`
	}
	return fmt.Sprintf(`import (
%s
)

%s`, identLines(strings.Join(imports, "\n")), code)
}

func PrintImportsAndCode(godsl GoDSL) ([]string, string) {
	caseExpression, caseStatement, caseTopLevelStatement := exhaustiveSwitch(godsl)
	if caseExpression != nil {
		caseFunctionCreation, caseFunctionInvocation, caseObjectCreation, caseObjectAccess, caseVariableReference, caseCast, caseLiteral := (*caseExpression).sealedExpressionCases()
		if caseFunctionCreation != nil {
			panic("TODO godsl Print caseFunctionCreation")
		} else if caseFunctionInvocation != nil {
			panic("TODO godsl Print caseFunctionInvocation")
		} else if caseObjectCreation != nil {
			panic("TODO godsl Print caseObjectCreation")
		} else if caseObjectAccess != nil {
			panic("TODO godsl Print caseObjectAccess")
		} else if caseVariableReference != nil {
			return printVariableReference(*caseVariableReference)
		} else if caseCast != nil {
			return printCast(*caseCast)
		} else if caseLiteral != nil {
			return printLiteral(*caseLiteral)
		} else {
			panic("godsl Print cases caseExpression")
		}
	} else if caseStatement != nil {
		caseVariableDeclaration, caseReturn, caseIf, caseNativeFunctionInvocation := (*caseStatement).sealedStatementCases()
		if caseVariableDeclaration != nil {
			panic("TODO godsl Print caseVariableDeclaration")
		} else if caseReturn != nil {
			return printReturn(*caseReturn)
		} else if caseIf != nil {
			panic("TODO godsl Print caseIf")
		} else if caseNativeFunctionInvocation != nil {
			return printNativeFunctionInvocation(*caseNativeFunctionInvocation)
		} else {
			panic("godsl Print cases caseStatement")
		}
	} else if caseTopLevelStatement != nil {
		caseNativeFunctionDeclaration := (*caseTopLevelStatement).sealedTopLevelStatementCases()
		if caseNativeFunctionDeclaration != nil {
			return printNativeFunctionDeclaration(*caseNativeFunctionDeclaration)
		} else {
			panic("godsl Print cases caseTopLevelStatement")
		}
	} else {
		panic("godsl Print cases godsl")
	}
}

func printType(g Type) string {
	return g.typeToString()
}

func printLiteral(g goLiteral) ([]string, string) {
	return []string{}, g.value
}

func printCast(g goCast) ([]string, string) {
	imports, code := PrintImportsAndCode(g.expression)
	code += ".(" + printType(g.toType) + ")"
	return imports, code
}

func printVariableReference(g goVariableReference) ([]string, string) {
	return []string{}, g.name
}

func printReturn(g goReturn) ([]string, string) {
	imports, result := PrintImportsAndCode(g.exp)
	result = "return " + result
	return imports, result
}

func printNativeFunctionInvocation(g goNativeFunctionInvocation) ([]string, string) {
	imports := []string{}
	result := ""
	if g.pkg != nil {
		imports = append(imports, *g.pkg)
		result += *g.pkg + "."
	}
	for i, v := range g.vars {
		result += v
		if i < len(g.vars)-1 {
			result += ", "
		} else {
			result += " := "
		}
	}
	params := ""
	for i, param := range g.params {
		if i > 0 {
			params += ","
		}
		impt, code := PrintImportsAndCode(param)
		imports = append(imports, impt...)
		params += code
	}
	result += fmt.Sprintf("%s(%s)", g.name, params)
	return imports, result
}

func printNativeFunctionDeclaration(g goNativeFunctionDeclaration) ([]string, string) {
	imports := []string{}
	params := strings.Join(g.params, ",")
	bodyStatements := []string{}
	for _, statement := range g.body {
		imp, s := PrintImportsAndCode(statement)
		imports = append(imports, imp...)
		bodyStatements = append(bodyStatements, s)
	}
	body := identLines(strings.Join(bodyStatements, "\n"))
	return imports, fmt.Sprintf(`func %s(%s) {
%s
}`, g.name, params, body)
}

func identLines(str string) string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = "	" + line
		} else {
			lines[i] = ""
		}

	}
	return strings.Join(lines, "\n")
}
