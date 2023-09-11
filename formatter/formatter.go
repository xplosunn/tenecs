package formatter

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"sort"
	"strconv"
	"strings"
)

func DisplayFileTopLevel(parsed parser.FileTopLevel) string {
	pkg, imports, declarations := parser.FileTopLevelFields(parsed)

	result := DisplayPackage(pkg)
	result += "\n\n"

	importsCode := []string{}
	for _, impt := range imports {
		importsCode = append(importsCode, DisplayImport(impt))
	}
	sort.Strings(importsCode)
	result += strings.Join(importsCode, "\n")
	if len(importsCode) > 0 {
		result += "\n"
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
		if line != "" {
			lines[i] = "  " + line
		} else {
			lines[i] = ""
		}

	}
	return strings.Join(lines, "\n")
}

func mapNameToString(collection []parser.Name) []string {
	return mapTo(collection, func(item parser.Name) string {
		return item.String
	})
}

func mapTo[T any, R any](collection []T, iteratee func(item T) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		res := iteratee(item)

		result[i] = res
	}

	return result
}

func DisplayPackage(pkg parser.Package) string {
	result := "package "
	for i, name := range pkg.DotSeparatedNames {
		if i > 0 {
			result += "."
		}
		result += name.String
	}
	return result
}

func DisplayImport(impt parser.Import) string {
	return fmt.Sprintf("import %s", strings.Join(mapNameToString(impt.DotSeparatedVars), "."))
}

func DisplayTopLevelDeclaration(topLevelDec parser.TopLevelDeclaration) string {
	var result string
	parser.TopLevelDeclarationExhaustiveSwitch(
		topLevelDec,
		func(topLevelDeclaration parser.Declaration) {
			result = DisplayDeclaration(topLevelDeclaration)
		},
		func(topLevelDeclaration parser.Interface) {
			result = DisplayInterface(topLevelDeclaration)
		},
		func(topLevelDeclaration parser.Struct) {
			result = DisplayStruct(topLevelDeclaration)
		},
	)
	return result
}

func DisplayStruct(struc parser.Struct) string {
	name, generics, variables := parser.StructFields(struc)
	result := "struct " + name.String
	if len(generics) > 0 {
		result += "<"
		for i, generic := range generics {
			if i > 0 {
				result += ", "
			}
			result += generic.String
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
	return name.String + ": " + DisplayTypeAnnotation(typeAnnotation)
}

func DisplayInterface(interf parser.Interface) string {
	name, genericNames, variables := parser.InterfaceFields(interf)
	generics := ""
	if len(genericNames) > 0 {
		generics = "<"
		for i, genericName := range genericNames {
			if i > 0 {
				generics += ", "
			}
			generics += genericName.String
		}
		generics += ">"
	}
	result := "interface " + name.String + generics + " {\n"
	for _, interfaceVariable := range variables {
		result += identLines(DisplayInterfaceVariable(interfaceVariable)) + "\n"
	}
	result += "}"
	return result
}

func DisplayInterfaceVariable(interfaceVariable parser.InterfaceVariable) string {
	name, typeAnnotation := parser.InterfaceVariableFields(interfaceVariable)
	return "public " + name.String + ": " + DisplayTypeAnnotation(typeAnnotation)
}

func DisplayModule(module parser.Module) string {
	implementing, generics, declarations := parser.ModuleFields(module)

	genericsStr := ""
	if len(generics) > 0 {
		genericsStr = "<"
		for i, generic := range generics {
			if i > 0 {
				genericsStr += ", "
			}
			genericsStr += DisplayTypeAnnotation(generic)
		}
		genericsStr += ">"
	}

	result := fmt.Sprintf("implement %s%s {\n", implementing.String, genericsStr)

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
	isPublic, name, typeAnnotation, expression := parser.ModuleDeclarationFields(moduleDeclaration)
	result := ""
	if isPublic {
		result += "public "
	}
	result += name.String
	if typeAnnotation != nil {
		result += ": " + DisplayTypeAnnotation(*typeAnnotation) + " = "
	} else {
		result += " := "
	}
	result += DisplayExpression(expression)
	return result
}

func DisplayExpression(expression parser.Expression) string {
	result := ""
	parser.ExpressionExhaustiveSwitch(
		expression,
		func(expression parser.Module) {
			result = DisplayModule(expression)
		},
		func(expression parser.LiteralExpression) {
			result = DisplayLiteralExpression(expression)
		},
		func(expression parser.ReferenceOrInvocation) {
			result = DisplayReferenceOrInvocation(expression)
		},
		func(expression parser.Lambda) {
			result = DisplayLambda(expression)
		},
		func(expression parser.Declaration) {
			result = DisplayDeclaration(expression)
		},
		func(expression parser.If) {
			result = DisplayIf(expression)
		},
		func(expression parser.Array) {
			result = DisplayArray(expression)
		},
		func(expression parser.When) {
			result = DisplayWhen(expression)
		},
	)
	return result
}

func DisplayArray(array parser.Array) string {
	result := "["
	if array.Generic != nil {
		result += DisplayTypeAnnotation(*array.Generic)
	}
	result += "]("
	for i, expressionBox := range array.Expressions {
		if i > 0 {
			result += ", "
		}
		result += DisplayExpressionBox(expressionBox)
	}
	result += ")"
	return result
}

func DisplayIf(parserIf parser.If) string {
	condition, thenBlock, elseIfs, elseBlock := parser.IfFields(parserIf)
	result := "if " + DisplayExpressionBox(condition) + " {\n"
	for _, expressionBox := range thenBlock {
		result += identLines(DisplayExpressionBox(expressionBox)) + "\n"
	}
	result += "}"
	for _, elseIf := range elseIfs {
		result += " else if " + DisplayExpressionBox(elseIf.Condition) + " {\n"
		for _, expressionBox := range elseIf.ThenBlock {
			result += identLines(DisplayExpressionBox(expressionBox)) + "\n"
		}
		result += "}"
	}
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
	name, typeAnnotation, expressionBox := parser.DeclarationFields(declaration)
	result := name.String
	if typeAnnotation != nil {
		result += ": " + DisplayTypeAnnotation(*typeAnnotation) + " = "
	} else {
		result += " := "
	}
	result += DisplayExpressionBox(expressionBox)
	return result
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
			result += generic.String
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
		for i, expressionBox := range block {
			if _, ok := expressionBox.Expression.(parser.Declaration); ok && i > 0 {
				result += "\n"
			}
			result += identLines(DisplayExpressionBox(expressionBox)) + "\n"
		}
		result += "}"
	}
	return result
}

func DisplayParameter(parameter parser.Parameter) string {
	name, typeAnnotationPtr := parser.ParameterFields(parameter)
	result := name.String
	if typeAnnotationPtr != nil {
		result += ": " + DisplayTypeAnnotation(*typeAnnotationPtr)
	}
	return result
}

func DisplayReferenceOrInvocation(referenceOrInvocation parser.ReferenceOrInvocation) string {
	varName, argumentsListPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)
	result := varName.String
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
				result += DisplayTypeAnnotation(generic)
			}
			result += ">"
		}
		result += "("

		arguments := []string{}
		lineSplitting := false
		for i, argument := range argumentsListPtr.Arguments {
			str := DisplayExpressionBox(argument)
			arguments = append(arguments, str)
			if i < len(argumentsListPtr.Arguments)-1 && strings.Contains(str, "\n") {
				lineSplitting = true
			}
		}

		for i, argument := range arguments {
			if i > 0 {
				if lineSplitting {
					result += ","
				} else {
					result += ", "
				}
			}
			if lineSplitting {
				result += "\n" + identLines(argument)
			} else {
				result += argument
			}
			if lineSplitting && i == len(arguments)-1 {
				result += "\n"
			}
		}
		result += ")"
	}
	return result
}

func DisplayExpressionBox(expressionBox parser.ExpressionBox) string {
	expression, accessOrInvocationChain := parser.ExpressionBoxFields(expressionBox)
	result := DisplayExpression(expression)
	for _, accessOrInvocation := range accessOrInvocationChain {
		result += "." + accessOrInvocation.VarName.String
		result += DisplayArgumentsList(accessOrInvocation.Arguments)
	}
	return result
}

func DisplayLiteralExpression(expression parser.LiteralExpression) string {
	result := ""
	parser.LiteralExhaustiveSwitch(
		expression.Literal,
		func(literal float64) { result = fmt.Sprintf("%f", literal) },
		func(literal int) { result = fmt.Sprintf("%d", literal) },
		func(literal string) { result = literal },
		func(literal bool) { result = strconv.FormatBool(literal) },
		func() { result = "null" },
	)
	return result
}

func DisplayTypeAnnotation(typeAnnotation parser.TypeAnnotation) string {
	result := ""
	for i, element := range typeAnnotation.OrTypes {
		if i > 0 {
			result += " | "
		}
		result += DisplayTypeAnnotationElement(element)
	}
	return result
}

func DisplayTypeAnnotationElement(typeAnnotationElement parser.TypeAnnotationElement) string {
	result := ""
	parser.TypeAnnotationElementExhaustiveSwitch(
		typeAnnotationElement,
		func(typeAnnotation parser.SingleNameType) {
			generics := ""
			if len(typeAnnotation.Generics) > 0 {
				generics = "<"
				for i, generic := range typeAnnotation.Generics {
					if i > 0 {
						generics += ", "
					}
					generics += DisplayTypeAnnotation(generic)
				}
				generics += ">"
			}
			result = typeAnnotation.TypeName.String + generics
		},
		func(typeAnnotation parser.FunctionType) {
			result = ""
			if len(typeAnnotation.Generics) > 0 {
				result += "<"
				for i, generic := range typeAnnotation.Generics {
					if i > 0 {
						result += ", "
					}
					result += generic.String
				}
				result += ">"
			}
			result += "("
			for i, argument := range typeAnnotation.Arguments {
				if i > 0 {
					result += ", "
				}
				result += DisplayTypeAnnotation(argument)
			}
			result += ") -> " + DisplayTypeAnnotation(typeAnnotation.ReturnType)
		},
	)
	return result
}

func DisplayWhen(when parser.When) string {
	result := "when "
	result += DisplayExpressionBox(when.Over)
	result += " {\n"

	resultCases := ""
	for i, is := range when.Is {
		resultCases += "is "
		if is.Name != nil {
			resultCases += is.Name.String + ": "
		}
		resultCases += DisplayTypeAnnotation(is.Type)
		resultCases += " => {\n"
		for _, thenExp := range is.ThenBlock {
			resultCases += identLines(DisplayExpressionBox(thenExp)) + "\n"
		}
		resultCases += "}"
		if i < len(when.Is)-1 {
			resultCases += "\n"
		}
	}
	if when.Other != nil {
		resultCases += "\n"
		if when.Other.Name != nil {
			resultCases += "other " + when.Other.Name.String + " => {\n"
		} else {
			resultCases += "other => {\n"
		}
		for _, thenExp := range when.Other.ThenBlock {
			resultCases += identLines(DisplayExpressionBox(thenExp)) + "\n"
		}
		resultCases += "}"
	}

	result += identLines(resultCases) + "\n"
	result += "}"

	return result
}
