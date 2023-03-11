package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_string = packageWith(
	withFunction("join", tenecs_string_join),
)

var tenecs_string_join = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "left",
			VariableType: &BasicTypeString,
		},
		types.FunctionArgument{
			Name:         "right",
			VariableType: &BasicTypeString,
		},
	},
	ReturnType: &BasicTypeString,
}
