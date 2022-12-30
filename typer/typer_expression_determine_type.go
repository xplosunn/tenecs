package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func determineVariableTypeOfExpression(variableName string, expression parser.Expression, universe binding.Universe) (binding.Universe, types.VariableType, *type_error.TypecheckError) {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := expression.Cases()
	if caseLiteralExp != nil {
		return universe, determineVariableTypeOfLiteral(caseLiteralExp.Literal), nil
	} else if caseReferenceOrInvocation != nil {
		varType, err := determineVariableTypeOfReferenceOrInvocation(*caseReferenceOrInvocation, universe)
		return universe, varType, err
	} else if caseLambda != nil {
		function := types.Function{
			Arguments:  []types.FunctionArgument{},
			ReturnType: nil,
		}
		parameters, annotatedReturnType, block := parser.LambdaFields(*caseLambda)
		_ = block
		for _, parameter := range parameters {
			if parameter.Type == nil {
				return universe, nil, type_error.PtrTypeCheckErrorf("parameter '%s' needs to be type annotated as the variable '%s' is not public", parameter.Name, variableName)
			}

			varType, err := validateTypeAnnotationInUniverse(*parameter.Type, universe)
			if err != nil {
				return universe, nil, err
			}
			function.Arguments = append(function.Arguments, types.FunctionArgument{
				Name:         parameter.Name,
				VariableType: varType,
			})
		}
		if annotatedReturnType == nil {
			return universe, nil, type_error.PtrTypeCheckErrorf("return type needs to be type annotated as the variable '%s' is not public", variableName)
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
		updatedUniverse, err = binding.CopyAddingVariable(updatedUniverse, fieldName, variableType)
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

func determineVariableTypeOfLiteral(literal parser.Literal) types.VariableType {
	return parser.LiteralFold(
		literal,
		func(arg float64) types.BasicType {
			return basicTypeFloat
		},
		func(arg int) types.BasicType {
			return basicTypeInt
		},
		func(arg string) types.BasicType {
			return basicTypeString
		},
		func(arg bool) types.BasicType {
			return basicTypeBoolean
		},
	)
}

func determineVariableTypeOfIf(caseIf parser.If, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	err := expectVariableTypeOfExpression(caseIf.Condition, basicTypeBoolean, universe)
	if err != nil {
		return nil, err
	}

	varTypeOfBlock := func(expressions []parser.Expression) (types.VariableType, *type_error.TypecheckError) {
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
			return nil, type_error.PtrTypeCheckErrorf("if and else blocks should yield the same type, but if is %s and then is %s", printableName(thenVarType), printableName(elseVarType))
		}
		return thenVarType, nil
	} else {
		return void, nil
	}
}

func determineVariableTypeOfReferenceOrInvocation(referenceOrInvocation parser.ReferenceOrInvocation, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	dotSeparatedVarName, argumentsPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)

	if len(referenceOrInvocation.DotSeparatedVars) == 1 {
		constructor, ok := binding.GetConstructorByName(universe, referenceOrInvocation.DotSeparatedVars[0])
		if ok {
			if argumentsPtr == nil {
				varType := types.Function{
					Arguments:  constructor.Arguments,
					ReturnType: constructor.ReturnType,
				}
				return varType, nil
			} else {
				arguments := *argumentsPtr
				if len(arguments) != len(constructor.Arguments) {
					return nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(constructor.Arguments), len(arguments))}
				}
				for i2, argument := range arguments {
					expectedType := constructor.Arguments[i2].VariableType
					err := expectVariableTypeOfExpression(argument, expectedType, universe)
					if err != nil {
						return nil, err
					}
				}
				return constructor.ReturnType, nil
			}
		}
	}

	currentUniverse := universe
	for i, varName := range dotSeparatedVarName {
		varType, ok := binding.GetTypeByVariableName(currentUniverse, varName)
		if !ok {
			return nil, &type_error.TypecheckError{Message: "not found in scope: " + varName}
		}

		if i < len(dotSeparatedVarName)-1 {
			caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			if caseInterface != nil {
				interfaceVariables, err := binding.GetGlobalInterfaceVariables(universe, *caseInterface)
				if err != nil {
					return nil, err
				}
				currentUniverse = binding.NewFromInterfaceVariables(interfaceVariables, universe)
			} else if caseFunction != nil {
				return nil, type_error.PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
			} else if caseBasicType != nil {
				return nil, type_error.PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
			} else if caseVoid != nil {
				return nil, type_error.PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
			} else {
				panic(fmt.Errorf("cases on %v", varType))
			}
		} else {
			caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			if caseInterface != nil {
				if argumentsPtr == nil {
					return *caseInterface, nil
				} else {
					return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else if caseFunction != nil {
				if argumentsPtr == nil {
					varType, ok := binding.GetTypeByVariableName(currentUniverse, varName)
					if !ok {
						return nil, &type_error.TypecheckError{Message: "not found in scope: " + varName}
					}
					return varType, nil
				} else {
					arguments := *argumentsPtr
					if len(arguments) != len(caseFunction.Arguments) {
						return nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(caseFunction.Arguments), len(arguments))}
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
					return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else if caseVoid != nil {
				if argumentsPtr == nil {
					return *caseVoid, nil
				} else {
					return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else {
				panic(fmt.Errorf("cases on %v", varType))
			}
		}
	}

	panic("empty dotSeparatedVarName")
}
