package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
)

func determineVariableTypeOfExpression(variableName string, expression parser.Expression, universe Universe) (VariableType, *TypecheckError) {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda := expression.Cases()
	if caseLiteralExp != nil {
		return determineVariableTypeOfLiteral(caseLiteralExp.Literal), nil
	} else if caseReferenceOrInvocation != nil {
		return determineVariableTypeOfReferenceOrInvocation(*caseReferenceOrInvocation, universe)
	} else if caseLambda != nil {
		function := Function{
			Arguments:  []FunctionArgument{},
			ReturnType: nil,
		}
		parameters, annotatedReturnType, block := parser.LambdaFields(*caseLambda)
		_ = block
		for _, parameter := range parameters {
			if parameter.Type == "" {
				return nil, PtrTypeCheckErrorf("parameter '%s' needs to be type annotated as the variable '%s' is not public", parameter.Name, variableName)
			}

			varType, ok := universe.TypeByTypeName.Get(parameter.Type)
			if !ok {
				return nil, PtrTypeCheckErrorf("not found type: %s", parameter.Type)
			}
			function.Arguments = append(function.Arguments, FunctionArgument{
				Name:         parameter.Name,
				VariableType: varType,
			})
		}
		if annotatedReturnType == "" {
			return nil, PtrTypeCheckErrorf("return type needs to be type annotated as the variable '%s' is not public", variableName)
		}
		varType, ok := universe.TypeByTypeName.Get(annotatedReturnType)
		if !ok {
			return nil, PtrTypeCheckErrorf("not found type: %s", annotatedReturnType)
		}
		function.ReturnType = varType
		return function, nil
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
