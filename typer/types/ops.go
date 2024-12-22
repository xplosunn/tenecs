package types

import (
	"fmt"
	"reflect"
)

func VariableTypeContainedIn(sub VariableType, super VariableType) bool {
	superOr, ok := super.(*OrVariableType)
	if !ok {
		return VariableTypeEq(sub, super)
	}
	subOr, ok := sub.(*OrVariableType)
	if !ok {
		for _, superElement := range superOr.Elements {
			if VariableTypeEq(sub, superElement) {
				return true
			}
		}
		return false
	}
	for _, subElement := range subOr.Elements {
		if !VariableTypeContainedIn(subElement, super) {
			return false
		}
	}
	return true
}

func VariableTypeAddToOr(varType VariableType, or *OrVariableType) {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	_ = caseTypeArgument
	_ = caseList
	_ = caseKnownType
	_ = caseFunction
	if caseOr != nil {
		for _, element := range caseOr.Elements {
			VariableTypeAddToOr(element, or)
		}
	} else {
		if !VariableTypeContainedIn(varType, or) {
			or.Elements = append(or.Elements, varType)
		}
	}
}

func VariableTypeCombine(v1 VariableType, v2 VariableType) VariableType {
	result := &OrVariableType{Elements: []VariableType{}}

	addAll := func(varType VariableType) {
		VariableTypeAddToOr(varType, result)
	}

	addAll(v1)
	addAll(v2)

	if len(result.Elements) == 1 {
		return result.Elements[0]
	}

	return result
}

func VariableTypeEq(v1 VariableType, v2 VariableType) bool {
	if v1 == nil || v2 == nil {
		panic(fmt.Errorf("trying to eq %v to %v", v1, v2))
	}
	v1CaseTypeArgument, v1CaseList, v1CaseKnownType, v1CaseFunction, v1CaseOr := v1.VariableTypeCases()
	v2CaseTypeArgument, v2CaseList, v2CaseKnownType, v2CaseFunction, v2CaseOr := v2.VariableTypeCases()
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
			if !VariableTypeEq(v1Generic, v2CaseKnownType.Generics[i]) {
				return false
			}
		}
		return true
	}
	if v1CaseTypeArgument != nil && v2CaseTypeArgument != nil {
		return v1CaseTypeArgument.Name == v2CaseTypeArgument.Name
	}
	if v1CaseList != nil && v2CaseList != nil {
		return VariableTypeEq(v1CaseList.Generic, v2CaseList.Generic)
	}
	if v1CaseFunction != nil && v2CaseFunction != nil {
		f1 := normalizeGenericsFunction(*v1CaseFunction, "", map[string]string{})
		f2 := normalizeGenericsFunction(*v2CaseFunction, "", map[string]string{})
		if len(f1.Arguments) != len(f2.Arguments) {
			return false
		}
		for i, f1Arg := range f1.Arguments {
			if !VariableTypeEq(f1Arg.VariableType, f2.Arguments[i].VariableType) {
				return false
			}
		}
		return VariableTypeEq(f1.ReturnType, f2.ReturnType)
	}
	if v1CaseOr != nil || v2CaseOr != nil {
		if v1CaseOr != nil && v2CaseOr != nil {
			for _, v1Element := range v1CaseOr.Elements {
				foundEq := false
				for _, v2Element := range v2CaseOr.Elements {
					if VariableTypeEq(v1Element, v2Element) {
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
					if VariableTypeEq(v1Element, v2Element) {
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
				if !VariableTypeEq(element, v2) {
					return false
				}
			}
			return true
		} else {
			return VariableTypeEq(v2, v1)
		}
	}
	return reflect.DeepEqual(v1, v2)
}

func normalizeGenericsVariableType(varType VariableType, prefix string, normalizationMap map[string]string) VariableType {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		return normalizeGenericsTypeArgument(*caseTypeArgument, prefix, normalizationMap)
	} else if caseList != nil {
		return caseList
	} else if caseKnownType != nil {
		return caseKnownType
	} else if caseFunction != nil {
		return normalizeGenericsFunction(*caseFunction, prefix, normalizationMap)
	} else if caseOr != nil {
		return normalizeGenericsOrVariableType(*caseOr, prefix, normalizationMap)
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}

func normalizeGenericsTypeArgument(typeArg TypeArgument, prefix string, normalizationMap map[string]string) *TypeArgument {
	newValue, ok := normalizationMap[typeArg.Name]
	if ok {
		return &TypeArgument{Name: newValue}
	}
	return &typeArg
}

func normalizeGenericsFunction(function Function, prefix string, previousNormalizationMap map[string]string) *Function {
	normalizationMap := map[string]string{}
	for k, v := range previousNormalizationMap {
		normalizationMap[k] = v
	}

	normalized := Function{
		Generics:   []string{},
		Arguments:  []FunctionArgument{},
		ReturnType: nil,
	}

	for i, generic := range function.Generics {
		newTypeName := fmt.Sprintf("%s%d", prefix, i)
		normalizationMap[generic] = newTypeName
		normalized.Generics = append(normalized.Generics, newTypeName)
	}
	for i, functionArgument := range normalized.Arguments {
		normalized.Arguments = append(normalized.Arguments, FunctionArgument{
			Name:         functionArgument.Name,
			VariableType: normalizeGenericsVariableType(functionArgument.VariableType, fmt.Sprintf("%s_%d", prefix, i), normalizationMap),
		})
	}
	normalized.ReturnType = normalizeGenericsVariableType(function.ReturnType, fmt.Sprintf("%s_R", prefix), normalizationMap)

	return &normalized
}

func normalizeGenericsOrVariableType(or OrVariableType, prefix string, normalizationMap map[string]string) *OrVariableType {
	result := OrVariableType{Elements: []VariableType{}}
	for _, element := range or.Elements {
		VariableTypeAddToOr(normalizeGenericsVariableType(element, prefix, normalizationMap), &result)
	}
	return &result
}
