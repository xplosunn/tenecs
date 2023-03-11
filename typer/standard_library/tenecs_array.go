package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_array = packageWith(
	withFunction("emptyArray", tenecs_array_emptyArray),
	withFunction("append", tenecs_array_append),
)

var tenecs_array_emptyArray = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{},
	ReturnType: &types.Array{
		OfType: &types.TypeArgument{
			Name: "T",
		},
	},
}

var tenecs_array_append = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: &types.Array{
				OfType: &types.TypeArgument{
					Name: "T",
				},
			},
		},
		types.FunctionArgument{
			Name: "newElement",
			VariableType: &types.TypeArgument{
				Name: "T",
			},
		},
	},
	ReturnType: &types.Array{
		OfType: &types.TypeArgument{
			Name: "T",
		},
	},
}
