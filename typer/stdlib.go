package typer

import (
	"github.com/xplosunn/tenecs/typer/types"
	"strings"
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
	"String":  &basicTypeString,
	"Float":   &basicTypeFloat,
	"Int":     &basicTypeInt,
	"Boolean": &basicTypeBoolean,
	"Void":    &void,
}

var basicTypeString = types.BasicType{Type: "String"}
var basicTypeFloat = types.BasicType{Type: "Float"}
var basicTypeInt = types.BasicType{Type: "Int"}
var basicTypeBoolean = types.BasicType{Type: "Boolean"}
var void = types.Void{}

var topLevelPackages = map[string]Package{
	"tenecs": packageWithPackages(map[string]Package{
		"os": packageWithInterfaces(map[string]*types.Interface{
			"Runtime": &runtimeInterface,
			"Console": &consoleInterface,
			"Main": {
				Package: "tenecs.os",
				Name:    "Main",
				Variables: map[string]types.VariableType{
					"main": &types.Function{
						Arguments: []types.FunctionArgument{
							{
								Name:         "runtime",
								VariableType: &runtimeInterface,
							},
						},
						ReturnType: &void,
					},
				},
			},
		}),
		"test": packageWithInterfaces(map[string]*types.Interface{
			"UnitTests": {
				Package: "tenecs.test",
				Name:    "UnitTests",
				Variables: map[string]types.VariableType{
					"tests": &types.Function{
						Arguments: []types.FunctionArgument{
							{
								Name:         "registry",
								VariableType: &unitTestRegistryInterface,
							},
						},
						ReturnType: &void,
					},
				},
			},
			"UnitTestRegistry": &unitTestRegistryInterface,
			"Assert":           &assertInterface,
		}),
	}),
}

var runtimeInterface = types.Interface{
	Package: "tenecs.os",
	Name:    "Runtime",
	Variables: map[string]types.VariableType{
		"console": &consoleInterface,
	},
}

var consoleInterface = types.Interface{
	Package: "tenecs.os",
	Name:    "Console",
	Variables: map[string]types.VariableType{
		"log": &types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "message",
					VariableType: &basicTypeString,
				},
			},
			ReturnType: &void,
		},
	},
}

var unitTestRegistryInterface = types.Interface{
	Package: "tenecs.test",
	Name:    "UnitTestRegistry",
	Variables: map[string]types.VariableType{
		"test": &types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "name",
					VariableType: &basicTypeString,
				},
				{
					Name: "theTest",
					VariableType: &types.Function{
						Arguments: []types.FunctionArgument{
							{
								Name:         "assert",
								VariableType: &assertInterface,
							},
						},
						ReturnType: &void,
					},
				},
			},
			ReturnType: &void,
		},
	},
}

var assertInterface = types.Interface{
	Package: "tenecs.test",
	Name:    "Assert",
	Variables: map[string]types.VariableType{
		"equal": &types.Function{
			Generics: []string{"T"},
			Arguments: []types.FunctionArgument{
				{
					Name:         "value",
					VariableType: &types.TypeArgument{Name: "T"},
				},
				{
					Name:         "expected",
					VariableType: &types.TypeArgument{Name: "T"},
				},
			},
			ReturnType: &void,
		},
	},
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

func StdLibGetOrPanic(ref string) *types.Interface {
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
		panic("StdLibGetOrPanic" + ref)
	}
	return pkg.Interfaces[finalName]
}
