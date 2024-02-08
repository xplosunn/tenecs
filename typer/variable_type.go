package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/types"
	"strings"
)

func PrintableNameWithoutPackage(varType types.VariableType) string {
	name := printableName(varType)
	split := strings.Split(name, ".")
	return split[len(split)-1]
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
