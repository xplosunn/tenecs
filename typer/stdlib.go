package typer

var StdLib = packageWithPackages(topLevelPackages)

var topLevelPackages = map[string]Package{
	"tenecs": packageWithPackages(map[string]Package{
		"os": packageWithInterfaces(map[string]Interface{
			"Runtime": runtimeInterface,
			"Main": {
				Variables: map[string]VariableType{
					"main": Function{
						ArgumentTypes: []VariableType{runtimeInterface},
						ReturnType:    Void{},
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
					ArgumentTypes: []VariableType{BasicType{Type: "String"}},
					ReturnType:    Void{},
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
