package binding

import (
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/fsamin/go-dump"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

type ResolutionError struct {
	VariableType *types.VariableType
	Problem      string
}

func ResolutionErrorCouldNotResolve(typeName string) *ResolutionError {
	return &ResolutionError{
		VariableType: nil,
		Problem:      "not found type: " + typeName,
	}
}

func ResolutionErrorWrongNumberOfGenerics(variableType types.VariableType, expected int, got int) *ResolutionError {
	return &ResolutionError{
		VariableType: &variableType,
		Problem:      fmt.Sprintf("wrong number of generics, expected %d but got %d", expected, got),
	}
}

func ResolutionErrorNotAValidGeneric(variableType types.VariableType) *ResolutionError {
	return &ResolutionError{
		VariableType: &variableType,
		Problem:      "not a valid generic",
	}
}

func (err *ResolutionError) Error() string {
	return err.Problem
}

type Universe interface {
	impl() *universeImpl
}

type universeImpl struct {
	TypeByTypeName     immutable.Map[string, types.VariableType]
	TypeByVariableName immutable.Map[string, types.VariableType]
}

func (u universeImpl) impl() *universeImpl {
	return &u
}

func PrettyPrint(u Universe, name string) {
	fmt.Printf("%s TypeByTypeName Keys: %v\n", name, mapKeys(u.impl().TypeByVariableName))
	fmt.Printf("%s TypeByVariableName Keys: %v\n", name, mapKeys(u.impl().TypeByVariableName))
	fmt.Printf("%s dump:\n", name)
	dump.Dump(u)
}

func mapKeys[V any](m immutable.Map[string, V]) []string {
	result := []string{}
	iterator := m.Iterator()
	for !iterator.Done() {
		key, _, _ := iterator.Next()
		result = append(result, key)
	}
	return result
}

func NewFromDefaults(defaultTypesWithoutImport map[string]types.VariableType) Universe {
	mapBuilder := immutable.NewMapBuilder[string, types.VariableType](nil)

	for key, value := range defaultTypesWithoutImport {
		mapBuilder.Set(key, value)
	}
	return universeImpl{
		TypeByTypeName:     *mapBuilder.Map(),
		TypeByVariableName: *immutable.NewMap[string, types.VariableType](nil),
	}
}

func NewFromInterfaceVariables(node parser.Node, interfaceVariables map[string]types.VariableType, universeToCopy Universe) (Universe, *type_error.TypecheckError) {
	universe := universeToCopy
	var err *type_error.TypecheckError
	for key, value := range interfaceVariables {
		universe, err = CopyAddingVariable(universe, parser.Name{
			Node:   node,
			String: key,
		}, value)
		if err != nil {
			return nil, err
		}
	}
	return universe, nil
}

func GetTypeByTypeName(universe Universe, typeName string, generics []types.StructFieldVariableType) (types.VariableType, *ResolutionError) {
	u := universe.impl()
	varType, ok := u.TypeByTypeName.Get(typeName)
	if !ok {
		return nil, ResolutionErrorCouldNotResolve(typeName)
	}

	return applyGenerics(varType, generics)
}

func applyGenerics(varType types.VariableType, generics []types.StructFieldVariableType) (types.VariableType, *ResolutionError) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid, caseArray, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		if len(generics) != 0 {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, 0, len(generics))
		}
		return varType, nil
	} else if caseStruct != nil {
		if len(generics) != len(caseStruct.Generics) {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, len(caseStruct.Generics), len(generics))
		}

		newFields := map[string]types.StructFieldVariableType{}
		for fieldName, fieldType := range caseStruct.Fields {
			newFields[fieldName] = fieldType
		}
		for i, genericName := range caseStruct.Generics {
			genericType := generics[i]
			for fieldName, fieldVariableType := range newFields {
				resolved, err := resolveGeneric(fieldVariableType, genericName, genericType)
				if err != nil {
					return nil, err
				}
				newFields[fieldName] = resolved
			}
		}
		return &types.Struct{
			Package:  caseStruct.Package,
			Name:     caseStruct.Name,
			Generics: caseStruct.Generics,
			Fields:   newFields,
		}, nil
	} else if caseInterface != nil {
		if len(generics) != 0 {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, 0, len(generics))
		}
		return varType, nil
	} else if caseFunction != nil {
		panic("TODO applyGenerics caseFunction")
	} else if caseBasicType != nil {
		if len(generics) != 0 {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, 0, len(generics))
		}
		return varType, nil
	} else if caseVoid != nil {
		if len(generics) != 0 {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, 0, len(generics))
		}
		return varType, nil
	} else if caseArray != nil {
		if len(generics) != 1 {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, 1, len(generics))
		}
		return &types.Array{
			OfType: generics[0],
		}, nil
	} else if caseOr != nil {
		panic("unexpected applyGenerics caseOr")
	} else {
		panic(fmt.Errorf("code on %v", varType))
	}
}

func resolveGeneric(over types.StructFieldVariableType, genericName string, resolveWith types.StructFieldVariableType) (types.StructFieldVariableType, *ResolutionError) {
	caseTypeArgument, caseStruct, caseBasicType, caseVoid, caseArray, caseOr := over.StructFieldVariableTypeCases()
	if caseTypeArgument != nil {
		if caseTypeArgument.Name == genericName {
			return resolveWith, nil
		}
		return caseTypeArgument, nil
	} else if caseStruct != nil {
		newStruct := &types.Struct{
			Package:  caseStruct.Package,
			Name:     caseStruct.Name,
			Generics: caseStruct.Generics,
			Fields:   caseStruct.Fields,
		}
		for fieldName, variableType := range caseStruct.Fields {
			newFieldType, err := resolveGeneric(variableType, genericName, resolveWith)
			if err != nil {
				return nil, err
			}
			newStruct.Fields[fieldName] = newFieldType.(types.StructFieldVariableType)
		}
		return newStruct, nil
	} else if caseBasicType != nil {
		return caseBasicType, nil
	} else if caseVoid != nil {
		return caseVoid, nil
	} else if caseArray != nil {
		newOfType, err := resolveGeneric(caseArray.OfType, genericName, resolveWith)
		if err != nil {
			return nil, err
		}
		return &types.Array{
			OfType: newOfType.(types.StructFieldVariableType),
		}, nil
		panic("todo resolveGeneric caseArray")
	} else if caseOr != nil {
		panic("todo resolveGeneric caseOr")
	} else {
		panic(fmt.Errorf("cases on %v", over))
	}
}

func GetTypeByVariableName(universe Universe, variableName string) (types.VariableType, bool) {
	u := universe.impl()
	return u.TypeByVariableName.Get(variableName)
}

func CopyAddingType(universe Universe, typeName parser.Name, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.TypeByTypeName.Get(typeName.String)
	if ok {
		return nil, type_error.PtrOnNodef(typeName.Node, "type already exists %s", typeName)
	}
	return universeImpl{
		TypeByTypeName:     *u.TypeByTypeName.Set(typeName.String, varType),
		TypeByVariableName: u.TypeByVariableName,
	}, nil
}

func CopyOverridingType(universe Universe, typeName string, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.TypeByTypeName.Get(typeName)
	if !ok {
		panic(fmt.Sprintf("cannot override %s in universe", typeName))
	}
	return universeImpl{
		TypeByTypeName:     *u.TypeByTypeName.Set(typeName, varType),
		TypeByVariableName: u.TypeByVariableName,
	}, nil
}

func CopyAddingVariable(universe Universe, variableName parser.Name, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.TypeByVariableName.Get(variableName.String)
	if ok {
		return nil, type_error.PtrOnNodef(variableName.Node, "duplicate variable '%s'", variableName.String)
	}
	return universeImpl{
		TypeByTypeName:     u.TypeByTypeName,
		TypeByVariableName: *u.TypeByVariableName.Set(variableName.String, varType),
	}, nil
}

func CopyOverridingVariableType(universe Universe, variableName string, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.TypeByVariableName.Get(variableName)
	if !ok {
		panic(fmt.Sprintf("cannot override %s in universe", variableName))
	}
	return universeImpl{
		TypeByTypeName:     u.TypeByTypeName,
		TypeByVariableName: *u.TypeByVariableName.Set(variableName, varType),
	}, nil
}

func CopyAddingFunctionArguments(universe Universe, functionArgumentNames []parser.Name, functionArgumentVariableTypes []types.VariableType) (Universe, *type_error.TypecheckError) {
	result := universe
	if len(functionArgumentNames) != len(functionArgumentVariableTypes) {
		panic("programatic err on CopyAddingFunctionArguments: len(functionArgumentNames) != len(functionArgumentVariableTypes)")
	}
	for i, name := range functionArgumentNames {
		updatedResult, err := CopyAddingVariable(result, name, functionArgumentVariableTypes[i])
		if err != nil {
			return result, err
		}
		result = updatedResult
	}
	return result, nil
}
