package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_test = packageWith(
	withInterface("Assert", &tenecs_test_Assert),
	withInterface("UnitTestRegistry", &tenecs_test_UnitTestRegistry),
	withInterface("UnitTests", &tenecs_test_UnitTests),
)

var tenecs_test_Assert = types.Interface{
	Package: "tenecs.test",
	Name:    "Assert",
	Variables: map[string]types.VariableType{
		"equal": &types.Function{
			Generics: []string{"T"},
			Arguments: []types.FunctionArgument{
				{
					Name:         "value",
					VariableType: &types.TypeArgument{Name: "T"},
				},
				{
					Name:         "expected",
					VariableType: &types.TypeArgument{Name: "T"},
				},
			},
			ReturnType: &Void,
		},
	},
}

var tenecs_test_UnitTestRegistry = types.Interface{
	Package: "tenecs.test",
	Name:    "UnitTestRegistry",
	Variables: map[string]types.VariableType{
		"test": &types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "name",
					VariableType: &BasicTypeString,
				},
				{
					Name: "theTest",
					VariableType: &types.Function{
						Arguments: []types.FunctionArgument{
							{
								Name:         "assert",
								VariableType: &tenecs_test_Assert,
							},
						},
						ReturnType: &Void,
					},
				},
			},
			ReturnType: &Void,
		},
	},
}

var tenecs_test_UnitTests = types.Interface{
	Package: "tenecs.test",
	Name:    "UnitTests",
	Variables: map[string]types.VariableType{
		"tests": &types.Function{
			Arguments: []types.FunctionArgument{
				{
					Name:         "registry",
					VariableType: &tenecs_test_UnitTestRegistry,
				},
			},
			ReturnType: &Void,
		},
	},
}
