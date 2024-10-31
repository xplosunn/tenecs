package standard_library

import (
	"fmt"
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
