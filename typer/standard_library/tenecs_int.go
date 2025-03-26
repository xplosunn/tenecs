package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_int = packageWith(
	withFunction("abs", tenecs_int_abs),
	withFunction("div", tenecs_int_div),
	withFunction("greaterThan", tenecs_int_greaterThan),
	withFunction("lessThan", tenecs_int_lessThan),
	withFunction("minus", tenecs_int_minus),
	withFunction("mod", tenecs_int_mod),
	withFunction("negate", tenecs_int_negate),
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

var tenecs_int_greaterThan = &types.Function{
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
	ReturnType: types.Boolean(),
}

var tenecs_int_lessThan = &types.Function{
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
	ReturnType: types.Boolean(),
}

var tenecs_int_abs = &types.Function{
	Generics: []string{},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "a",
			VariableType: types.Int(),
		},
	},
	ReturnType: types.Int(),
}

var tenecs_int_negate = &types.Function{
	Generics: []string{},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "a",
			VariableType: types.Int(),
		},
	},
	ReturnType: types.Int(),
}