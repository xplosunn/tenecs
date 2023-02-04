package testgen

import (
	"errors"
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"regexp"
	"strings"
)

func Generate(program ast.Program, targetFunctionName string) (*parser.Module, error) {
	targetFunction, err := findFunctionInProgram(program, targetFunctionName)
	if err != nil {
		return nil, err
	}
	testCases, err := generateTestCases(targetFunction)

	singleTypeNameToTypeAnnotation := func(typeName string) *parser.TypeAnnotation {
		var typeAnnotation parser.TypeAnnotation = parser.SingleNameType{
			TypeName: typeName,
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
		Declarations: declarations,
	}, nil
}

func findFunctionInProgram(program ast.Program, functionName string) (*ast.Function, error) {
	var expression ast.Expression
	for _, declaration := range program.Declarations {
		if declaration.Name == functionName {
			expression = declaration.Expression
			break
		}
	}
	if expression == nil {
		return nil, fmt.Errorf("not found function %s", functionName)
	}
	function, ok := expression.(*ast.Function)
	if !ok {
		return nil, fmt.Errorf("%s is not a function", functionName)
	}
	return function, nil

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
	caseModule, caseLiteral, caseInvocation, caseAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseModule != nil {
		return nil, errors.New("todo testForExpression")
	} else if caseLiteral != nil {
		expectedOutputType := parser.LiteralFold(
			caseLiteral.Literal,
			func(arg float64) string { return "Float" },
			func(arg int) string { return "Int" },
			func(arg string) string { return "String" },
			func(arg bool) string { return "Boolean" },
		)
		name := parser.LiteralToString(caseLiteral.Literal)
		if !strings.HasPrefix(name, "\"") {
			name = fmt.Sprintf("\"%s\"", name)
		}
		return []testCase{
			{
				name:               name,
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
