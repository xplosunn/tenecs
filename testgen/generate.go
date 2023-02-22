package testgen

import (
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/interpreter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strconv"
	"strings"
)

func Generate(program ast.Program, targetFunctionName string) (*parser.Module, error) {
	targetFunction, err := findFunctionInProgram(program, targetFunctionName)
	if err != nil {
		return nil, err
	}
	testCases, err := generateTestCases(targetFunction)
	if err != nil {
		return nil, err
	}

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
										Value: "\"" + testCase.name + "\"",
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

		resultArgs := []parser.ExpressionBox{}
		for _, functionArgument := range testCase.functionArguments {
			resultArgs = append(resultArgs, parser.ExpressionBox{
				Expression:              functionArgument,
				AccessOrInvocationChain: nil,
			})
		}
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
								Arguments: resultArgs,
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

func nameOfFunctionForTestCase(test printableTestCase) string {
	suffix := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(test.name, "")
	suffix = cases.Title(language.English, cases.Compact).String(suffix)
	return "testCase" + suffix
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
	name              string
	functionArguments []ast.Expression
	expectedOutput    interpreter.Value
}

type printableTestCase struct {
	name               string
	functionArguments  []parser.Expression
	expectedOutput     parser.Expression
	expectedOutputType string
}

func generateTestCases(function *ast.Function) ([]printableTestCase, error) {
	testCases := []*testCase{}

	constraintsForTestCases, err := findConstraints(function)
	if err != nil {
		return nil, err
	}
	if len(constraintsForTestCases) == 0 {
		constraintsForTestCases = []testCaseConstraints{
			testCaseConstraints{
				argsConstraints: immutable.NewMap[string, []valueConstraint](nil),
			},
		}
	}
	for _, constraints := range constraintsForTestCases {
		test := testCase{}
		for _, functionArgument := range function.VariableType.Arguments {
			constraints, ok := constraints.argsConstraints.Get(functionArgument.Name)
			if !ok {
				constraints = []valueConstraint{}
			}
			value, err := satisfy(functionArgument.Name, functionArgument.VariableType, constraints)
			if err != nil {
				return nil, err
			}
			test.functionArguments = append(test.functionArguments, value)
		}
		testCases = append(testCases, &test)
	}
	for _, test := range testCases {
		_, value, err := interpreter.EvalBlock(
			interpreter.NewScope(),
			[]ast.Expression{
				ast.Declaration{
					VariableType: types.Void{},
					Name:         "target",
					Expression:   function,
				},
				ast.ReferenceAndMaybeInvocation{
					VariableType: function.VariableType,
					Name:         "target",
					ArgumentsList: &ast.ArgumentsList{
						Arguments: test.functionArguments,
					},
				},
			},
		)
		if err != nil {
			return nil, err
		}
		test.expectedOutput = value
	}

	generateTestNames(testCases)

	printableTests := []printableTestCase{}
	for _, test := range testCases {
		printableTest, err := makePrintable(*test)
		if err != nil {
			return nil, err
		}
		printableTests = append(printableTests, printableTest)
	}
	return printableTests, nil
}

func makePrintable(test testCase) (printableTestCase, error) {

	functionArgs := []parser.Expression{}
	for _, argument := range test.functionArguments {
		functionArgs = append(functionArgs, astExpressionToParserExpression(argument))
	}

	expectedOutputAst := valueToAstExpression(test.expectedOutput)
	expectedOutput := astExpressionToParserExpression(expectedOutputAst)

	expectedOutputType := typeNameOfValue(test.expectedOutput)

	return printableTestCase{
		name:               test.name,
		functionArguments:  functionArgs,
		expectedOutput:     expectedOutput,
		expectedOutputType: expectedOutputType,
	}, nil
}

func generateTestNames(tests []*testCase) {
	existingTestNames := map[string]bool{}
	for _, test := range tests {
		name := ""
		interpreter.ValueExhaustiveSwitch(
			test.expectedOutput,
			func(value interpreter.ValueVoid) {
				name = "void"
			},
			func(value interpreter.ValueBoolean) {
				name = strconv.FormatBool(value.Bool)
			},
			func(value interpreter.ValueFloat) {
				name = fmt.Sprintf("%f", value.Float)
			},
			func(value interpreter.ValueInt) {
				name = fmt.Sprintf("%d", value.Int)
			},
			func(value interpreter.ValueString) {
				name = strings.TrimPrefix(strings.TrimSuffix(value.String, "\""), "\"")
			},
			func(value interpreter.ValueFunction) {
				panic("TODO generateTestNames function")
			},
		)
		test.name = name
		if _, ok := existingTestNames[test.name]; ok {
			test.name = name + " again"
		}
		existingTestNames[test.name] = true
	}
}

func astExpressionToParserExpression(expression ast.Expression) parser.Expression {
	caseModule, caseLiteral, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseModule != nil {
		declarations := []parser.ModuleDeclaration{}
		for _, _ = range caseModule.Variables {
			panic("TODO astExpressionToParserExpression caseModule.Variables")
		}
		return parser.Module{
			Implementing: caseModule.Implements.Name,
			Declarations: declarations,
		}
	} else if caseLiteral != nil {
		return parser.LiteralExpression{
			Literal: caseLiteral.Literal,
		}
	} else if caseReferenceAndMaybeInvocation != nil {
		panic("TODO astExpressionToParserExpression caseReferenceAndMaybeInvocation")
	} else if caseWithAccessAndMaybeInvocation != nil {
		panic("TODO astExpressionToParserExpression caseWithAccessAndMaybeInvocation")
	} else if caseFunction != nil {
		parameters := []parser.Parameter{}
		for i, _ := range caseFunction.VariableType.Arguments {
			parameters = append(parameters, parser.Parameter{
				Name: fmt.Sprintf("arg%d", i),
				Type: nil,
			})
		}
		block := []parser.ExpressionBox{}
		for _, exp := range caseFunction.Block {
			block = append(block, parser.ExpressionBox{Expression: astExpressionToParserExpression(exp)})
		}
		return parser.Lambda{
			Generics:   caseFunction.VariableType.Generics,
			Parameters: parameters,
			ReturnType: nil,
			Block:      block,
		}
		panic("TODO astExpressionToParserExpression caseFunction")
	} else if caseDeclaration != nil {
		panic("TODO astExpressionToParserExpression caseDeclaration")
	} else if caseIf != nil {
		panic("TODO astExpressionToParserExpression caseIf")
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func typeNameOfValue(value interpreter.Value) string {
	var result string
	interpreter.ValueExhaustiveSwitch(
		value,
		func(value interpreter.ValueVoid) {
			panic("TODO typeNameOfValue string")
		},
		func(value interpreter.ValueBoolean) {
			result = "Boolean"
		},
		func(value interpreter.ValueFloat) {
			result = "Float"
		},
		func(value interpreter.ValueInt) {
			result = "Int"
		},
		func(value interpreter.ValueString) {
			result = "String"
		},
		func(value interpreter.ValueFunction) {
			panic("TODO typeNameOfValue function")
		},
	)
	return result
}

func valueToAstExpression(value interpreter.Value) ast.Expression {
	var result ast.Expression
	interpreter.ValueExhaustiveSwitch(
		value,
		func(value interpreter.ValueVoid) {
			panic("TODO ValueToAstExpression void")
		},
		func(value interpreter.ValueBoolean) {
			result = ast.Literal{
				VariableType: types.BasicType{
					Type: "Boolean",
				},
				Literal: parser.LiteralBool{
					Value: value.Bool,
				},
			}
		},
		func(value interpreter.ValueFloat) {
			result = ast.Literal{
				VariableType: types.BasicType{
					Type: "Float",
				},
				Literal: parser.LiteralFloat{
					Value: value.Float,
				},
			}
		},
		func(value interpreter.ValueInt) {
			result = ast.Literal{
				VariableType: types.BasicType{
					Type: "Int",
				},
				Literal: parser.LiteralInt{
					Value: value.Int,
				},
			}
		},
		func(value interpreter.ValueString) {
			result = ast.Literal{
				VariableType: types.BasicType{
					Type: "String",
				},
				Literal: parser.LiteralString{
					Value: value.String,
				},
			}
		},
		func(value interpreter.ValueFunction) {
			result = value.AstFunction
		},
	)
	return result
}
