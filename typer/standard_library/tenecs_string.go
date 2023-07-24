package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_string = packageWith(
	withFunction("join", tenecs_string_join),
)

var tenecs_string_join = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "left",
			VariableType: types.String(),
		},
		types.FunctionArgument{
			Name:         "right",
			VariableType: types.String(),
		},
	},
	ReturnType: types.String(),
}
