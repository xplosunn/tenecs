package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_http = packageWith(
	withStruct("RuntimeServer", tenecs_http_RuntimeServer, tenecs_http_RuntimeServer_Fields...),
	withStruct("Server", tenecs_http_Server, tenecs_http_Server_Fields...),
	withStruct("ServerError", tenecs_http_ServerError, tenecs_http_ServerError_Fields...),
	withFunction("newServer", tenecs_http_newServer),
)

var tenecs_http_RuntimeServer = types.Interface(
	"tenecs.http",
	"RuntimeServer",
	nil,
)

var tenecs_http_RuntimeServer_Fields = []func(fields *StructWithFields){
	structField("serve", &types.Function{
		Arguments: []types.FunctionArgument{
			types.FunctionArgument{
				Name:         "server",
				VariableType: tenecs_http_Server,
			},
			types.FunctionArgument{
				Name:         "address",
				VariableType: types.String(),
			},
		},
		ReturnType: tenecs_http_ServerError,
	}),
}

var tenecs_http_Server = types.Interface(
	"tenecs.http",
	"Server",
	nil,
)

var tenecs_http_Server_Fields = []func(fields *StructWithFields){
	structField("restHandlerGet", &types.Function{
		Generics: []string{"ResponseBody"},
		Arguments: []types.FunctionArgument{
			types.FunctionArgument{
				Name:         "toJson",
				VariableType: tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "ResponseBody"}),
			},
			types.FunctionArgument{
				Name:         "route",
				VariableType: types.String(),
			},
			types.FunctionArgument{
				Name: "handler",
				VariableType: &types.Function{
					Arguments: []types.FunctionArgument{
						types.FunctionArgument{
							Name:         "responseStatus",
							VariableType: types.UncheckedApplyGenerics(tenecs_ref_Ref, []types.VariableType{types.Int()}),
						},
					},
					ReturnType: &types.TypeArgument{Name: "ResponseBody"},
				},
			},
		},
		ReturnType: types.Void(),
	}),
	structField("restHandlerPost", &types.Function{
		Generics: []string{"RequestBody", "ResponseBody"},
		Arguments: []types.FunctionArgument{
			types.FunctionArgument{
				Name:         "fromJson",
				VariableType: tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "RequestBody"}),
			},
			types.FunctionArgument{
				Name:         "toJson",
				VariableType: tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "ResponseBody"}),
			},
			types.FunctionArgument{
				Name:         "route",
				VariableType: types.String(),
			},
			types.FunctionArgument{
				Name: "handler",
				VariableType: &types.Function{
					Arguments: []types.FunctionArgument{
						types.FunctionArgument{
							Name:         "requestBody",
							VariableType: &types.TypeArgument{Name: "RequestBody"},
						},
						types.FunctionArgument{
							Name:         "responseStatus",
							VariableType: types.UncheckedApplyGenerics(tenecs_ref_Ref, []types.VariableType{types.Int()}),
						},
					},
					ReturnType: &types.TypeArgument{Name: "ResponseBody"},
				},
			},
		},
		ReturnType: types.Void(),
	}),
	structField("runRestPostWithBody", &types.Function{
		Arguments: []types.FunctionArgument{
			types.FunctionArgument{
				Name:         "route",
				VariableType: types.String(),
			},
			types.FunctionArgument{
				Name:         "requestBody",
				VariableType: types.String(),
			},
		},
		ReturnType: types.String(),
	}),
}

var tenecs_http_ServerError = types.Struct(
	"tenecs.http",
	"ServerError",
	nil,
)

var tenecs_http_ServerError_Fields = []func(fields *StructWithFields){
	structField("message", types.String()),
}

var tenecs_http_newServer = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "refCreator",
			VariableType: tenecs_ref_RefCreator,
		},
	},
	ReturnType: tenecs_http_Server,
}
