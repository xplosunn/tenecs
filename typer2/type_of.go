package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func typeOfExpressionBox(expressionBox parser.ExpressionBox, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	expression, accessOrInvocations := parser.ExpressionBoxFields(expressionBox)

	varType, err := typeOfExpression(expression, universe)
	if err != nil {
		return nil, err
	}
	if len(accessOrInvocations) == 0 {
		return varType, nil
	}

	for _, accessOrInvocation := range accessOrInvocations {
		varType, err = typeOfAccess(varType, accessOrInvocation.VarName)
		if err != nil {
			return nil, err
		}
		if accessOrInvocation.Arguments != nil {
			function, ok := varType.(*types.Function)
			if !ok {
				return nil, type_error.PtrOnNodef(accessOrInvocation.Arguments.Node, "Expected a function in order to invoke")
			}
			varType, err = typeOfInvocation(function, *accessOrInvocation.Arguments, universe)
		}
	}

	return varType, nil
}

func typeOfBlock(block []parser.ExpressionBox, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	for i, exp := range block {
		if i == len(block)-1 {
			return typeOfExpressionBox(exp, universe)
		}
		dec, ok := exp.Expression.(parser.Declaration)
		if !ok {
			continue
		}
		decType, err := typeOfExpressionBox(dec.ExpressionBox, universe)
		if err != nil {
			return nil, err
		}
		universe, err = binding.CopyAddingVariable(universe, dec.Name, decType)
		if err != nil {
			return nil, err
		}
	}
	return &types.Void{}, nil
}

func typeOfExpression(expression parser.Expression, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	var varType types.VariableType
	var err *type_error.TypecheckError
	parser.ExpressionExhaustiveSwitch(
		expression,
		func(expression parser.Module) {
			var ok bool
			varType, ok = binding.GetTypeByTypeName(universe, expression.Implementing.String)
			if !ok {
				err = type_error.PtrOnNodef(expression.Implementing.Node, "No module found named %s", expression.Implementing.String)
			}
		},
		func(expression parser.LiteralExpression) {
			parser.LiteralExhaustiveSwitch(
				expression.Literal,
				func(literal float64) {
					varType = &standard_library.BasicTypeFloat
				},
				func(literal int) {
					varType = &standard_library.BasicTypeInt
				},
				func(literal string) {
					varType = &standard_library.BasicTypeString
				},
				func(literal bool) {
					varType = &standard_library.BasicTypeBoolean
				},
			)
		},
		func(expression parser.ReferenceOrInvocation) {
			var ok bool
			varType, ok = binding.GetTypeByVariableName(universe, expression.Var.String)
			if !ok {
				err = type_error.PtrOnNodef(expression.Var.Node, "Reference not found: %s", expression.Var.String)
				return
			}
			if expression.Arguments != nil {
				function, ok := varType.(*types.Function)
				if !ok {
					err = type_error.PtrOnNodef(expression.Var.Node, "Needs to be a function for invocation: %s", expression.Var.String)
					return
				}
				varType, err = typeOfInvocation(function, *expression.Arguments, universe)
			}
		},
		func(expression parser.Lambda) {
			localUniverse := universe

			generics := []string{}

			for _, generic := range expression.Generics {
				localUniverse, err = binding.CopyAddingType(localUniverse, generic, &types.TypeArgument{Name: generic.String})
				if err != nil {
					return
				}
				generics = append(generics, generic.String)
			}
			if len(generics) == 0 {
				generics = nil
			}

			arguments := []types.FunctionArgument{}
			for _, argument := range expression.Parameters {
				if argument.Type == nil {
					err = type_error.PtrOnNodef(argument.Name.Node, "Type annotation required for %s", argument.Name)
				}
				varType, err2 := validateTypeAnnotationInUniverse(*argument.Type, localUniverse)
				if err2 != nil {
					err = err2
					return
				}
				arguments = append(arguments, types.FunctionArgument{
					Name:         argument.Name.String,
					VariableType: varType,
				})
			}

			if expression.ReturnType == nil {
				err = type_error.PtrOnNodef(expression.Node, "Return yype annotation required")
				return
			}

			returnVarType, err2 := validateTypeAnnotationInUniverse(*expression.ReturnType, localUniverse)
			if err2 != nil {
				err = err2
				return
			}

			varType = &types.Function{
				Generics:   generics,
				Arguments:  arguments,
				ReturnType: returnVarType,
			}
		},
		func(expression parser.Declaration) {
			varType = &types.Void{}
		},
		func(expression parser.If) {
			if len(expression.ElseBlock) == 0 {
				varType = &types.Void{}
				return
			}
			typeOfThen, err2 := typeOfBlock(expression.ThenBlock, universe)
			if err2 != nil {
				err = err2
				return
			}
			typeOfElse, err2 := typeOfBlock(expression.ElseBlock, universe)
			if err2 != nil {
				err = err2
				return
			}
			varType = variableTypeCombine(typeOfThen, typeOfElse)
		},
		func(expression parser.Array) {
			if expression.Generic == nil {
				err = type_error.PtrOnNodef(expression.Node, "Missing generic")
				return
			}
			varType, err = validateTypeAnnotationInUniverse(*expression.Generic, universe)
			if err != nil {
				return
			}
			structVarType, ok := types.StructFieldVariableTypeFromVariableType(varType)
			if !ok {
				err = type_error.PtrOnNodef(expression.Node, "Not a valid generic: %s", printableName(varType))
			}
			varType = &types.Array{OfType: structVarType}
		},
		func(expression parser.When) {
			for _, whenIs := range expression.Is {
				t, err2 := typeOfBlock(whenIs.ThenBlock, universe)
				if err2 != nil {
					err = err2
					return
				}
				if varType == nil {
					varType = t
				} else {
					varType = variableTypeCombine(t, varType)
				}
			}
			if expression.Other != nil {
				t, err2 := typeOfBlock(expression.Other.ThenBlock, universe)
				if err2 != nil {
					err = err2
					return
				}
				if varType == nil {
					varType = t
				} else {
					varType = variableTypeCombine(t, varType)
				}
			}
			if varType == nil {
				err = type_error.PtrOnNodef(expression.Node, "Empty when")
			}
		},
	)
	return varType, err
}

func typeOfAccess(over types.VariableType, access parser.Name) (types.VariableType, *type_error.TypecheckError) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid, caseArray, caseOr := over.VariableTypeCases()
	if caseTypeArgument != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", printableName(over))
	} else if caseStruct != nil {
		varType, ok := caseStruct.Fields[access.String]
		if !ok {
			return nil, type_error.PtrOnNodef(access.Node, "no field named %s on struct %s", access.String, printableName(over))
		}
		return types.VariableTypeFromStructFieldVariableType(varType), nil
	} else if caseInterface != nil {
		varType, ok := caseInterface.Variables[access.String]
		if !ok {
			return nil, type_error.PtrOnNodef(access.Node, "no field named %s on interface %s", access.String, printableName(over))
		}
		return varType, nil
	} else if caseFunction != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", printableName(over))
	} else if caseBasicType != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", printableName(over))
	} else if caseVoid != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", printableName(over))
	} else if caseArray != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", printableName(over))
	} else if caseOr != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", printableName(over))
	} else {
		panic(fmt.Errorf("cases on %v", over))
	}
}

func typeOfInvocation(function *types.Function, argumentsList parser.ArgumentsList, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	resolvedGenericsFunction, _, err := resolveFunctionGenerics(argumentsList.Node, function, argumentsList.Generics, universe)
	if err != nil {
		return nil, err
	}

	return resolvedGenericsFunction.ReturnType, nil
}

func resolveGeneric(over types.VariableType, genericName string, resolveWith types.StructFieldVariableType) (types.VariableType, *type_error.TypecheckError) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid, caseArray, caseOr := over.VariableTypeCases()
	if caseTypeArgument != nil {
		if caseTypeArgument.Name == genericName {
			return types.VariableTypeFromStructFieldVariableType(resolveWith), nil
		}
		return caseTypeArgument, nil
	} else if caseStruct != nil {
		newStruct := &types.Struct{
			Package: caseStruct.Package,
			Name:    caseStruct.Name,
			Fields:  caseStruct.Fields,
		}
		for fieldName, variableType := range caseStruct.Fields {
			newFieldType, err := resolveGeneric(types.VariableTypeFromStructFieldVariableType(variableType), genericName, resolveWith)
			if err != nil {
				return nil, err
			}
			newStruct.Fields[fieldName] = newFieldType.(types.StructFieldVariableType)
		}
		return newStruct, nil
	} else if caseInterface != nil {
		panic("todo resolveGeneric caseInterface")
	} else if caseFunction != nil {
		panic("todo resolveGeneric caseFunction")
	} else if caseBasicType != nil {
		return caseBasicType, nil
	} else if caseVoid != nil {
		return caseVoid, nil
	} else if caseArray != nil {
		newOfType, err := resolveGeneric(types.VariableTypeFromStructFieldVariableType(caseArray.OfType), genericName, resolveWith)
		if err != nil {
			return nil, err
		}
		return &types.Array{
			OfType: newOfType.(types.StructFieldVariableType),
		}, nil
		panic("todo resolveGeneric caseArray")
	} else if caseOr != nil {
		panic("todo resolveGeneric caseOr")
	} else {
		panic(fmt.Errorf("cases on %v", over))
	}
}
