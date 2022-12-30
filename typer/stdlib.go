package typer

var StdLib = Package{
	Packages:   topLevelPackages,
	Interfaces: nil,
}

var StdLibInterfaceVariables = map[string]map[string]VariableType{
	"tenecs.os.Main": map[string]VariableType{
		"main": Function{
			Arguments: []FunctionArgument{
				{
					Name:         "runtime",
					VariableType: runtimeInterface,
				},
			},
			ReturnType: void,
		},
	},
	"tenecs.os.Runtime": map[string]VariableType{
		"console": consoleInterface,
	},
	"tenecs.os.Console": map[string]VariableType{
		"log": Function{
			Arguments: []FunctionArgument{
				{
					Name:         "message",
					VariableType: basicTypeString,
				},
			},
			ReturnType: void,
		},
	},
}

type Package struct {
	Packages   map[string]Package
	Interfaces map[string]Interface
}

var DefaultTypesAvailableWithoutImport = map[string]VariableType{
	"String":  basicTypeString,
	"Float":   basicTypeFloat,
	"Int":     basicTypeInt,
	"Boolean": basicTypeBoolean,
	"Void":    void,
}

var basicTypeString = BasicType{Type: "String"}
var basicTypeFloat = BasicType{Type: "Float"}
var basicTypeInt = BasicType{Type: "Int"}
var basicTypeBoolean = BasicType{Type: "Boolean"}
var void = Void{}

var topLevelPackages = map[string]Package{
	"tenecs": packageWithPackages(map[string]Package{
		"os": packageWithInterfaces(map[string]Interface{
			"Runtime": runtimeInterface,
			"Console": consoleInterface,
			"Main": {
				Package: "tenecs.os",
				Name:    "Main",
			},
		}),
	}),
}

var runtimeInterface = Interface{
	Package: "tenecs.os",
	Name:    "Runtime",
}

var consoleInterface = Interface{
	Package: "tenecs.os",
	Name:    "Console",
}

func packageWithPackages(packages map[string]Package) Package {
	return Package{
		Packages:   packages,
		Interfaces: map[string]Interface{},
	}
}

func packageWithInterfaces(interfaces map[string]Interface) Package {
	return Package{
		Packages:   map[string]Package{},
		Interfaces: interfaces,
	}
}
