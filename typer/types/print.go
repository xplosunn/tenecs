package types

import (
	"fmt"
	"strings"
)

func PrintableNameWithoutPackage(varType VariableType) string {
	name := PrintableName(varType)
	split := strings.Split(name, ".")
	return split[len(split)-1]
}

func PrintableName(varType VariableType) string {
	if varType == nil {
		panic("PrintableName nil")
	}
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		//TODO FIXME drop the "<" and ">"
		return "<" + caseTypeArgument.Name + ">"
	} else if caseList != nil {
		return "List<" + PrintableName(caseList.Generic) + ">"
	} else if caseKnownType != nil {
		generics := ""
		if len(caseKnownType.Generics) > 0 {
			generics = "<"
			for i, generic := range caseKnownType.Generics {
				if i > 0 {
					generics += ", "
				}
				generics += PrintableName(generic)
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
			result = result + PrintableName(argumentType.VariableType)
		}
		return result + ") ~> " + PrintableName(caseFunction.ReturnType)
	} else if caseOr != nil {
		result := ""
		for i, element := range caseOr.Elements {
			if i > 0 {
				result += " | "
			}
			result += PrintableName(element)
		}
		return result
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}
