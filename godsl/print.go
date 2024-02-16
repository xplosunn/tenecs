package godsl

import (
	"fmt"
	"strings"
)

func Print(godsl GoDSL) string {
	imports, code := print(godsl)
	for i, imp := range imports {
		imports[i] = `"` + imp + `"`
	}
	return fmt.Sprintf(`import (
%s
)

%s`, identLines(strings.Join(imports, "\n")), code)
}

func print(godsl GoDSL) ([]string, string) {
	caseExpression, caseStatement, caseTopLevelStatement := exhaustiveSwitch(godsl)
	if caseExpression != nil {
		caseFunctionCreation, caseFunctionInvocation, caseObjectCreation, caseObjectAccess := (*caseExpression).sealedExpressionCases()
		if caseFunctionCreation != nil {
			panic("TODO godsl Print caseFunctionCreation")
		} else if caseFunctionInvocation != nil {
			panic("TODO godsl Print caseFunctionInvocation")
		} else if caseObjectCreation != nil {
			panic("TODO godsl Print caseObjectCreation")
		} else if caseObjectAccess != nil {
			panic("TODO godsl Print caseObjectAccess")
		} else {
			panic("godsl Print cases caseExpression")
		}
	} else if caseStatement != nil {
		caseVariableDeclaration, caseReturn, caseIf, caseNativeFunctionInvocation := (*caseStatement).sealedStatementCases()
		if caseVariableDeclaration != nil {
			panic("TODO godsl Print caseVariableDeclaration")
		} else if caseReturn != nil {
			panic("TODO godsl Print caseReturn")
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
	params := strings.Join(g.params, ", ")
	result += fmt.Sprintf("%s(%s)", g.name, params)
	return imports, result
}

func printNativeFunctionDeclaration(g goNativeFunctionDeclaration) ([]string, string) {
	imports := []string{}
	params := strings.Join(g.params, ",")
	bodyStatements := []string{}
	for _, statement := range g.body {
		imp, s := print(statement)
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
