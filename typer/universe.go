package typer

import (
	"encoding/json"
	"github.com/benbjohnson/immutable"
)

type Universe struct {
	TypeByTypeName     immutable.Map[string, VariableType]
	TypeByVariableName immutable.Map[string, VariableType]
}

func NewUniverseFromDefaults() Universe {
	mapBuilder := immutable.NewMapBuilder[string, VariableType](nil)

	for key, value := range DefaultTypesAvailableWithoutImport {
		mapBuilder.Set(key, value)
	}
	return Universe{
		TypeByTypeName:     *mapBuilder.Map(),
		TypeByVariableName: *immutable.NewMap[string, VariableType](nil),
	}
}

func NewUniverseFromInterface(interf Interface) Universe {
	mapBuilder := immutable.NewMapBuilder[string, VariableType](nil)

	for key, value := range interf.Variables {
		mapBuilder.Set(key, value)
	}
	return Universe{
		TypeByTypeName:     *immutable.NewMap[string, VariableType](nil),
		TypeByVariableName: *mapBuilder.Map(),
	}
}

func copyUniverseAddingType(universe Universe, typeName string, varType VariableType) (Universe, *TypecheckError) {
	_, ok := universe.TypeByTypeName.Get(typeName)
	if ok {
		bytes, err := json.Marshal(universe.TypeByTypeName)
		if err != nil {
			panic(err)
		}
		return universe, PtrTypeCheckErrorf("type already exists %s in %s", typeName, string(bytes))
	}
	return Universe{
		TypeByTypeName:     *universe.TypeByTypeName.Set(typeName, varType),
		TypeByVariableName: universe.TypeByVariableName,
	}, nil
}

func copyUniverseAddingVariable(universe Universe, variableName string, varType VariableType) (Universe, *TypecheckError) {
	_, ok := universe.TypeByVariableName.Get(variableName)
	if ok {
		bytes, err := json.Marshal(universe.TypeByVariableName)
		if err != nil {
			panic(err)
		}
		return universe, PtrTypeCheckErrorf("variable already exists %s in %s", variableName, string(bytes))
	}
	return Universe{
		TypeByTypeName:     universe.TypeByTypeName,
		TypeByVariableName: *universe.TypeByVariableName.Set(variableName, varType),
	}, nil
}

func copyUniverseAddingVariables(universe Universe, variables map[string]VariableType) (Universe, *TypecheckError) {
	result := universe
	for name, varType := range variables {
		updatedResult, err := copyUniverseAddingVariable(result, name, varType)
		if err != nil {
			return result, err
		}
		result = updatedResult
	}
	return result, nil
}

func copyUniverseAddingFunctionArguments(universe Universe, functionArguments []FunctionArgument) (Universe, *TypecheckError) {
	result := universe
	for _, argument := range functionArguments {
		updatedResult, err := copyUniverseAddingVariable(result, argument.Name, argument.VariableType)
		if err != nil {
			return result, err
		}
		result = updatedResult
	}
	return result, nil
}
