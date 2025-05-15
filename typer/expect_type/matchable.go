package expect_type

import (
	"errors"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/types"
)

type MatchableVariableType interface {
	sealedMatchableVariableType()
	MatchableVariableTypeCases() (*MatchableList, *MatchableKnownType, *OrMatchableVariableType)
}

type MatchableList struct {
	Generic MatchableVariableType
}

func (l *MatchableList) sealedMatchableVariableType() {}
func (l *MatchableList) MatchableVariableTypeCases() (*MatchableList, *MatchableKnownType, *OrMatchableVariableType) {
	return l, nil, nil
}

type MatchableKnownType struct {
	Package                  string
	Name                     string
	DeclaredGenerics         []string
	Generics                 []MatchableVariableType
	GenericsMatchableByField []string
}

func (k *MatchableKnownType) sealedMatchableVariableType() {}
func (k *MatchableKnownType) MatchableVariableTypeCases() (*MatchableList, *MatchableKnownType, *OrMatchableVariableType) {
	return nil, k, nil
}

type OrMatchableVariableType struct {
	Elements []MatchableVariableType
}

func (o *OrMatchableVariableType) sealedMatchableVariableType() {}
func (o *OrMatchableVariableType) MatchableVariableTypeCases() (*MatchableList, *MatchableKnownType, *OrMatchableVariableType) {
	return nil, nil, o
}

func KnownTypeGenericsMatchByField(caseKnownType types.KnownType, resolveStructFields map[binding.Ref]map[string]types.VariableType) ([]string, error) {
	genericsMatchableByField := []string{}
	structFields := resolveStructFields[binding.Ref{
		Package: caseKnownType.Package,
		Name:    caseKnownType.Name,
	}]
	for _, declaredGenericName := range caseKnownType.DeclaredGenerics {
		if structFields == nil {
			panic("no fields")
		}
		var foundMatchingField *string
		for fieldName, structFieldVarType := range structFields {
			if types.VariableTypeEq(structFieldVarType, &types.TypeArgument{Name: declaredGenericName}) {
				foundMatchingField = &fieldName
				break
			}
		}
		if foundMatchingField == nil {
			return nil, errors.New("matching on a struct with generics requires the struct to have one field of that type")
		}
		genericsMatchableByField = append(genericsMatchableByField, *foundMatchingField)
	}
	return genericsMatchableByField, nil
}

func AsMatchable(varType types.VariableType, resolveStructFields map[binding.Ref]map[string]types.VariableType) (MatchableVariableType, error) {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		return nil, errors.New("can't match on generic")
	} else if caseList != nil {
		panic("TODO AsMatchable caseList")
	} else if caseKnownType != nil {
		generics := []MatchableVariableType{}
		for _, generic := range caseKnownType.Generics {
			matchable, err := AsMatchable(generic, resolveStructFields)
			if err != nil {
				return nil, err
			}
			_, _, caseMatchableOr := matchable.MatchableVariableTypeCases()
			if caseMatchableOr != nil {
				return nil, errors.New("can't match on or in generic position")
			}
			generics = append(generics, matchable)
		}
		genericsMatchableByField, err := KnownTypeGenericsMatchByField(*caseKnownType, resolveStructFields)
		if err != nil {
			return nil, err
		}
		if caseKnownType.Generics == nil {
			generics = nil
		}
		return &MatchableKnownType{
			Package:                  caseKnownType.Package,
			Name:                     caseKnownType.Name,
			DeclaredGenerics:         caseKnownType.DeclaredGenerics,
			Generics:                 generics,
			GenericsMatchableByField: genericsMatchableByField,
		}, nil
	} else if caseFunction != nil {
		return nil, errors.New("can't match on function")
	} else if caseOr != nil {
		result := &OrMatchableVariableType{
			Elements: []MatchableVariableType{},
		}
		for _, element := range caseOr.Elements {
			toAdd, err := AsMatchable(element, resolveStructFields)
			if err != nil {
				return nil, err
			}
			result.Elements = append(result.Elements, toAdd)
		}
		return result, nil
	} else {
		panic("cases on VariableTypeCases")
	}
}
