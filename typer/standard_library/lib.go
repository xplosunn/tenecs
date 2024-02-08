package standard_library

import (
	"github.com/xplosunn/tenecs/typer/types"
	"strings"
	"testing"
)

var StdLib = Package{
	Packages:   topLevelPackages,
	Interfaces: nil,
}

var DefaultTypesAvailableWithoutImport = map[string]types.VariableType{
	"String":  types.String(),
	"Float":   types.Float(),
	"Int":     types.Int(),
	"Boolean": types.Boolean(),
	"Void":    types.Void(),
	"Array": types.Array(&types.TypeArgument{
		Name: "T",
	}),
}

var topLevelPackages = map[string]Package{
	"tenecs": packageWith(
		withPackage("array", tenecs_array),
		withPackage("boolean", tenecs_boolean),
		withPackage("compare", tenecs_compare),
		withPackage("http", tenecs_http),
		withPackage("int", tenecs_int),
		withPackage("json", tenecs_json),
		withPackage("os", tenecs_os),
		withPackage("ref", tenecs_ref),
		withPackage("string", tenecs_string),
		withPackage("test", tenecs_test),
	),
}

func StdLibGetOrPanic(t *testing.T, ref string) *types.KnownType {
	pkg := StdLib
	split := strings.Split(ref, ".")
	var finalName string
	for i, name := range split {
		if i < len(split)-1 {
			pkg = pkg.Packages[name]
		} else {
			finalName = name
		}
	}
	if pkg.Interfaces[finalName] == nil {
		t.Fatal("StdLibGetOrPanic" + ref)
	}
	return pkg.Interfaces[finalName].Interface
}
