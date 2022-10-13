package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_boolean = packageWith(
	withFunction("and", tenecs_boolean_and),
	withFunction("not", tenecs_boolean_not),
	withFunction("or", tenecs_boolean_or),
)

var tenecs_boolean_and = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "a",
			VariableType: types.Boolean(),
		},
		types.FunctionArgument{
			Name:         "b",
			VariableType: types.Boolean(),
		},
	},
	ReturnType: types.Boolean(),
}

var tenecs_boolean_not = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "b",
			VariableType: types.Boolean(),
		},
	},
	ReturnType: types.Boolean(),
}

var tenecs_boolean_or = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "a",
			VariableType: types.Boolean(),
		},
		types.FunctionArgument{
			Name:         "b",
			VariableType: types.Boolean(),
		},
	},
	ReturnType: types.Boolean(),
}
