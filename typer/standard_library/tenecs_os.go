package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_os = packageWith(
	withInterface("Console", &tenecs_os_Console),
	withInterface("Main", &tenecs_os_Main),
	withInterface("Runtime", &tenecs_os_Runtime),
)

var tenecs_os_Console = types.Interface{
	Package: "tenecs.os",
	Name:    "Console",
	Variables: map[string]types.VariableType{
		"log": &types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "message",
					VariableType: &BasicTypeString,
				},
			},
			ReturnType: &Void,
		},
	},
}

var tenecs_os_Main = types.Interface{
	Package: "tenecs.os",
	Name:    "Main",
	Variables: map[string]types.VariableType{
		"main": &types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "runtime",
					VariableType: &tenecs_os_Runtime,
				},
			},
			ReturnType: &Void,
		},
	},
}

var tenecs_os_Runtime = types.Interface{
	Package: "tenecs.os",
	Name:    "Runtime",
	Variables: map[string]types.VariableType{
		"console": &tenecs_os_Console,
	},
}
