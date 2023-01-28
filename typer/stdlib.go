package typer

import "github.com/xplosunn/tenecs/typer/types"

var StdLib = Package{
	Packages:   topLevelPackages,
	Interfaces: nil,
}

var StdLibInterfaceVariables = map[string]map[string]types.VariableType{
	"tenecs.os.Main": map[string]types.VariableType{
		"main": types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "runtime",
					VariableType: runtimeInterface,
				},
			},
			ReturnType: void,
		},
	},
	"tenecs.os.Runtime": map[string]types.VariableType{
		"console": consoleInterface,
	},
	"tenecs.os.Console": map[string]types.VariableType{
		"log": types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "message",
					VariableType: basicTypeString,
				},
			},
			ReturnType: void,
		},
	},
	"tenecs.test.UnitTests": map[string]types.VariableType{
		"tests": types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "registry",
					VariableType: unitTestRegistryInterface,
				},
			},
			ReturnType: void,
		},
	},
	"tenecs.test.UnitTestRegistry": map[string]types.VariableType{
		"suite": types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "name",
					VariableType: basicTypeString,
				},
				{
					Name: "theSuite",
					VariableType: types.Function{
						Arguments: []types.FunctionArgument{
							{
								Name:         "registry",
								VariableType: unitTestRegistryInterface,
							},
						},
						ReturnType: void,
					},
				},
			},
			ReturnType: void,
		},
		"test": types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "name",
					VariableType: basicTypeString,
				},
				{
					Name: "theTest",
					VariableType: types.Function{
						Arguments: []types.FunctionArgument{
							{
								Name:         "assert",
								VariableType: assertInterface,
							},
						},
						ReturnType: void,
					},
				},
			},
			ReturnType: void,
		},
	},
	"tenecs.test.Assert": map[string]types.VariableType{
		"equal": types.Function{
			Generics: []string{"T"},
			Arguments: []types.FunctionArgument{
				{
					Name:         "value",
					VariableType: types.TypeArgument{Name: "T"},
				},
				{
					Name:         "expected",
					VariableType: types.TypeArgument{Name: "T"},
				},
			},
			ReturnType: void,
		},
	},
}

type Package struct {
	Packages   map[string]Package
	Interfaces map[string]types.Interface
}

var DefaultTypesAvailableWithoutImport = map[string]types.VariableType{
	"String":  basicTypeString,
	"Float":   basicTypeFloat,
	"Int":     basicTypeInt,
	"Boolean": basicTypeBoolean,
	"Void":    void,
}

var basicTypeString = types.BasicType{Type: "String"}
var basicTypeFloat = types.BasicType{Type: "Float"}
var basicTypeInt = types.BasicType{Type: "Int"}
var basicTypeBoolean = types.BasicType{Type: "Boolean"}
var void = types.Void{}

var topLevelPackages = map[string]Package{
	"tenecs": packageWithPackages(map[string]Package{
		"os": packageWithInterfaces(map[string]types.Interface{
			"Runtime": runtimeInterface,
			"Console": consoleInterface,
			"Main": {
				Package: "tenecs.os",
				Name:    "Main",
			},
		}),
		"test": packageWithInterfaces(map[string]types.Interface{
			"UnitTests": {
				Package: "tenecs.test",
				Name:    "UnitTests",
			},
			"UnitTestRegistry": unitTestRegistryInterface,
			"Assert":           assertInterface,
		}),
	}),
}

var runtimeInterface = types.Interface{
	Package: "tenecs.os",
	Name:    "Runtime",
}

var consoleInterface = types.Interface{
	Package: "tenecs.os",
	Name:    "Console",
}

var unitTestRegistryInterface = types.Interface{
	Package: "tenecs.test",
	Name:    "UnitTestRegistry",
}

var assertInterface = types.Interface{
	Package: "tenecs.test",
	Name:    "Assert",
}

func packageWithPackages(packages map[string]Package) Package {
	return Package{
		Packages:   packages,
		Interfaces: map[string]types.Interface{},
	}
}

func packageWithInterfaces(interfaces map[string]types.Interface) Package {
	return Package{
		Packages:   map[string]Package{},
		Interfaces: interfaces,
	}
}
