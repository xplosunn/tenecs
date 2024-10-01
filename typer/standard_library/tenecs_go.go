package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_go = packageWith(
	withStruct("Console", &tenecs_go_Console, tenecs_go_Console_Fields...),
	withStruct("Main", &tenecs_go_Main, tenecs_go_Main_Fields...),
	withStruct("Runtime", &tenecs_go_Runtime, tenecs_go_Runtime_Fields...),
)

var tenecs_go_Console = types.KnownType{
	Package: "tenecs.go",
	Name:    "Console",
}

var tenecs_go_Console_Fields = []func(fields *StructWithFields){
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

var tenecs_go_Main = types.KnownType{
	Package: "tenecs.go",
	Name:    "Main",
}

var tenecs_go_Main_Fields = []func(fields *StructWithFields){
	structField("main", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "runtime",
				VariableType: &tenecs_go_Runtime,
			},
		},
		ReturnType: types.Void(),
	}),
}

var tenecs_go_Runtime = types.KnownType{
	Package: "tenecs.go",
	Name:    "Runtime",
}

var tenecs_go_Runtime_Fields = []func(fields *StructWithFields){
	structField("console", &tenecs_go_Console),
	structField("http", tenecs_http_RuntimeServer),
	structField("ref", tenecs_ref_RefCreator),
}
