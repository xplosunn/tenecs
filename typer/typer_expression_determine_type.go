package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/program"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func determineVariableTypeOfExpression(variableName string, expression parser.Expression, universe binding.Universe) (binding.Universe, program.Expression, *type_error.TypecheckError) {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := expression.Cases()
	if caseLiteralExp != nil {
		return universe, determineVariableTypeOfLiteral(caseLiteralExp.Literal), nil
	} else if caseReferenceOrInvocation != nil {
		varType, err := determineVariableTypeOfReferenceOrInvocation(*caseReferenceOrInvocation, universe)
		return universe, varType, err
	} else if caseLambda != nil {
		var functionUniqueId string
		functionUniqueId, universe = binding.CopyAddingParserFunctionGeneratingUniqueId(universe, *caseLambda)
		function := types.Function{
			Arguments:  []types.FunctionArgument{},
			ReturnType: nil,
		}
		parameters, annotatedReturnType, block := parser.LambdaFields(*caseLambda)
		_ = block
		for _, parameter := range parameters {
			if parameter.Type == nil {
				return nil, nil, type_error.PtrTypeCheckErrorf("parameter '%s' needs to be type annotated as the variable '%s' is not public", parameter.Name, variableName)
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
			return nil, nil, type_error.PtrTypeCheckErrorf("return type needs to be type annotated as the variable '%s' is not public", variableName)
		}
		varType, err := validateTypeAnnotationInUniverse(*annotatedReturnType, universe)
		if err != nil {
			return nil, nil, err
		}
		function.ReturnType = varType
		programExp := program.Function{
			UniqueId:     functionUniqueId,
			VariableType: function,
			Block:        nil,
		}
		return universe, programExp, nil
	} else if caseDeclaration != nil {
		fieldName, fieldExpression := parser.DeclarationFields(*caseDeclaration)
		updatedUniverse, programExp, err := determineVariableTypeOfExpression(fieldName, fieldExpression, universe)
		if err != nil {
			return nil, nil, err
		}
		varType := program.VariableTypeOfExpression(programExp)
		updatedUniverse, err = binding.CopyAddingVariable(updatedUniverse, fieldName, varType)
		if err != nil {
			return nil, nil, err
		}
		declarationProgramExp := program.Declaration{
			VariableType: void,
			Name:         fieldName,
			Expression:   programExp,
		}
		return updatedUniverse, declarationProgramExp, nil
	} else if caseIf != nil {
		updatedUniverse, programExp, err := determineVariableTypeOfIf(*caseIf, universe)
		if err != nil {
			return nil, nil, err
		}
		return updatedUniverse, programExp, nil
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func determineVariableTypeOfLiteral(literal parser.Literal) program.Expression {
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
	return program.Literal{
		VariableType: varType,
		Literal:      literal,
	}
}

func determineVariableTypeOfIf(caseIf parser.If, universe binding.Universe) (binding.Universe, program.Expression, *type_error.TypecheckError) {
	u2, conditionProgramExp, err := expectVariableTypeOfExpression(caseIf.Condition, basicTypeBoolean, universe)
	if err != nil {
		return nil, nil, err
	}
	universe = u2

	varTypeOfBlock := func(expressions []parser.Expression, universe binding.Universe) (binding.Universe, []program.Expression, types.VariableType, *type_error.TypecheckError) {
		if len(expressions) == 0 {
			return universe, []program.Expression{}, void, nil
		}
		localUniverse := universe
		programExpressions := []program.Expression{}
		for i, exp := range expressions {
			u, programExp, err := determineVariableTypeOfExpression("//", exp, localUniverse)
			if err != nil {
				return nil, nil, nil, err
			}
			localUniverse = u
			varType := program.VariableTypeOfExpression(programExp)
			universe, err = binding.ImportParserFunctionsFrom(universe, localUniverse)
			if err != nil {
				return nil, nil, nil, err
			}
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
		return universe, program.If{
			VariableType: thenVarType,
			Condition:    conditionProgramExp,
			ThenBlock:    thenProgramExpressions,
			ElseBlock:    elseProgramExpressions,
		}, nil
	} else {
		return universe, program.If{
			VariableType: void,
			Condition:    conditionProgramExp,
			ThenBlock:    thenProgramExpressions,
			ElseBlock:    []program.Expression{},
		}, nil
	}
}

func determineVariableTypeOfReferenceOrInvocation(referenceOrInvocation parser.ReferenceOrInvocation, universe binding.Universe) (program.Expression, *type_error.TypecheckError) {
	dotSeparatedVarName, argumentsPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)

	if len(referenceOrInvocation.DotSeparatedVars) == 1 {
		constructor, ok := binding.GetConstructorByName(universe, referenceOrInvocation.DotSeparatedVars[0])
		if ok {
			if argumentsPtr == nil {
				varType := types.Function{
					Arguments:  constructor.Arguments,
					ReturnType: constructor.ReturnType,
				}
				programExp := program.ReferenceOrInvocation{
					VariableType:     varType,
					DotSeparatedVars: dotSeparatedVarName,
					Arguments:        nil,
				}
				return programExp, nil
			} else {
				arguments := *argumentsPtr
				if len(arguments) != len(constructor.Arguments) {
					return nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(constructor.Arguments), len(arguments))}
				}
				argumentProgramExpressions := []program.Expression{}
				for i2, argument := range arguments {
					expectedType := constructor.Arguments[i2].VariableType
					_, programExp, err := expectVariableTypeOfExpression(argument, expectedType, universe)
					if err != nil {
						return nil, err
					}
					argumentProgramExpressions = append(argumentProgramExpressions, programExp)
				}
				programExp := program.ReferenceOrInvocation{
					VariableType:     constructor.ReturnType,
					DotSeparatedVars: dotSeparatedVarName,
					Arguments: &program.ArgumentsList{
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
					programExp := program.ReferenceOrInvocation{
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
					programExp := program.ReferenceOrInvocation{
						VariableType:     varType,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments:        nil,
					}
					return programExp, nil
				} else {
					arguments := *argumentsPtr
					if len(arguments) != len(caseFunction.Arguments) {
						return nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(caseFunction.Arguments), len(arguments))}
					}
					argumentProgramExpressions := []program.Expression{}
					for i2, argument := range arguments {
						expectedType := caseFunction.Arguments[i2].VariableType
						_, programExp, err := expectVariableTypeOfExpression(argument, expectedType, universe)
						if err != nil {
							return nil, err
						}
						argumentProgramExpressions = append(argumentProgramExpressions, programExp)
					}
					programExp := program.ReferenceOrInvocation{
						VariableType:     caseFunction.ReturnType,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments: &program.ArgumentsList{
							Arguments: argumentProgramExpressions,
						},
					}
					return programExp, nil
				}
			} else if caseBasicType != nil {
				if argumentsPtr == nil {
					programExp := program.ReferenceOrInvocation{
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
					programExp := program.ReferenceOrInvocation{
						VariableType:     *caseVoid,
						DotSeparatedVars: dotSeparatedVarName,
						Arguments:        nil,
					}
					return programExp, nil
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
