package typer

var StdLib = Package{
	Packages:   topLevelPackages,
	Interfaces: nil,
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
			"Main": {
				Package: "tenecs.os",
				Name:    "Main",
				Variables: map[string]VariableType{
					"main": Function{
						ArgumentTypes: []VariableType{runtimeInterface},
						ReturnType:    void,
					},
				},
			},
		}),
	}),
}

var runtimeInterface = Interface{
	Package: "tenecs.os",
	Name:    "Runtime",
	Variables: map[string]VariableType{
		"console": Interface{
			Variables: map[string]VariableType{
				"log": Function{
					ArgumentTypes: []VariableType{basicTypeString},
					ReturnType:    void,
				},
			},
		},
	},
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
