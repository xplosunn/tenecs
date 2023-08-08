package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_http = packageWith(
	withInterface("Server", tenecs_http_Server, tenecs_http_Server_Fields),
	withInterface("ServerError", tenecs_http_ServerError, tenecs_http_ServerError_Fields),
	withFunction("newServer", tenecs_http_newServer),
)

var tenecs_http_Server = types.Interface(
	"tenecs.http",
	"Server",
	nil,
)

var tenecs_http_Server_Fields = map[string]types.VariableType{
	"restHandlerGet": &types.Function{
		Generics: []string{"ResponseBody"},
		Arguments: []types.FunctionArgument{
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
	},
	"restHandlerPost": &types.Function{
		Generics: []string{"RequestBody", "ResponseBody"},
		Arguments: []types.FunctionArgument{
			types.FunctionArgument{
				Name:         "fromJson",
				VariableType: tenecs_json_FromJson_Of(&types.TypeArgument{Name: "RequestBody"}),
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
	},
	"serve": &types.Function{
		Arguments: []types.FunctionArgument{
			types.FunctionArgument{
				Name:         "address",
				VariableType: types.String(),
			},
			types.FunctionArgument{
				Name:         "blocker",
				VariableType: tenecs_execution_Blocker,
			},
		},
		ReturnType: tenecs_http_ServerError,
	},
}

var tenecs_http_ServerError = types.Interface(
	"tenecs.http",
	"ServerError",
	nil,
)

var tenecs_http_ServerError_Fields = map[string]types.VariableType{
	"message": types.String(),
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
