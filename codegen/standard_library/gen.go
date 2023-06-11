package standard_library

import (
	"fmt"
	"strings"
)

//go:generate go run ../standard_library_generate/main.go

type Function struct {
	Imports []string
	Code    string
}

type RuntimeFunction struct {
	Imports []string
	Params  []string
	Body    []string
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
		params += param + " any"
	}

	body := strings.Join(f.Body, "\n")

	return Function{
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

func body(b ...string) func(*RuntimeFunction) {
	return func(runtimeFunction *RuntimeFunction) {
		runtimeFunction.Body = b
	}
}
