package binding

import (
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/fsamin/go-dump"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

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

func NewFromStructVariables(node parser.Node, interfaceVariables map[string]types.StructFieldVariableType, universeToCopy Universe) (Universe, *type_error.TypecheckError) {
	universe := universeToCopy
	var err *type_error.TypecheckError
	for key, value := range interfaceVariables {
		universe, err = CopyAddingVariable(universe, parser.Name{
			Node:   node,
			String: key,
		}, types.VariableTypeFromStructFieldVariableType(value))
		if err != nil {
			return nil, err
		}
	}
	return universe, nil
}

func GetTypeByTypeName(universe Universe, typeName string) (types.VariableType, bool) {
	u := universe.impl()
	return u.TypeByTypeName.Get(typeName)
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
