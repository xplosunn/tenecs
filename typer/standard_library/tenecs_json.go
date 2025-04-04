package standard_library

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/types"
)

var tenecs_json = packageWith(
	withStruct(Tenecs_json_JsonConverter),
	withStruct(Tenecs_json_JsonField),
	withFunction("jsonList", tenecs_json_jsonList),
	withFunction("jsonBoolean", tenecs_json_jsonBoolean),
	withFunction("jsonInt", tenecs_json_jsonInt),
	withFunction("jsonObject0", tenecs_json_jsonObject0),
	withFunctions(tenecs_json_jsonObject),
	withFunction("jsonOr", tenecs_json_jsonOr),
	withFunction("jsonString", tenecs_json_jsonString),
)

var Tenecs_json_JsonConverter = structWithFields("JsonConverter", tenecs_json_JsonConverter, tenecs_json_FromJson_Fields...)

var tenecs_json_JsonConverter = types.Struct(
	"tenecs.json",
	"JsonConverter",
	[]string{"T"},
)

func tenecs_json_JsonConverter_Of(varType types.VariableType) *types.KnownType {
	return types.UncheckedApplyGenerics(tenecs_json_JsonConverter, []types.VariableType{varType})
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
	structField("Converter", tenecs_json_JsonConverter_Of(&types.TypeArgument{Name: "Field"})),
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

var tenecs_json_jsonList = functionFromType("<T>(of: JsonConverter<T>) ~> JsonConverter<List<T>>", Tenecs_json_JsonConverter)

var tenecs_json_jsonBoolean = functionFromType("() ~> JsonConverter<Boolean>", Tenecs_json_JsonConverter)

var tenecs_json_jsonInt = functionFromType("() ~> JsonConverter<Int>", Tenecs_json_JsonConverter)

var tenecs_json_jsonObject0 = functionFromType("<R>(build: () ~> R) ~> JsonConverter<R>", Tenecs_json_JsonConverter)

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
				Name: fmt.Sprintf("jsonConverterFieldI%d", j),
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
				ReturnType: tenecs_json_JsonConverter_Of(&types.TypeArgument{Name: "R"}),
			},
		})
	}
	return result
}()

var tenecs_json_jsonOr = functionFromType("<A, B>(ConverterA: JsonConverter<A>, ConverterB: JsonConverter<B>, toJsonConverterPicker: (A | B) ~> JsonConverter<A> | JsonConverter<B>) ~> JsonConverter<A | B>", Tenecs_json_JsonConverter)

var tenecs_json_jsonString = functionFromType("() ~> JsonConverter<String>", Tenecs_json_JsonConverter)
