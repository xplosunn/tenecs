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
			VariableType: &BasicTypeInt,
		},
		types.FunctionArgument{
			Name:         "b",
			VariableType: &BasicTypeInt,
		},
	},
	ReturnType: &BasicTypeInt,
}

var tenecs_int_plus = &types.Function{
	Generics: []string{},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "a",
			VariableType: &BasicTypeInt,
		},
		types.FunctionArgument{
			Name:         "b",
			VariableType: &BasicTypeInt,
		},
	},
	ReturnType: &BasicTypeInt,
}

var tenecs_int_times = &types.Function{
	Generics: []string{},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "a",
			VariableType: &BasicTypeInt,
		},
		types.FunctionArgument{
			Name:         "b",
			VariableType: &BasicTypeInt,
		},
	},
	ReturnType: &BasicTypeInt,
}