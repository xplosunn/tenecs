package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"reflect"
	"strings"
)

func expectVariableTypeOfExpression(exp parser.Expression, expectedType VariableType, universe Universe) *TypecheckError {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda := exp.Cases()
	if caseLiteralExp != nil {
		varType := determineVariableTypeOfLiteral(caseLiteralExp.Literal)
		if !variableTypeEq(varType, expectedType) {
			return PtrTypeCheckErrorf("expected type %s but found %s", printableName(expectedType), printableName(varType))
		}
		return nil
	} else if caseReferenceOrInvocation != nil {
		varType, err := determineVariableTypeOfReferenceOrInvocation(*caseReferenceOrInvocation, universe)
		if err != nil {
			return err
		}
		if !variableTypeEq(varType, expectedType) {
			return PtrTypeCheckErrorf("in expression '%s' expected %s but found %s", strings.Join(caseReferenceOrInvocation.DotSeparatedVars, "."), printableName(expectedType), printableName(varType))
		}
		return nil
	} else if caseLambda != nil {
		return expectVariableTypeOfLambdaSignature(*caseLambda, expectedType, universe)
	} else {
		panic(fmt.Errorf("cases on %v", exp))
	}
}

func expectVariableTypeOfLambdaSignature(lambda parser.Lambda, expectedType VariableType, universe Universe) *TypecheckError {
	var expectedFunction Function
	caseInterface, caseFunction, caseBasicType, caseVoid := expectedType.Cases()
	if caseInterface != nil {
		return PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else if caseFunction != nil {
		expectedFunction = *caseFunction
	} else if caseBasicType != nil {
		return PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else if caseVoid != nil {
		return PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else {
		panic(fmt.Errorf("cases on %v", expectedType))
	}

	parameters, annotatedReturnType, block := parser.LambdaFields(lambda)
	_ = block
	if len(parameters) != len(expectedFunction.Arguments) {
		return PtrTypeCheckErrorf("expected same number of arguments as interface variable (%d) but found %d", len(expectedFunction.Arguments), len(parameters))
	}
	for i, parameter := range parameters {
		if parameter.Type == "" {
			continue
		}

		varType, ok := universe.TypeByTypeName.Get(parameter.Type)
		if !ok {
			return PtrTypeCheckErrorf("not found type: %s", parameter.Type)
		}

		if !variableTypeEq(varType, expectedFunction.Arguments[i].VariableType) {
			return PtrTypeCheckErrorf("in parameter position %d expected type %s but you have annotated %s", i, printableName(expectedFunction.Arguments[i].VariableType), parameter.Type)
		}
	}

	if annotatedReturnType == "" {
		return nil
	}
	varType, ok := universe.TypeByTypeName.Get(annotatedReturnType)
	if !ok {
		return PtrTypeCheckErrorf("not found type: %s", annotatedReturnType)
	}

	if !variableTypeEq(varType, expectedFunction.ReturnType) {
		return PtrTypeCheckErrorf("in return type expected type %s but you have annotated %s", printableName(expectedFunction.ReturnType), annotatedReturnType)
	}
	return nil
}

func variableTypeEq(v1 VariableType, v2 VariableType) bool {
	return reflect.DeepEqual(v1, v2)
}
