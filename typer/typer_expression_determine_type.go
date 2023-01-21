package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func determineTypeOfExpression(validateFunctionBlock bool, variableName string, expression parser.Expression, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := expression.Cases()
	if caseLiteralExp != nil {
		return universe, determineTypeOfLiteral(caseLiteralExp.Literal), nil
	} else if caseReferenceOrInvocation != nil {
		varType, err := determineTypeOfReferenceOrInvocation(validateFunctionBlock, *caseReferenceOrInvocation, universe)
		return universe, varType, err
	} else if caseLambda != nil {
		localUniverse := universe
		generics, parameters, annotatedReturnType, block := parser.LambdaFields(*caseLambda)
		_ = block
		function := types.Function{
			Generics:   generics,
			Arguments:  []types.FunctionArgument{},
			ReturnType: nil,
		}
		for _, generic := range generics {
			u, err := binding.CopyAddingType(localUniverse, generic, types.TypeArgument{Name: generic})
			if err != nil {
				return nil, nil, err
			}
			localUniverse = u
		}
		for _, parameter := range parameters {
			if parameter.Type == nil {
				return nil, nil, type_error.PtrTypeCheckErrorf("parameter '%s' needs to be type annotated as the variable '%s' is not public", parameter.Name, variableName)
			}

			varType, err := validateTypeAnnotationInUniverse(*parameter.Type, localUniverse)
			if err != nil {
				return nil, nil, err
			}
			function.Arguments = append(function.Arguments, types.FunctionArgument{
				Name:         parameter.Name,
				VariableType: varType,
			})
		}
		if annotatedReturnType == nil {
			return nil, nil, type_error.PtrTypeCheckErrorf("return type needs to be type annotated as the variable '%s' is not public", variableName)
		}
		varType, err := validateTypeAnnotationInUniverse(*annotatedReturnType, localUniverse)
		if err != nil {
			return nil, nil, err
		}
		function.ReturnType = varType

		localUniverse, err = binding.CopyAddingFunctionArguments(localUniverse, function.Arguments)
		if err != nil {
			return nil, nil, err
		}

		functionBlock := []ast.Expression{}
		if validateFunctionBlock {
			if function.ReturnType != void && len(block) == 0 {
				return nil, nil, type_error.PtrTypeCheckErrorf("Function has return type of %s but has empty body", printableName(function.ReturnType))
			}
			for i, blockExp := range block {
				if i < len(block)-1 {
					u, astExp, err := determineTypeOfExpression(true, "===", blockExp, localUniverse)
					if err != nil {
						return nil, nil, err
					}
					functionBlock = append(functionBlock, astExp)
					localUniverse = u
				} else {
					_, astExp, err := expectTypeOfExpression(true, blockExp, varType, localUniverse)
					if err != nil {
						return nil, nil, err
					}
					functionBlock = append(functionBlock, astExp)
				}
			}
		}
		programExp := ast.Function{
			VariableType: function,
			Block:        functionBlock,
		}
		return universe, programExp, nil
	} else if caseDeclaration != nil {
		fieldName, fieldExpression := parser.DeclarationFields(*caseDeclaration)
		updatedUniverse, programExp, err := determineTypeOfExpression(validateFunctionBlock, fieldName, fieldExpression, universe)
		if err != nil {
			return nil, nil, err
		}
		varType := ast.VariableTypeOfExpression(programExp)
		updatedUniverse, err = binding.CopyAddingVariable(updatedUniverse, fieldName, varType)
		if err != nil {
			return nil, nil, err
		}
		declarationProgramExp := ast.Declaration{
			VariableType: void,
			Name:         fieldName,
			Expression:   programExp,
		}
		return updatedUniverse, declarationProgramExp, nil
	} else if caseIf != nil {
		updatedUniverse, programExp, err := determineTypeOfIf(validateFunctionBlock, *caseIf, universe)
		if err != nil {
			return nil, nil, err
		}
		return updatedUniverse, programExp, nil
	} else {
		panic(fmt.Errorf("code on %v", expression))
	}
}

func determineTypeOfLiteral(literal parser.Literal) ast.Expression {
	varType := parser.LiteralFold(
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
	return ast.Literal{
		VariableType: varType,
		Literal:      literal,
	}
}

func determineTypeOfIf(validateFunctionBlock bool, caseIf parser.If, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	u2, conditionProgramExp, err := expectTypeOfExpression(validateFunctionBlock, caseIf.Condition, basicTypeBoolean, universe)
	if err != nil {
		return nil, nil, err
	}
	universe = u2

	varTypeOfBlock := func(expressions []parser.Expression, universe binding.Universe) (binding.Universe, []ast.Expression, types.VariableType, *type_error.TypecheckError) {
		if len(expressions) == 0 {
			return universe, []ast.Expression{}, void, nil
		}
		localUniverse := universe
		programExpressions := []ast.Expression{}
		for i, exp := range expressions {
			u, programExp, err := determineTypeOfExpression(validateFunctionBlock, "//", exp, localUniverse)
			if err != nil {
				return nil, nil, nil, err
			}
			localUniverse = u
			varType := ast.VariableTypeOfExpression(programExp)
			programExpressions = append(programExpressions, programExp)
			if i == len(expressions)-1 {
				return universe, programExpressions, varType, nil
			}
		}
		panic("should have returned before")
	}
	u2, thenProgramExpressions, thenVarType, err := varTypeOfBlock(caseIf.ThenBlock, universe)
	if err != nil {
		return nil, nil, err
	}
	universe = u2
	if len(caseIf.ElseBlock) > 0 {
		u2, elseProgramExpressions, elseVarType, err := varTypeOfBlock(caseIf.ThenBlock, universe)
		if err != nil {
			return nil, nil, err
		}
		universe = u2
		if !variableTypeEq(thenVarType, elseVarType) {
			return nil, nil, type_error.PtrTypeCheckErrorf("if and else blocks should yield the same type, but if is %s and then is %s", printableName(thenVarType), printableName(elseVarType))
		}
		return universe, ast.If{
			VariableType: thenVarType,
			Condition:    conditionProgramExp,
			ThenBlock:    thenProgramExpressions,
			ElseBlock:    elseProgramExpressions,
		}, nil
	} else {
		return universe, ast.If{
			VariableType: void,
			Condition:    conditionProgramExp,
			ThenBlock:    thenProgramExpressions,
			ElseBlock:    []ast.Expression{},
		}, nil
	}
}

func determineTypeOfReferenceOrInvocation(validateFunctionBlock bool, referenceOrInvocation parser.ReferenceOrInvocation, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	dotSeparatedVarName, argumentsPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)

	if len(referenceOrInvocation.DotSeparatedVars) == 1 {
		constructor, ok := binding.GetConstructorByName(universe, referenceOrInvocation.DotSeparatedVars[0])
		if ok {
			if argumentsPtr == nil {
				varType := types.Function{
					Generics:   constructor.Generics,
					Arguments:  constructor.Arguments,
					ReturnType: types.VariableTypeFromConstructableVariableType(constructor.ReturnType),
				}
				programExp := ast.ReferenceOrInvocation{
					VariableType:     varType,
					DotSeparatedVars: dotSeparatedVarName,
					Arguments:        nil,
				}
				return programExp, nil
			} else {
				argumentsList := *argumentsPtr
				if len(argumentsList.Arguments) != len(constructor.Arguments) {
					return nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(constructor.Arguments), len(argumentsList.Arguments))}
				}
				if len(argumentsList.Generics) != len(constructor.Generics) {
					return nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d generics annotated but got %d", len(constructor.Generics), len(argumentsList.Generics))}
				}
				argumentProgramExpressions := []ast.Expression{}
				for i2, argument := range argumentsList.Arguments {
					expectedType := constructor.Arguments[i2].VariableType
					expectedTypeArg, isGeneric := expectedType.(types.TypeArgument)
					if isGeneric {
						caseFunctionGenericIndex := -1
						for index, functionGeneric := range constructor.Generics {
							if functionGeneric == expectedTypeArg.Name {
								caseFunctionGenericIndex = index
								break
							}
						}
						if caseFunctionGenericIndex == -1 {
							return nil, &type_error.TypecheckError{Message: fmt.Sprintf("unexpected error not found generic %s", expectedTypeArg.Name)}
						}
						invocationGeneric := argumentsList.Generics[caseFunctionGenericIndex]
						newExpectedType, err := validateTypeAnnotationInUniverse(parser.SingleNameType{TypeName: invocationGeneric}, universe)
						if err != nil {
							return nil, &type_error.TypecheckError{Message: fmt.Sprintf("not found annotated generic type %s", invocationGeneric)}
						}
						expectedType = newExpectedType
					}
					_, programExp, err := expectTypeOfExpression(validateFunctionBlock, argument, expectedType, universe)
					if err != nil {
						return nil, err
					}
					argumentProgramExpressions = append(argumentProgramExpressions, programExp)
				}
				returnType := constructor.ReturnType
				caseStruct, caseInterface := returnType.ConstructableCases()
				_ = caseInterface
				if caseStruct != nil && len(constructor.Generics) > 0 {
					if caseStruct.ResolvedTypeArguments == nil {
						caseStruct.ResolvedTypeArguments = []types.ResolvedTypeArgument{}
					}
					for i, generic := range argumentsList.Generics {
						genericVarType, err := validateTypeAnnotationInUniverse(parser.SingleNameType{TypeName: generic}, universe)
						if err != nil {
							return nil, &type_error.TypecheckError{Message: fmt.Sprintf("not found annotated generic type %s", generic)}
						}
						structVarType, ok := types.StructVariableTypeFromVariableType(genericVarType)
						if !ok {
							return nil, &type_error.TypecheckError{Message: fmt.Sprintf("not a valid annotated generic type %s", generic)}
						}
						caseStruct.ResolvedTypeArguments = append(caseStruct.ResolvedTypeArguments, types.ResolvedTypeArgument{
							Name:               constructor.Generics[i],
							StructVariableType: structVarType,
						})
					}
					returnType = caseStruct
				}
				programExp := ast.ReferenceOrInvocation{
					VariableType:     types.VariableTypeFromConstructableVariableType(returnType),
					DotSeparatedVars: dotSeparatedVarName,
					Arguments: &ast.ArgumentsList{
						Arguments: argumentProgramExpressions,
					},
				}
				return programExp, nil
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
			caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			if caseTypeArgument != nil {
				return nil, type_error.PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
			} else if caseStruct != nil {
				structVariables, err := binding.GetGlobalStructVariables(universe, *caseStruct)
				if err != nil {
					return nil, err
				}
				currentUniverse = binding.NewFromStructVariables(structVariables, universe)
			} else if caseInterface != nil {
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
				panic(fmt.Errorf("code on %v", varType))
			}
		} else {
			caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			if caseTypeArgument != nil {
				if argumentsPtr == nil {
					programExp := ast.ReferenceOrInvocation{
						VariableType:     *caseTypeArgument,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments:        nil,
					}
					return programExp, nil
				} else {
					return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else if caseStruct != nil {
				if argumentsPtr == nil {
					programExp := ast.ReferenceOrInvocation{
						VariableType:     *caseStruct,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments:        nil,
					}
					return programExp, nil
				} else {
					return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else if caseInterface != nil {
				if argumentsPtr == nil {
					programExp := ast.ReferenceOrInvocation{
						VariableType:     *caseInterface,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments:        nil,
					}
					return programExp, nil
				} else {
					return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else if caseFunction != nil {
				if argumentsPtr == nil {
					varType, ok := binding.GetTypeByVariableName(currentUniverse, varName)
					if !ok {
						return nil, &type_error.TypecheckError{Message: "not found in scope: " + varName}
					}
					programExp := ast.ReferenceOrInvocation{
						VariableType:     varType,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments:        nil,
					}
					return programExp, nil
				} else {
					argumentsList := *argumentsPtr
					if len(argumentsList.Arguments) != len(caseFunction.Arguments) {
						return nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(caseFunction.Arguments), len(argumentsList.Arguments))}
					}
					if len(argumentsList.Generics) != len(caseFunction.Generics) {
						return nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d generics annotated but got %d", len(caseFunction.Generics), len(argumentsList.Generics))}
					}
					argumentProgramExpressions := []ast.Expression{}
					for i2, argument := range argumentsList.Arguments {
						expectedType := caseFunction.Arguments[i2].VariableType
						expectedTypeArg, isGeneric := expectedType.(types.TypeArgument)
						if isGeneric {
							caseFunctionGenericIndex := -1
							for index, functionGeneric := range caseFunction.Generics {
								if functionGeneric == expectedTypeArg.Name {
									caseFunctionGenericIndex = index
									break
								}
							}
							if caseFunctionGenericIndex == -1 {
								return nil, &type_error.TypecheckError{Message: fmt.Sprintf("unexpected error not found generic %s", expectedTypeArg.Name)}
							}
							invocationGeneric := argumentsList.Generics[caseFunctionGenericIndex]
							newExpectedType, err := validateTypeAnnotationInUniverse(parser.SingleNameType{TypeName: invocationGeneric}, universe)
							if err != nil {
								return nil, &type_error.TypecheckError{Message: fmt.Sprintf("not found annotated generic type %s", invocationGeneric)}
							}
							expectedType = newExpectedType
						}
						_, programExp, err := expectTypeOfExpression(validateFunctionBlock, argument, expectedType, universe)
						if err != nil {
							return nil, err
						}
						argumentProgramExpressions = append(argumentProgramExpressions, programExp)
					}
					returnType := caseFunction.ReturnType
					returnTypeArg, isGeneric := returnType.(types.TypeArgument)
					if isGeneric {
						caseFunctionGenericIndex := -1
						for index, functionGeneric := range caseFunction.Generics {
							if functionGeneric == returnTypeArg.Name {
								caseFunctionGenericIndex = index
								break
							}
						}
						if caseFunctionGenericIndex == -1 {
							return nil, &type_error.TypecheckError{Message: fmt.Sprintf("unexpected error not found return generic %s", returnTypeArg.Name)}
						}
						invocationGeneric := argumentsList.Generics[caseFunctionGenericIndex]
						newReturnType, err := validateTypeAnnotationInUniverse(parser.SingleNameType{TypeName: invocationGeneric}, universe)
						if err != nil {
							return nil, &type_error.TypecheckError{Message: fmt.Sprintf("not found return generic type %s", invocationGeneric)}
						}
						returnType = newReturnType
					}
					programExp := ast.ReferenceOrInvocation{
						VariableType:     returnType,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments: &ast.ArgumentsList{
							Arguments: argumentProgramExpressions,
						},
					}
					return programExp, nil
				}
			} else if caseBasicType != nil {
				if argumentsPtr == nil {
					programExp := ast.ReferenceOrInvocation{
						VariableType:     *caseBasicType,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments:        nil,
					}
					return programExp, nil
				} else {
					return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else if caseVoid != nil {
				if argumentsPtr == nil {
					programExp := ast.ReferenceOrInvocation{
						VariableType:     *caseVoid,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments:        nil,
					}
					return programExp, nil
				} else {
					return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				}
			} else {
				panic(fmt.Errorf("code on %v", varType))
			}
		}
	}

	panic("empty dotSeparatedVarName")
}
