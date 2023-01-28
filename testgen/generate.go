package testgen

import (
	"errors"
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"regexp"
)

func Generate(program ast.Program, constructorName string, targetFunctionName string) (*parser.Module, error) {
	module, err := findModuleInProgram(program, constructorName)
	if err != nil {
		return nil, err
	}
	targetFunction, err := findFunctionInModule(module, targetFunctionName)
	if err != nil {
		return nil, err
	}
	testCases, err := generateTestCases(targetFunction)

	singleTypeNameToTypeAnnotation := func(typeName string) *parser.TypeAnnotation {
		var typeAnnotation parser.TypeAnnotation = parser.SingleNameType{
			TypeName: "UnitTestRegistry",
		}
		return &typeAnnotation
	}

	testsBlock := []parser.ExpressionBox{}

	for _, testCase := range testCases {
		testsBlock = append(testsBlock, parser.ExpressionBox{
			Expression: parser.ReferenceOrInvocation{
				Var:       "registry",
				Arguments: nil,
			},
			AccessOrInvocationChain: []parser.AccessOrInvocation{
				{
					VarName: "test",
					Arguments: &parser.ArgumentsList{
						Arguments: []parser.ExpressionBox{
							{
								Expression: parser.LiteralExpression{
									Literal: parser.LiteralString{
										Value: testCase.name,
									},
								},
							},
							{
								Expression: parser.ReferenceOrInvocation{
									Var: nameOfFunctionForTestCase(testCase),
								},
							},
						},
					},
				},
			},
		})
	}

	declarations := []parser.ModuleDeclaration{
		{
			Public: true,
			Name:   "tests",
			Expression: parser.Lambda{
				Parameters: []parser.Parameter{
					{
						Name: "registry",
						Type: singleTypeNameToTypeAnnotation("UnitTestRegistry"),
					},
				},
				ReturnType: singleTypeNameToTypeAnnotation("Void"),
				Block:      testsBlock,
			},
		},
	}

	for _, testCase := range testCases {
		block := []parser.ExpressionBox{}

		block = append(block, parser.ExpressionBox{
			Expression: parser.Declaration{
				Name: "module",
				ExpressionBox: parser.ExpressionBox{
					Expression: parser.ReferenceOrInvocation{
						Var: constructorName,
						Arguments: &parser.ArgumentsList{
							Arguments: []parser.ExpressionBox{},
						},
					},
				},
			},
		})

		block = append(block, parser.ExpressionBox{
			Expression: parser.Declaration{
				Name: "result",
				ExpressionBox: parser.ExpressionBox{
					Expression: parser.ReferenceOrInvocation{
						Var: "module",
					},
					AccessOrInvocationChain: []parser.AccessOrInvocation{
						{
							VarName: targetFunctionName,
							Arguments: &parser.ArgumentsList{
								Arguments: testCase.functionArguments,
							},
						},
					},
				},
			},
		})

		block = append(block, parser.ExpressionBox{
			Expression: parser.Declaration{
				Name: "expected",
				ExpressionBox: parser.ExpressionBox{
					Expression: testCase.expectedOutput,
				},
			},
		})

		block = append(block, parser.ExpressionBox{
			Expression: parser.ReferenceOrInvocation{
				Var: "assert",
			},
			AccessOrInvocationChain: []parser.AccessOrInvocation{
				{
					VarName: "equal",
					Arguments: &parser.ArgumentsList{
						Generics: []string{testCase.expectedOutputType},
						Arguments: []parser.ExpressionBox{
							{
								Expression: parser.ReferenceOrInvocation{
									Var: "result",
								},
							},
							{
								Expression: parser.ReferenceOrInvocation{
									Var: "expected",
								},
							},
						},
					},
				},
			},
		})

		declarations = append(declarations, parser.ModuleDeclaration{
			Name: nameOfFunctionForTestCase(testCase),
			Expression: parser.Lambda{
				Parameters: []parser.Parameter{
					{
						Name: "assert",
						Type: singleTypeNameToTypeAnnotation("Assert"),
					},
				},
				ReturnType: singleTypeNameToTypeAnnotation("Void"),
				Block:      block,
			},
		})
	}

	return &parser.Module{
		Implementing: "UnitTests",
		Name:         "generated",
		Declarations: declarations,
	}, nil
}

func findModuleInProgram(program ast.Program, constructorName string) (*ast.Module, error) {
	for _, module := range program.Modules {
		if module.Name == constructorName {
			return module, nil
		}
	}
	return nil, fmt.Errorf("not found: %s", constructorName)
}

func findFunctionInModule(module *ast.Module, functionName string) (*ast.Function, error) {
	variable := module.Variables[functionName]
	function, ok := variable.(*ast.Function)
	if ok {
		return function, nil
	}
	return nil, fmt.Errorf("not found function %s in %s", functionName, module.Name)
}

type testCase struct {
	name               string
	functionArguments  []parser.ExpressionBox
	expectedOutput     parser.Expression
	expectedOutputType string
}

func nameOfFunctionForTestCase(test testCase) string {
	return "testCase" + regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(test.name, "")
}

func generateTestCases(function *ast.Function) ([]testCase, error) {
	if len(function.Block) != 1 {
		return nil, errors.New("todo != 1")
	}

	return testForExpression(function.Block[0])
}

func testForExpression(expression ast.Expression) ([]testCase, error) {
	caseLiteral, caseInvocation, caseAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseLiteral != nil {
		expectedOutputType := parser.LiteralFold(
			caseLiteral.Literal,
			func(arg float64) string { return "Float" },
			func(arg int) string { return "Int" },
			func(arg string) string { return "String" },
			func(arg bool) string { return "Boolean" },
		)
		return []testCase{
			{
				name:               parser.LiteralToString(caseLiteral.Literal),
				functionArguments:  []parser.ExpressionBox{},
				expectedOutput:     parser.LiteralExpression{Literal: caseLiteral.Literal},
				expectedOutputType: expectedOutputType,
			},
		}, nil
	} else if caseInvocation != nil {
		return nil, errors.New("todo testForExpression")
	} else if caseAccessAndMaybeInvocation != nil {
		return nil, errors.New("todo testForExpression")
	} else if caseFunction != nil {
		return nil, errors.New("todo testForExpression")
	} else if caseDeclaration != nil {
		return nil, errors.New("todo testForExpression")
	} else if caseIf != nil {
		return nil, errors.New("todo testForExpression")
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}
