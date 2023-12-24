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

func nameFromString(name string) parser.Name {
	return parser.Name{
		Node:   parser.Node{},
		String: name,
	}
}

func ptrNameFromString(name string) *parser.Name {
	return &parser.Name{
		Node:   parser.Node{},
		String: name,
	}
}

func Generate(program ast.Program, targetFunctionName string) (*parser.Implementation, error) {
	targetFunction, err := findFunctionInProgram(program, targetFunctionName)
	if err != nil {
		return nil, err
	}
	testCases, err := generateTestCases(program, targetFunction)
	if err != nil {
		return nil, err
	}

	singleTypeNameToTypeAnnotation := func(typeName string) *parser.TypeAnnotation {
		return &parser.TypeAnnotation{
			OrTypes: []parser.TypeAnnotationElement{
				parser.SingleNameType{
					TypeName: nameFromString(typeName),
				},
			},
		}
	}

	testsBlock := []parser.ExpressionBox{}

	for _, testCase := range testCases {
		testsBlock = append(testsBlock, parser.ExpressionBox{
			Expression: parser.ReferenceOrInvocation{
				Var:       nameFromString("registry"),
				Arguments: nil,
			},
			AccessOrInvocationChain: []parser.AccessOrInvocation{
				{
					VarName: ptrNameFromString("test"),
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
									Var: nameFromString(nameOfFunctionForTestCase(testCase)),
								},
							},
						},
					},
				},
			},
		})
	}

	declarations := []parser.ImplementationDeclaration{
		{
			Public: true,
			Name:   nameFromString("tests"),
			Expression: parser.Lambda{
				Parameters: []parser.Parameter{
					{
						Name: nameFromString("registry"),
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
				Name: nameFromString("result"),
				ExpressionBox: parser.ExpressionBox{
					Expression: parser.ReferenceOrInvocation{
						Var: nameFromString(targetFunctionName),
						Arguments: &parser.ArgumentsList{
							Arguments: resultArgs,
						},
					},
					AccessOrInvocationChain: nil,
				},
			},
		})

		block = append(block, parser.ExpressionBox{
			Expression: parser.Declaration{
				Name: nameFromString("expected"),
				ExpressionBox: parser.ExpressionBox{
					Expression: testCase.expectedOutput,
				},
			},
		})

		block = append(block, parser.ExpressionBox{
			Expression: parser.ReferenceOrInvocation{
				Var: nameFromString("testkit"),
			},
			AccessOrInvocationChain: []parser.AccessOrInvocation{
				{
					VarName:   ptrNameFromString("assert"),
					Arguments: nil,
				},
				{
					VarName: ptrNameFromString("equal"),
					Arguments: &parser.ArgumentsList{
						Generics: []parser.TypeAnnotation{
							parser.TypeAnnotation{
								OrTypes: []parser.TypeAnnotationElement{
									parser.SingleNameType{
										TypeName: nameFromString(testCase.expectedOutputType),
									},
								},
							},
						},
						Arguments: []parser.ExpressionBox{
							{
								Expression: parser.ReferenceOrInvocation{
									Var: nameFromString("result"),
								},
							},
							{
								Expression: parser.ReferenceOrInvocation{
									Var: nameFromString("expected"),
								},
							},
						},
					},
				},
			},
		})

		declarations = append(declarations, parser.ImplementationDeclaration{
			Name: nameFromString(nameOfFunctionForTestCase(testCase)),
			Expression: parser.Lambda{
				Parameters: []parser.Parameter{
					{
						Name: nameFromString("testkit"),
						Type: singleTypeNameToTypeAnnotation("UnitTestKit"),
					},
				},
				ReturnType: singleTypeNameToTypeAnnotation("Void"),
				Block:      block,
			},
		})
	}

	return &parser.Implementation{
		Implementing: nameFromString("UnitTests"),
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

func generateTestCases(program ast.Program, function *ast.Function) ([]printableTestCase, error) {
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

	satisfier := NewSatisfier(program)

	for _, constraints := range constraintsForTestCases {
		test := testCase{}
		for _, functionArgument := range function.VariableType.Arguments {
			constraints, ok := constraints.argsConstraints.Get(functionArgument.Name)
			if !ok {
				constraints = []valueConstraint{}
			}
			value, err := satisfy(satisfier, functionArgument.Name, functionArgument.VariableType, constraints)
			if err != nil {
				return nil, err
			}
			test.functionArguments = append(test.functionArguments, value)
		}
		testCases = append(testCases, &test)
	}
	for _, test := range testCases {
		scope, err := interpreter.NewScope(program)
		if err != nil {
			return nil, err
		}
		_, value, err := interpreter.EvalBlock(
			scope,
			[]ast.Expression{
				ast.Declaration{
					Name:       "target",
					Expression: function,
				},
				ast.Invocation{
					VariableType: function.VariableType.ReturnType,
					Over: ast.Reference{
						VariableType: function.VariableType,
						Name:         "target",
					},
					Generics:  nil,
					Arguments: test.functionArguments,
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

func makeName(outputValue interpreter.Value) string {
	name := ""
	interpreter.ValueExhaustiveSwitch(
		outputValue,
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
		func(value interpreter.ValueNativeFunction) {
			panic("TODO generateTestNames native function")
		},
		func(value interpreter.ValueStructFunction) {
			panic("TODO generateTestNames struct function")
		},
		func(value interpreter.ValueStruct) {
			for _, v := range value.KeyValues {
				name = name + makeName(v)
			}
		},
		func(value interpreter.ValueArray) {
			name = "["
			for i, v := range value.Values {
				if i > 0 {
					name += ","
				}
				name += makeName(v)
			}
			name += "]"
		},
	)
	return name
}

func generateTestNames(tests []*testCase) {
	existingTestNames := map[string]bool{}
	for _, test := range tests {
		name := makeName(test.expectedOutput)
		test.name = name
		if _, ok := existingTestNames[test.name]; ok {
			test.name = name + " again"
		}
		existingTestNames[test.name] = true
	}
}

func astExpressionToParserExpression(expression ast.Expression) parser.Expression {
	caseImplementation, caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseArray, caseWhen := expression.ExpressionCases()
	if caseImplementation != nil {
		declarations := []parser.ImplementationDeclaration{}
		for _, _ = range caseImplementation.Variables {
			panic("TODO astExpressionToParserExpression caseImplementation.Variables")
		}
		return parser.Implementation{
			Implementing: nameFromString(caseImplementation.Implements.Name),
			Declarations: declarations,
		}
	} else if caseLiteral != nil {
		return parser.LiteralExpression{
			Literal: caseLiteral.Literal,
		}
	} else if caseReference != nil {
		var args *parser.ArgumentsList
		return parser.ReferenceOrInvocation{
			Var:       nameFromString(caseReference.Name),
			Arguments: args,
		}
	} else if caseAccess != nil {
		panic("TODO astExpressionToParserExpression caseWithAccessAndMaybeInvocation")
	} else if caseInvocation != nil {
		if ref, ok := caseInvocation.Over.(ast.Reference); ok {
			generics := []parser.TypeAnnotation{}
			for _, generic := range caseInvocation.Generics {
				generics = append(generics, typeAnnotationOfVariableType(generic))
			}
			arguments := []parser.ExpressionBox{}
			for _, argumentExp := range caseInvocation.Arguments {
				arguments = append(arguments, parser.ExpressionBox{
					Expression: astExpressionToParserExpression(argumentExp),
				})
			}
			args := &parser.ArgumentsList{
				Generics:  generics,
				Arguments: arguments,
			}
			return parser.ReferenceOrInvocation{
				Var:       nameFromString(ref.Name),
				Arguments: args,
			}
		}
		panic("TODO astExpressionToParserExpression caseInvocation: " + fmt.Sprintf("%T", caseInvocation.Over))
	} else if caseFunction != nil {
		parameters := []parser.Parameter{}
		for i, _ := range caseFunction.VariableType.Arguments {
			parameters = append(parameters, parser.Parameter{
				Name: nameFromString(fmt.Sprintf("arg%d", i)),
				Type: nil,
			})
		}
		block := []parser.ExpressionBox{}
		for _, exp := range caseFunction.Block {
			block = append(block, parser.ExpressionBox{Expression: astExpressionToParserExpression(exp)})
		}
		genericNames := []parser.Name{}
		for _, generic := range caseFunction.VariableType.Generics {
			genericNames = append(genericNames, nameFromString(generic))
		}
		if caseFunction.VariableType.Generics == nil {
			genericNames = nil
		}
		return parser.Lambda{
			Generics:   genericNames,
			Parameters: parameters,
			ReturnType: nil,
			Block:      block,
		}
	} else if caseDeclaration != nil {
		panic("TODO astExpressionToParserExpression caseDeclaration")
	} else if caseIf != nil {
		panic("TODO astExpressionToParserExpression caseIf")
	} else if caseArray != nil {
		genericTypeAnnotation := typeAnnotationOfVariableType(caseArray.ContainedVariableType)

		expressions := []parser.ExpressionBox{}

		for _, argument := range caseArray.Arguments {
			expressions = append(expressions, parser.ExpressionBox{
				Expression:              astExpressionToParserExpression(argument),
				AccessOrInvocationChain: nil,
			})
		}

		return parser.Array{
			Generic:     &genericTypeAnnotation,
			Expressions: expressions,
		}
	} else if caseWhen != nil {
		panic("TODO astExpressionToParserExpression caseWhen")
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
		func(value interpreter.ValueNativeFunction) {
			panic("TODO typeNameOfValue native function")
		},
		func(value interpreter.ValueStructFunction) {
			panic("TODO typeNameOfValue struct function")
		},
		func(value interpreter.ValueStruct) {
			result = value.StructName
		},
		func(value interpreter.ValueArray) {
			result = "Array<" + typeNameOfVariableType(value.Type) + ">"
		},
	)
	return result
}

func typeNameOfVariableType(varType types.VariableType) string {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO typeNameOfVariableType caseTypeArgument")
	} else if caseKnownType != nil {
		generics := ""
		if len(caseKnownType.Generics) > 0 {
			generics = "<"
			for i, generic := range caseKnownType.Generics {
				if i > 0 {
					generics += ", "
				}
				generics += typeNameOfVariableType(generic)
			}
			generics += ">"
		}
		pkg := caseKnownType.Package
		if pkg != "" {
			pkg += "."
		}
		return pkg + caseKnownType.Name + generics
	} else if caseFunction != nil {
		panic("TODO typeNameOfVariableType caseFunction")
	} else if caseOr != nil {
		panic("TODO typeNameOfVariableType caseOr")
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}

func typeAnnotationOfVariableType(variableType types.VariableType) parser.TypeAnnotation {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO typeAnnotationOfVariableType caseTypeArgument")
	} else if caseKnownType != nil {
		generics := []parser.TypeAnnotation{}
		for _, generic := range caseKnownType.Generics {
			generics = append(generics, typeAnnotationOfVariableType(generic))
		}
		if len(generics) == 0 {
			generics = nil
		}
		return parser.TypeAnnotation{
			OrTypes: []parser.TypeAnnotationElement{
				parser.SingleNameType{
					TypeName: nameFromString(caseKnownType.Name),
					Generics: generics,
				},
			},
		}
	} else if caseFunction != nil {
		panic("TODO typeAnnotationOfVariableType caseFunction")
	} else if caseOr != nil {
		panic("TODO typeAnnotationOfVariableType caseOr")
	} else {
		panic(fmt.Errorf("cases on %v", variableType))
	}
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
				VariableType: types.Boolean(),
				Literal: parser.LiteralBool{
					Value: value.Bool,
				},
			}
		},
		func(value interpreter.ValueFloat) {
			result = ast.Literal{
				VariableType: types.Float(),
				Literal: parser.LiteralFloat{
					Value: value.Float,
				},
			}
		},
		func(value interpreter.ValueInt) {
			result = ast.Literal{
				VariableType: types.Int(),
				Literal: parser.LiteralInt{
					Value: value.Int,
				},
			}
		},
		func(value interpreter.ValueString) {
			result = ast.Literal{
				VariableType: types.String(),
				Literal: parser.LiteralString{
					Value: value.String,
				},
			}
		},
		func(value interpreter.ValueFunction) {
			result = value.AstFunction
		},
		func(value interpreter.ValueNativeFunction) {
			panic("TODO valueToAstExpression ValueNativeFunction")
		},
		func(value interpreter.ValueStructFunction) {
			panic("TODO valueToAstExpression ValueStructFunction")
		},
		func(value interpreter.ValueStruct) {
			args := []ast.Expression{}
			for _, value := range value.OrderedValues {
				args = append(args, valueToAstExpression(value))
			}
			result = ast.Invocation{
				VariableType: &types.KnownType{
					Package: "",
					Name:    value.StructName,
				},
				Over: ast.Reference{
					VariableType: nil,
					Name:         value.StructName,
				},
				Generics:  nil,
				Arguments: args,
			}
		},
		func(value interpreter.ValueArray) {
			if len(value.Values) == 0 {
				result = ast.Array{
					ContainedVariableType: value.Type,
					Arguments:             []ast.Expression{},
				}
			} else {
				panic("TODO valueToAstExpression ValueArray")
			}
		},
	)
	return result
}
