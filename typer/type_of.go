package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
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
		if accessOrInvocation.VarName != nil {
			varType, err = typeOfAccess(varType, *accessOrInvocation.VarName, universe)
			if err != nil {
				return nil, err
			}
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
	return types.Void(), nil
}

func typeOfExpression(expression parser.Expression, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	var varType types.VariableType
	var err *type_error.TypecheckError
	parser.ExpressionExhaustiveSwitch(
		expression,
		func(expression parser.Implementation) {
			generics := []types.VariableType{}
			for _, generic := range expression.Generics {
				varType, err2 := validateTypeAnnotationInUniverse(generic, universe)
				if err2 != nil {
					err = err2
					return
				}
				generics = append(generics, varType)
			}
			varType2, err2 := binding.GetTypeByTypeName(universe, expression.Implementing.String, generics)
			varType = varType2
			err = TypecheckErrorFromResolutionError(expression.Node, err2)
		},
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
					err = type_error.PtrOnNodef(argument.Name.Node, "Type annotation required for %s", argument.Name.String)
					return
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
				err = type_error.PtrOnNodef(expression.Node, "Return type annotation required")
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
			varType = types.Void()
		},
		func(expression parser.If) {
			if len(expression.ElseBlock) == 0 {
				varType = types.Void()
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
			varType = types.VariableTypeCombine(typeOfThen, typeOfElse)
			for _, elseIf := range expression.ElseIfs {
				typeOfElseIf, err2 := typeOfBlock(elseIf.ThenBlock, universe)
				if err2 != nil {
					err = err2
					return
				}
				varType = types.VariableTypeCombine(varType, typeOfElseIf)
			}
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
			array, ok := types.Array(varType)
			if !ok {
				err = type_error.PtrOnNodef(expression.Node, "Not a valid generic: %s", printableName(varType))
				return
			}
			varType = array
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
					varType = types.VariableTypeCombine(t, varType)
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

func typeOfAccess(over types.VariableType, access parser.Name, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := over.VariableTypeCases()
	if caseTypeArgument != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", printableName(over))
	} else if caseKnownType != nil {
		fields, resolutionErr := binding.GetFields(universe, caseKnownType)
		if resolutionErr != nil {
			return nil, TypecheckErrorFromResolutionError(access.Node, resolutionErr)
		}
		varType, ok := fields[access.String]
		if !ok {
			return nil, type_error.PtrOnNodef(access.Node, "no field named %s on %s", access.String, printableName(over))
		}
		return varType, nil
	} else if caseFunction != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", printableName(over))
	} else if caseOr != nil {
		return nil, type_error.PtrOnNodef(access.Node, "can't access over %s", printableName(over))
	} else {
		panic(fmt.Errorf("cases on %v", over))
	}
}

func typeOfInvocation(function *types.Function, argumentsList parser.ArgumentsList, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	resolvedGenericsFunction, _, _, err := resolveFunctionGenerics(argumentsList.Node, function, argumentsList.Generics, argumentsList.Arguments, universe)
	if err != nil {
		return nil, err
	}

	return resolvedGenericsFunction.ReturnType, nil
}
