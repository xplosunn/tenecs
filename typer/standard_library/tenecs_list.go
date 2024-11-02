package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_list = packageWith(
	withFunction("append", tenecs_list_append),
	withStruct(Tenecs_list_Break),
	withFunction("filter", tenecs_list_filter),
	withFunction("flatMap", tenecs_list_flatMap),
	withFunction("fold", tenecs_list_fold),
	withFunction("forEach", tenecs_list_forEach),
	withFunction("length", tenecs_list_length),
	withFunction("map", tenecs_list_map),
	withFunction("mapUntil", tenecs_list_mapUntil),
	withFunction("mapNotNull", tenecs_list_mapNotNull),
	withFunction("repeat", tenecs_list_repeat),
)

var tenecs_list_append = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "list",
			VariableType: types.List(&types.TypeArgument{
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
	ReturnType: types.List(&types.TypeArgument{
		Name: "T",
	}),
}

var tenecs_list_filter = &types.Function{
	Generics: []string{
		"A",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "list",
			VariableType: types.List(&types.TypeArgument{
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
	ReturnType: types.List(&types.TypeArgument{
		Name: "A",
	}),
}

var tenecs_list_flatMap = &types.Function{
	Generics: []string{
		"A",
		"B",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "list",
			VariableType: types.List(&types.TypeArgument{
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
				ReturnType: types.List(&types.TypeArgument{
					Name: "B",
				}),
			},
		},
	},
	ReturnType: types.List(&types.TypeArgument{
		Name: "B",
	}),
}

var tenecs_list_fold = &types.Function{
	Generics: []string{
		"A",
		"Acc",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "list",
			VariableType: types.List(&types.TypeArgument{
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
var tenecs_list_forEach = &types.Function{
	Generics: []string{
		"A",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "list",
			VariableType: types.List(&types.TypeArgument{
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
	ReturnType: types.Void(),
}

var tenecs_list_length = &types.Function{
	Generics: []string{
		"T",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "list",
			VariableType: types.List(&types.TypeArgument{
				Name: "T",
			}),
		},
	},
	ReturnType: types.Int(),
}

var tenecs_list_map = &types.Function{
	Generics: []string{
		"A",
		"B",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "list",
			VariableType: types.List(&types.TypeArgument{
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
	ReturnType: types.List(&types.TypeArgument{
		Name: "B",
	}),
}

var tenecs_list_mapUntil = &types.Function{
	Generics: []string{
		"A",
		"B",
		"S", // same as Break
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "list",
			VariableType: types.List(&types.TypeArgument{
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
						tenecs_list_Break,
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
			types.List(&types.TypeArgument{
				Name: "B",
			}),
		},
	},
}

var tenecs_list_mapNotNull = &types.Function{
	Generics: []string{
		"A",
		"B",
	},
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name: "list",
			VariableType: types.List(&types.TypeArgument{
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
	ReturnType: types.List(&types.TypeArgument{
		Name: "B",
	}),
}

var tenecs_list_repeat = &types.Function{
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
	ReturnType: types.List(&types.TypeArgument{
		Name: "A",
	}),
}

var Tenecs_list_Break = structWithFields("Break", tenecs_list_Break, tenecs_list_Break_Fields...)

var tenecs_list_Break = types.Struct("tenecs.list", "Break", []string{"S"})

var tenecs_list_Break_Fields = []func(fields *StructWithFields){
	structField("value", &types.TypeArgument{Name: "S"}),
}
