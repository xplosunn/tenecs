package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_int = packageWith(
	withFunction("minus", tenecs_int_minus),
	withFunction("plus", tenecs_int_plus),
	withFunction("times", tenecs_int_times),
)

var tenecs_int_minus = &types.Function{
	Generics: []string{},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "a",
			VariableType: types.Int(),
		},
		types.FunctionArgument{
			Name:         "b",
			VariableType: types.Int(),
		},
	},
	ReturnType: types.Int(),
}

var tenecs_int_plus = &types.Function{
	Generics: []string{},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "a",
			VariableType: types.Int(),
		},
		types.FunctionArgument{
			Name:         "b",
			VariableType: types.Int(),
		},
	},
	ReturnType: types.Int(),
}

var tenecs_int_times = &types.Function{
	Generics: []string{},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "a",
			VariableType: types.Int(),
		},
		types.FunctionArgument{
			Name:         "b",
			VariableType: types.Int(),
		},
	},
	ReturnType: types.Int(),
}
