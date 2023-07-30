package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_json = packageWith(
	withInterface("FromJson", tenecs_json_FromJson, tenecs_json_FromJson_Fields),
	withInterface("JsonError", tenecs_json_JsonError, tenecs_json_JsonError_Fields),
	withFunction("jsonError", tenecs_json_jsonError),
	withFunction("parseArray", tenecs_json_parseArray),
	withFunction("parseBoolean", tenecs_json_parseBoolean),
	withFunction("parseInt", tenecs_json_parseInt),
	withFunction("parseString", tenecs_json_parseString),
	withFunction("toJson", tenecs_json_toJson),
)

var tenecs_json_FromJson = types.Interface(
	"tenecs.json",
	"FromJson",
	[]string{"T"},
)

func tenecs_json_FromJson_Of(varType types.VariableType) *types.KnownType {
	fromJson := types.Interface(
		"tenecs.json",
		"FromJson",
		[]string{"T"},
	)
	fromJson.Generics = []types.VariableType{varType}
	return fromJson
}

var tenecs_json_FromJson_Fields = map[string]types.VariableType{
	"parse": &types.Function{
		Arguments: []types.FunctionArgument{
			types.FunctionArgument{
				Name:         "json",
				VariableType: types.String(),
			},
		},
		ReturnType: &types.OrVariableType{
			Elements: []types.VariableType{
				&types.TypeArgument{
					Name: "T",
				},
				tenecs_json_JsonError,
			},
		},
	},
}

var tenecs_json_JsonError = types.Interface(
	"tenecs.json",
	"JsonError",
	nil,
)

var tenecs_json_JsonError_Fields = map[string]types.VariableType{
	"message": types.String(),
}

var tenecs_json_jsonError = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "message",
			VariableType: types.String(),
		},
	},
	ReturnType: tenecs_json_JsonError,
}

var tenecs_json_parseArray = &types.Function{
	Generics: []string{"T"},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "of",
			VariableType: tenecs_json_FromJson_Of(&types.TypeArgument{Name: "T"}),
		},
	},
	ReturnType: tenecs_json_FromJson_Of(types.UncheckedArray(&types.TypeArgument{Name: "T"})),
}

var tenecs_json_parseBoolean = &types.Function{
	Arguments:  []types.FunctionArgument{},
	ReturnType: tenecs_json_FromJson_Of(types.Boolean()),
}

var tenecs_json_parseInt = &types.Function{
	Arguments:  []types.FunctionArgument{},
	ReturnType: tenecs_json_FromJson_Of(types.Int()),
}

var tenecs_json_parseString = &types.Function{
	Arguments:  []types.FunctionArgument{},
	ReturnType: tenecs_json_FromJson_Of(types.String()),
}

var tenecs_json_toJson = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "input",
			VariableType: &types.TypeArgument{
				Name: "T",
			},
		},
	},
	ReturnType: types.String(),
}
