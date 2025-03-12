package ast

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/exp/maps"
)

type RefDependencies map[Ref]Set[Ref]

func DetermineRefDependencies(program Program) RefDependencies {
	result := RefDependencies{}
	for ref, exp := range program.Declarations {
		if result[ref] == nil {
			result[ref] = Set[Ref]{}
		}
		result[ref].PutAll(refDependenciesOfExpression(exp))
	}
	for ref, varType := range program.StructFunctions {
		if result[ref] == nil {
			result[ref] = Set[Ref]{}
		}
		result[ref].PutAll(refDependenciesOfVariableType(varType))
	}
	for ref, varType := range program.NativeFunctions {
		if result[ref] == nil {
			result[ref] = Set[Ref]{}
		}
		result[ref].PutAll(refDependenciesOfVariableType(varType))
	}
	for ref, fieldMap := range program.FieldsByType {
		if result[ref] == nil {
			result[ref] = Set[Ref]{}
		}
		for _, varType := range fieldMap {
			result[ref].PutAll(refDependenciesOfVariableType(varType))
		}
	}
	for _, ref := range maps.Keys(result) {
		result[ref].Remove(ref)
	}
	return result
}

func refDependenciesOfExpression(expression Expression) []Ref {
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
		return refDependenciesOfLiteral(*caseLiteral)
	} else if caseReference != nil {
		return refDependenciesOfReference(*caseReference)
	} else if caseAccess != nil {
		return refDependenciesOfAccess(*caseAccess)
	} else if caseInvocation != nil {
		return refDependenciesOfInvocation(*caseInvocation)
	} else if caseFunction != nil {
		return refDependenciesOfFunction(*caseFunction)
	} else if caseDeclaration != nil {
		return refDependenciesOfDeclaration(caseDeclaration)
	} else if caseIf != nil {
		return refDependenciesOfIf(*caseIf)
	} else if caseList != nil {
		return refDependenciesOfList(*caseList)
	} else if caseWhen != nil {
		return refDependenciesOfWhen(*caseWhen)
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func refDependenciesOfWhen(when When) []Ref {
	result := refDependenciesOfVariableType(when.VariableType)
	result = append(result, refDependenciesOfExpression(when.Over)...)
	for _, whenCase := range when.Cases {
		result = append(result, refDependenciesOfVariableType(whenCase.VariableType)...)
		for _, expression := range whenCase.Block {
			result = append(result, refDependenciesOfExpression(expression)...)
		}
	}
	for _, expression := range when.OtherCase {
		result = append(result, refDependenciesOfExpression(expression)...)
	}
	return result
}

func refDependenciesOfList(list List) []Ref {
	result := refDependenciesOfVariableType(list.ContainedVariableType)
	for _, expression := range list.Arguments {
		result = append(result, refDependenciesOfExpression(expression)...)
	}
	return result
}

func refDependenciesOfIf(exp If) []Ref {
	result := refDependenciesOfVariableType(exp.VariableType)
	result = append(result, refDependenciesOfExpression(exp.Condition)...)
	for _, expression := range exp.ThenBlock {
		result = append(result, refDependenciesOfExpression(expression)...)
	}
	for _, expression := range exp.ElseBlock {
		result = append(result, refDependenciesOfExpression(expression)...)
	}
	return result
}

func refDependenciesOfDeclaration(declaration *Declaration) []Ref {
	return refDependenciesOfExpression(declaration.Expression)
}

func refDependenciesOfFunction(function Function) []Ref {
	result := refDependenciesOfVariableType(function.VariableType)
	for _, expression := range function.Block {
		result = append(result, refDependenciesOfExpression(expression)...)
	}
	return result
}

func refDependenciesOfInvocation(invocation Invocation) []Ref {
	result := refDependenciesOfVariableType(invocation.VariableType)
	result = append(result, refDependenciesOfExpression(invocation.Over)...)
	for _, expression := range invocation.Arguments {
		result = append(result, refDependenciesOfExpression(expression)...)
	}
	return result
}

func refDependenciesOfAccess(access Access) []Ref {
	result := refDependenciesOfVariableType(access.VariableType)
	result = append(result, refDependenciesOfExpression(access.Over)...)
	return result
}

func refDependenciesOfReference(reference Reference) []Ref {
	result := refDependenciesOfVariableType(reference.VariableType)
	if reference.PackageName != nil {
		result = append(result, Ref{
			Package: *reference.PackageName,
			Name:    reference.Name,
		})
	}
	return result
}

func refDependenciesOfLiteral(literal Literal) []Ref {
	return []Ref{}
}

func refDependenciesOfVariableType(variableType types.VariableType) []Ref {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		return []Ref{}
	} else if caseList != nil {
		return refDependenciesOfVariableType(caseList.Generic)
	} else if caseKnownType != nil {
		if caseKnownType.Package == "" {
			return []Ref{}
		}
		return []Ref{
			Ref{
				Package: caseKnownType.Package,
				Name:    caseKnownType.Name,
			},
		}
	} else if caseFunction != nil {
		result := refDependenciesOfVariableType(caseFunction.ReturnType)
		for _, functionArgument := range caseFunction.Arguments {
			result = append(result, refDependenciesOfVariableType(functionArgument.VariableType)...)
		}
		return result
	} else if caseOr != nil {
		result := []Ref{}
		for _, element := range caseOr.Elements {
			result = append(result, refDependenciesOfVariableType(element)...)
		}
		return result
	} else {
		panic(fmt.Errorf("cases on %v", variableType))
	}
}
