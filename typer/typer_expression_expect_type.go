package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"reflect"
	"strings"
)

func expectVariableTypeOfExpression(exp parser.Expression, expectedType VariableType, universe Universe) *TypecheckError {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := exp.Cases()
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
	} else if caseDeclaration != nil {
		if !variableTypeEq(expectedType, void) {
			return PtrTypeCheckErrorf("expected type %s but found Void (variable declarations return void)", printableName(expectedType))
		}
		return nil
	} else if caseIf != nil {
		varType, err := determineVariableTypeOfIf(*caseIf, universe)
		if err != nil {
			return err
		}
		if !variableTypeEq(varType, expectedType) {
			return PtrTypeCheckErrorf("expected type %s but found %s", printableName(expectedType), printableName(varType))
		}
		return nil
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
		if parameter.Type == nil {
			continue
		}

		varType, err := validateTypeAnnotationInUniverse(*parameter.Type, universe)
		if err != nil {
			return err
		}

		if !variableTypeEq(varType, expectedFunction.Arguments[i].VariableType) {
			return PtrTypeCheckErrorf("in parameter position %d expected type %s but you have annotated %s", i, printableName(expectedFunction.Arguments[i].VariableType), printableNameOfTypeAnnotation(*parameter.Type))
		}
	}

	if annotatedReturnType == nil {
		return nil
	}
	varType, err := validateTypeAnnotationInUniverse(*annotatedReturnType, universe)
	if err != nil {
		return err
	}

	if !variableTypeEq(varType, expectedFunction.ReturnType) {
		return PtrTypeCheckErrorf("in return type expected type %s but you have annotated %s", printableName(expectedFunction.ReturnType), printableNameOfTypeAnnotation(*annotatedReturnType))
	}
	return nil
}

func variableTypeEq(v1 VariableType, v2 VariableType) bool {
	v1CaseInterface, v1CaseFunction, v1CaseBasicType, v1CaseVoid := v1.Cases()
	_ = v1CaseInterface
	_ = v1CaseBasicType
	_ = v1CaseVoid
	v2CaseInterface, v2CaseFunction, v2CaseBasicType, v2CaseVoid := v2.Cases()
	_ = v2CaseInterface
	_ = v2CaseBasicType
	_ = v2CaseVoid
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
