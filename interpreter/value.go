package interpreter

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
)

type Value interface {
	sealedValue()
}

type ValueVoid struct {
}

func (v ValueVoid) sealedValue() {}

type ValueBoolean struct {
	Bool bool
}

func (v ValueBoolean) sealedValue() {}

type ValueFloat struct {
	Float float64
}

func (v ValueFloat) sealedValue() {}

type ValueInt struct {
	Int int
}

func (v ValueInt) sealedValue() {}

type ValueString struct {
	String string
}

func (v ValueString) sealedValue() {}

type ValueFunction struct {
	Scope       Scope
	AstFunction ast.Function
}

func (v ValueFunction) sealedValue() {}

type ValueNativeFunction struct {
	Scope    Scope
	Function *types.Function
	Invoke   func(values []Value) Value
}

func (v ValueNativeFunction) sealedValue() {}

type ValueStructFunction struct {
	Scope  Scope
	Create func(values []Value) ValueStruct
}

func (v ValueStructFunction) sealedValue() {}

type ValueStruct struct {
	Scope         Scope
	StructName    string
	KeyValues     map[string]Value
	OrderedValues []Value
}

func (v ValueStruct) sealedValue() {}

func ValueExpect[V Value](value Value) (V, bool) {
	result, ok := value.(V)
	return result, ok
}

func ValueExhaustiveSwitch(
	value Value,
	caseVoid func(value ValueVoid),
	caseBoolean func(value ValueBoolean),
	caseFloat func(value ValueFloat),
	caseInt func(value ValueInt),
	caseString func(value ValueString),
	caseFunction func(value ValueFunction),
	caseNativeFunction func(value ValueNativeFunction),
	caseStructFunction func(value ValueStructFunction),
	caseStruct func(value ValueStruct),
) {
	valueVoid, ok := value.(ValueVoid)
	if ok {
		caseVoid(valueVoid)
		return
	}
	valueBoolean, ok := value.(ValueBoolean)
	if ok {
		caseBoolean(valueBoolean)
		return
	}
	valueFloat, ok := value.(ValueFloat)
	if ok {
		caseFloat(valueFloat)
		return
	}
	valueInt, ok := value.(ValueInt)
	if ok {
		caseInt(valueInt)
		return
	}
	valueString, ok := value.(ValueString)
	if ok {
		caseString(valueString)
		return
	}
	valueFunction, ok := value.(ValueFunction)
	if ok {
		caseFunction(valueFunction)
		return
	}
	valueNativeFunction, ok := value.(ValueNativeFunction)
	if ok {
		caseNativeFunction(valueNativeFunction)
		return
	}
	valueStructFunction, ok := value.(ValueStructFunction)
	if ok {
		caseStructFunction(valueStructFunction)
		return
	}
	valueStruct, ok := value.(ValueStruct)
	if ok {
		caseStruct(valueStruct)
		return
	}
	panic(fmt.Errorf("ValueExhaustiveSwitch not implemented for %T", value))
}
