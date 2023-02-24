package interpreter

import (
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/typer/ast"
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
	for _, declaration := range program.Declarations {
		_, value, err := EvalExpression(firstScope, declaration.Expression)
		if err != nil {
			return firstScope, err
		}
		scopeValueByName.Set(declaration.Name, value)
	}
	return scopeImpl{
		ValueByName: scopeValueByName.Map(),
	}, nil
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
