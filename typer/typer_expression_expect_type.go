package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
	"reflect"
)

func expectTypeOfExpressionBox(validateFunctionBlock bool, expressionBox parser.ExpressionBox, expectedType types.VariableType, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	noInvocation := expressionBox.AccessOrInvocationChain == nil
	caseExpectFunction, expectsFunction := expectedType.(*types.Function)
	caseLambda, expressionIsLambda := expressionBox.Expression.(parser.Lambda)
	if noInvocation && expectsFunction && expressionIsLambda {
		return expectTypeOfLambda(validateFunctionBlock, caseLambda, caseExpectFunction, universe)
	}
	universe, astExp, err := determineTypeOfExpressionBox(validateFunctionBlock, expressionBox, universe)
	if err != nil {
		return nil, nil, err
	}
	varType := ast.VariableTypeOfExpression(astExp)
	if !variableTypeEq(varType, expectedType) {
		return nil, nil, type_error.PtrOnNodef(parser.GetExpressionNode(expressionBox.Expression), "expected type %s but found %s", printableName(expectedType), printableName(varType))
	}
	return universe, astExp, nil
}

func expectTypeOfExpression(validateFunctionBlock bool, exp parser.Expression, expectedType types.VariableType, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	resultUniverse := universe
	var resultExpression ast.Expression
	var err *type_error.TypecheckError
	parser.ExpressionExhaustiveSwitch(
		exp,
		func(expression parser.Module) {
			resultUniverse, resultExpression, err = determineTypeOfModule(validateFunctionBlock, expression, universe)
		},
		func(expression parser.LiteralExpression) {
			programExp := determineTypeOfLiteral(expression.Literal)
			varType := ast.VariableTypeOfExpression(programExp)
			if !variableTypeEq(varType, expectedType) {
				err = type_error.PtrOnNodef(parser.GetExpressionNode(exp), "expected type %s but found %s", printableName(expectedType), printableName(varType))
				return
			}
			resultExpression = programExp
		},
		func(expression parser.ReferenceOrInvocation) {
			programExp, err2 := determineTypeOfReferenceOrInvocation(validateFunctionBlock, expression, universe)
			if err2 != nil {
				err = err2
				return
			}
			varType := ast.VariableTypeOfExpression(programExp)
			if !variableTypeEq(varType, expectedType) {
				err = type_error.PtrOnNodef(expression.Var.Node, "in expression '%s' expected %s but found %s", expression.Var.String, printableName(expectedType), printableName(varType))
				return
			}
			resultExpression = programExp
		},
		func(expression parser.Lambda) {
			resultUniverse, resultExpression, err = expectTypeOfLambda(validateFunctionBlock, expression, expectedType, universe)
		},
		func(expression parser.Declaration) {
			u, programExp, err2 := determineTypeOfExpressionBox(validateFunctionBlock, expression.ExpressionBox, universe)
			if err2 != nil {
				err = err2
				return
			}
			if !variableTypeEq(expectedType, &standard_library.Void) {
				err = type_error.PtrOnNodef(expression.Name.Node, "expected type %s but found Void (variable declarations return void)", printableName(expectedType))
				return
			}
			universe = u
			resultExpression = programExp
		},
		func(expression parser.If) {
			u, programExp, err2 := determineTypeOfIf(validateFunctionBlock, expression, universe)
			if err2 != nil {
				err = err2
				return
			}
			varType := ast.VariableTypeOfExpression(programExp)
			if !variableTypeEq(varType, expectedType) {
				err = type_error.PtrOnNodef(expression.Node, "expected type %s but found %s", printableName(expectedType), printableName(varType))
				return
			}
			universe = u
			resultExpression = programExp
		},
	)
	return resultUniverse, resultExpression, err
}

func expectTypeOfLambda(validateFunctionBlock bool, lambda parser.Lambda, expectedType types.VariableType, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	var expectedFunction types.Function
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := expectedType.VariableTypeCases()
	if caseTypeArgument != nil {
		return nil, nil, type_error.PtrOnNodef(lambda.Node, "expected type %s but found a Function", printableName(expectedType))
	} else if caseStruct != nil {
		return nil, nil, type_error.PtrOnNodef(lambda.Node, "expected type %s but found a Function", printableName(expectedType))
	} else if caseInterface != nil {
		return nil, nil, type_error.PtrOnNodef(lambda.Node, "expected type %s but found a Function", printableName(expectedType))
	} else if caseFunction != nil {
		expectedFunction = *caseFunction
	} else if caseBasicType != nil {
		return nil, nil, type_error.PtrOnNodef(lambda.Node, "expected type %s but found a Function", printableName(expectedType))
	} else if caseVoid != nil {
		return nil, nil, type_error.PtrOnNodef(lambda.Node, "expected type %s but found a Function", printableName(expectedType))
	} else {
		panic(fmt.Errorf("code on %v", expectedType))
	}

	functionArgs := []types.FunctionArgument{}
	generics, parameters, annotatedReturnType, block := parser.LambdaFields(lambda)
	_ = block
	if len(generics) != len(expectedFunction.Generics) {
		return nil, nil, type_error.PtrOnNodef(lambda.Node, "expected same number of generics as interface variable (%d) but found %d", len(expectedFunction.Generics), len(generics))
	}
	if len(parameters) != len(expectedFunction.Arguments) {
		return nil, nil, type_error.PtrOnNodef(lambda.Node, "expected same number of arguments as interface variable (%d) but found %d", len(expectedFunction.Arguments), len(parameters))
	}
	localUniverse := universe
	for _, generic := range generics {
		u, err := binding.CopyAddingType(localUniverse, generic, &types.TypeArgument{Name: generic.String})
		if err != nil {
			return nil, nil, err
		}
		localUniverse = u
	}
	for i, parameter := range parameters {
		if parameter.Type == nil {
			functionArgs = append(functionArgs, types.FunctionArgument{
				Name:         parameter.Name.String,
				VariableType: expectedFunction.Arguments[i].VariableType,
			})
			continue
		}

		varType, err := validateTypeAnnotationInUniverse(*parameter.Type, localUniverse)
		if err != nil {
			return nil, nil, err
		}

		if !variableTypeEq(varType, expectedFunction.Arguments[i].VariableType) {
			return nil, nil, type_error.PtrOnNodef(parameter.Name.Node, "in parameter position %d expected type %s but you have annotated %s", i, printableName(expectedFunction.Arguments[i].VariableType), printableNameOfTypeAnnotation(*parameter.Type))
		}

		functionArgs = append(functionArgs, types.FunctionArgument{
			Name:         parameter.Name.String,
			VariableType: varType,
		})
	}
	functionArgumentNames := []parser.Name{}
	for _, parameter := range lambda.Parameters {
		functionArgumentNames = append(functionArgumentNames, parameter.Name)
	}
	functionArgumentVariableTypes := []types.VariableType{}
	for _, argument := range functionArgs {
		functionArgumentVariableTypes = append(functionArgumentVariableTypes, argument.VariableType)
	}
	localUniverse, err := binding.CopyAddingFunctionArguments(localUniverse, functionArgumentNames, functionArgumentVariableTypes)
	if err != nil {
		return nil, nil, err
	}

	functionBlock := []ast.Expression{}
	if validateFunctionBlock {
		_, hasReturnTypeVoid := expectedFunction.ReturnType.(*types.Void)
		if !hasReturnTypeVoid && len(block) == 0 {
			return nil, nil, type_error.PtrOnNodef(lambda.Node, "Function has return type of %s but has empty body", printableName(expectedFunction.ReturnType))
		}
		for i, blockExp := range block {
			if i < len(block)-1 {
				u, astExp, err := determineTypeOfExpressionBox(true, blockExp, localUniverse)
				if err != nil {
					return nil, nil, err
				}
				functionBlock = append(functionBlock, astExp)
				localUniverse = u
			} else {
				_, astExp, err := expectTypeOfExpressionBox(true, blockExp, caseFunction.ReturnType, localUniverse)
				if err != nil {
					return nil, nil, err
				}
				functionBlock = append(functionBlock, astExp)
			}
		}
	}
	programExp := ast.Function{
		VariableType: caseFunction,
		Block:        functionBlock,
	}

	if annotatedReturnType == nil {
		return universe, programExp, nil
	}
	varType, err := validateTypeAnnotationInUniverse(*annotatedReturnType, localUniverse)
	if err != nil {
		return nil, nil, err
	}

	if !variableTypeEq(varType, expectedFunction.ReturnType) {
		return nil, nil, type_error.PtrOnNodef(lambda.Node, "in return type expected type %s but you have annotated %s", printableName(expectedFunction.ReturnType), printableNameOfTypeAnnotation(*annotatedReturnType))
	}
	return universe, programExp, nil
}

func variableTypeEq(v1 types.VariableType, v2 types.VariableType) bool {
	if v1 == nil || v2 == nil {
		panic(fmt.Errorf("trying to eq %v to %v", v1, v2))
	}
	v1CaseTypeArgument, v1CaseStruct, v1CaseInterface, v1CaseFunction, v1CaseBasicType, v1CaseVoid := v1.VariableTypeCases()
	_ = v1CaseStruct
	_ = v1CaseInterface
	_ = v1CaseBasicType
	_ = v1CaseVoid
	v2CaseTypeArgument, v2CaseStruct, v2CaseInterface, v2CaseFunction, v2CaseBasicType, v2CaseVoid := v2.VariableTypeCases()
	_ = v2CaseStruct
	_ = v2CaseInterface
	_ = v2CaseBasicType
	_ = v2CaseVoid
	if v1CaseStruct != nil && v2CaseStruct != nil {
		return v1CaseStruct.Name == v2CaseStruct.Name
	}
	if v1CaseTypeArgument != nil && v2CaseTypeArgument != nil {
		return v1CaseTypeArgument.Name == v2CaseTypeArgument.Name
	}
	if v1CaseFunction != nil && v2CaseFunction != nil {
		f1 := *v1CaseFunction
		f2 := *v2CaseFunction
		if len(f1.Arguments) != len(f2.Arguments) {
			return false
		}
		for i, f1Arg := range f1.Arguments {
			if !variableTypeEq(f1Arg.VariableType, f2.Arguments[i].VariableType) {
				return false
			}
		}
		return variableTypeEq(f1.ReturnType, f2.ReturnType)
	}
	return reflect.DeepEqual(v1, v2)
}
