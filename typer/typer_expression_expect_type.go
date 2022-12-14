package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
	"reflect"
	"strings"
)

func expectTypeOfExpression(exp parser.Expression, expectedType types.VariableType, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := exp.Cases()
	if caseLiteralExp != nil {
		programExp := determineTypeOfLiteral(caseLiteralExp.Literal)
		varType := ast.VariableTypeOfExpression(programExp)
		if !variableTypeEq(varType, expectedType) {
			return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found %s", printableName(expectedType), printableName(varType))
		}
		return universe, programExp, nil
	} else if caseReferenceOrInvocation != nil {
		programExp, err := determineTypeOfReferenceOrInvocation(*caseReferenceOrInvocation, universe)
		if err != nil {
			return nil, nil, err
		}
		varType := ast.VariableTypeOfExpression(programExp)
		if !variableTypeEq(varType, expectedType) {
			return nil, nil, type_error.PtrTypeCheckErrorf("in expression '%s' expected %s but found %s", strings.Join(caseReferenceOrInvocation.DotSeparatedVars, "."), printableName(expectedType), printableName(varType))
		}
		return universe, programExp, nil
	} else if caseLambda != nil {
		return expectTypeOfLambdaSignature(*caseLambda, expectedType, universe)
	} else if caseDeclaration != nil {
		universe, programExp, err := determineTypeOfExpression("%%", caseDeclaration.Expression, universe)
		if err != nil {
			return nil, nil, err
		}
		if !variableTypeEq(expectedType, void) {
			return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found Void (variable declarations return void)", printableName(expectedType))
		}
		return universe, programExp, nil
	} else if caseIf != nil {
		universe, programExp, err := determineTypeOfIf(*caseIf, universe)
		if err != nil {
			return nil, nil, err
		}
		varType := ast.VariableTypeOfExpression(programExp)
		if !variableTypeEq(varType, expectedType) {
			return nil, nil, type_error.PtrTypeCheckErrorf("expected type %s but found %s", printableName(expectedType), printableName(varType))
		}
		return universe, programExp, nil
	} else {
		panic(fmt.Errorf("cases on %v", exp))
	}
}

func expectTypeOfLambdaSignature(lambda parser.Lambda, expectedType types.VariableType, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	var functionUniqueId string
	functionUniqueId, universe = binding.CopyAddingParserFunctionGeneratingUniqueId(universe, lambda)

	var expectedFunction types.Function
	caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := expectedType.Cases()
	if caseStruct != nil {
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
		panic(fmt.Errorf("cases on %v", expectedType))
	}

	parameters, annotatedReturnType, block := parser.LambdaFields(lambda)
	_ = block
	if len(parameters) != len(expectedFunction.Arguments) {
		return nil, nil, type_error.PtrTypeCheckErrorf("expected same number of arguments as interface variable (%d) but found %d", len(expectedFunction.Arguments), len(parameters))
	}
	for i, parameter := range parameters {
		if parameter.Type == nil {
			continue
		}

		varType, err := validateTypeAnnotationInUniverse(*parameter.Type, universe)
		if err != nil {
			return nil, nil, err
		}

		if !variableTypeEq(varType, expectedFunction.Arguments[i].VariableType) {
			return nil, nil, type_error.PtrTypeCheckErrorf("in parameter position %d expected type %s but you have annotated %s", i, printableName(expectedFunction.Arguments[i].VariableType), printableNameOfTypeAnnotation(*parameter.Type))
		}
	}

	programExp := ast.Function{
		UniqueId:     functionUniqueId,
		VariableType: *caseFunction,
		Block:        nil,
	}

	if annotatedReturnType == nil {
		return universe, programExp, nil
	}
	varType, err := validateTypeAnnotationInUniverse(*annotatedReturnType, universe)
	if err != nil {
		return nil, nil, err
	}

	if !variableTypeEq(varType, expectedFunction.ReturnType) {
		return nil, nil, type_error.PtrTypeCheckErrorf("in return type expected type %s but you have annotated %s", printableName(expectedFunction.ReturnType), printableNameOfTypeAnnotation(*annotatedReturnType))
	}
	return universe, programExp, nil
}

func variableTypeEq(v1 types.VariableType, v2 types.VariableType) bool {
	v1CaseStruct, v1CaseInterface, v1CaseFunction, v1CaseBasicType, v1CaseVoid := v1.Cases()
	_ = v1CaseStruct
	_ = v1CaseInterface
	_ = v1CaseBasicType
	_ = v1CaseVoid
	v2CaseStruct, v2CaseInterface, v2CaseFunction, v2CaseBasicType, v2CaseVoid := v2.Cases()
	_ = v2CaseStruct
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
