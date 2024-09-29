package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_test = packageWith(
	withStruct("Assert", &tenecs_test_Assert, tenecs_test_Assert_Fields...),
	withStruct("UnitTestKit", &tenecs_test_UnitTestKit, tenecs_test_UnitTestKit_Fields...),
	withStruct("UnitTestRegistry", &tenecs_test_UnitTestRegistry, tenecs_test_UnitTestRegistry_Fields...),
	withStruct("UnitTests", &tenecs_test_UnitTests, tenecs_test_UnitTests_Fields...),
)

var tenecs_test_Assert = types.KnownType{
	Package: "tenecs.test",
	Name:    "Assert",
}

var tenecs_test_Assert_Fields = []func(fields *StructWithFields){
	structField("equal", &types.Function{
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
		ReturnType: types.Void(),
	}),
	structField("fail", &types.Function{
		Generics: []string{"T"},
		Arguments: []types.FunctionArgument{
			{
				Name:         "message",
				VariableType: types.String(),
			},
		},
		ReturnType: &types.TypeArgument{Name: "T"},
	}),
}

var tenecs_test_UnitTestKit = types.KnownType{
	Package: "tenecs.test",
	Name:    "UnitTestKit",
}

var tenecs_test_UnitTestKit_Fields = []func(fields *StructWithFields){
	structField("assert", &tenecs_test_Assert),
	structField("runtime", &tenecs_os_Runtime),
}

var tenecs_test_UnitTestRegistry = types.KnownType{
	Package: "tenecs.test",
	Name:    "UnitTestRegistry",
}

var tenecs_test_UnitTestRegistry_Fields = []func(fields *StructWithFields){
	structField("test", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "name",
				VariableType: types.String(),
			},
			{
				Name: "theTest",
				VariableType: &types.Function{
					Arguments: []types.FunctionArgument{
						{
							Name:         "testkit",
							VariableType: &tenecs_test_UnitTestKit,
						},
					},
					ReturnType: types.Void(),
				},
			},
		},
		ReturnType: types.Void(),
	}),
}

var tenecs_test_UnitTests = types.KnownType{
	Package: "tenecs.test",
	Name:    "UnitTests",
}

var tenecs_test_UnitTests_Fields = []func(fields *StructWithFields){
	structField("tests", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "registry",
				VariableType: &tenecs_test_UnitTestRegistry,
			},
		},
		ReturnType: types.Void(),
	}),
}
