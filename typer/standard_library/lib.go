package standard_library

import (
	"github.com/xplosunn/tenecs/typer/types"
	"strings"
	"testing"
)

var StdLib = Package{
	Packages: topLevelPackages,
}

var DefaultTypesAvailableWithoutImport = map[string]types.VariableType{
	"String":  types.String(),
	"Float":   types.Float(),
	"Int":     types.Int(),
	"Boolean": types.Boolean(),
	"Void":    types.Void(),
	"List": types.List(&types.TypeArgument{
		Name: "T",
	}),
}

var topLevelPackages = map[string]Package{
	"tenecs": packageWith(
		withPackage("list", tenecs_list),
		withPackage("boolean", tenecs_boolean),
		withPackage("compare", tenecs_compare),
		withPackage("error", tenecs_error),
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
	if pkg.Structs[finalName] == nil {
		t.Fatal("StdLibGetOrPanic" + ref)
	}
	return pkg.Structs[finalName].Struct
}

func StdLibGetFunctionOrPanic(t *testing.T, ref string) *types.Function {
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
	if pkg.Structs[finalName] == nil {
		t.Fatal("StdLibGetOrPanic" + ref)
	}
	arguments := []types.FunctionArgument{}
	for _, fieldName := range pkg.Structs[finalName].FieldNamesSorted {
		arguments = append(arguments, types.FunctionArgument{
			Name:         fieldName,
			VariableType: pkg.Structs[finalName].Fields[fieldName],
		})
	}
	if len(pkg.Structs[finalName].Struct.Generics) > 0 {
		panic("todo StdLibGetFunctionOrPanic with generics")
	}
	return &types.Function{
		Generics:   nil,
		Arguments:  arguments,
		ReturnType: pkg.Structs[finalName].Struct,
	}
}
