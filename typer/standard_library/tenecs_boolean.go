package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_boolean = packageWith(
	withFunction("not", tenecs_boolean_not),
)

var tenecs_boolean_not = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "b",
			VariableType: types.Boolean(),
		},
	},
	ReturnType: types.Boolean(),
}
