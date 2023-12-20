package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_array = packageWith(
	withFunction("append", tenecs_array_append),
	withFunction("filter", tenecs_array_filter),
	withFunction("flatMap", tenecs_array_flatMap),
	withFunction("length", tenecs_array_length),
	withFunction("map", tenecs_array_map),
	withFunction("repeat", tenecs_array_repeat),
)

var tenecs_array_append = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.UncheckedArray(&types.TypeArgument{
				Name: "T",
			}),
		},
		types.FunctionArgument{
			Name: "newElement",
			VariableType: &types.TypeArgument{
				Name: "T",
			},
		},
	},
	ReturnType: types.UncheckedArray(&types.TypeArgument{
		Name: "T",
	}),
}

var tenecs_array_filter = &types.Function{
	Generics: []string{
		"A",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.UncheckedArray(&types.TypeArgument{
				Name: "A",
			}),
		},
		types.FunctionArgument{
			Name: "keep",
			VariableType: &types.Function{
				Arguments: []types.FunctionArgument{
					types.FunctionArgument{
						Name: "a",
						VariableType: &types.TypeArgument{
							Name: "A",
						},
					},
				},
				ReturnType: types.Boolean(),
			},
		},
	},
	ReturnType: types.UncheckedArray(&types.TypeArgument{
		Name: "A",
	}),
}

var tenecs_array_flatMap = &types.Function{
	Generics: []string{
		"A",
		"B",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.UncheckedArray(&types.TypeArgument{
				Name: "A",
			}),
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
				ReturnType: types.UncheckedArray(&types.TypeArgument{
					Name: "B",
				}),
			},
		},
	},
	ReturnType: types.UncheckedArray(&types.TypeArgument{
		Name: "B",
	}),
}

var tenecs_array_length = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.UncheckedArray(&types.TypeArgument{
				Name: "T",
			}),
		},
	},
	ReturnType: types.Int(),
}

var tenecs_array_map = &types.Function{
	Generics: []string{
		"A",
		"B",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.UncheckedArray(&types.TypeArgument{
				Name: "A",
			}),
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
	ReturnType: types.UncheckedArray(&types.TypeArgument{
		Name: "B",
	}),
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
			Name:         "times",
			VariableType: types.Int(),
		},
	},
	ReturnType: types.UncheckedArray(&types.TypeArgument{
		Name: "A",
	}),
}
