package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_ref = packageWith(
	withInterface("Ref", tenecs_ref_Ref, tenecs_ref_Ref_Fields),
	withInterface("RefCreator", tenecs_ref_RefCreator, tenecs_ref_RefCreator_Fields),
)

var tenecs_ref_Ref = types.Interface(
	"tenecs.ref",
	"Ref",
	[]string{"T"},
)

var tenecs_ref_Ref_Fields = map[string]types.VariableType{
	"get": &types.Function{
		Arguments:  []types.FunctionArgument{},
		ReturnType: &types.TypeArgument{Name: "T"},
	},
	"set": &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "value",
				VariableType: &types.TypeArgument{Name: "T"},
			},
		},
		ReturnType: types.Void(),
	},
}

var tenecs_ref_RefCreator = types.Interface(
	"tenecs.ref",
	"RefCreator",
	nil,
)

var tenecs_ref_RefCreator_Fields = map[string]types.VariableType{
	"new": &types.Function{
		Generics: []string{"T"},
		Arguments: []types.FunctionArgument{
			{
				Name:         "value",
				VariableType: &types.TypeArgument{Name: "T"},
			},
		},
		ReturnType: tenecs_ref_Ref,
	},
}
