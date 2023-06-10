package codegen

import (
	"fmt"
	"strings"
)

func GenerateRuntime() ([]Import, string) {
	imports := []Import{}

	imports = append(imports, "fmt")
	console := ofMap(map[string]string{
		"log": function(params("Pmessage"), body(`fmt.Println(Pmessage)`)),
	})

	runtime := ofMap(map[string]string{
		"console": console,
	})

	return imports, runtime
}

func ofMap(m map[string]string) string {
	result := "map[string]any{"
	for k, v := range m {
		result += "\n" + fmt.Sprintf(`"%s": %s,`, k, v)
	}
	result += "\n}"

	return result
}

type RuntimeFunction struct {
	Params []string
	Body   []string
}

func function(opts ...func(*RuntimeFunction)) string {
	f := &RuntimeFunction{}
	for _, opt := range opts {
		opt(f)
	}

	params := ""
	for i, param := range f.Params {
		if i > 0 {
			params += ","
		}
		params += param + " any"
	}

	body := strings.Join(f.Body, "\n")

	return fmt.Sprintf(`func (%s) any {
%s
return nil
}`, params, body)
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
