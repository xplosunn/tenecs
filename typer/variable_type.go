package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/types"
	"reflect"
)

func variableTypeContainedIn(sub types.VariableType, super types.VariableType) bool {
	superOr, ok := super.(*types.OrVariableType)
	if !ok {
		return variableTypeEq(sub, super)
	}
	subOr, ok := sub.(*types.OrVariableType)
	if !ok {
		for _, superElement := range superOr.Elements {
			if variableTypeEq(sub, superElement) {
				return true
			}
		}
		return false
	}
	for _, subElement := range subOr.Elements {
		if !variableTypeContainedIn(subElement, super) {
			return false
		}
	}
	return true
}

func variableTypeAddToOr(varType types.VariableType, or *types.OrVariableType) {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	_ = caseTypeArgument
	_ = caseKnownType
	_ = caseFunction
	if caseOr != nil {
		for _, element := range caseOr.Elements {
			variableTypeAddToOr(element, or)
		}
	} else {
		if !variableTypeContainedIn(varType, or) {
			or.Elements = append(or.Elements, varType)
		}
	}
}

func variableTypeCombine(v1 types.VariableType, v2 types.VariableType) types.VariableType {
	result := &types.OrVariableType{Elements: []types.VariableType{}}

	addAll := func(varType types.VariableType) {
		variableTypeAddToOr(varType, result)
	}

	addAll(v1)
	addAll(v2)

	if len(result.Elements) == 1 {
		return result.Elements[0]
	}

	return result
}

func variableTypeEq(v1 types.VariableType, v2 types.VariableType) bool {
	if v1 == nil || v2 == nil {
		panic(fmt.Errorf("trying to eq %v to %v", v1, v2))
	}
	v1CaseTypeArgument, v1CaseKnownType, v1CaseFunction, v1CaseOr := v1.VariableTypeCases()
	v2CaseTypeArgument, v2CaseKnownType, v2CaseFunction, v2CaseOr := v2.VariableTypeCases()
	if v1CaseKnownType != nil && v2CaseKnownType != nil {
		if v1CaseKnownType.Package != v2CaseKnownType.Package {
			return false
		}
		if v1CaseKnownType.Name != v2CaseKnownType.Name {
			return false
		}
		if len(v1CaseKnownType.Generics) != len(v2CaseKnownType.Generics) {
			panic("unexpected diff len(generics)")
		}
		for i, v1Generic := range v1CaseKnownType.Generics {
			if !variableTypeEq(v1Generic, v2CaseKnownType.Generics[i]) {
				return false
			}
		}
		return true
	}
	if v1CaseTypeArgument != nil && v2CaseTypeArgument != nil {
		return v1CaseTypeArgument.Name == v2CaseTypeArgument.Name
	}
	if v1CaseFunction != nil && v2CaseFunction != nil {
		f1 := *v1CaseFunction
		f2 := *v2CaseFunction
		if len(f1.Arguments) != len(f2.Arguments) {
			return false
		}
		for i, f1Arg := range f1.Arguments {
			if !variableTypeEq(f1Arg.VariableType, f2.Arguments[i].VariableType) {
				return false
			}
		}
		return variableTypeEq(f1.ReturnType, f2.ReturnType)
	}
	if v1CaseOr != nil || v2CaseOr != nil {
		if v1CaseOr != nil && v2CaseOr != nil {
			for _, v1Element := range v1CaseOr.Elements {
				foundEq := false
				for _, v2Element := range v2CaseOr.Elements {
					if variableTypeEq(v1Element, v2Element) {
						foundEq = true
						break
					}
				}
				if !foundEq {
					return false
				}
			}
			for _, v2Element := range v2CaseOr.Elements {
				foundEq := false
				for _, v1Element := range v1CaseOr.Elements {
					if variableTypeEq(v1Element, v2Element) {
						foundEq = true
						break
					}
				}
				if !foundEq {
					return false
				}
			}
			return true
		} else if v1CaseOr != nil {
			for _, element := range v1CaseOr.Elements {
				if !variableTypeEq(element, v2) {
					return false
				}
				return true
			}
		} else {
			return variableTypeEq(v2, v1)
		}
	}
	return reflect.DeepEqual(v1, v2)
}

func printableNameOfTypeAnnotation(typeAnnotation parser.TypeAnnotation) string {
	var result string
	for i, typeAnnotationElement := range typeAnnotation.OrTypes {
		if i > 0 {
			result += " | "
		}
		parser.TypeAnnotationElementExhaustiveSwitch(
			typeAnnotationElement,
			func(typeAnnotation parser.SingleNameType) {
				result = typeAnnotation.TypeName.String
			},
			func(typeAnnotation parser.FunctionType) {
				result = "("
				for i, argument := range typeAnnotation.Arguments {
					if i > 0 {
						result += ", "
					}
					result += printableNameOfTypeAnnotation(argument)
				}
				result = result + ") -> " + printableNameOfTypeAnnotation(typeAnnotation.ReturnType)
			},
		)
	}
	return result
}

func printableName(varType types.VariableType) string {
	if varType == nil {
		return "(nil!)"
	}
	caseTypeArgument, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		return "<" + caseTypeArgument.Name + ">"
	} else if caseKnownType != nil {
		generics := ""
		if len(caseKnownType.Generics) > 0 {
			generics = "<"
			for i, generic := range caseKnownType.Generics {
				if i > 0 {
					generics += ", "
				}
				generics += printableName(generic)
			}
			generics += ">"
		}
		pkg := caseKnownType.Package
		if pkg != "" {
			pkg += "."
		}
		return pkg + caseKnownType.Name + generics
	} else if caseFunction != nil {
		result := "("
		for i, argumentType := range caseFunction.Arguments {
			if i > 0 {
				result = result + ", "
			}
			result = result + printableName(argumentType.VariableType)
		}
		return result + ") -> " + printableName(caseFunction.ReturnType)
	} else if caseOr != nil {
		result := ""
		for i, element := range caseOr.Elements {
			if i > 0 {
				result += " | "
			}
			result += printableName(element)
		}
		return result
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}
