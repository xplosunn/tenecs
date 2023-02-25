package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
	"reflect"
)

func expectTypeOfExpressionBox(validateFunctionBlock bool, expressionBox parser.ExpressionBox, expectedType types.VariableType, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	noInvocation := expressionBox.AccessOrInvocationChain == nil
	caseExpectFunction, expectsFunction := expectedType.(types.Function)
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
		return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found %s", printableName(expectedType), printableName(varType))
	}
	return universe, astExp, nil
}

func expectTypeOfExpression(validateFunctionBlock bool, exp parser.Expression, expectedType types.VariableType, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	caseModule, caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := exp.ExpressionCases()
	if caseModule != nil {
		return determineTypeOfModule(validateFunctionBlock, *caseModule, universe)
	} else if caseLiteralExp != nil {
		programExp := determineTypeOfLiteral(caseLiteralExp.Literal)
		varType := ast.VariableTypeOfExpression(programExp)
		if !variableTypeEq(varType, expectedType) {
			return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found %s", printableName(expectedType), printableName(varType))
		}
		return universe, programExp, nil
	} else if caseReferenceOrInvocation != nil {
		programExp, err := determineTypeOfReferenceOrInvocation(validateFunctionBlock, *caseReferenceOrInvocation, universe)
		if err != nil {
			return nil, nil, err
		}
		varType := ast.VariableTypeOfExpression(programExp)
		if !variableTypeEq(varType, expectedType) {
			return nil, nil, type_error.PtrTypeCheckErrorf("in expression '%s' expected %s but found %s", caseReferenceOrInvocation.Var, printableName(expectedType), printableName(varType))
		}
		return universe, programExp, nil
	} else if caseLambda != nil {
		return expectTypeOfLambda(validateFunctionBlock, *caseLambda, expectedType, universe)
	} else if caseDeclaration != nil {
		universe, programExp, err := determineTypeOfExpressionBox(validateFunctionBlock, caseDeclaration.ExpressionBox, universe)
		if err != nil {
			return nil, nil, err
		}
		if !variableTypeEq(expectedType, void) {
			return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found Void (variable declarations return void)", printableName(expectedType))
		}
		return universe, programExp, nil
	} else if caseIf != nil {
		universe, programExp, err := determineTypeOfIf(validateFunctionBlock, *caseIf, universe)
		if err != nil {
			return nil, nil, err
		}
		varType := ast.VariableTypeOfExpression(programExp)
		if !variableTypeEq(varType, expectedType) {
			return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found %s", printableName(expectedType), printableName(varType))
		}
		return universe, programExp, nil
	} else {
		panic(fmt.Errorf("code on %v", exp))
	}
}

func expectTypeOfLambda(validateFunctionBlock bool, lambda parser.Lambda, expectedType types.VariableType, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	var expectedFunction types.Function
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := expectedType.VariableTypeCases()
	if caseTypeArgument != nil {
		return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else if caseStruct != nil {
		return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else if caseInterface != nil {
		return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else if caseFunction != nil {
		expectedFunction = *caseFunction
	} else if caseBasicType != nil {
		return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else if caseVoid != nil {
		return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else {
		panic(fmt.Errorf("code on %v", expectedType))
	}

	functionArgs := []types.FunctionArgument{}
	generics, parameters, annotatedReturnType, block := parser.LambdaFields(lambda)
	_ = block
	if len(generics) != len(expectedFunction.Generics) {
		return nil, nil, type_error.PtrTypeCheckErrorf("expected same number of generics as interface variable (%d) but found %d", len(expectedFunction.Generics), len(generics))
	}
	if len(parameters) != len(expectedFunction.Arguments) {
		return nil, nil, type_error.PtrTypeCheckErrorf("expected same number of arguments as interface variable (%d) but found %d", len(expectedFunction.Arguments), len(parameters))
	}
	localUniverse := universe
	for _, generic := range generics {
		u, err := binding.CopyAddingType(localUniverse, generic, types.TypeArgument{Name: generic})
		if err != nil {
			return nil, nil, err
		}
		localUniverse = u
	}
	for i, parameter := range parameters {
		if parameter.Type == nil {
			functionArgs = append(functionArgs, types.FunctionArgument{
				Name:         parameter.Name,
				VariableType: expectedFunction.Arguments[i].VariableType,
			})
			continue
		}

		varType, err := validateTypeAnnotationInUniverse(*parameter.Type, localUniverse)
		if err != nil {
			return nil, nil, err
		}

		if !variableTypeEq(varType, expectedFunction.Arguments[i].VariableType) {
			return nil, nil, type_error.PtrTypeCheckErrorf("in parameter position %d expected type %s but you have annotated %s", i, printableName(expectedFunction.Arguments[i].VariableType), printableNameOfTypeAnnotation(*parameter.Type))
		}

		functionArgs = append(functionArgs, types.FunctionArgument{
			Name:         parameter.Name,
			VariableType: varType,
		})
	}
	localUniverse, err := binding.CopyAddingFunctionArguments(localUniverse, functionArgs)
	if err != nil {
		return nil, nil, err
	}

	functionBlock := []ast.Expression{}
	if validateFunctionBlock {
		if expectedFunction.ReturnType != void && len(block) == 0 {
			return nil, nil, type_error.PtrTypeCheckErrorf("Function has return type of %s but has empty body", printableName(expectedFunction.ReturnType))
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
		VariableType: *caseFunction,
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
		return nil, nil, type_error.PtrTypeCheckErrorf("in return type expected type %s but you have annotated %s", printableName(expectedFunction.ReturnType), printableNameOfTypeAnnotation(*annotatedReturnType))
	}
	return universe, programExp, nil
}

func variableTypeEq(v1 types.VariableType, v2 types.VariableType) bool {
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
