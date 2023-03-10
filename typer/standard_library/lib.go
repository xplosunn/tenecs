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

type Package struct {
	Packages   map[string]Package
	Interfaces map[string]*types.Interface
}

var DefaultTypesAvailableWithoutImport = map[string]types.VariableType{
	"String":  &BasicTypeString,
	"Float":   &BasicTypeFloat,
	"Int":     &BasicTypeInt,
	"Boolean": &BasicTypeBoolean,
	"Void":    &Void,
}

var BasicTypeString = types.BasicType{Type: "String"}
var BasicTypeFloat = types.BasicType{Type: "Float"}
var BasicTypeInt = types.BasicType{Type: "Int"}
var BasicTypeBoolean = types.BasicType{Type: "Boolean"}
var Void = types.Void{}

var topLevelPackages = map[string]Package{
	"tenecs": packageWithPackages(map[string]Package{
		"os":   tenecs_os,
		"test": tenecs_test,
	}),
}

func packageWithPackages(packages map[string]Package) Package {
	return Package{
		Packages:   packages,
		Interfaces: map[string]*types.Interface{},
	}
}

func packageWithInterfaces(interfaces map[string]*types.Interface) Package {
	return Package{
		Packages:   map[string]Package{},
		Interfaces: interfaces,
	}
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
