package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_os = packageWith(
	withInterface("Console", &tenecs_os_Console, tenecs_os_Console_Fields),
	withInterface("Main", &Tenecs_os_Main, Tenecs_os_Main_Fields),
	withInterface("Runtime", &tenecs_os_Runtime, tenecs_os_Runtime_Fields),
	withInterface("RuntimeExecution", tenecs_os_RuntimeExecution, tenecs_os_RuntimeExecution_Fields),
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

var Tenecs_os_Main = types.KnownType{
	Package: "tenecs.os",
	Name:    "Main",
}

var Tenecs_os_Main_Fields = map[string]types.VariableType{
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
	"console":   &tenecs_os_Console,
	"execution": tenecs_os_RuntimeExecution,
	"ref":       tenecs_ref_RefCreator,
}

var tenecs_os_RuntimeExecution = types.Interface(
	"tenecs.os",
	"RuntimeExecution",
	nil,
)

var tenecs_os_RuntimeExecution_Fields = map[string]types.VariableType{
	"runBlocking": &types.Function{
		Generics: []string{"R"},
		Arguments: []types.FunctionArgument{
			{
				Name:         "operation",
				VariableType: tenecs_execution_BlockingOperation,
			},
		},
		ReturnType: &types.TypeArgument{Name: "R"},
	},
}
