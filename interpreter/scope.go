package interpreter

import (
	"fmt"
	"github.com/benbjohnson/immutable"
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

func NewScope() Scope {
	return scopeImpl{
		ValueByName: immutable.NewMap[string, Value](nil),
	}
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
