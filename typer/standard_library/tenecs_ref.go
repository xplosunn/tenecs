package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_ref = packageWith(
	withStruct(Tenecs_ref_Ref),
	withStruct(Tenecs_ref_RefCreator),
)

var Tenecs_ref_Ref = structWithFields("Ref", tenecs_ref_Ref, tenecs_ref_Ref_Fields...)

var tenecs_ref_Ref = types.Struct(
	"tenecs.ref",
	"Ref",
	[]string{"T"},
)

var tenecs_ref_Ref_Fields = []func(fields *StructWithFields){
	structField("get", &types.Function{
		Arguments:  []types.FunctionArgument{},
		ReturnType: &types.TypeArgument{Name: "T"},
	}),
	structField("set", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "value",
				VariableType: &types.TypeArgument{Name: "T"},
			},
		},
		ReturnType: types.Void(),
	}),
	structField("modify", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name: "f",
				VariableType: &types.Function{
					Generics: []string{},
					Arguments: []types.FunctionArgument{
						types.FunctionArgument{
							Name:         "value",
							VariableType: &types.TypeArgument{Name: "T"},
						},
					},
					ReturnType: &types.TypeArgument{Name: "T"},
				},
			},
		},
		ReturnType: types.Void(),
	}),
}

var Tenecs_ref_RefCreator = structWithFields("RefCreator", tenecs_ref_RefCreator, tenecs_ref_RefCreator_Fields...)

var tenecs_ref_RefCreator = types.Struct(
	"tenecs.ref",
	"RefCreator",
	nil,
)

var tenecs_ref_RefCreator_Fields = []func(fields *StructWithFields){
	structField("new", &types.Function{
		Generics: []string{"T"},
		Arguments: []types.FunctionArgument{
			{
				Name:         "value",
				VariableType: &types.TypeArgument{Name: "T"},
			},
		},
		ReturnType: tenecs_ref_Ref,
	}),
}
