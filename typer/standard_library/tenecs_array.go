package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_array = packageWith(
	withFunction("emptyArray", tenecs_array_emptyArray),
	withFunction("append", tenecs_array_append),
	withFunction("map", tenecs_array_map),
	withFunction("repeat", tenecs_array_repeat),
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

var tenecs_array_append = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: &types.Array{
				OfType: &types.TypeArgument{
					Name: "T",
				},
			},
		},
		types.FunctionArgument{
			Name: "newElement",
			VariableType: &types.TypeArgument{
				Name: "T",
			},
		},
	},
	ReturnType: &types.Array{
		OfType: &types.TypeArgument{
			Name: "T",
		},
	},
}

var tenecs_array_map = &types.Function{
	Generics: []string{
		"A",
		"B",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: &types.Array{
				OfType: &types.TypeArgument{
					Name: "A",
				},
			},
		},
		types.FunctionArgument{
			Name: "f",
			VariableType: &types.Function{
				Arguments: []types.FunctionArgument{
					types.FunctionArgument{
						Name: "a",
						VariableType: &types.TypeArgument{
							Name: "A",
						},
					},
				},
				ReturnType: &types.TypeArgument{
					Name: "B",
				},
			},
		},
	},
	ReturnType: &types.Array{
		OfType: &types.TypeArgument{
			Name: "B",
		},
	},
}

var tenecs_array_repeat = &types.Function{
	Generics: []string{
		"A",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "elem",
			VariableType: &types.TypeArgument{
				Name: "A",
			},
		},
		types.FunctionArgument{
			Name: "times",
			VariableType: &types.BasicType{
				Type: "Int",
			},
		},
	},
	ReturnType: &types.Array{
		OfType: &types.TypeArgument{
			Name: "A",
		},
	},
}
