package testgen

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"strconv"
	"strings"
	"testing"
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

func Generate(parsedProgram parser.FileTopLevel, program ast.Program, targetFunctionName ast.Ref) ([]parser.Declaration, error) {
	return generate(golang.RunCodeBlockingAndReturningOutputWhenFinished, parsedProgram, program, targetFunctionName)
}

func GenerateCached(t *testing.T, parsedProgram parser.FileTopLevel, program ast.Program, targetFunctionName ast.Ref) ([]parser.Declaration, error) {
	return generate(func(code string) (string, error) {
		return golang.RunCodeUnlessCached(t, code), nil
	}, parsedProgram, program, targetFunctionName)
}

func generate(runCode func(string) (string, error), parsedProgram parser.FileTopLevel, program ast.Program, targetFunctionName ast.Ref) ([]parser.Declaration, error) {
	targetFunction, err := findFunctionInProgram(program, targetFunctionName)
	if err != nil {
		return nil, err
	}
	testCases, err := generateTestCases(runCode, parsedProgram, program, targetFunctionName, targetFunction)
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

	declarations := []parser.Declaration{}

	for _, testCase := range testCases {
		block := []parser.ExpressionBox{}

		resultArgs := []parser.NamedArgument{}
		for _, functionArgument := range testCase.functionArguments {
			resultArgs = append(resultArgs, parser.NamedArgument{
				Argument: parser.ExpressionBox{
					Expression:              functionArgument,
					AccessOrInvocationChain: nil,
				},
			})
		}
		block = append(block, parser.ExpressionBox{
			Expression: parser.Declaration{
				Name: nameFromString("result"),
				ExpressionBox: parser.ExpressionBox{
					Expression: parser.ReferenceOrInvocation{
						Var: nameFromString(targetFunctionName.Name),
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
					DotOrArrowName: &parser.DotOrArrowName{
						Dot:     true,
						Arrow:   false,
						VarName: nameFromString("assert"),
					},
					Arguments: nil,
				},
				{
					DotOrArrowName: &parser.DotOrArrowName{
						Dot:     true,
						Arrow:   false,
						VarName: nameFromString("equal"),
					},
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
						Arguments: []parser.NamedArgument{
							{
								Argument: parser.ExpressionBox{
									Expression: parser.ReferenceOrInvocation{
										Var: nameFromString("result"),
									},
								},
							},
							{
								Argument: parser.ExpressionBox{
									Expression: parser.ReferenceOrInvocation{
										Var: nameFromString("expected"),
									},
								},
							},
						},
					},
				},
			},
		})

		declarations = append(declarations, parser.Declaration{
			Name: nameFromString("_"),
			ExpressionBox: parser.ExpressionBox{
				Expression: parser.ReferenceOrInvocation{
					Var: nameFromString("UnitTest"),
					Arguments: &parser.ArgumentsList{
						Arguments: []parser.NamedArgument{
							parser.NamedArgument{
								Argument: parser.ExpressionBox{
									Expression: parser.LiteralExpression{
										Literal: parser.LiteralString{
											Value: fmt.Sprintf("\"%s\"", testCase.name),
										},
									},
								},
							},
							parser.NamedArgument{
								Argument: parser.ExpressionBox{
									Expression: parser.LambdaOrList{
										Generics: nil,
										List:     nil,
										Lambda: &parser.Lambda{
											Signature: parser.LambdaSignature{
												Parameters: []parser.Parameter{
													{
														Name: nameFromString("testkit"),
														Type: singleTypeNameToTypeAnnotation("UnitTestKit"),
													},
												},
												ReturnType: singleTypeNameToTypeAnnotation("Void"),
											},
											Block: block,
										},
									},
								},
							},
						},
					},
				},
			},
		})
	}

	return declarations, nil
}

func findFunctionInProgram(program ast.Program, functionName ast.Ref) (*ast.Function, error) {
	var expression ast.Expression
	for decName, decExp := range program.Declarations {
		if decName == functionName {
			expression = decExp
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

type Json string

type testCase struct {
	name               string
	functionArguments  []parser.Expression
	functionReturnType types.VariableType
	expectedOutput     Json
}

type printableTestCase struct {
	name               string
	functionArguments  []parser.Expression
	expectedOutput     parser.Expression
	expectedOutputType string
}

func generateTestCases(runCode func(string) (string, error), parsedProgram parser.FileTopLevel, program ast.Program, functionName ast.Ref, function *ast.Function) ([]printableTestCase, error) {
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
			test.functionArguments = append(test.functionArguments, astExpressionToParserExpression(value))
		}
		test.functionReturnType = function.VariableType.ReturnType
		testCases = append(testCases, &test)
	}
	for _, test := range testCases {
		err := determineExpectedOutput(runCode, test, parsedProgram, program, functionName)
		if err != nil {
			return nil, err
		}
	}

	generateTestNames(testCases)

	printableTests := []printableTestCase{}
	for _, test := range testCases {
		printableTest, err := makePrintable(*test, function.VariableType, program)
		if err != nil {
			return nil, err
		}
		printableTests = append(printableTests, printableTest)
	}
	return printableTests, nil
}

func generateToJsonFunction(program ast.Program, variableType types.VariableType, functionName string) ([]parser.Import, string, error) {
	importFrom := func(vars []string, alias *string) parser.Import {
		names := []parser.Name{}
		for _, s := range vars {
			names = append(names, parser.Name{String: s})
		}
		var as *parser.Name
		if alias != nil {
			as = &parser.Name{String: *alias}
		}
		return parser.Import{DotSeparatedVars: names, As: as}
	}
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		return nil, "", errors.New("can't do generateToJsonFunction for type argument")
	} else if caseList != nil {
		ofImports, ofFunctionCode, err := generateToJsonFunction(program, caseList.Generic, fmt.Sprintf("%s_of", functionName))
		if err != nil {
			return nil, "", err
		}
		imports := append(
			ofImports,
			importFrom([]string{"tenecs", "json", "jsonList"}, nil),
			importFrom([]string{"tenecs", "json", "JsonConverter"}, nil),
		)
		ofTypeName := types.PrintableNameWithoutPackage(caseList.Generic)
		code := ofFunctionCode + fmt.Sprintf(`
%s := (): JsonConverter<List<%s>> => {
	jsonList(%s())
}
`, functionName, ofTypeName, fmt.Sprintf("%s_of", functionName))
		return imports, code, nil
	} else if caseKnownType != nil {
		if caseKnownType.Package == "" {
			switch caseKnownType.Name {
			case "Void":
				panic("TODO generateToJsonFunction Void")
			case "String":
				return []parser.Import{importFrom([]string{"tenecs", "json", "jsonString"}, &functionName)}, "", nil
			case "Int":
				return []parser.Import{importFrom([]string{"tenecs", "json", "jsonInt"}, &functionName)}, "", nil
			case "Boolean":
				return []parser.Import{importFrom([]string{"tenecs", "json", "jsonBoolean"}, &functionName)}, "", nil
			default:
				panic("unknown in generateToJsonFunction caseKnownType base " + caseKnownType.Name)
			}
		} else {
			if len(caseKnownType.Generics) > 0 {
				panic("TODO generateToJsonFunction caseKnownType with generics " + caseKnownType.Name)
			}

			imports := []parser.Import{
				importFrom([]string{"tenecs", "json", "JsonConverter"}, nil),
				importFrom([]string{"tenecs", "json", "JsonField"}, nil),
			}
			result := ""
			fields := program.FieldsByType[ast.Ref{
				Package: caseKnownType.Package,
				Name:    caseKnownType.Name,
			}]
			for fieldName, fieldVarType := range fields {
				functionImports, functionCode, err := generateToJsonFunction(program, fieldVarType, fmt.Sprintf("%s_%s", functionName, fieldName))
				if err != nil {
					return nil, "", err
				}
				imports = append(imports, functionImports...)
				result += functionCode + "\n"
			}
			result += fmt.Sprintf("%s := (): JsonConverter<%s> => {\n", functionName, types.PrintableNameWithoutPackage(variableType))
			constructorFunc := program.StructFunctions[ast.Ref{
				Package: caseKnownType.Package,
				Name:    caseKnownType.Name,
			}]
			if constructorFunc == nil {
				panic("nil constructorFunc")
			}
			result += fmt.Sprintf("jsonObject%d(\n", len(constructorFunc.Arguments))
			imports = append(imports, importFrom([]string{"tenecs", "json", fmt.Sprintf("jsonObject%d", len(constructorFunc.Arguments))}, nil))
			result += caseKnownType.Name + ",\n"
			for i, argument := range constructorFunc.Arguments {
				result += fmt.Sprintf(`JsonField("%s", %s(), (obj: %s) => obj.%s)`,
					argument.Name, fmt.Sprintf("%s_%s", functionName, argument.Name), types.PrintableNameWithoutPackage(variableType), argument.Name)
				if i < len(constructorFunc.Arguments)-1 {
					result += ","
				}
				result += "\n"
			}
			result += ")\n"
			result += "}"
			return imports, result, nil
		}
	} else if caseFunction != nil {
		return nil, "", errors.New("can't do generateToJsonFunction for function")
	} else if caseOr != nil {
		panic("TODO generateToJsonFunction caseOr")
	} else {
		panic("cases on variableType")
	}
}

func determineExpectedOutput(runCode func(string) (string, error), test *testCase, originalParsed parser.FileTopLevel, program ast.Program, targetFunctionName ast.Ref) error {
	tmpMain := "tmp_Main_qwertyuiopasdfghjklzxcvbnm"
	tmpToJson := "tmp_toJson_qwertyuiopasdfghjklzxcvbnm"

	parsedWithAddedImports, err := func() (parser.FileTopLevel, error) {
		parsed := parser.FileTopLevel{
			Tokens:               nil,
			Package:              originalParsed.Package,
			Imports:              []parser.Import{},
			TopLevelDeclarations: originalParsed.TopLevelDeclarations,
		}
		parsed.Imports = append(parsed.Imports, originalParsed.Imports...)
		parsed.Imports = append(parsed.Imports, parser.Import{
			DotSeparatedVars: []parser.Name{{String: "tenecs"}, {String: "go"}, {String: "Main"}},
			As:               &parser.Name{String: tmpMain},
		})

		programStr := formatter.DisplayFileTopLevelIgnoringComments(parsed)
		result, err := parser.ParseString(programStr)
		if err != nil {
			return parser.FileTopLevel{}, err
		}
		return *result, err
	}()
	if err != nil {
		return err
	}

	toJsonImports, toJsonCode, err := generateToJsonFunction(program, test.functionReturnType, tmpToJson)
	if err != nil {
		return err
	}
	parsedWithAddedImports.Imports = append(parsedWithAddedImports.Imports, toJsonImports...)

	tmpFunctionName := "tmp_function_test_qwertyuiopasdfghjklzxcvbnm"

	tmpProgramStr := func() string {
		invocationStr := targetFunctionName.Name
		invocationStr += "("
		for i, argument := range test.functionArguments {
			if i > 0 {
				invocationStr += ", "
			}
			invocationStr += formatter.DisplayExpression(argument)
		}
		invocationStr += ")"

		invocationStr = tmpToJson + "().toJson(" + invocationStr + ")"
		invocationStr = "runtime.console.log(" + invocationStr + ")"

		tmpProgramStr := formatter.DisplayFileTopLevel(parsedWithAddedImports)
		tmpProgramStr += toJsonCode
		tmpProgramStr += fmt.Sprintf(`
%s := %s(
  (runtime) => {
    %s
  }
)
`, tmpFunctionName, tmpMain, invocationStr)
		return tmpProgramStr
	}()

	runOutput, err := func() (string, error) {
		fileTopLevel, err := parser.ParseString(tmpProgramStr)
		if err != nil {
			return "", err
		}
		desugared, err := desugar.Desugar(*fileTopLevel)
		if err != nil {
			return "", err
		}
		program, err := typer.TypecheckSingleFile(desugared)
		if err != nil {
			return "", err
		}
		pkgName := ""
		for i, name := range originalParsed.Package.DotSeparatedNames {
			if i > 0 {
				pkgName += "."
			}
			pkgName += name.String
		}
		generatedProgram := codegen_golang.GenerateProgramMain(program, ast.Ref{
			Package: pkgName,
			Name:    tmpFunctionName,
		})

		return runCode(generatedProgram)
	}()
	if err != nil {
		return err
	}
	test.expectedOutput = Json(runOutput)
	return nil
}

func makePrintable(test testCase, function *types.Function, program ast.Program) (printableTestCase, error) {
	name := test.name
	if name == "" {
		name = "(empty)"
	}

	functionArgs := []parser.Expression{}
	for _, argument := range test.functionArguments {
		functionArgs = append(functionArgs, argument)
	}

	expectedOutputAst, err := parseJsonAsInstanceOfType(test.expectedOutput, function.ReturnType, program)
	if err != nil {
		return printableTestCase{}, err
	}
	expectedOutput := astExpressionToParserExpression(expectedOutputAst)

	expectedOutputType := typeNameOfVariableType(function.ReturnType)
	split := strings.Split(expectedOutputType, ".")
	expectedOutputType = split[len(split)-1]
	return printableTestCase{
		name:               name,
		functionArguments:  functionArgs,
		expectedOutput:     expectedOutput,
		expectedOutputType: expectedOutputType,
	}, nil
}

func generateTestNames(tests []*testCase) {
	existingTestNames := map[string]bool{}
	for _, test := range tests {
		name := string(test.expectedOutput)
		name = strings.ReplaceAll(name, "\n", "")
		name = strings.ReplaceAll(name, "\"", "")
		test.name = name
		if _, ok := existingTestNames[test.name]; ok {
			test.name = name + " again"
		}
		existingTestNames[test.name] = true
	}
}

func astExpressionToParserExpression(expression ast.Expression) parser.Expression {
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
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
			arguments := []parser.NamedArgument{}
			for _, argumentExp := range caseInvocation.Arguments {
				arguments = append(arguments, parser.NamedArgument{
					Argument: parser.ExpressionBox{
						Expression: astExpressionToParserExpression(argumentExp),
					},
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
		genericNames := []parser.TypeAnnotation{}
		for _, generic := range caseFunction.VariableType.Generics {
			genericNames = append(genericNames, parser.TypeAnnotation{
				OrTypes: []parser.TypeAnnotationElement{
					parser.SingleNameType{
						TypeName: nameFromString(generic),
					},
				},
			})
		}
		generics := &parser.LambdaOrListGenerics{
			Generics: genericNames,
		}
		if caseFunction.VariableType.Generics == nil {
			generics = nil
		}
		return parser.LambdaOrList{
			Generics: generics,
			List:     nil,
			Lambda: &parser.Lambda{
				Signature: parser.LambdaSignature{
					Parameters: parameters,
					ReturnType: nil,
				},
				Block: block,
			},
		}
	} else if caseDeclaration != nil {
		panic("TODO astExpressionToParserExpression caseDeclaration")
	} else if caseIf != nil {
		panic("TODO astExpressionToParserExpression caseIf")
	} else if caseList != nil {
		genericTypeAnnotation := typeAnnotationOfVariableType(caseList.ContainedVariableType)

		expressions := []parser.ExpressionBox{}

		for _, argument := range caseList.Arguments {
			expressions = append(expressions, parser.ExpressionBox{
				Expression:              astExpressionToParserExpression(argument),
				AccessOrInvocationChain: nil,
			})
		}

		return parser.LambdaOrList{
			Generics: &parser.LambdaOrListGenerics{
				Generics: []parser.TypeAnnotation{
					genericTypeAnnotation,
				},
			},
			List: &parser.List{
				Expressions: expressions,
			},
			Lambda: nil,
		}
	} else if caseWhen != nil {
		panic("TODO astExpressionToParserExpression caseWhen")
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func typeNameOfVariableType(varType types.VariableType) string {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO typeNameOfVariableType caseTypeArgument")
	} else if caseList != nil {
		return "List<" + typeNameOfVariableType(caseList.Generic) + ">"
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
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO typeAnnotationOfVariableType caseTypeArgument")
	} else if caseList != nil {
		return parser.TypeAnnotation{
			OrTypes: []parser.TypeAnnotationElement{
				parser.SingleNameType{
					TypeName: nameFromString("List"),
					Generics: []parser.TypeAnnotation{typeAnnotationOfVariableType(caseList.Generic)},
				},
			},
		}
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

func parseJsonAsInstanceOfType(value Json, variableType types.VariableType, program ast.Program) (ast.Expression, error) {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		return nil, errors.New("TODO parseJsonAsInstanceOfType caseTypeArgument")
	} else if caseList != nil {
		var preResult []json.RawMessage
		err := json.Unmarshal([]byte(value), &preResult)
		if err != nil {
			return nil, err
		}
		result := ast.List{
			ContainedVariableType: caseList.Generic,
			Arguments:             []ast.Expression{},
		}
		for _, elemJson := range preResult {
			elem, err := parseJsonAsInstanceOfType(Json(elemJson), caseList.Generic, program)
			if err != nil {
				return nil, err
			}
			result.Arguments = append(result.Arguments, elem)
		}
		return result, nil
	} else if caseKnownType != nil {
		if caseKnownType.Package == "" {
			switch caseKnownType.Name {
			case "Boolean":
				var result bool
				err := json.Unmarshal([]byte(value), &result)
				if err != nil {
					return nil, err
				}
				return ast.Literal{
					VariableType: types.Boolean(),
					Literal: parser.LiteralString{
						Value: strconv.FormatBool(result),
					},
				}, nil
			case "String":
				var result string
				err := json.Unmarshal([]byte(value), &result)
				if err != nil {
					return nil, err
				}
				return ast.Literal{
					VariableType: types.String(),
					Literal: parser.LiteralString{
						Value: strings.TrimSpace(string(value)),
					},
				}, nil
			default:
				return nil, errors.New("TODO parseJsonAsInstanceOfType caseKnownType: " + caseKnownType.Name)
			}
		}
		if len(caseKnownType.Generics) > 0 {
			return nil, errors.New("TODO parseJsonAsInstanceOfType caseKnownType with generics")
		}

		var preResult map[string]json.RawMessage
		err := json.Unmarshal([]byte(value), &preResult)
		if err != nil {
			return nil, err
		}
		resultArguments := []ast.Expression{}
		constructorFunc := program.StructFunctions[ast.Ref{
			Package: caseKnownType.Package,
			Name:    caseKnownType.Name,
		}]
		for _, argument := range constructorFunc.Arguments {
			elemJson := preResult[argument.Name]
			elem, err := parseJsonAsInstanceOfType(Json(elemJson), argument.VariableType, program)
			if err != nil {
				return nil, err
			}
			resultArguments = append(resultArguments, elem)
		}
		return ast.Invocation{
			VariableType: variableType,
			Over: ast.Reference{
				VariableType: constructorFunc,
				PackageName:  &caseKnownType.Package,
				Name:         caseKnownType.Name,
			},
			Generics:  nil,
			Arguments: resultArguments,
		}, nil

	} else if caseFunction != nil {
		return nil, errors.New("TODO parseJsonAsInstanceOfType caseFunction")
	} else if caseOr != nil {
		return nil, errors.New("TODO parseJsonAsInstanceOfType caseOr")
	} else {
		panic(fmt.Errorf("cases on %v", variableType))
	}
}
