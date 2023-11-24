package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_os = packageWith(
	withInterface("Console", &tenecs_os_Console, tenecs_os_Console_Fields),
	withInterface("Main", &tenecs_os_Main, tenecs_os_Main_Fields),
	withInterface("Runtime", &tenecs_os_Runtime, tenecs_os_Runtime_Fields),
)

var tenecs_os_Console = types.KnownType{
	Package: "tenecs.os",
	Name:    "Console",
}

var tenecs_os_Console_Fields = map[string]types.VariableType{
	"log": &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "message",
				VariableType: types.String(),
			},
		},
		ReturnType: types.Void(),
	},
}

var tenecs_os_Main = types.KnownType{
	Package: "tenecs.os",
	Name:    "Main",
}

var tenecs_os_Main_Fields = map[string]types.VariableType{
	"main": &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "runtime",
				VariableType: &tenecs_os_Runtime,
			},
		},
		ReturnType: types.Void(),
	},
}

var tenecs_os_Runtime = types.KnownType{
	Package: "tenecs.os",
	Name:    "Runtime",
}

var tenecs_os_Runtime_Fields = map[string]types.VariableType{
	"console": &tenecs_os_Console,
	"http":    tenecs_http_RuntimeServer,
	"ref":     tenecs_ref_RefCreator,
}
