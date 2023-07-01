package interpreter

import (
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"strings"
)

type Scope interface {
	impl() scopeImpl
}

type scopeImpl struct {
	ValueByName *immutable.Map[string, Value]
}

func (scope scopeImpl) impl() scopeImpl {
	return scope
}

func NewScope(program ast.Program) (Scope, error) {
	firstScopeValueByName := immutable.NewMap[string, Value](nil)

	firstScope := scopeImpl{
		ValueByName: firstScopeValueByName,
	}
	scopeValueByName := immutable.NewMapBuilder[string, Value](nil)
	for structName, structFunction := range program.StructFunctions {
		scopeValueByName.Set(structName, ValueStructFunction{
			Scope: firstScope,
			Create: func(values []Value) ValueStruct {
				if len(structFunction.Arguments) != len(values) {
					panic(fmt.Sprintf("ValueStructFunction Create len(%v) != len(%v)", structFunction.Arguments, values))
				}
				structKeyValues := map[string]Value{}
				for i, argument := range structFunction.Arguments {
					structKeyValues[argument.Name] = values[i]
				}
				return ValueStruct{
					Scope:         firstScope,
					StructName:    structName,
					KeyValues:     structKeyValues,
					OrderedValues: values,
				}
			},
		})
	}
	for functionName, function := range program.NativeFunctions {
		var invoke func(passedGenerics []types.StructFieldVariableType, values []Value) Value
		if functionName == "join" {
			invoke = func(passedGenerics []types.StructFieldVariableType, values []Value) Value {
				return ValueString{
					String: strings.TrimSuffix(values[0].(ValueString).String, "\"") + strings.TrimPrefix(values[1].(ValueString).String, "\""),
				}
			}
		} else {
			panic("todo NewScope NativeFunction " + functionName)
		}
		scopeValueByName.Set(functionName, ValueNativeFunction{
			Scope:    firstScope,
			Function: function,
			Invoke:   invoke,
		})
	}
	scope := &scopeImpl{
		ValueByName: scopeValueByName.Map(),
	}
	for _, declaration := range program.Declarations {
		_, value, err := EvalExpression(*scope, declaration.Expression)
		if err != nil {
			return firstScope, err
		}
		scope = &scopeImpl{
			ValueByName: scope.ValueByName.Set(declaration.Name, value),
		}
	}

	return scope, nil
}

func Resolve(scope Scope, name string) (Value, error) {
	value, ok := scope.impl().ValueByName.Get(name)
	if !ok {
		return nil, fmt.Errorf("couldn't find %s in Scope", name)
	}
	return value, nil
}

func CopyAdding(scope Scope, name string, value Value) Scope {
	return scopeImpl{ValueByName: scope.impl().ValueByName.Set(name, value)}
}
