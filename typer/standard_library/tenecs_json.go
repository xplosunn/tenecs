package standard_library

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/types"
)

var tenecs_json = packageWith(
	withFunction("field", tenecs_json_field),
	withInterface("FromJson", tenecs_json_FromJson, tenecs_json_FromJson_Fields),
	withInterface("FromJsonField", tenecs_json_FromJsonField, tenecs_json_FromJsonField_Fields),
	withStruct("JsonError", tenecs_json_JsonError, tenecs_json_JsonError_Fields...),
	withFunction("parseArray", tenecs_json_parseArray),
	withFunction("parseBoolean", tenecs_json_parseBoolean),
	withFunction("parseInt", tenecs_json_parseInt),
	withFunction("parseObject0", tenecs_json_parseObject0),
	withFunctions(tenecs_json_parseObject),
	withFunction("parseOr", tenecs_json_parseOr),
	withFunction("parseString", tenecs_json_parseString),
	withFunction("toJson", tenecs_json_toJson),
)

var tenecs_json_field = &types.Function{
	Generics: []string{"T"},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "name",
			VariableType: types.String(),
		},
		types.FunctionArgument{
			Name:         "fromJson",
			VariableType: tenecs_json_FromJson_Of(&types.TypeArgument{Name: "T"}),
		},
	},
	ReturnType: tenecs_json_FromJsonField,
}

var tenecs_json_FromJson = types.Interface(
	"tenecs.json",
	"FromJson",
	[]string{"T"},
)

func tenecs_json_FromJson_Of(varType types.VariableType) *types.KnownType {
	return types.UncheckedApplyGenerics(tenecs_json_FromJson, []types.VariableType{varType})
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

var tenecs_json_FromJsonField = types.Interface(
	"tenecs.json",
	"FromJsonField",
	[]string{"T"},
)

func tenecs_json_FromJsonField_Of(varType types.VariableType) *types.KnownType {
	return types.UncheckedApplyGenerics(tenecs_json_FromJsonField, []types.VariableType{varType})
}

var tenecs_json_FromJsonField_Fields = map[string]types.VariableType{
	"name":     types.String(),
	"fromJson": tenecs_json_FromJson_Of(&types.TypeArgument{Name: "T"}),
}

var tenecs_json_JsonError = types.Struct(
	"tenecs.json",
	"JsonError",
	nil,
)

var tenecs_json_JsonError_Fields = []func(fields *StructWithFields){
	structField("message", types.String()),
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

var tenecs_json_parseObject0 = &types.Function{
	Generics: []string{"R"},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "build",
			VariableType: &types.Function{
				Arguments:  []types.FunctionArgument{},
				ReturnType: &types.TypeArgument{Name: "R"},
			},
		},
	},
	ReturnType: tenecs_json_FromJson_Of(&types.TypeArgument{Name: "R"}),
}

var tenecs_json_parseObject = func() []NamedFunction {
	result := []NamedFunction{}
	for i := 1; i < 23; i++ {
		generics := []string{"R"}
		buildArguments := []types.FunctionArgument{}
		argumentsAfterBuild := []types.FunctionArgument{}
		for j := 0; j < i; j++ {
			generics = append(generics, fmt.Sprintf("I%d", j))
			buildArguments = append(buildArguments, types.FunctionArgument{
				Name:         fmt.Sprintf("i%d", j),
				VariableType: &types.TypeArgument{Name: fmt.Sprintf("I%d", j)},
			})
			argumentsAfterBuild = append(argumentsAfterBuild, types.FunctionArgument{
				Name:         fmt.Sprintf("fromJsonI%d", j),
				VariableType: tenecs_json_FromJsonField_Of(&types.TypeArgument{Name: fmt.Sprintf("I%d", j)}),
			})
		}

		result = append(result, NamedFunction{
			name: fmt.Sprintf("parseObject%d", i),
			function: &types.Function{
				Generics: generics,
				Arguments: append([]types.FunctionArgument{
					types.FunctionArgument{
						Name: "build",
						VariableType: &types.Function{
							Arguments:  buildArguments,
							ReturnType: &types.TypeArgument{Name: "R"},
						},
					},
				}, argumentsAfterBuild...),
				ReturnType: tenecs_json_FromJson_Of(&types.TypeArgument{Name: "R"}),
			},
		})
	}
	return result
}()

var tenecs_json_parseOr = &types.Function{
	Generics: []string{"A", "B"},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "fromA",
			VariableType: tenecs_json_FromJson_Of(&types.TypeArgument{Name: "A"}),
		},
		types.FunctionArgument{
			Name:         "fromB",
			VariableType: tenecs_json_FromJson_Of(&types.TypeArgument{Name: "B"}),
		},
	},
	ReturnType: tenecs_json_FromJson_Of(&types.OrVariableType{
		Elements: []types.VariableType{
			&types.TypeArgument{Name: "A"},
			&types.TypeArgument{Name: "B"},
		},
	}),
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
