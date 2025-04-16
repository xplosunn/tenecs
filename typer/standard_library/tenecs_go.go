package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_go = packageWith(
	withStruct(Tenecs_go_Console),
	withStruct(Tenecs_go_Main),
	withStruct(Tenecs_go_Runtime),
	withStruct(Tenecs_go_Time),
)

var Tenecs_go_Console = structWithFields("Console", &tenecs_go_Console, tenecs_go_Console_Fields...)

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

var Tenecs_go_Main = structWithFields("Main", &tenecs_go_Main, tenecs_go_Main_Fields...)

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

var Tenecs_go_Runtime = structWithFields("Runtime", &tenecs_go_Runtime, tenecs_go_Runtime_Fields...)

var tenecs_go_Runtime = types.KnownType{
	Package: "tenecs.go",
	Name:    "Runtime",
}

var tenecs_go_Runtime_Fields = []func(fields *StructWithFields){
	structField("console", &tenecs_go_Console),
	structField("ref", tenecs_ref_RefCreator),
	structField("time", &tenecs_go_Time),
}

var Tenecs_go_Time = structWithFields("Time", &tenecs_go_Time, tenecs_go_Time_Fields...)

var tenecs_go_Time = types.KnownType{
	Package: "tenecs.go",
	Name:    "Time",
}

var tenecs_go_Time_Fields = []func(fields *StructWithFields){
	structField("today", functionFromType("() ~> Date", Tenecs_time_Date)),
}
