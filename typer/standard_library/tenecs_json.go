package standard_library

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/types"
)

var tenecs_json = packageWith(
	withStruct(Tenecs_json_JsonSchema),
	withStruct(Tenecs_json_JsonField),
	withFunction("jsonList", tenecs_json_jsonList),
	withFunction("jsonBoolean", tenecs_json_jsonBoolean),
	withFunction("jsonInt", tenecs_json_jsonInt),
	withFunction("jsonObject0", tenecs_json_jsonObject0),
	withFunctions(tenecs_json_jsonObject),
	withFunction("jsonOr", tenecs_json_jsonOr),
	withFunction("jsonString", tenecs_json_jsonString),
)

var Tenecs_json_JsonSchema = structWithFields("JsonSchema", tenecs_json_JsonSchema, tenecs_json_FromJson_Fields...)

var tenecs_json_JsonSchema = types.Struct(
	"tenecs.json",
	"JsonSchema",
	[]string{"T"},
)

func tenecs_json_JsonSchema_Of(varType types.VariableType) *types.KnownType {
	return types.UncheckedApplyGenerics(tenecs_json_JsonSchema, []types.VariableType{varType})
}

var tenecs_json_FromJson_Fields = []func(fields *StructWithFields){
	structField("fromJson", &types.Function{
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
				tenecs_error_Error,
			},
		},
	}),
	structField("toJson", &types.Function{
		Arguments: []types.FunctionArgument{
			types.FunctionArgument{
				Name: "value",
				VariableType: &types.TypeArgument{
					Name: "T",
				},
			},
		},
		ReturnType: types.String(),
	}),
}

var Tenecs_json_JsonField = structWithFields("JsonField", tenecs_json_JsonField, tenecs_json_JsonField_Fields...)

var tenecs_json_JsonField = types.Struct(
	"tenecs.json",
	"JsonField",
	[]string{"Record", "Field"},
)

func tenecs_json_JsonField_Of(recordVarType types.VariableType, fieldVarType types.VariableType) *types.KnownType {
	return types.UncheckedApplyGenerics(tenecs_json_JsonField, []types.VariableType{recordVarType, fieldVarType})
}

var tenecs_json_JsonField_Fields = []func(fields *StructWithFields){
	structField("name", types.String()),
	structField("schema", tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "Field"})),
	structField("access", &types.Function{
		Arguments: []types.FunctionArgument{
			types.FunctionArgument{
				Name:         "record",
				VariableType: &types.TypeArgument{Name: "Record"},
			},
		},
		ReturnType: &types.TypeArgument{Name: "Field"},
	}),
}

var tenecs_json_jsonList = &types.Function{
	Generics: []string{"T"},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "of",
			VariableType: tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "T"}),
		},
	},
	ReturnType: tenecs_json_JsonSchema_Of(types.List(&types.TypeArgument{Name: "T"})),
}

var tenecs_json_jsonBoolean = &types.Function{
	Arguments:  []types.FunctionArgument{},
	ReturnType: tenecs_json_JsonSchema_Of(types.Boolean()),
}

var tenecs_json_jsonInt = &types.Function{
	Arguments:  []types.FunctionArgument{},
	ReturnType: tenecs_json_JsonSchema_Of(types.Int()),
}

var tenecs_json_jsonObject0 = &types.Function{
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
	ReturnType: tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "R"}),
}

var tenecs_json_jsonObject = func() []NamedFunction {
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
				Name: fmt.Sprintf("jsonSchemaFieldI%d", j),
				VariableType: tenecs_json_JsonField_Of(
					&types.TypeArgument{Name: "R"},
					&types.TypeArgument{Name: fmt.Sprintf("I%d", j)},
				),
			})
		}

		result = append(result, NamedFunction{
			name: fmt.Sprintf("jsonObject%d", i),
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
				ReturnType: tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "R"}),
			},
		})
	}
	return result
}()

var tenecs_json_jsonOr = &types.Function{
	Generics: []string{"A", "B"},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "schemaA",
			VariableType: tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "A"}),
		},
		types.FunctionArgument{
			Name:         "schemaB",
			VariableType: tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "B"}),
		},
		types.FunctionArgument{
			Name: "toJsonSchemaPicker",
			VariableType: &types.Function{
				Arguments: []types.FunctionArgument{
					types.FunctionArgument{
						Name: "either",
						VariableType: &types.OrVariableType{
							Elements: []types.VariableType{
								&types.TypeArgument{Name: "A"},
								&types.TypeArgument{Name: "B"},
							},
						},
					},
				},
				ReturnType: &types.OrVariableType{
					Elements: []types.VariableType{
						tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "A"}),
						tenecs_json_JsonSchema_Of(&types.TypeArgument{Name: "B"}),
					},
				},
			},
		},
	},
	ReturnType: tenecs_json_JsonSchema_Of(&types.OrVariableType{
		Elements: []types.VariableType{
			&types.TypeArgument{Name: "A"},
			&types.TypeArgument{Name: "B"},
		},
	}),
}

var tenecs_json_jsonString = &types.Function{
	Arguments:  []types.FunctionArgument{},
	ReturnType: tenecs_json_JsonSchema_Of(types.String()),
}
