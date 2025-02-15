package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_int = packageWith(
	withFunction("div", tenecs_int_div),
	withFunction("minus", tenecs_int_minus),
	withFunction("mod", tenecs_int_mod),
	withFunction("plus", tenecs_int_plus),
	withFunction("ponyDiv", tenecs_int_ponyDiv),
	withFunction("ponyMod", tenecs_int_ponyMod),
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

var tenecs_int_div = &types.Function{
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
	ReturnType: &types.OrVariableType{
		Elements: []types.VariableType{
			types.Int(),
			tenecs_error_Error,
		},
	},
}

var tenecs_int_ponyDiv = &types.Function{
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

var tenecs_int_mod = &types.Function{
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
	ReturnType: &types.OrVariableType{
		Elements: []types.VariableType{
			types.Int(),
			tenecs_error_Error,
		},
	},
}

var tenecs_int_ponyMod = &types.Function{
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
