package typer

import (
	"encoding/json"
	"github.com/benbjohnson/immutable"
)

type Universe struct {
	TypeByTypeName           immutable.Map[string, VariableType]
	TypeByVariableName       immutable.Map[string, VariableType]
	Constructors             immutable.Map[string, Constructor]
	GlobalInterfaceVariables immutable.Map[string, map[string]VariableType]
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
		TypeByTypeName:           *mapBuilder.Map(),
		TypeByVariableName:       *immutable.NewMap[string, VariableType](nil),
		Constructors:             *immutable.NewMap[string, Constructor](nil),
		GlobalInterfaceVariables: *immutable.NewMap[string, map[string]VariableType](nil),
	}
}

func NewUniverseFromInterfaceVariables(interfaceVariables map[string]VariableType, globalInterfaceVariables immutable.Map[string, map[string]VariableType]) Universe {
	mapBuilder := immutable.NewMapBuilder[string, VariableType](nil)

	for key, value := range interfaceVariables {
		mapBuilder.Set(key, value)
	}
	return Universe{
		TypeByTypeName:           *immutable.NewMap[string, VariableType](nil),
		TypeByVariableName:       *mapBuilder.Map(),
		Constructors:             *immutable.NewMap[string, Constructor](nil),
		GlobalInterfaceVariables: globalInterfaceVariables,
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
		TypeByTypeName:           *universe.TypeByTypeName.Set(typeName, varType),
		TypeByVariableName:       universe.TypeByVariableName,
		Constructors:             universe.Constructors,
		GlobalInterfaceVariables: universe.GlobalInterfaceVariables,
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
		TypeByTypeName:           universe.TypeByTypeName,
		TypeByVariableName:       *universe.TypeByVariableName.Set(variableName, varType),
		Constructors:             universe.Constructors,
		GlobalInterfaceVariables: universe.GlobalInterfaceVariables,
	}, nil
}

func copyUniverseAddingGlobalInterfaceVariables(universe Universe, interf Interface, variables map[string]VariableType) (Universe, *TypecheckError) {
	interfaceRef := interf.Package + "." + interf.Name
	return copyUniverseAddingGlobalInterfaceRefVariables(universe, interfaceRef, variables)
}

func copyUniverseAddingGlobalInterfaceRefVariables(universe Universe, interfaceRef string, variables map[string]VariableType) (Universe, *TypecheckError) {
	_, ok := universe.GlobalInterfaceVariables.Get(interfaceRef)
	if ok {
		bytes, err := json.Marshal(universe.TypeByVariableName)
		if err != nil {
			panic(err)
		}
		return universe, PtrTypeCheckErrorf("variable already exists %s in %s", interfaceRef, string(bytes))
	}
	return Universe{
		TypeByTypeName:           universe.TypeByTypeName,
		TypeByVariableName:       universe.TypeByVariableName,
		Constructors:             universe.Constructors,
		GlobalInterfaceVariables: *universe.GlobalInterfaceVariables.Set(interfaceRef, variables),
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
		TypeByTypeName:           universe.TypeByTypeName,
		TypeByVariableName:       universe.TypeByVariableName,
		Constructors:             *universe.Constructors.Set(moduleName, constructor),
		GlobalInterfaceVariables: universe.GlobalInterfaceVariables,
	}, nil
}

func GetGlobalInterfaceVariables(universe Universe, interf Interface) (map[string]VariableType, *TypecheckError) {
	interfaceRef := interf.Package + "." + interf.Name
	variables, ok := universe.GlobalInterfaceVariables.Get(interfaceRef)
	if !ok {
		bytes, err := json.Marshal(universe.GlobalInterfaceVariables)
		if err != nil {
			panic(err)
		}
		return nil, PtrTypeCheckErrorf("not found %s in GlobalInterfaceVariables %s", interfaceRef, string(bytes))
	}
	return variables, nil
}
