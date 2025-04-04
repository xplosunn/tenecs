package standard_library

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/types"
)

//go:generate go run ../standard_library_generate/main.go

type Function interface {
	sealedFunction()
	FunctionCases() (*NativeFunction, *StructFunction)
}

type NativeFunction struct {
	Imports []string
	Code    string
}

func (f NativeFunction) sealedFunction() {}

func (f NativeFunction) FunctionCases() (*NativeFunction, *StructFunction) {
	return &f, nil
}

type StructFunction struct {
	Struct           *types.KnownType
	Fields           map[string]types.VariableType
	FieldNamesSorted []string
}

func (f StructFunction) sealedFunction() {}

func (f StructFunction) FunctionCases() (*NativeFunction, *StructFunction) {
	return nil, &f
}

type RuntimeFunction struct {
	Imports []string
	Params  []string
	Body    string
}

func function(opts ...func(*RuntimeFunction)) NativeFunction {
	f := &RuntimeFunction{}
	for _, opt := range opts {
		opt(f)
	}

	params := ""
	for i, param := range f.Params {
		if i > 0 {
			params += ", "
		}
		params += param + " any"
	}

	body := f.Body

	return NativeFunction{
		Imports: f.Imports,
		Code: fmt.Sprintf(`func (%s) any {
%s
return nil
}`, params, body),
	}
}

func imports(i ...string) func(*RuntimeFunction) {
	return func(runtimeFunction *RuntimeFunction) {
		runtimeFunction.Imports = i
	}
}

func params(p ...string) func(*RuntimeFunction) {
	return func(runtimeFunction *RuntimeFunction) {
		runtimeFunction.Params = p
	}
}

func body(b string) func(*RuntimeFunction) {
	return func(runtimeFunction *RuntimeFunction) {
		runtimeFunction.Body = b
	}
}

func structFunction(structWithFields *standard_library.StructWithFields) Function {
	return StructFunction{
		Struct:           structWithFields.Struct,
		Fields:           structWithFields.Fields,
		FieldNamesSorted: structWithFields.FieldNamesSorted,
	}
}
