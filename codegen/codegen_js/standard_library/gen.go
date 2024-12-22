package standard_library

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/standard_library"
)

//go:generate go run ../standard_library_generate/main.go

type Function struct {
	Code string
}

type RuntimeFunction struct {
	Params []string
	Body   string
}

func function(opts ...func(*RuntimeFunction)) Function {
	f := &RuntimeFunction{}
	for _, opt := range opts {
		opt(f)
	}

	params := ""
	for i, param := range f.Params {
		if i > 0 {
			params += ", "
		}
		params += param
	}

	body := f.Body

	return Function{
		Code: fmt.Sprintf(`(%s) {
%s
return null
}`, params, body),
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
	bodyStr := "return ({\n"
	bodyStr += fmt.Sprintf(`  "$type": "%s",`, structWithFields.Struct.Name) + "\n"
	for _, fieldName := range structWithFields.FieldNamesSorted {
		bodyStr += fmt.Sprintf(`  "%s": %s,`, fieldName, fieldName) + "\n"
	}
	bodyStr += "})"
	return function(
		params(structWithFields.FieldNamesSorted...),
		body(bodyStr),
	)
}
