package binding

import (
	"encoding/json"
	"github.com/benbjohnson/immutable"
	"github.com/segmentio/ksuid"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
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
	ParserFunctionByUniqueId immutable.Map[string, parser.Lambda]
}

func (u universeImpl) impl() *universeImpl {
	return &u
}

type Constructor struct {
	Arguments  []types.FunctionArgument
	ReturnType types.Interface
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
		ParserFunctionByUniqueId: *immutable.NewMap[string, parser.Lambda](nil),
	}
}

func NewFromInterfaceVariables(interfaceVariables map[string]types.VariableType, universeToCopyGlobalInterfaceVariables Universe) Universe {
	mapBuilder := immutable.NewMapBuilder[string, types.VariableType](nil)

	for key, value := range interfaceVariables {
		mapBuilder.Set(key, value)
	}
	return universeImpl{
		TypeByTypeName:           *immutable.NewMap[string, types.VariableType](nil),
		TypeByVariableName:       *mapBuilder.Map(),
		Constructors:             *immutable.NewMap[string, Constructor](nil),
		GlobalInterfaceVariables: universeToCopyGlobalInterfaceVariables.impl().GlobalInterfaceVariables,
		ParserFunctionByUniqueId: *immutable.NewMap[string, parser.Lambda](nil),
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

func GetParserFunctionByUniqueId(universe Universe, id string) (parser.Lambda, *type_error.TypecheckError) {
	u := universe.impl()
	lambda, ok := u.ParserFunctionByUniqueId.Get(id)
	if !ok {
		bytes, err := json.Marshal(u.GlobalInterfaceVariables)
		if err != nil {
			panic(err)
		}
		return parser.Lambda{}, type_error.PtrTypeCheckErrorf("not found %s in ParserFunctionByUniqueId %s", id, string(bytes))
	}
	return lambda, nil
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
		ParserFunctionByUniqueId: u.ParserFunctionByUniqueId,
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
		ParserFunctionByUniqueId: u.ParserFunctionByUniqueId,
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
		bytes, err := json.Marshal(u.TypeByVariableName)
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
		ParserFunctionByUniqueId: u.ParserFunctionByUniqueId,
	}, nil
}

func CopyAddingVariables(universe Universe, variables map[string]ast.Expression) (Universe, *type_error.TypecheckError) {
	result := universe
	for name, programExp := range variables {
		varType := ast.VariableTypeOfExpression(programExp)
		updatedResult, err := CopyAddingVariable(result, name, varType)
		if err != nil {
			return result, err
		}
		result = updatedResult
	}
	return result, nil
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

func CopyAddingConstructor(universe Universe, moduleName string, constructor Constructor) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.Constructors.Get(moduleName)
	if ok {
		bytes, err := json.Marshal(u.Constructors)
		if err != nil {
			panic(err)
		}
		return nil, type_error.PtrTypeCheckErrorf("constructor already exists %s in %s", moduleName, string(bytes))
	}
	return universeImpl{
		TypeByTypeName:           u.TypeByTypeName,
		TypeByVariableName:       u.TypeByVariableName,
		Constructors:             *u.Constructors.Set(moduleName, constructor),
		GlobalInterfaceVariables: u.GlobalInterfaceVariables,
		ParserFunctionByUniqueId: u.ParserFunctionByUniqueId,
	}, nil
}

func CopyAddingParserFunctionByUniqueId(universe Universe, uniqueId string, parserFunction parser.Lambda) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.ParserFunctionByUniqueId.Get(uniqueId)
	if ok {
		bytes, err := json.Marshal(u.ParserFunctionByUniqueId)
		if err != nil {
			panic(err)
		}
		return nil, type_error.PtrTypeCheckErrorf("parser function already exists %s in %s", uniqueId, string(bytes))
	}
	return universeImpl{
		TypeByTypeName:           u.TypeByTypeName,
		TypeByVariableName:       u.TypeByVariableName,
		Constructors:             u.Constructors,
		GlobalInterfaceVariables: u.GlobalInterfaceVariables,
		ParserFunctionByUniqueId: *u.ParserFunctionByUniqueId.Set(uniqueId, parserFunction),
	}, nil
}

func CopyAddingParserFunctionGeneratingUniqueId(universe Universe, parserFunction parser.Lambda) (string, Universe) {
	id := ksuid.New().String()
	u, err := CopyAddingParserFunctionByUniqueId(universe, id, parserFunction)
	if err != nil {
		panic("CopyAddingParserFunctionGeneratingUniqueId: " + err.Error())
	}
	return id, u
}

func ImportParserFunctionsFrom(universeToAddTo Universe, universeToTakeFrom Universe) (Universe, *type_error.TypecheckError) {
	var err *type_error.TypecheckError
	uToTake := universeToTakeFrom.impl()

	iterator := uToTake.ParserFunctionByUniqueId.Iterator()
	for uniqueId, parserFunction, hasNext := iterator.Next(); hasNext; {
		universeToAddTo, err = CopyAddingParserFunctionByUniqueId(universeToAddTo, uniqueId, parserFunction)
		if err != nil {
			return universeToAddTo, err
		}
	}
	return universeToAddTo, nil
}
