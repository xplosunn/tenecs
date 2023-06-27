package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_compare = packageWith(
	withFunction("eq", tenecs_compare_eq),
)

var tenecs_compare_eq = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "first",
			VariableType: &types.TypeArgument{
				Name: "T",
			},
		},
		types.FunctionArgument{
			Name: "second",
			VariableType: &types.TypeArgument{
				Name: "T",
			},
		},
	},
	ReturnType: &types.Array{
		OfType: &types.BasicType{
			Type: "Boolean",
		},
	},
}
