package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_string = packageWith(
	withFunction("join", tenecs_string_join),
	withFunction("hasPrefix", tenecs_string_hasPrefix),
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

var tenecs_string_hasPrefix = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "str",
			VariableType: types.String(),
		},
		types.FunctionArgument{
			Name:         "prefix",
			VariableType: types.String(),
		},
	},
	ReturnType: types.Boolean(),
}
