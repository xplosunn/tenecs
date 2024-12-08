package formatter

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/xplosunn/tenecs/parser"
	"sort"
	"strconv"
	"strings"
	"text/scanner"
)

func displayFileTopLevel(parsed parser.FileTopLevel, ignoreComments bool) string {
	pkg, imports, declarations := parser.FileTopLevelFields(parsed)

	result, tokens := displayPackage(pkg, parsed.Tokens, ignoreComments)
	result += "\n\n"

	importsCode := []string{}
	for _, impt := range imports {
		r, t := displayImport(impt, tokens, ignoreComments)
		tokens = t
		importsCode = append(importsCode, r)
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
		r, t := displayTopLevelDeclaration(topLevelDeclaration, tokens, ignoreComments)
		tokens = t
		result += r + "\n"
	}

	result += displayRemainingComments(tokens, ignoreComments)

	return result
}

func indentLines(str string) string {
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

func displayRemainingCommentsBeforeNodeWithValue(value string, tokens []lexer.Token, ignoreComments bool) (string, []lexer.Token) {
	if ignoreComments {
		return "", tokens
	}
	for _, token := range tokens {
		if token.Value == value {
			return displayRemainingCommentsBeforeNode(parser.Node{
				Pos:    token.Pos,
				EndPos: token.Pos,
			}, tokens, ignoreComments)
		}
	}
	panic("Reached end of loop on displayRemainingCommentsBeforeNodeWithValue: " + value)
}

func displayRemainingCommentsBeforeNode(node parser.Node, tokens []lexer.Token, ignoreComments bool) (string, []lexer.Token) {
	if ignoreComments {
		return "", tokens
	}
	comments := ""
	for i, token := range tokens {
		if token.Pos.Line >= node.Pos.Line && token.Pos.Column >= node.Pos.Column {
			if i < len(tokens)-1 {
				return comments, tokens[i+1:]
			} else {
				return comments, []lexer.Token{}
			}
		}
		if token.Type == scanner.Comment {
			comments += token.Value + "\n"
		}
	}
	return comments, []lexer.Token{}
}

func displayRemainingComments(tokens []lexer.Token, ignoreComments bool) string {
	if ignoreComments {
		return ""
	}
	comments := ""
	for _, token := range tokens {
		if token.Type == scanner.Comment {
			comments += token.Value + "\n"
		}
	}
	return comments
}

func lastOfNonEmptySlice[T any](s []T) T {
	if len(s) == 0 {
		panic("didn't expect empty slice")
	}
	return s[len(s)-1]
}

func displayPackage(pkg parser.Package, tokens []lexer.Token, ignoreComments bool) (string, []lexer.Token) {
	result, tokens := displayRemainingCommentsBeforeNode(lastOfNonEmptySlice(pkg.DotSeparatedNames).Node, tokens, ignoreComments)
	result += "package "
	for i, name := range pkg.DotSeparatedNames {
		if i > 0 {
			result += "."
		}
		result += name.String
	}
	return result, tokens
}

func displayImport(impt parser.Import, tokens []lexer.Token, ignoreComments bool) (string, []lexer.Token) {
	result, tokens := displayRemainingCommentsBeforeNode(lastOfNonEmptySlice(impt.DotSeparatedVars).Node, tokens, ignoreComments)
	result += fmt.Sprintf("import %s", strings.Join(mapNameToString(impt.DotSeparatedVars), "."))
	if impt.As != nil {
		result += " as " + impt.As.String
	}
	return result, tokens
}

func displayTopLevelDeclaration(topLevelDec parser.TopLevelDeclaration, tokens []lexer.Token, ignoreComments bool) (string, []lexer.Token) {
	var result string
	parser.TopLevelDeclarationExhaustiveSwitch(
		topLevelDec,
		func(topLevelDeclaration parser.Declaration) {
			result, tokens = displayRemainingCommentsBeforeNode(topLevelDeclaration.ExpressionBox.Node, tokens, ignoreComments)
			result += displayDeclaration(topLevelDeclaration)
		},
		func(topLevelDeclaration parser.Struct) {
			result, tokens = displayRemainingCommentsBeforeNode(topLevelDeclaration.Name.Node, tokens, ignoreComments)
			r, t := displayStruct(topLevelDeclaration, tokens, ignoreComments)
			tokens = t
			result += r
		},
		func(topLevelDeclaration parser.TypeAlias) {
			result, tokens = displayRemainingCommentsBeforeNode(topLevelDeclaration.Name.Node, tokens, ignoreComments)
			r, t := displayTypeAlias(topLevelDeclaration, tokens, ignoreComments)
			tokens = t
			result += r
		},
	)
	return result, tokens
}

func displayTypeAlias(typeAlias parser.TypeAlias, tokens []lexer.Token, ignoreComments bool) (string, []lexer.Token) {
	name, generics, typ := parser.TypeAliasFields(typeAlias)
	result := "typealias " + name.String
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
	r, t := displayRemainingCommentsBeforeNode(typ.Node, tokens, ignoreComments)
	result += r
	tokens = t
	result += " = " + displayTypeAnnotation(typ)
	return result, tokens
}

func displayStruct(struc parser.Struct, tokens []lexer.Token, ignoreComments bool) (string, []lexer.Token) {
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
		r, t := displayRemainingCommentsBeforeNode(structVariable.Type.Node, tokens, ignoreComments)
		tokens = t
		result += indentLines(r)
		result += indentLines(displayStructVariable(structVariable))
		if i < len(variables)-1 {
			result += ",\n"
		} else {
			result += "\n"
		}
	}
	r, t := displayRemainingCommentsBeforeNodeWithValue(")", tokens, ignoreComments)
	tokens = t
	result += indentLines(r)
	result += ")"
	return result, tokens
}

func displayStructVariable(structVariable parser.StructVariable) string {
	name, typeAnnotation := parser.StructVariableFields(structVariable)
	return name.String + ": " + displayTypeAnnotation(typeAnnotation)
}

func displayExpression(expression parser.Expression) string {
	result := ""
	parser.ExpressionExhaustiveSwitch(
		expression,
		func(expression parser.LiteralExpression) {
			result = displayLiteralExpression(expression)
		},
		func(expression parser.ReferenceOrInvocation) {
			result = displayReferenceOrInvocation(expression)
		},
		func(expression parser.Lambda) {
			result = displayLambda(expression)
		},
		func(expression parser.Declaration) {
			result = displayDeclaration(expression)
		},
		func(expression parser.If) {
			result = displayIf(expression)
		},
		func(expression parser.List) {
			result = displayList(expression)
		},
		func(expression parser.When) {
			result = displayWhen(expression)
		},
	)
	return result
}

func displayList(list parser.List) string {
	result := "["
	if list.Generic != nil {
		result += displayTypeAnnotation(*list.Generic)
	}
	result += "]("
	for i, expressionBox := range list.Expressions {
		if i > 0 {
			result += ", "
		}
		result += displayExpressionBox(expressionBox)
	}
	result += ")"
	return result
}

func displayIf(parserIf parser.If) string {
	condition, thenBlock, elseIfs, elseBlock := parser.IfFields(parserIf)
	result := "if " + displayExpressionBox(condition) + " {\n"
	for _, expressionBox := range thenBlock {
		result += indentLines(displayExpressionBox(expressionBox)) + "\n"
	}
	result += "}"
	for _, elseIf := range elseIfs {
		result += " else if " + displayExpressionBox(elseIf.Condition) + " {\n"
		for _, expressionBox := range elseIf.ThenBlock {
			result += indentLines(displayExpressionBox(expressionBox)) + "\n"
		}
		result += "}"
	}
	if len(elseBlock) > 0 {
		result += " else {\n"
		for _, expressionBox := range elseBlock {
			result += indentLines(displayExpressionBox(expressionBox)) + "\n"
		}
		result += "}"
	}
	return result
}

func displayDeclaration(declaration parser.Declaration) string {
	name, typeAnnotation, shortcircuit, expressionBox := parser.DeclarationFields(declaration)
	result := name.String
	if shortcircuit != nil {
		if shortcircuit.TypeAnnotation != nil {
			if typeAnnotation != nil {
				result += ": " + displayTypeAnnotation(*typeAnnotation) + " ? " + displayTypeAnnotation(*shortcircuit.TypeAnnotation) + " = "
			} else {
				result += " :? " + displayTypeAnnotation(*shortcircuit.TypeAnnotation) + " = "
			}
		} else {
			if typeAnnotation != nil {
				result += ": " + displayTypeAnnotation(*typeAnnotation) + " ?= "
			} else {
				result += " :?= "
			}
		}
	} else {
		if typeAnnotation != nil {
			result += ": " + displayTypeAnnotation(*typeAnnotation) + " = "
		} else {
			result += " := "
		}
	}
	result += displayExpressionBox(expressionBox)
	return result
}

func displayLambda(lambda parser.Lambda) string {
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
		result += displayParameter(parameter)
	}
	result += ")"
	if returnTypePtr != nil {
		result += ": " + displayTypeAnnotation(*returnTypePtr)
	}
	result += " => {\n"
	for i, expressionBox := range block {
		if _, ok := expressionBox.Expression.(parser.Declaration); ok && i > 0 {
			result += "\n"
		}
		result += indentLines(displayExpressionBox(expressionBox)) + "\n"
	}
	result += "}"
	return result
}

func displayParameter(parameter parser.Parameter) string {
	name, typeAnnotationPtr := parser.ParameterFields(parameter)
	result := name.String
	if typeAnnotationPtr != nil {
		result += ": " + displayTypeAnnotation(*typeAnnotationPtr)
	}
	return result
}

func displayReferenceOrInvocation(referenceOrInvocation parser.ReferenceOrInvocation) string {
	varName, argumentsListPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)
	result := varName.String
	result += displayArgumentsList(argumentsListPtr)
	return result
}

func displayArgumentsList(argumentsListPtr *parser.ArgumentsList) string {
	result := ""
	if argumentsListPtr != nil {
		if len(argumentsListPtr.Generics) > 0 {
			result += "<"
			for i, generic := range argumentsListPtr.Generics {
				if i > 0 {
					result += ", "
				}
				result += displayTypeAnnotation(generic)
			}
			result += ">"
		}
		result += "("

		arguments := []string{}
		lineSplitting := false
		for i, argument := range argumentsListPtr.Arguments {
			if argument.Name != nil {
				lineSplitting = true
			}
			str := displayNamedArgument(argument)
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
				result += "\n" + indentLines(argument)
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

func displayNamedArgument(namedArgument parser.NamedArgument) string {
	name, expressionBox := parser.NamedArgumentFields(namedArgument)
	result := displayExpressionBox(expressionBox)
	if name != nil {
		result = name.String + " = " + result
	}
	return result
}

func displayExpressionBox(expressionBox parser.ExpressionBox) string {
	expression, accessOrInvocationChain := parser.ExpressionBoxFields(expressionBox)
	result := displayExpression(expression)
	for _, accessOrInvocation := range accessOrInvocationChain {
		if accessOrInvocation.DotOrArrowName != nil {
			separator := "."
			if accessOrInvocation.DotOrArrowName.Arrow {
				separator = "->"
			}
			result += separator + accessOrInvocation.DotOrArrowName.VarName.String
		}
		if accessOrInvocation.Arguments != nil {
			result += displayArgumentsList(accessOrInvocation.Arguments)
		}
	}
	return result
}

func displayLiteralExpression(expression parser.LiteralExpression) string {
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

func displayTypeAnnotation(typeAnnotation parser.TypeAnnotation) string {
	result := ""
	for i, element := range typeAnnotation.OrTypes {
		if i > 0 {
			result += " | "
		}
		result += displayTypeAnnotationElement(element)
	}
	return result
}

func displayTypeAnnotationElement(typeAnnotationElement parser.TypeAnnotationElement) string {
	result := ""
	parser.TypeAnnotationElementExhaustiveSwitch(
		typeAnnotationElement,
		func(underscoreTypeAnnotation parser.SingleNameType) {
			result = "_"
		},
		func(typeAnnotation parser.SingleNameType) {
			generics := ""
			if len(typeAnnotation.Generics) > 0 {
				generics = "<"
				for i, generic := range typeAnnotation.Generics {
					if i > 0 {
						generics += ", "
					}
					generics += displayTypeAnnotation(generic)
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
				if argument.Name != nil && argument.Name.String != "_" {
					result += argument.Name.String + ": "
				}
				result += displayTypeAnnotation(argument.Type)
			}
			result += ") ~> " + displayTypeAnnotation(typeAnnotation.ReturnType)
		},
	)
	return result
}

func displayWhen(when parser.When) string {
	result := "when "
	result += displayExpressionBox(when.Over)
	result += " {\n"

	resultCases := ""
	for i, is := range when.Is {
		resultCases += "is "
		if is.Name != nil {
			resultCases += is.Name.String + ": "
		}
		resultCases += displayTypeAnnotation(is.Type)
		resultCases += " => {\n"
		for _, thenExp := range is.ThenBlock {
			resultCases += indentLines(displayExpressionBox(thenExp)) + "\n"
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
			resultCases += indentLines(displayExpressionBox(thenExp)) + "\n"
		}
		resultCases += "}"
	}

	result += indentLines(resultCases) + "\n"
	result += "}"

	return result
}
