package codegen_golang

import (
	"fmt"
	"golang.org/x/exp/maps"
	"sort"
	"strings"
)

func GenerateRuntime() ([]Import, string) {
	imports := []Import{}

	imports = append(imports, "fmt", "time")
	console := ofMap("tenecs_go_Console", map[string]string{
		"_log": function(params("Pmessage"), body(`fmt.Println(Pmessage)`)),
	})

	time := ofMap("tenecs_go_Time", map[string]string{
		"_today": function(params(), body(`t := time.Now()
return tenecs_time_Date{
  _year: t.Year(),
  _month: int(t.Month()),
  _day: t.Day(),
}`)),
	})

	runtime := ofMap("tenecs_go_Runtime", map[string]string{
		"_console": console,
		"_ref":     runtimeRefCreator(),
		"_time":    time,
	})

	return imports, runtime
}

func runtimeRefCreator() string {
	return ofMap("tenecs_ref_RefCreator", map[string]string{
		"_new": function(
			params("Pvalue"),
			body(`var ref any = Pvalue
return tenecs_ref_Ref{
_get: func()any {
return ref
},
_set: func(value any)any {
ref = value
return nil
},
_modify: func(f any) any {
ref = f.(func(any)any)(ref)
return nil
},
}
`),
		),
	})
}

func ofMap(structName string, m map[string]string) string {
	result := structName + "{"
	keys := maps.Keys(m)
	sort.Strings(keys)
	for _, k := range keys {
		result += "\n" + fmt.Sprintf(`%s: %s,`, k, m[k])
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
