package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_os = packageWith(
	withStruct("Console", &tenecs_os_Console, tenecs_os_Console_Fields...),
	withStruct("Main", &Tenecs_os_Main, Tenecs_os_Main_Fields...),
	withStruct("Runtime", &tenecs_os_Runtime, tenecs_os_Runtime_Fields...),
)

var tenecs_os_Console = types.KnownType{
	Package: "tenecs.os",
	Name:    "Console",
}

var tenecs_os_Console_Fields = []func(fields *StructWithFields){
	structField("log", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "message",
				VariableType: types.String(),
			},
		},
		ReturnType: types.Void(),
	}),
}

var Tenecs_os_Main = types.KnownType{
	Package: "tenecs.os",
	Name:    "Main",
}

var Tenecs_os_Main_Fields = []func(fields *StructWithFields){
	structField("main", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "runtime",
				VariableType: &tenecs_os_Runtime,
			},
		},
		ReturnType: types.Void(),
	}),
}

var tenecs_os_Runtime = types.KnownType{
	Package: "tenecs.os",
	Name:    "Runtime",
}

var tenecs_os_Runtime_Fields = []func(fields *StructWithFields){
	structField("console", &tenecs_os_Console),
	structField("http", tenecs_http_RuntimeServer),
	structField("ref", tenecs_ref_RefCreator),
}
