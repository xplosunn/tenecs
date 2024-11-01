package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_test = packageWith(
	withStruct(Tenecs_test_Assert),
	withStruct(Tenecs_test_UnitTest),
	withStruct(Tenecs_test_UnitTestKit),
	withStruct(Tenecs_test_UnitTestRegistry),
	withStruct(Tenecs_test_UnitTestSuite),
)

var Tenecs_test_Assert = structWithFields("Assert", &tenecs_test_Assert, tenecs_test_Assert_Fields...)

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

var Tenecs_test_UnitTestKit = structWithFields("UnitTestKit", &tenecs_test_UnitTestKit, tenecs_test_UnitTestKit_Fields...)

var tenecs_test_UnitTestKit = types.KnownType{
	Package: "tenecs.test",
	Name:    "UnitTestKit",
}

var tenecs_test_UnitTestKit_Fields = []func(fields *StructWithFields){
	structField("assert", &tenecs_test_Assert),
	structField("ref", tenecs_ref_RefCreator),
}

var Tenecs_test_UnitTestRegistry = structWithFields("UnitTestRegistry", &tenecs_test_UnitTestRegistry, tenecs_test_UnitTestRegistry_Fields...)

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

var Tenecs_test_UnitTest = structWithFields("UnitTest", &tenecs_test_UnitTest, tenecs_test_UnitTest_Fields...)

var tenecs_test_UnitTest = types.KnownType{
	Package: "tenecs.test",
	Name:    "UnitTest",
}

var tenecs_test_UnitTest_Fields = []func(fields *StructWithFields){
	structField("name", types.String()),
	structField("theTest", &types.Function{
		Arguments: []types.FunctionArgument{
			{
				Name:         "testkit",
				VariableType: &tenecs_test_UnitTestKit,
			},
		},
		ReturnType: types.Void(),
	}),
}

var Tenecs_test_UnitTestSuite = structWithFields("UnitTestSuite", &tenecs_test_UnitTestSuite, tenecs_test_UnitTestSuite_Fields...)

var tenecs_test_UnitTestSuite = types.KnownType{
	Package: "tenecs.test",
	Name:    "UnitTestSuite",
}

var tenecs_test_UnitTestSuite_Fields = []func(fields *StructWithFields){
	structField("name", types.String()),
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
