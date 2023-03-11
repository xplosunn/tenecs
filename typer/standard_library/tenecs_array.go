package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_array = packageWith(
	withFunction("emptyArray", tenecs_array_emptyArray),
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
