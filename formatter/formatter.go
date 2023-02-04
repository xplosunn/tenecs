package formatter

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"strconv"
	"strings"
)

func DisplayFileTopLevel(parsed parser.FileTopLevel) string {
	pkg, imports, declarations := parser.FileTopLevelFields(parsed)

	result := DisplayPackage(pkg)
	result += "\n\n"
	for _, impt := range imports {
		result += DisplayImport(impt) + "\n"
	}
	result += "\n"
	for i, topLevelDeclaration := range declarations {
		if i > 0 {
			result += "\n"
		}
		result += DisplayTopLevelDeclaration(topLevelDeclaration) + "\n"
	}
	return result
}

func identLines(str string) string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = "  " + line
	}
	return strings.Join(lines, "\n")
}

func DisplayPackage(pkg parser.Package) string {
	return fmt.Sprintf("package %s", pkg.Identifier)
}

func DisplayImport(impt parser.Import) string {
	return fmt.Sprintf("import %s", strings.Join(impt.DotSeparatedVars, "."))
}

func DisplayTopLevelDeclaration(topLevelDec parser.TopLevelDeclaration) string {
	caseDeclaration, caseInterface, caseStruct := topLevelDec.TopLevelDeclarationCases()
	if caseDeclaration != nil {
		return DisplayDeclaration(*caseDeclaration)
	} else if caseInterface != nil {
		return DisplayInterface(*caseInterface)
	} else if caseStruct != nil {
		return DisplayStruct(*caseStruct)
	} else {
		panic(fmt.Errorf("cases on %v", topLevelDec))
	}
}

func DisplayStruct(struc parser.Struct) string {
	name, generics, variables := parser.StructFields(struc)
	result := "struct " + name
	if len(generics) > 0 {
		result += "<"
		for i, generic := range generics {
			if i > 0 {
				result += ", "
			}
			result += generic
		}
		result += ">"
	}
	result += "(\n"
	for i, structVariable := range variables {
		result += identLines(DisplayStructVariable(structVariable))
		if i < len(variables)-1 {
			result += ",\n"
		} else {
			result += "\n"
		}
	}
	result += ")"
	return result
}

func DisplayStructVariable(structVariable parser.StructVariable) string {
	name, typeAnnotation := parser.StructVariableFields(structVariable)
	return name + ": " + DisplayTypeAnnotation(typeAnnotation)
}

func DisplayInterface(interf parser.Interface) string {
	name, variables := parser.InterfaceFields(interf)
	result := "interface " + name + " {\n"
	for _, interfaceVariable := range variables {
		result += identLines(DisplayInterfaceVariable(interfaceVariable)) + "\n"
	}
	result += "}"
	return result
}

func DisplayInterfaceVariable(interfaceVariable parser.InterfaceVariable) string {
	name, typeAnnotation := parser.InterfaceVariableFields(interfaceVariable)
	return name + ": " + DisplayTypeAnnotation(typeAnnotation)
}

func DisplayModule(module parser.Module) string {
	implementing, declarations := parser.ModuleFields(module)

	result := fmt.Sprintf("implement %s {\n", implementing)

	for i, moduleDeclaration := range declarations {
		result += identLines(DisplayModuleDeclaration(moduleDeclaration)) + "\n"
		if i < len(declarations)-1 {
			result += "\n"
		}
	}

	result += "}"
	return result
}

func DisplayModuleDeclaration(moduleDeclaration parser.ModuleDeclaration) string {
	isPublic, name, expression := parser.ModuleDeclarationFields(moduleDeclaration)
	result := ""
	if isPublic {
		result += "public "
	}
	result += name + " := "
	result += DisplayExpression(expression)
	return result
}

func DisplayExpression(expression parser.Expression) string {
	caseModule, caseLiteralExpression, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseModule != nil {
		return DisplayModule(*caseModule)
	} else if caseLiteralExpression != nil {
		return DisplayLiteralExpression(*caseLiteralExpression)
	} else if caseReferenceOrInvocation != nil {
		return DisplayReferenceOrInvocation(*caseReferenceOrInvocation)
	} else if caseLambda != nil {
		return DisplayLambda(*caseLambda)
	} else if caseDeclaration != nil {
		return DisplayDeclaration(*caseDeclaration)
	} else if caseIf != nil {
		return DisplayIf(*caseIf)
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func DisplayIf(parserIf parser.If) string {
	condition, thenBlock, elseBlock := parser.IfFields(parserIf)
	result := "if " + DisplayExpressionBox(condition) + " {\n"
	for _, expressionBox := range thenBlock {
		result += identLines(DisplayExpressionBox(expressionBox)) + "\n"
	}
	result += "}"
	if len(elseBlock) > 0 {
		result += " else {\n"
		for _, expressionBox := range elseBlock {
			result += identLines(DisplayExpressionBox(expressionBox)) + "\n"
		}
		result += "}"
	}
	return result
}

func DisplayDeclaration(declaration parser.Declaration) string {
	name, expressionBox := parser.DeclarationFields(declaration)
	return name + " := " + DisplayExpressionBox(expressionBox)
}

func DisplayLambda(lambda parser.Lambda) string {
	generics, parameters, returnTypePtr, block := parser.LambdaFields(lambda)
	result := ""
	if len(generics) > 0 {
		result += "<"
		for i, generic := range generics {
			if i > 0 {
				result += ", "
			}
			result += generic
		}
		result += ">"
	}
	result += "("
	for i, parameter := range parameters {
		if i > 0 {
			result += ", "
		}
		result += DisplayParameter(parameter)
	}
	result += ")"
	if returnTypePtr != nil {
		result += ": " + DisplayTypeAnnotation(*returnTypePtr)
	}
	if len(block) == 1 {
		expressionBox := block[0]
		noAccessOrInvocations := expressionBox.AccessOrInvocationChain == nil || len(expressionBox.AccessOrInvocationChain) == 0
		_, isModule := expressionBox.Expression.(parser.Module)
		if noAccessOrInvocations && isModule {
			result += " => "
			for _, expressionBox := range block {
				result += DisplayExpressionBox(expressionBox)
			}
		} else {
			result += " => {\n"
			for _, expressionBox := range block {
				result += identLines(DisplayExpressionBox(expressionBox)) + "\n"
			}
			result += "}"
		}
	} else {
		result += " => {\n"
		for _, expressionBox := range block {
			result += identLines(DisplayExpressionBox(expressionBox)) + "\n"
		}
		result += "}"
	}
	return result
}

func DisplayParameter(parameter parser.Parameter) string {
	name, typeAnnotationPtr := parser.ParameterFields(parameter)
	result := name
	if typeAnnotationPtr != nil {
		result += ": " + DisplayTypeAnnotation(*typeAnnotationPtr)
	}
	return result
}

func DisplayReferenceOrInvocation(referenceOrInvocation parser.ReferenceOrInvocation) string {
	varName, argumentsListPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)
	result := varName
	result += DisplayArgumentsList(argumentsListPtr)
	return result
}

func DisplayArgumentsList(argumentsListPtr *parser.ArgumentsList) string {
	result := ""
	if argumentsListPtr != nil {
		if len(argumentsListPtr.Generics) > 0 {
			result += "<"
			for i, generic := range argumentsListPtr.Generics {
				if i > 0 {
					result += ", "
				}
				result += generic
			}
			result += ">"
		}
		result += "("
		for i, argument := range argumentsListPtr.Arguments {
			if i > 0 {
				result += ", "
			}
			result += DisplayExpressionBox(argument)
		}
		result += ")"
	}
	return result
}

func DisplayExpressionBox(expressionBox parser.ExpressionBox) string {
	expression, accessOrInvocationChain := parser.ExpressionBoxFields(expressionBox)
	result := DisplayExpression(expression)
	for _, accessOrInvocation := range accessOrInvocationChain {
		result += "." + accessOrInvocation.VarName
		result += DisplayArgumentsList(accessOrInvocation.Arguments)
	}
	return result
}

func DisplayLiteralExpression(expression parser.LiteralExpression) string {
	return parser.LiteralFold(
		expression.Literal,
		func(arg float64) string { return fmt.Sprintf("%f", arg) },
		func(arg int) string { return fmt.Sprintf("%d", arg) },
		func(arg string) string { return arg },
		func(arg bool) string { return strconv.FormatBool(arg) },
	)
}

func DisplayTypeAnnotation(typeAnnotation parser.TypeAnnotation) string {
	caseSingleName, caseFunctionType := typeAnnotation.TypeAnnotationCases()
	if caseSingleName != nil {
		return caseSingleName.TypeName
	} else if caseFunctionType != nil {
		result := "("
		for i, argument := range caseFunctionType.Arguments {
			if i > 0 {
				result += ", "
			}
			result += DisplayTypeAnnotation(argument)
		}
		return result + ") -> " + DisplayTypeAnnotation(caseFunctionType.ReturnType)
	} else {
		panic(fmt.Errorf("cases on %v", typeAnnotation))
	}
}
