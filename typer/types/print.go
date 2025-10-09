package types

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
)

func PrintableNameWithoutPackage(varType VariableType) string {
	name := PrintableName(varType)
	split := strings.Split(name, ".")
	return split[len(split)-1]
}

func PrintableNameSorted(varType VariableType) string {
	return printableName(varType, true)
}

func PrintableName(varType VariableType) string {
	return printableName(varType, false)
}

func PrintableNameKnownTypeWithoutGenerics(caseKnownType *KnownType) string {
	pkg := caseKnownType.Package
	if pkg != "" {
		pkg += "."
	}
	return pkg + caseKnownType.Name
}

func printableName(varType VariableType, sort bool) string {
	if varType == nil {
		panic("PrintableName nil")
	}
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := FlattenOr(varType).VariableTypeCases()
	if caseTypeArgument != nil {
		return caseTypeArgument.Name
	} else if caseList != nil {
		return "List<" + printableName(caseList.Generic, sort) + ">"
	} else if caseKnownType != nil {
		generics := ""
		if len(caseKnownType.Generics) > 0 {
			generics = "<"
			for i, generic := range caseKnownType.Generics {
				if i > 0 {
					generics += ", "
				}
				generics += printableName(generic, sort)
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
			result = result + printableName(argumentType.VariableType, sort)
		}
		return result + ") ~> " + printableName(caseFunction.ReturnType, sort)
	} else if caseOr != nil {
		elements := []string{}
		for _, element := range caseOr.Elements {
			elements = append(elements, printableName(element, sort))
		}
		if sort {
			slices.Sort(elements)
		}
		return strings.Join(elements, " | ")
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}
