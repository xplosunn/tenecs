package binding

import (
	"encoding/json"
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

type Universe interface {
	impl() *universeImpl
}

type universeImpl struct {
	TypeByTypeName           immutable.Map[string, types.VariableType]
	TypeByVariableName       immutable.Map[string, types.VariableType]
	Constructors             immutable.Map[string, Constructor]
	GlobalInterfaceVariables immutable.Map[string, map[string]types.VariableType]
	GlobalStructVariables    immutable.Map[string, map[string]types.StructVariableType]
}

func (u universeImpl) impl() *universeImpl {
	return &u
}

func PrettyPrint(u Universe, name string) {
	fmt.Printf("%s TypeByTypeName Keys: %v\n", name, mapKeys(u.impl().TypeByVariableName))
	fmt.Printf("%s TypeByVariableName Keys: %v\n", name, mapKeys(u.impl().TypeByVariableName))
	fmt.Printf("%s Constructors Keys: %v\n", name, mapKeys(u.impl().Constructors))
	fmt.Printf("%s GlobalInterfaceVariables Keys: %v\n", name, mapKeys(u.impl().GlobalInterfaceVariables))
	fmt.Printf("%s GlobalStructVariables Keys: %v\n", name, mapKeys(u.impl().GlobalStructVariables))
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

type Constructor struct {
	Arguments  []types.FunctionArgument
	ReturnType types.ConstructableVariableType
}

func NewFromDefaults(defaultTypesWithoutImport map[string]types.VariableType) Universe {
	mapBuilder := immutable.NewMapBuilder[string, types.VariableType](nil)

	for key, value := range defaultTypesWithoutImport {
		mapBuilder.Set(key, value)
	}
	return universeImpl{
		TypeByTypeName:           *mapBuilder.Map(),
		TypeByVariableName:       *immutable.NewMap[string, types.VariableType](nil),
		Constructors:             *immutable.NewMap[string, Constructor](nil),
		GlobalInterfaceVariables: *immutable.NewMap[string, map[string]types.VariableType](nil),
		GlobalStructVariables:    *immutable.NewMap[string, map[string]types.StructVariableType](nil),
	}
}

func NewFromInterfaceVariables(interfaceVariables map[string]types.VariableType, universeToCopyGlobalVariables Universe) Universe {
	mapBuilder := immutable.NewMapBuilder[string, types.VariableType](nil)

	for key, value := range interfaceVariables {
		mapBuilder.Set(key, value)
	}
	return universeImpl{
		TypeByTypeName:           *immutable.NewMap[string, types.VariableType](nil),
		TypeByVariableName:       *mapBuilder.Map(),
		Constructors:             *immutable.NewMap[string, Constructor](nil),
		GlobalInterfaceVariables: universeToCopyGlobalVariables.impl().GlobalInterfaceVariables,
		GlobalStructVariables:    universeToCopyGlobalVariables.impl().GlobalStructVariables,
	}
}

func NewFromStructVariables(interfaceVariables map[string]types.StructVariableType, universeToCopyGlobalVariables Universe) Universe {
	mapBuilder := immutable.NewMapBuilder[string, types.VariableType](nil)

	for key, value := range interfaceVariables {
		mapBuilder.Set(key, types.VariableTypeFromStructVariableType(value))
	}
	return universeImpl{
		TypeByTypeName:           *immutable.NewMap[string, types.VariableType](nil),
		TypeByVariableName:       *mapBuilder.Map(),
		Constructors:             *immutable.NewMap[string, Constructor](nil),
		GlobalInterfaceVariables: universeToCopyGlobalVariables.impl().GlobalInterfaceVariables,
		GlobalStructVariables:    universeToCopyGlobalVariables.impl().GlobalStructVariables,
	}
}

func GetTypeByTypeName(universe Universe, typeName string) (types.VariableType, bool) {
	u := universe.impl()
	return u.TypeByTypeName.Get(typeName)
}

func GetTypeByVariableName(universe Universe, variableName string) (types.VariableType, bool) {
	u := universe.impl()
	return u.TypeByVariableName.Get(variableName)
}

func GetConstructorByName(universe Universe, name string) (Constructor, bool) {
	u := universe.impl()
	return u.Constructors.Get(name)
}

func GetGlobalInterfaceVariables(universe Universe, interf types.Interface) (map[string]types.VariableType, *type_error.TypecheckError) {
	u := universe.impl()
	interfaceRef := interf.Package + "." + interf.Name
	variables, ok := u.GlobalInterfaceVariables.Get(interfaceRef)
	if !ok {
		bytes, err := json.Marshal(u.GlobalInterfaceVariables)
		if err != nil {
			panic(err)
		}
		return nil, type_error.PtrTypeCheckErrorf("not found %s in GlobalInterfaceVariables %s", interfaceRef, string(bytes))
	}
	return variables, nil
}

func GetGlobalStructVariables(universe Universe, struc types.Struct) (map[string]types.StructVariableType, *type_error.TypecheckError) {
	u := universe.impl()
	structRef := struc.Package + "." + struc.Name
	variables, ok := u.GlobalStructVariables.Get(structRef)
	if !ok {
		bytes, err := json.Marshal(u.GlobalStructVariables)
		if err != nil {
			panic(err)
		}
		return nil, type_error.PtrTypeCheckErrorf("not found %s in GlobalStructVariables %s", structRef, string(bytes))
	}
	return variables, nil
}

func CopyAddingType(universe Universe, typeName string, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.TypeByTypeName.Get(typeName)
	if ok {
		bytes, err := json.Marshal(u.TypeByTypeName)
		if err != nil {
			panic(err)
		}
		return nil, type_error.PtrTypeCheckErrorf("type already exists %s in %s", typeName, string(bytes))
	}
	return universeImpl{
		TypeByTypeName:           *u.TypeByTypeName.Set(typeName, varType),
		TypeByVariableName:       u.TypeByVariableName,
		Constructors:             u.Constructors,
		GlobalInterfaceVariables: u.GlobalInterfaceVariables,
		GlobalStructVariables:    u.GlobalStructVariables,
	}, nil
}

func CopyAddingVariable(universe Universe, variableName string, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.TypeByVariableName.Get(variableName)
	if ok {
		bytes, err := json.Marshal(u.TypeByVariableName)
		if err != nil {
			panic(err)
		}
		return nil, type_error.PtrTypeCheckErrorf("variable already exists %s in %s", variableName, string(bytes))
	}
	return universeImpl{
		TypeByTypeName:           u.TypeByTypeName,
		TypeByVariableName:       *u.TypeByVariableName.Set(variableName, varType),
		Constructors:             u.Constructors,
		GlobalInterfaceVariables: u.GlobalInterfaceVariables,
		GlobalStructVariables:    u.GlobalStructVariables,
	}, nil
}

func CopyAddingGlobalInterfaceVariables(universe Universe, interf types.Interface, variables map[string]types.VariableType) (Universe, *type_error.TypecheckError) {
	interfaceRef := interf.Package + "." + interf.Name
	return CopyAddingGlobalInterfaceRefVariables(universe, interfaceRef, variables)
}

func CopyAddingGlobalInterfaceRefVariables(universe Universe, interfaceRef string, variables map[string]types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.GlobalInterfaceVariables.Get(interfaceRef)
	if ok {
		bytes, err := json.Marshal(u.GlobalInterfaceVariables)
		if err != nil {
			panic(err)
		}
		return nil, type_error.PtrTypeCheckErrorf("variable already exists %s in %s", interfaceRef, string(bytes))
	}
	return universeImpl{
		TypeByTypeName:           u.TypeByTypeName,
		TypeByVariableName:       u.TypeByVariableName,
		Constructors:             u.Constructors,
		GlobalInterfaceVariables: *u.GlobalInterfaceVariables.Set(interfaceRef, variables),
		GlobalStructVariables:    u.GlobalStructVariables,
	}, nil
}

func CopyAddingGlobalStructVariables(universe Universe, struc types.Struct, variables map[string]types.StructVariableType) (Universe, *type_error.TypecheckError) {
	structRef := struc.Package + "." + struc.Name
	return CopyAddingGlobalStructRefVariables(universe, structRef, variables)
}

func CopyAddingGlobalStructRefVariables(universe Universe, structRef string, variables map[string]types.StructVariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.GlobalStructVariables.Get(structRef)
	if ok {
		bytes, err := json.Marshal(u.GlobalStructVariables)
		if err != nil {
			panic(err)
		}
		return nil, type_error.PtrTypeCheckErrorf("variable already exists %s in %s", structRef, string(bytes))
	}
	return universeImpl{
		TypeByTypeName:           u.TypeByTypeName,
		TypeByVariableName:       u.TypeByVariableName,
		Constructors:             u.Constructors,
		GlobalInterfaceVariables: u.GlobalInterfaceVariables,
		GlobalStructVariables:    *u.GlobalStructVariables.Set(structRef, variables),
	}, nil
}

func CopyAddingFunctionArguments(universe Universe, functionArguments []types.FunctionArgument) (Universe, *type_error.TypecheckError) {
	result := universe
	for _, argument := range functionArguments {
		updatedResult, err := CopyAddingVariable(result, argument.Name, argument.VariableType)
		if err != nil {
			return result, err
		}
		result = updatedResult
	}
	return result, nil
}

func CopyAddingConstructor(universe Universe, typeName string, constructor Constructor) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.Constructors.Get(typeName)
	if ok {
		bytes, err := json.Marshal(u.Constructors)
		if err != nil {
			panic(err)
		}
		return nil, type_error.PtrTypeCheckErrorf("constructor already exists %s in %s", typeName, string(bytes))
	}
	return universeImpl{
		TypeByTypeName:           u.TypeByTypeName,
		TypeByVariableName:       u.TypeByVariableName,
		Constructors:             *u.Constructors.Set(typeName, constructor),
		GlobalInterfaceVariables: u.GlobalInterfaceVariables,
		GlobalStructVariables:    u.GlobalStructVariables,
	}, nil
}
