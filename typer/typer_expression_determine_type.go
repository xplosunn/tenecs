package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
)

func determineVariableTypeOfExpression(variableName string, expression parser.Expression, universe Universe) (Universe, VariableType, *TypecheckError) {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := expression.Cases()
	if caseLiteralExp != nil {
		return universe, determineVariableTypeOfLiteral(caseLiteralExp.Literal), nil
	} else if caseReferenceOrInvocation != nil {
		varType, err := determineVariableTypeOfReferenceOrInvocation(*caseReferenceOrInvocation, universe)
		return universe, varType, err
	} else if caseLambda != nil {
		function := Function{
			Arguments:  []FunctionArgument{},
			ReturnType: nil,
		}
		parameters, annotatedReturnType, block := parser.LambdaFields(*caseLambda)
		_ = block
		for _, parameter := range parameters {
			if parameter.Type == nil {
				return universe, nil, PtrTypeCheckErrorf("parameter '%s' needs to be type annotated as the variable '%s' is not public", parameter.Name, variableName)
			}

			varType, err := validateTypeAnnotationInUniverse(*parameter.Type, universe)
			if err != nil {
				return universe, nil, err
			}
			function.Arguments = append(function.Arguments, FunctionArgument{
				Name:         parameter.Name,
				VariableType: varType,
			})
		}
		if annotatedReturnType == nil {
			return universe, nil, PtrTypeCheckErrorf("return type needs to be type annotated as the variable '%s' is not public", variableName)
		}
		varType, err := validateTypeAnnotationInUniverse(*annotatedReturnType, universe)
		if err != nil {
			return universe, nil, err
		}
		function.ReturnType = varType
		return universe, function, nil
	} else if caseDeclaration != nil {
		fieldName, fieldExpression := parser.DeclarationFields(*caseDeclaration)
		updatedUniverse, variableType, err := determineVariableTypeOfExpression(fieldName, fieldExpression, universe)
		if err != nil {
			return universe, nil, err
		}
		updatedUniverse, err = copyUniverseAddingVariable(updatedUniverse, fieldName, variableType)
		if err != nil {
			return universe, nil, err
		}
		return updatedUniverse, void, nil
	} else if caseIf != nil {
		varType, err := determineVariableTypeOfIf(*caseIf, universe)
		if err != nil {
			return universe, nil, err
		}
		return universe, varType, nil
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func determineVariableTypeOfLiteral(literal parser.Literal) VariableType {
	return parser.LiteralFold(
		literal,
		func(arg float64) BasicType {
			return basicTypeFloat
		},
		func(arg int) BasicType {
			return basicTypeInt
		},
		func(arg string) BasicType {
			return basicTypeString
		},
		func(arg bool) BasicType {
			return basicTypeBoolean
		},
	)
}

func determineVariableTypeOfIf(caseIf parser.If, universe Universe) (VariableType, *TypecheckError) {
	err := expectVariableTypeOfExpression(caseIf.Condition, basicTypeBoolean, universe)
	if err != nil {
		return nil, err
	}

	varTypeOfBlock := func(expressions []parser.Expression) (VariableType, *TypecheckError) {
		if len(expressions) == 0 {
			return void, nil
		}
		localUniverse := universe
		for i, exp := range expressions {
			u, varType, err := determineVariableTypeOfExpression("//", exp, localUniverse)
			if err != nil {
				return nil, err
			}
			localUniverse = u
			if i == len(expressions)-1 {
				return varType, nil
			}
		}
		panic("should have returned before")
	}
	thenVarType, err := varTypeOfBlock(caseIf.ThenBlock)
	if err != nil {
		return nil, err
	}
	if len(caseIf.ElseBlock) > 0 {
		elseVarType, err := varTypeOfBlock(caseIf.ThenBlock)
		if err != nil {
			return nil, err
		}
		if !variableTypeEq(thenVarType, elseVarType) {
			return nil, PtrTypeCheckErrorf("if and else blocks should yield the same type, but if is %s and then is %s", printableName(thenVarType), printableName(elseVarType))
		}
		return thenVarType, nil
	} else {
		return void, nil
	}
}

func determineVariableTypeOfReferenceOrInvocation(referenceOrInvocation parser.ReferenceOrInvocation, universe Universe) (VariableType, *TypecheckError) {
	dotSeparatedVarName, argumentsPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)

	currentUniverse := universe
	for i, varName := range dotSeparatedVarName {
		varType, ok := currentUniverse.TypeByVariableName.Get(varName)
		if !ok {
			return nil, &TypecheckError{Message: "not found in scope: " + varName}
		}

		if i < len(dotSeparatedVarName)-1 {
			caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			if caseInterface != nil {
				currentUniverse = NewUniverseFromInterface(*caseInterface)
			} else if caseFunction != nil {
				return nil, PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
			} else if caseBasicType != nil {
				return nil, PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
			} else if caseVoid != nil {
				return nil, PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
			} else {
				panic(fmt.Errorf("cases on %v", varType))
			}
		} else {
			caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			if caseInterface != nil {
				if argumentsPtr == nil {
					return *caseInterface, nil
				} else {
					return nil, PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else if caseFunction != nil {
				if argumentsPtr == nil {
					varType, ok := currentUniverse.TypeByVariableName.Get(varName)
					if !ok {
						return nil, &TypecheckError{Message: "not found in scope: " + varName}
					}
					return varType, nil
				} else {
					arguments := *argumentsPtr
					if len(arguments) != len(caseFunction.Arguments) {
						return nil, &TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(caseFunction.Arguments), len(arguments))}
					}
					for i2, argument := range arguments {
						expectedType := caseFunction.Arguments[i2].VariableType
						err := expectVariableTypeOfExpression(argument, expectedType, universe)
						if err != nil {
							return nil, err
						}
					}
					return caseFunction.ReturnType, nil
				}
			} else if caseBasicType != nil {
				if argumentsPtr == nil {
					return *caseBasicType, nil
				} else {
					return nil, PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else if caseVoid != nil {
				if argumentsPtr == nil {
					return *caseVoid, nil
				} else {
					return nil, PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else {
				panic(fmt.Errorf("cases on %v", varType))
			}
		}
	}

	panic("empty dotSeparatedVarName")
}
