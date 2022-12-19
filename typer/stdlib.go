package typer

var StdLib = Package{
	Packages:   topLevelPackages,
	Interfaces: nil,
}

var topLevelPackages = map[string]Package{
	"tenecs": packageWithPackages(map[string]Package{
		"os": packageWithInterfaces(map[string]Interface{
			"Runtime": runtimeInterface,
			"Main": {
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

var basicTypeString = BasicType{Type: "String"}
var void = Void{}

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
