package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_array = packageWith(
	withFunction("append", tenecs_array_append),
	withStruct("Break", tenecs_array_Break, tenecs_array_Break_Fields...),
	withFunction("filter", tenecs_array_filter),
	withFunction("flatMap", tenecs_array_flatMap),
	withFunction("fold", tenecs_array_fold),
	withFunction("forEach", tenecs_array_forEach),
	withFunction("length", tenecs_array_length),
	withFunction("map", tenecs_array_map),
	withFunction("mapUntil", tenecs_array_mapUntil),
	withFunction("mapNotNull", tenecs_array_mapNotNull),
	withFunction("repeat", tenecs_array_repeat),
)

var tenecs_array_append = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.Array(&types.TypeArgument{
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
	ReturnType: types.Array(&types.TypeArgument{
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
			VariableType: types.Array(&types.TypeArgument{
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
	ReturnType: types.Array(&types.TypeArgument{
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
			VariableType: types.Array(&types.TypeArgument{
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
				ReturnType: types.Array(&types.TypeArgument{
					Name: "B",
				}),
			},
		},
	},
	ReturnType: types.Array(&types.TypeArgument{
		Name: "B",
	}),
}

var tenecs_array_fold = &types.Function{
	Generics: []string{
		"A",
		"Acc",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.Array(&types.TypeArgument{
				Name: "A",
			}),
		},
		types.FunctionArgument{
			Name: "zero",
			VariableType: &types.TypeArgument{
				Name: "Acc",
			},
		},
		types.FunctionArgument{
			Name: "f",
			VariableType: &types.Function{
				Arguments: []types.FunctionArgument{
					types.FunctionArgument{
						Name: "acc",
						VariableType: &types.TypeArgument{
							Name: "Acc",
						},
					},
					types.FunctionArgument{
						Name: "a",
						VariableType: &types.TypeArgument{
							Name: "A",
						},
					},
				},
				ReturnType: &types.TypeArgument{
					Name: "Acc",
				},
			},
		},
	},
	ReturnType: &types.TypeArgument{
		Name: "Acc",
	},
}
var tenecs_array_forEach = &types.Function{
	Generics: []string{
		"A",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.Array(&types.TypeArgument{
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
				ReturnType: types.Void(),
			},
		},
	},
	ReturnType: &types.TypeArgument{
		Name: "Acc",
	},
}

var tenecs_array_length = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.Array(&types.TypeArgument{
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
			VariableType: types.Array(&types.TypeArgument{
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
	ReturnType: types.Array(&types.TypeArgument{
		Name: "B",
	}),
}

var tenecs_array_mapUntil = &types.Function{
	Generics: []string{
		"A",
		"B",
		"S", // same as Break
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.Array(&types.TypeArgument{
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
				ReturnType: &types.OrVariableType{
					Elements: []types.VariableType{
						tenecs_array_Break,
						&types.TypeArgument{
							Name: "B",
						},
					},
				},
			},
		},
	},
	ReturnType: &types.OrVariableType{
		Elements: []types.VariableType{
			&types.TypeArgument{
				Name: "S",
			},
			types.Array(&types.TypeArgument{
				Name: "B",
			}),
		},
	},
}

var tenecs_array_mapNotNull = &types.Function{
	Generics: []string{
		"A",
		"B",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "array",
			VariableType: types.Array(&types.TypeArgument{
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
				ReturnType: &types.OrVariableType{
					Elements: []types.VariableType{
						&types.TypeArgument{
							Name: "B",
						},
						types.Void(),
					},
				},
			},
		},
	},
	ReturnType: types.Array(&types.TypeArgument{
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
	ReturnType: types.Array(&types.TypeArgument{
		Name: "A",
	}),
}

var tenecs_array_Break = types.Struct("tenecs.array", "Break", []string{"S"})

var tenecs_array_Break_Fields = []func(fields *StructWithFields){
	structField("values", &types.TypeArgument{Name: "S"}),
}
