package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func typeOfExpressionBox(expressionBox parser.ExpressionBox, file string, scope binding.Scope) (types.VariableType, *type_error.TypecheckError) {
	expression, accessOrInvocations := parser.ExpressionBoxFields(expressionBox)

	varType, err := typeOfExpression(expression, file, scope)
	if err != nil {
		return nil, err
	}
	if len(accessOrInvocations) == 0 {
		return varType, nil
	}

	for _, accessOrInvocation := range accessOrInvocations {
		if accessOrInvocation.DotOrArrowName != nil {
			varType, err = typeOfAccess(varType, accessOrInvocation.DotOrArrowName.VarName, scope)
			if err != nil {
				return nil, err
			}
		}
		if accessOrInvocation.Arguments != nil {
			function, ok := varType.(*types.Function)
			if !ok {
				return nil, type_error.PtrOnNodef(accessOrInvocation.Arguments.Node, "Expected a function in order to invoke")
			}
			varType, err = typeOfInvocation(function, *accessOrInvocation.Arguments, file, scope)
		}
	}

	return varType, nil
}

func typeOfBlock(block []parser.ExpressionBox, file string, scope binding.Scope) (types.VariableType, *type_error.TypecheckError) {
	for i, exp := range block {
		if i == len(block)-1 {
			return typeOfExpressionBox(exp, file, scope)
		}
		dec, ok := exp.Expression.(parser.Declaration)
		if !ok {
			continue
		}
		decType, err := typeOfExpressionBox(dec.ExpressionBox, file, scope)
		if err != nil {
			return nil, err
		}
		scope, err = binding.CopyAddingLocalVariable(scope, dec.Name, decType)
		if err != nil {
			return nil, err
		}
	}
	return types.Void(), nil
}

func typeOfExpression(expression parser.Expression, file string, scope binding.Scope) (types.VariableType, *type_error.TypecheckError) {
	var varType types.VariableType
	var err *type_error.TypecheckError
	parser.ExpressionExhaustiveSwitch(
		expression,
		func(expression parser.LiteralExpression) {
			parser.LiteralExhaustiveSwitch(
				expression.Literal,
				func(literal float64) {
					varType = types.Float()
				},
				func(literal int) {
					varType = types.Int()
				},
				func(literal string) {
					varType = types.String()
				},
				func(literal bool) {
					varType = types.Boolean()
				},
				func() {
					varType = types.Void()
				},
			)
		},
		func(expression parser.ReferenceOrInvocation) {
			var ok bool
			varType, ok = binding.GetTypeByVariableName(scope, expression.Var.String)
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
				varType, err = typeOfInvocation(function, *expression.Arguments, file, scope)
			}
		},
		func(expression parser.Lambda) {
			localScope := scope

			generics := []string{}

			for _, generic := range expression.Generics {
				localScope, err = binding.CopyAddingTypeToAllFiles(localScope, generic, &types.TypeArgument{Name: generic.String})
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
					err = type_error.PtrOnNodef(argument.Name.Node, "Type annotation required for %s", argument.Name.String)
					return
				}
				varType, err2 := validateTypeAnnotationInScope(*argument.Type, file, localScope)
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
				err = type_error.PtrOnNodef(expression.Node, "Return type annotation required")
				return
			}

			returnVarType, err2 := validateTypeAnnotationInScope(*expression.ReturnType, file, localScope)
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
			varType = types.Void()
		},
		func(expression parser.If) {
			if len(expression.ElseBlock) == 0 {
				varType = types.Void()
				return
			}
			typeOfThen, err2 := typeOfBlock(expression.ThenBlock, file, scope)
			if err2 != nil {
				err = err2
				return
			}
			typeOfElse, err2 := typeOfBlock(expression.ElseBlock, file, scope)
			if err2 != nil {
				err = err2
				return
			}
			varType = types.VariableTypeCombine(typeOfThen, typeOfElse)
			for _, elseIf := range expression.ElseIfs {
				typeOfElseIf, err2 := typeOfBlock(elseIf.ThenBlock, file, scope)
				if err2 != nil {
					err = err2
					return
				}
				varType = types.VariableTypeCombine(varType, typeOfElseIf)
			}
		},
		func(expression parser.List) {
			if expression.Generic == nil {
				if len(expression.Expressions) > 0 {
					varTypeOr := &types.OrVariableType{Elements: []types.VariableType{}}
					for _, expressionBox := range expression.Expressions {
						varType2, err2 := typeOfExpressionBox(expressionBox, file, scope)
						if err2 != nil {
							err = err2
							return
						}
						types.VariableTypeAddToOr(varType2, varTypeOr)
					}
					if len(varTypeOr.Elements) == 1 {
						varType = varTypeOr.Elements[0]
					} else {
						varType = varTypeOr
					}
					varType = types.List(varType)
					return
				} else {
					err = type_error.PtrOnNodef(expression.Node, "Missing generic")
					return
				}
			}
			varType, err = validateTypeAnnotationInScope(*expression.Generic, file, scope)
			if err != nil {
				return
			}
			varType = types.List(varType)
		},
		func(expression parser.When) {
			for _, whenIs := range expression.Is {
				t, err2 := typeOfBlock(whenIs.ThenBlock, file, scope)
				if err2 != nil {
					err = err2
					return
				}
				if varType == nil {
					varType = t
				} else {
					varType = types.VariableTypeCombine(t, varType)
				}
			}
			if expression.Other != nil {
				t, err2 := typeOfBlock(expression.Other.ThenBlock, file, scope)
				if err2 != nil {
					err = err2
					return
				}
				if varType == nil {
					varType = t
				} else {
					varType = types.VariableTypeCombine(t, varType)
				}
			}
			if varType == nil {
				err = type_error.PtrOnNodef(expression.Node, "Empty when")
			}
		},
	)
	return varType, err
}

func typeOfAccess(over types.VariableType, access parser.Name, scope binding.Scope) (types.VariableType, *type_error.TypecheckError) {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := over.VariableTypeCases()
	if caseTypeArgument != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", types.PrintableName(over))
	} else if caseKnownType != nil {
		fields, resolutionErr := binding.GetFields(scope, caseKnownType)
		if resolutionErr != nil {
			return nil, TypecheckErrorFromResolutionError(access.Node, resolutionErr)
		}
		varType, ok := fields[access.String]
		if !ok {
			return nil, type_error.PtrOnNodef(access.Node, "no field named %s on %s", access.String, types.PrintableName(over))
		}
		return varType, nil
	} else if caseFunction != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", types.PrintableName(over))
	} else if caseOr != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", types.PrintableName(over))
	} else {
		panic(fmt.Errorf("cases on %v", over))
	}
}

func typeOfInvocation(function *types.Function, argumentsList parser.ArgumentsList, file string, scope binding.Scope) (types.VariableType, *type_error.TypecheckError) {
	if len(function.Generics) == 0 {
		return function.ReturnType, nil
	}
	resolvedReturnType, err := typeOfReturnedByFunctionAfterResolvingGenerics(argumentsList.Node, function, argumentsList.Generics, argumentsList.Arguments, file, scope)
	if err != nil {
		return nil, err
	}

	return resolvedReturnType, nil
}

func typeOfReturnedByFunctionAfterResolvingGenerics(node parser.Node, function *types.Function, genericsPassed []parser.TypeAnnotation, argumentsPassed []parser.NamedArgument, file string, scope binding.Scope) (types.VariableType, *type_error.TypecheckError) {
	if len(genericsPassed) > 0 && len(function.Generics) != len(genericsPassed) {
		return nil, type_error.PtrOnNodef(node, "wrong number of generics, expected %d but got %d", len(function.Generics), len(genericsPassed))
	}
	resolve := map[string]types.VariableType{}
	inferredGenerics, err := attemptGenericInference(node, function, argumentsPassed, genericsPassed, nil, file, scope)
	if err != nil {
		return nil, err
	}
	for i, genericName := range function.Generics {
		resolve[genericName] = inferredGenerics[i]
	}

	return typeOfResolvingGeneric(node, function.ReturnType, resolve)
}

func typeOfResolvingGeneric(node parser.Node, varType types.VariableType, resolve map[string]types.VariableType) (types.VariableType, *type_error.TypecheckError) {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		resolved, ok := resolve[caseTypeArgument.Name]
		if ok {
			return resolved, nil
		} else {
			return nil, type_error.PtrOnNodef(node, "failed to determine generics (a type annotation might be required)")
		}
	} else if caseKnownType != nil {
		resultGenerics := []types.VariableType{}
		for _, generic := range caseKnownType.Generics {
			resolved, err := typeOfResolvingGeneric(node, generic, resolve)
			if err != nil {
				return nil, err
			}
			resultGenerics = append(resultGenerics, resolved)
		}
		return &types.KnownType{
			Package:          caseKnownType.Package,
			Name:             caseKnownType.Name,
			DeclaredGenerics: caseKnownType.DeclaredGenerics,
			Generics:         resultGenerics,
		}, nil
	} else if caseFunction != nil {
		arguments := []types.FunctionArgument{}
		for _, argument := range caseFunction.Arguments {
			varType, err := typeOfResolvingGeneric(node, argument.VariableType, resolve)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, types.FunctionArgument{
				Name:         argument.Name,
				VariableType: varType,
			})
		}
		returnVarType, err := typeOfResolvingGeneric(node, caseFunction.ReturnType, resolve)
		if err != nil {
			return nil, err
		}
		return &types.Function{
			Generics:   caseFunction.Generics,
			Arguments:  arguments,
			ReturnType: returnVarType,
		}, nil
	} else if caseOr != nil {
		resultElements := []types.VariableType{}
		for _, element := range caseOr.Elements {
			resolved, err := typeOfResolvingGeneric(node, element, resolve)
			if err != nil {
				return nil, err
			}
			resultElements = append(resultElements, resolved)
		}
		return &types.OrVariableType{Elements: resultElements}, nil
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}
