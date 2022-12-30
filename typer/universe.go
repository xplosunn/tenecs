package typer

import (
	"encoding/json"
	"github.com/benbjohnson/immutable"
)

type Universe struct {
	TypeByTypeName     immutable.Map[string, VariableType]
	TypeByVariableName immutable.Map[string, VariableType]
	Constructors       immutable.Map[string, Constructor]
}

type Constructor struct {
	Arguments  []FunctionArgument
	ReturnType Interface
}

func NewUniverseFromDefaults() Universe {
	mapBuilder := immutable.NewMapBuilder[string, VariableType](nil)

	for key, value := range DefaultTypesAvailableWithoutImport {
		mapBuilder.Set(key, value)
	}
	return Universe{
		TypeByTypeName:     *mapBuilder.Map(),
		TypeByVariableName: *immutable.NewMap[string, VariableType](nil),
		Constructors:       *immutable.NewMap[string, Constructor](nil),
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
		Constructors:       *immutable.NewMap[string, Constructor](nil),
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
		Constructors:       universe.Constructors,
	}, nil
}

func copyUniverseOverridingType(universe Universe, typeName string, varType VariableType) (Universe, *TypecheckError) {
	_, ok := universe.TypeByTypeName.Get(typeName)
	if !ok {
		bytes, err := json.Marshal(universe.TypeByTypeName)
		if err != nil {
			panic(err)
		}
		return universe, PtrTypeCheckErrorf("type %s not found in %s", typeName, string(bytes))
	}
	return Universe{
		TypeByTypeName:     *universe.TypeByTypeName.Set(typeName, varType),
		TypeByVariableName: universe.TypeByVariableName,
		Constructors:       universe.Constructors,
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
		Constructors:       universe.Constructors,
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

func copyUniverseAddingConstructor(universe Universe, moduleName string, constructor Constructor) (Universe, *TypecheckError) {
	_, ok := universe.Constructors.Get(moduleName)
	if ok {
		bytes, err := json.Marshal(universe.Constructors)
		if err != nil {
			panic(err)
		}
		return universe, PtrTypeCheckErrorf("constructor already exists %s in %s", moduleName, string(bytes))
	}
	return Universe{
		TypeByTypeName:     universe.TypeByTypeName,
		TypeByVariableName: universe.TypeByVariableName,
		Constructors:       *universe.Constructors.Set(moduleName, constructor),
	}, nil
}
