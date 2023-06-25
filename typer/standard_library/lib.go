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
	"String":  &BasicTypeString,
	"Float":   &BasicTypeFloat,
	"Int":     &BasicTypeInt,
	"Boolean": &BasicTypeBoolean,
	"Void":    &Void,
	"Array":   &Array,
}

var BasicTypeString = types.BasicType{Type: "String"}
var BasicTypeFloat = types.BasicType{Type: "Float"}
var BasicTypeInt = types.BasicType{Type: "Int"}
var BasicTypeBoolean = types.BasicType{Type: "Boolean"}
var Void = types.Void{}
var Array = types.Array{OfType: &types.TypeArgument{Name: "T"}}

var topLevelPackages = map[string]Package{
	"tenecs": packageWith(
		withPackage("array", tenecs_array),
		withPackage("json", tenecs_json),
		withPackage("os", tenecs_os),
		withPackage("string", tenecs_string),
		withPackage("test", tenecs_test),
	),
}

func StdLibGetOrPanic(t *testing.T, ref string) *types.Interface {
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
	return pkg.Interfaces[finalName]
}
