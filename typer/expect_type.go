package typer

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func expectTypeOfExpressionBox(expectedType types.VariableType, expressionBox parser.ExpressionBox, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	if len(expressionBox.AccessOrInvocationChain) == 0 {
		return expectTypeOfExpression(expectedType, expressionBox.Expression, file, scope)
	}

	varType, err := typeOfExpression(expressionBox.Expression, file, scope)
	if err != nil {
		return nil, err
	}

	astExp, err := expectTypeOfExpression(varType, expressionBox.Expression, file, scope)
	if err != nil {
		return nil, err
	}

	for i, accessOrInvocation := range expressionBox.AccessOrInvocationChain {
		if i < len(expressionBox.AccessOrInvocationChain)-1 {
			astExp, err = determineTypeOfAccessOrInvocation(astExp, accessOrInvocation, nil, file, scope)
			if err != nil {
				return nil, err
			}
		} else {
			astExp, err = determineTypeOfAccessOrInvocation(astExp, accessOrInvocation, &expectedType, file, scope)
			if err != nil {
				return nil, err
			}
			gotVarType := ast.VariableTypeOfExpression(astExp)
			if !types.VariableTypeContainedIn(gotVarType, expectedType) {
				return nil, type_error.PtrOnNodef(accessOrInvocation.Node, "Expected %s but got %s", types.PrintableName(expectedType), types.PrintableName(gotVarType))
			}
		}
	}

	return astExp, nil
}

func determineTypeOfAccessOrInvocation(over ast.Expression, accessOrInvocation parser.AccessOrInvocation, expectedReturnType *types.VariableType, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	lhsVarType := ast.VariableTypeOfExpression(over)
	astExp := over
	var err *type_error.TypecheckError
	if accessOrInvocation.DotOrArrowName != nil {
		lhsVarType, err = typeOfAccess(lhsVarType, accessOrInvocation.DotOrArrowName.VarName, scope)
		if err != nil {
			return nil, err
		}

		astExp = ast.Access{
			VariableType: lhsVarType,
			Over:         over,
			Access:       accessOrInvocation.DotOrArrowName.VarName.String,
		}
	}
	if accessOrInvocation.Arguments != nil {
		function, ok := lhsVarType.(*types.Function)
		if !ok {
			return nil, type_error.PtrOnNodef(accessOrInvocation.Arguments.Node, "Should be a function in order to be invoked but is %s", types.PrintableName(lhsVarType))
		}
		if len(function.Arguments) != len(accessOrInvocation.Arguments.Arguments) {
			return nil, type_error.PtrOnNodef(accessOrInvocation.Arguments.Node, "Invoked with wrong number of arguments, expected %d but got %d", len(function.Arguments), len(accessOrInvocation.Arguments.Arguments))
		}

		resolvedGenericsFunction, generics, arguments, err := resolveFunctionGenerics(
			accessOrInvocation.Arguments.Node,
			function,
			accessOrInvocation.Arguments.Generics,
			accessOrInvocation.Arguments.Arguments,
			expectedReturnType,
			file,
			scope,
		)
		if err != nil {
			return nil, err
		}

		astExp = ast.Invocation{
			VariableType: resolvedGenericsFunction.ReturnType,
			Over:         astExp,
			Generics:     generics,
			Arguments:    arguments,
		}
	}

	return astExp, nil
}

func expectTypeOfExpression(expectedType types.VariableType, expression parser.Expression, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	var astExp ast.Expression
	var err *type_error.TypecheckError
	parser.ExpressionExhaustiveSwitch(
		expression,
		func(expression parser.LiteralExpression) {
			astExp, err = expectTypeOfLiteral(expectedType, expression, file, scope)
		},
		func(expression parser.ReferenceOrInvocation) {
			astExp, err = expectTypeOfReferenceOrInvocation(expectedType, expression, file, scope)
		},
		func(expression parser.Lambda) {
			astExp, err = expectTypeOfLambda(expectedType, expression, file, scope)
		},
		func(expression parser.Declaration) {
			astExp, err = expectTypeOfDeclaration(expectedType, expression, file, scope)
		},
		func(expression parser.If) {
			astExp, err = expectTypeOfIf(expectedType, expression, file, scope)
		},
		func(expression parser.List) {
			astExp, err = expectTypeOfList(expectedType, expression, file, scope)
		},
		func(expression parser.When) {
			astExp, err = expectTypeOfWhen(expectedType, expression, file, scope)
		},
	)
	return astExp, err
}

func expectTypeOfWhen(expectedType types.VariableType, expression parser.When, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	typeOfOver, err := typeOfExpressionBox(expression.Over, file, scope)
	if err != nil {
		return nil, err
	}
	_, _, _, typeOverOr := typeOfOver.VariableTypeCases()
	if typeOverOr == nil {
		typeOverOr = &types.OrVariableType{
			Elements: []types.VariableType{typeOfOver},
		}
	}

	missingCases := map[string]types.VariableType{}
	for _, varType := range typeOverOr.Elements {
		missingCases[types.PrintableName(varType)] = varType
	}

	astOver, err := expectTypeOfExpressionBox(typeOfOver, expression.Over, file, scope)
	if err != nil {
		return nil, err
	}

	cases := map[types.VariableType][]ast.Expression{}
	caseNames := map[types.VariableType]*string{}

	for _, whenIs := range expression.Is {
		varType, err := scopecheck.ValidateTypeAnnotationInScope(whenIs.Type, file, scope)
		if err != nil {
			return nil, TypecheckErrorFromScopeCheckError(err)
		}
		if missingCases[types.PrintableName(varType)] != nil {
			delete(missingCases, types.PrintableName(varType))
			localScope := scope
			if whenIs.Name != nil {
				var err *binding.ResolutionError
				localScope, err = binding.CopyAddingLocalVariable(localScope, *whenIs.Name, varType)
				if err != nil {
					return nil, TypecheckErrorFromResolutionError(whenIs.Name.Node, err)
				}
			}
			astThen, err := expectTypeOfBlock(expectedType, whenIs.Node, whenIs.ThenBlock, file, localScope)
			if err != nil {
				return nil, err
			}
			cases[varType] = astThen
			if whenIs.Name != nil {
				caseNames[varType] = &whenIs.Name.String
			}
		} else {
			return nil, type_error.PtrOnNodef(whenIs.Node, "no matching for %s in %s", types.PrintableName(varType), types.PrintableName(typeOfOver))
		}
	}

	var otherCase []ast.Expression = nil
	var otherCaseName *string = nil
	if expression.Other != nil {
		orCases := []types.VariableType{}
		for _, variableType := range missingCases {
			orCases = append(orCases, variableType)
		}
		missingCases = nil
		varType := &types.OrVariableType{Elements: orCases}
		localScope := scope
		if expression.Other.Name != nil {
			var err *binding.ResolutionError
			localScope, err = binding.CopyAddingLocalVariable(scope, *expression.Other.Name, varType)
			if err != nil {
				return nil, TypecheckErrorFromResolutionError(expression.Other.Name.Node, err)
			}
		}
		astThen, err := expectTypeOfBlock(expectedType, expression.Other.Node, expression.Other.ThenBlock, file, localScope)
		if err != nil {
			return nil, err
		}
		otherCase = astThen
		if expression.Other.Name != nil {
			otherCaseName = &expression.Other.Name.String
		}
	}

	if len(missingCases) > 0 {
		varTypeNames := ""
		for _, varType := range missingCases {
			if varTypeNames != "" {
				varTypeNames += ", "
			}
			varTypeNames += types.PrintableName(varType)
		}
		return nil, type_error.PtrOnNodef(expression.Node, "missing cases for %s", varTypeNames)
	}

	return ast.When{
		VariableType:  expectedType,
		Over:          astOver,
		Cases:         cases,
		CaseNames:     caseNames,
		OtherCase:     otherCase,
		OtherCaseName: otherCaseName,
	}, nil
}

func expectTypeOfList(expectedType types.VariableType, expression parser.List, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	var expectedListOf types.VariableType

	if expression.Generic != nil {
		varType, err := scopecheck.ValidateTypeAnnotationInScope(*expression.Generic, file, scope)
		if err != nil {
			return nil, TypecheckErrorFromScopeCheckError(err)
		}
		expectedListOf = varType
	} else if len(expression.Expressions) == 0 {
		_, caseKnownType, _, _ := expectedType.VariableTypeCases()
		if caseKnownType != nil && caseKnownType.Package == "" && caseKnownType.Name == "List" {
			return ast.List{
				ContainedVariableType: caseKnownType.Generics[0],
				Arguments:             []ast.Expression{},
			}, nil
		} else {
			return nil, type_error.PtrOnNodef(expression.Node, "Could not infer list generic, please annotate it")
		}
	} else {
		or := &types.OrVariableType{
			Elements: []types.VariableType{},
		}
		for _, expressionBox := range expression.Expressions {
			varType, err := typeOfExpressionBox(expressionBox, file, scope)
			if err != nil {
				return nil, err
			}
			types.VariableTypeAddToOr(varType, or)
		}
		if len(or.Elements) == 0 {
			panic("TODO expectTypeOfList invalid")
		} else if len(or.Elements) == 1 {
			expectedListOf = or.Elements[0]
		} else {
			expectedListOf = or
		}
	}

	expectedList := types.List(expectedListOf)
	if !types.VariableTypeContainedIn(expectedList, expectedType) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected %s but got %s", types.PrintableName(expectedType), types.PrintableName(expectedList))
	}

	astArguments := []ast.Expression{}
	for _, expressionBox := range expression.Expressions {
		astExp, err := expectTypeOfExpressionBox(expectedListOf, expressionBox, file, scope)
		if err != nil {
			return nil, err
		}
		astArguments = append(astArguments, astExp)
	}

	return ast.List{
		ContainedVariableType: expectedListOf,
		Arguments:             astArguments,
	}, nil
}

func expectTypeOfIf(expectedType types.VariableType, expression parser.If, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	astCondition, err := expectTypeOfExpressionBox(types.Boolean(), expression.Condition, file, scope)
	if err != nil {
		return nil, err
	}

	thenBlock, err := expectTypeOfBlock(expectedType, expression.Node, expression.ThenBlock, file, scope)
	if err != nil {
		return nil, err
	}
	for len(expression.ElseIfs) > 0 {
		elem := expression.ElseIfs[len(expression.ElseIfs)-1]
		expression.ElseIfs = expression.ElseIfs[:len(expression.ElseIfs)-1]
		expression.ElseBlock = []parser.ExpressionBox{
			parser.ExpressionBox{
				Expression: parser.If{
					Node:      elem.Node,
					Condition: elem.Condition,
					ThenBlock: elem.ThenBlock,
					ElseBlock: expression.ElseBlock,
				},
			},
		}
	}
	var elseBlock []ast.Expression = nil
	if len(expression.ElseBlock) != 0 {
		block, err := expectTypeOfBlock(expectedType, expression.Node, expression.ElseBlock, file, scope)
		if err != nil {
			return nil, err
		}
		elseBlock = block
	}

	return ast.If{
		VariableType: expectedType,
		Condition:    astCondition,
		ThenBlock:    thenBlock,
		ElseBlock:    elseBlock,
	}, nil
}

func expectTypeOfDeclaration(expectedDeclarationType types.VariableType, expression parser.Declaration, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	if expression.ShortCircuit != nil {
		panic("failed to desugar before expectTypeOfDeclaration")
	}
	if !types.VariableTypeEq(expectedDeclarationType, types.Void()) {
		return nil, type_error.PtrOnNodef(expression.Name.Node, "Expected type %s but got void", types.PrintableName(expectedDeclarationType))
	}

	var expectedType types.VariableType
	var err *type_error.TypecheckError
	if expression.TypeAnnotation != nil {
		var err2 scopecheck.ScopeCheckError
		expectedType, err2 = scopecheck.ValidateTypeAnnotationInScope(*expression.TypeAnnotation, file, scope)
		err = TypecheckErrorFromScopeCheckError(err2)
	} else {
		expectedType, err = typeOfExpressionBox(expression.ExpressionBox, file, scope)
	}
	if err != nil {
		return nil, err
	}
	astExp, err := expectTypeOfExpressionBox(expectedType, expression.ExpressionBox, file, scope)
	if err != nil {
		return nil, err
	}
	return ast.Declaration{
		Name:       expression.Name.String,
		Expression: astExp,
	}, nil
}

func expectTypeOfLambda(expectedType types.VariableType, expression parser.Lambda, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	_, _, expectedFunction, expectedOr := expectedType.VariableTypeCases()
	if expectedOr != nil {
		for _, element := range expectedOr.Elements {
			f, ok := element.(*types.Function)
			if ok {
				if expectedFunction != nil {
					panic("TODO expectTypeOfLambda or with multiple functions")
				}
				expectedFunction = f
			}
		}
	}
	if expectedFunction == nil {
		return nil, type_error.PtrOnNodef(expression.Node, "Expected %s but got a function", types.PrintableName(expectedType))
	}

	if len(expression.Signature.Generics) != len(expectedFunction.Generics) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected %d generics but got %d", len(expectedFunction.Generics), len(expression.Signature.Generics))
	}

	localScope := scope
	for _, generic := range expression.Signature.Generics {
		var err *binding.ResolutionError
		localScope, err = binding.CopyAddingTypeToAllFiles(localScope, generic, &types.TypeArgument{Name: generic.String})
		if err != nil {
			return nil, TypecheckErrorFromResolutionError(generic.Node, err)
		}
	}

	if len(expression.Signature.Parameters) != len(expectedFunction.Arguments) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected %d params but got %d", len(expectedFunction.Arguments), len(expression.Signature.Parameters))
	}
	for i, parameter := range expression.Signature.Parameters {
		if parameter.Type != nil {
			var err scopecheck.ScopeCheckError
			paramType, err := scopecheck.ValidateTypeAnnotationInScope(*parameter.Type, file, localScope)
			if err != nil {
				return nil, TypecheckErrorFromScopeCheckError(err)
			}
			if !types.VariableTypeContainedIn(expectedFunction.Arguments[i].VariableType, paramType) {
				return nil, type_error.PtrOnNodef(expression.Node, "in parameter position %d expected type %s but you have annotated %s", i, types.PrintableName(expectedFunction.Arguments[i].VariableType), types.PrintableName(paramType))
			}
		}
		var err *binding.ResolutionError
		localScope, err = binding.CopyAddingLocalVariable(localScope, parameter.Name, expectedFunction.Arguments[i].VariableType)
		if err != nil {
			return nil, TypecheckErrorFromResolutionError(parameter.Name.Node, err)
		}
	}

	expectedTypeOfBlock := expectedFunction.ReturnType
	if expression.Signature.ReturnType != nil {
		returnType, err := scopecheck.ValidateTypeAnnotationInScope(*expression.Signature.ReturnType, file, localScope)
		if err != nil {
			return nil, TypecheckErrorFromScopeCheckError(err)
		}
		if !types.VariableTypeContainedIn(returnType, expectedFunction.ReturnType) {
			return nil, type_error.PtrOnNodef(expression.Node, "in return type expected type %s but you have annotated %s", types.PrintableName(expectedFunction.ReturnType), types.PrintableName(returnType))
		}
		expectedTypeOfBlock = returnType
	}

	astBlock, err := expectTypeOfBlock(expectedTypeOfBlock, expression.Node, expression.Block, file, localScope)
	if err != nil {
		return nil, err
	}

	varType := &types.Function{
		Generics:   expectedFunction.Generics,
		Arguments:  []types.FunctionArgument{},
		ReturnType: expectedFunction.ReturnType,
	}
	for i, arg := range expectedFunction.Arguments {
		varType.Arguments = append(varType.Arguments, types.FunctionArgument{
			Name:         expression.Signature.Parameters[i].Name.String,
			VariableType: arg.VariableType,
		})
	}

	return &ast.Function{
		VariableType: varType,
		Block:        astBlock,
	}, nil
}

func expectTypeOfBlock(expectedType types.VariableType, node parser.Node, block []parser.ExpressionBox, file string, scope binding.Scope) ([]ast.Expression, *type_error.TypecheckError) {
	result := []ast.Expression{}

	if len(block) == 0 {
		return nil, type_error.PtrOnNodef(node, "empty function block not allowed (maybe you want to return null?)")
	}

	localScope := scope
	for i, expressionBox := range block {
		var expectedTypeOfExpressionBox = expectedType
		var err *type_error.TypecheckError
		if i < len(block)-1 {
			expectedTypeOfExpressionBox, err = typeOfExpressionBox(expressionBox, file, localScope)
			if err != nil {
				return nil, err
			}
		}
		astExp, err := expectTypeOfExpressionBox(expectedTypeOfExpressionBox, expressionBox, file, localScope)
		if err != nil {
			return nil, err
		}
		result = append(result, astExp)
		astDec, isDec := astExp.(ast.Declaration)
		if isDec {
			var err *binding.ResolutionError
			localScope, err = binding.CopyAddingLocalVariable(localScope, parser.Name{
				String: astDec.Name,
			}, ast.VariableTypeOfExpression(astDec.Expression))
			if err != nil {
				// TODO FIXME shouldn't convert with an empty Node
				return nil, TypecheckErrorFromResolutionError(parser.Node{}, err)
			}
		}
	}

	return result, nil
}

func resolveFunctionGenerics(node parser.Node, function *types.Function, genericsPassed []parser.TypeAnnotation, argumentsPassed []parser.NamedArgument, expectedReturnType *types.VariableType, file string, scope binding.Scope) (*types.Function, []types.VariableType, []ast.Expression, *type_error.TypecheckError) {
	generics := []types.VariableType{}

	genericsPassedContainsUnderscore := false
	var err *type_error.TypecheckError
	for _, passed := range genericsPassed {
		for _, element := range passed.OrTypes {
			parser.TypeAnnotationElementExhaustiveSwitch(
				element,
				func(underscoreTypeAnnotation parser.SingleNameType) {
					if len(passed.OrTypes) > 1 {
						err = type_error.PtrOnNodef(underscoreTypeAnnotation.Node, "Cannot infer part of an or type")
						return
					}
					genericsPassedContainsUnderscore = true
				},
				func(typeAnnotation parser.SingleNameType) {},
				func(typeAnnotation parser.FunctionType) {},
			)
		}
	}
	if err != nil {
		return nil, nil, nil, err
	}

	if genericsPassedContainsUnderscore || (len(genericsPassed) == 0 && len(function.Generics) > 0) {
		inferredGenerics, err := attemptGenericInference(node, function, argumentsPassed, genericsPassed, expectedReturnType, file, scope)
		if err != nil {
			return nil, nil, nil, err
		}
		generics = inferredGenerics
	} else {
		if len(genericsPassed) != len(function.Generics) {
			return nil, nil, nil, type_error.PtrOnNodef(node, "expected %d generics but got %d", len(function.Generics), len(genericsPassed))
		}
		for _, generic := range genericsPassed {
			varType, err := scopecheck.ValidateTypeAnnotationInScope(generic, file, scope)
			if err != nil {
				return nil, nil, nil, TypecheckErrorFromScopeCheckError(err)
			}
			generics = append(generics, varType)
		}
	}
	genericsMap := map[string]types.VariableType{}
	for i, varType := range generics {
		genericsMap[function.Generics[i]] = varType
	}

	if len(argumentsPassed) != len(function.Arguments) {
		return nil, nil, nil, type_error.PtrOnNodef(node, "expected %d arguments but got %d", len(function.Arguments), len(argumentsPassed))
	}
	arguments := []types.FunctionArgument{}
	for _, argument := range function.Arguments {
		arguments = append(arguments, argument)
	}
	for i := 0; i < len(arguments); i++ {
		for genericName, resolveTo := range genericsMap {
			newVarType, err := binding.ResolveGeneric(arguments[i].VariableType, genericName, resolveTo)
			if err != nil {
				return nil, nil, nil, TypecheckErrorFromResolutionError(node, err)
			}
			arguments[i].VariableType = newVarType
		}
	}
	astArguments := []ast.Expression{}
	for i, argument := range argumentsPassed {
		if argument.Name != nil && argument.Name.String != arguments[i].Name {
			return nil, nil, nil, type_error.PtrOnNodef(argument.Name.Node, "name of argument should be '%s'", arguments[i].Name)
		}
		expectedArgType := arguments[i].VariableType
		astArg, err := expectTypeOfExpressionBox(expectedArgType, argument.Argument, file, scope)
		if err != nil {
			return nil, nil, nil, err
		}
		astArguments = append(astArguments, astArg)
	}

	returnType := function.ReturnType
	for genericName, resolveTo := range genericsMap {
		newVarType, err := binding.ResolveGeneric(returnType, genericName, resolveTo)
		if err != nil {
			return nil, nil, nil, TypecheckErrorFromResolutionError(node, err)
		}
		returnType = newVarType
	}

	return &types.Function{
		Generics:   nil,
		Arguments:  arguments,
		ReturnType: returnType,
	}, generics, astArguments, nil
}

func attemptGenericInference(node parser.Node, function *types.Function, argumentsPassed []parser.NamedArgument, genericsPassed []parser.TypeAnnotation, expectedReturnType *types.VariableType, file string, scope binding.Scope) ([]types.VariableType, *type_error.TypecheckError) {
	resolvedGenerics := []types.VariableType{}
	for genericIndex, functionGenericName := range function.Generics {
		if len(genericsPassed) > 0 {
			shouldInfer := false
			passed := genericsPassed[genericIndex]
			for _, element := range passed.OrTypes {
				var err *type_error.TypecheckError
				parser.TypeAnnotationElementExhaustiveSwitch(
					element,
					func(underscoreTypeAnnotation parser.SingleNameType) {
						if len(passed.OrTypes) > 1 {
							err = type_error.PtrOnNodef(underscoreTypeAnnotation.Node, "Cannot infer part of an or type")
							return
						}
						shouldInfer = true
					},
					func(typeAnnotation parser.SingleNameType) {},
					func(typeAnnotation parser.FunctionType) {},
				)
				if err != nil {
					return nil, err
				}
			}
			if !shouldInfer {
				varType, err := scopecheck.ValidateTypeAnnotationInScope(passed, file, scope)
				if err != nil {
					return nil, TypecheckErrorFromScopeCheckError(err)
				}
				resolvedGenerics = append(resolvedGenerics, varType)
				continue
			}
		}

		var found types.VariableType
		for i, arg := range argumentsPassed {
			var typeOfArgFunction types.VariableType
			_, _, caseParameterFunction, _ := function.Arguments[i].VariableType.VariableTypeCases()
			if caseParameterFunction != nil {
				if len(arg.Argument.AccessOrInvocationChain) == 0 {
					lambda, ok := arg.Argument.Expression.(parser.Lambda)
					if ok {
						if len(lambda.Signature.Generics) == 0 {
							argumentTypes, ok, err := tryToDetermineFunctionArgumentTypes(resolvedGenerics, lambda, function, caseParameterFunction, file, scope)
							if err != nil {
								return nil, err
							}
							if !ok {
								continue
							}
							localScope := scope
							for i, argType := range argumentTypes {
								var err *binding.ResolutionError
								localScope, err = binding.CopyAddingLocalVariable(localScope, lambda.Signature.Parameters[i].Name, argType)
								if err != nil {
									return nil, TypecheckErrorFromResolutionError(lambda.Signature.Parameters[i].Name.Node, err)
								}
							}
							var returnType types.VariableType
							if lambda.Signature.ReturnType != nil {
								rType, err := scopecheck.ValidateTypeAnnotationInScope(*lambda.Signature.ReturnType, file, scope)
								if err != nil {
									return nil, TypecheckErrorFromScopeCheckError(err)
								}
								returnType = rType
							} else {
								rType, err := typeOfBlock(lambda.Block, file, localScope)
								if err != nil {
									return nil, err
								}
								returnType = rType
							}
							arguments := []types.FunctionArgument{}
							for i, variableType := range argumentTypes {
								arguments = append(arguments, types.FunctionArgument{
									Name:         lambda.Signature.Parameters[i].Name.String,
									VariableType: variableType,
								})
							}
							typeOfArgFunction = &types.Function{
								Generics:   nil,
								Arguments:  arguments,
								ReturnType: returnType,
							}
						}
					}
				}
			}
			typeOfArg := typeOfArgFunction
			if typeOfArg == nil {
				typeOfThisArg, err := typeOfExpressionBox(arg.Argument, file, scope)
				if err != nil {
					continue
				}
				typeOfArg = typeOfThisArg
			}
			maybeInferred, ok := tryToInferGeneric(functionGenericName, function.Arguments[i].VariableType, typeOfArg)
			if !ok {
				return nil, type_error.PtrOnNodef(node, "Could not infer generics, please annotate them")
			}
			if maybeInferred != nil {
				if found == nil || types.VariableTypeContainedIn(found, maybeInferred) {
					found = maybeInferred
				} else {
					return nil, type_error.PtrOnNodef(node, "Could not infer generics, please annotate them")
				}
			}
		}
		if found == nil && expectedReturnType != nil {
			caseTypeArgument, _, _, _ := function.ReturnType.VariableTypeCases()
			if caseTypeArgument != nil && caseTypeArgument.Name == functionGenericName {
				found = *expectedReturnType
			}
		}
		if found == nil {
			return nil, type_error.PtrOnNodef(node, "Could not infer generics, please annotate them")
		}
		resolvedGenerics = append(resolvedGenerics, found)
	}
	if len(resolvedGenerics) == len(function.Generics) {
		return resolvedGenerics, nil
	} else {
		return nil, type_error.PtrOnNodef(node, "Could not infer generics, please annotate them")
	}
}

func tryToDetermineFunctionArgumentTypes(
	resolvedGenerics []types.VariableType,
	lambda parser.Lambda,
	function *types.Function,
	caseParameterFunction *types.Function,
	file string,
	scope binding.Scope,
) ([]types.VariableType, bool, *type_error.TypecheckError) {
	if len(lambda.Signature.Generics) > 0 {
		return nil, false, nil
	}
	arguments := []types.VariableType{}
	successInArguments := true
	for i, parameter := range lambda.Signature.Parameters {
		if parameter.Type == nil {
			typeOfParam, ok := tryToDetermineFunctionArgumentType(resolvedGenerics, function.Generics, caseParameterFunction.Arguments[i].VariableType)
			if !ok {
				return nil, false, nil
			}
			arguments = append(arguments, typeOfParam)
		} else {
			typeOfParam, err := scopecheck.ValidateTypeAnnotationInScope(*parameter.Type, file, scope)
			if err != nil {
				return nil, false, TypecheckErrorFromScopeCheckError(err)
			}
			arguments = append(arguments, typeOfParam)
		}
	}
	if !successInArguments {
		return nil, false, nil
	}
	return arguments, true, nil
}

func tryToDetermineFunctionArgumentType(
	resolvedGenerics []types.VariableType,
	functionGenerics []string,
	argumentVariableType types.VariableType,
) (types.VariableType, bool) {
	caseTypeArg, caseKnownType, _, _ := argumentVariableType.VariableTypeCases()
	if caseTypeArg != nil {
		for i, generic := range functionGenerics {
			if generic == caseTypeArg.Name {
				if len(resolvedGenerics) > i {
					return resolvedGenerics[i], true
				}
			}
		}
		return nil, false
	} else if caseKnownType != nil {
		return caseKnownType, true
	} else {
		return nil, false
	}
}

func tryToInferGeneric(genericName string, functionVarType types.VariableType, argVarType types.VariableType) (types.VariableType, bool) {
	funcCaseTypeArgument, funcCaseKnownType, funcCaseFunction, funcCaseOr := functionVarType.VariableTypeCases()
	if funcCaseTypeArgument != nil {
		if funcCaseTypeArgument.Name == genericName {
			return argVarType, true
		}
		return nil, true
	} else if funcCaseKnownType != nil {
		argKnownType, ok := argVarType.(*types.KnownType)
		if ok && len(funcCaseKnownType.Generics) == len(argKnownType.Generics) {
			for i, _ := range funcCaseKnownType.Generics {
				inferred, ok := tryToInferGeneric(genericName, funcCaseKnownType.Generics[i], argKnownType.Generics[i])
				if inferred != nil || !ok {
					return inferred, ok
				}
			}
		}
		return nil, true
	} else if funcCaseFunction != nil {
		for _, generic := range funcCaseFunction.Generics {
			if generic == genericName {
				return nil, true
			}
		}
		argFunction, ok := argVarType.(*types.Function)
		if !ok {
			return nil, false
		}
		if len(funcCaseFunction.Arguments) != len(argFunction.Arguments) {
			return nil, false
		}
		var found types.VariableType
		for i, _ := range funcCaseFunction.Arguments {
			maybeInferred, ok := tryToInferGeneric(genericName, funcCaseFunction.Arguments[i].VariableType, argFunction.Arguments[i].VariableType)
			if !ok {
				return nil, false
			}
			if maybeInferred != nil {
				if found == nil || types.VariableTypeContainedIn(found, maybeInferred) {
					found = maybeInferred
				} else {
					return nil, false
				}
			}
		}
		maybeInferred, ok := tryToInferGeneric(genericName, funcCaseFunction.ReturnType, argFunction.ReturnType)
		if !ok {
			return nil, false
		}
		if maybeInferred != nil {
			if found == nil || types.VariableTypeContainedIn(found, maybeInferred) {
				found = maybeInferred
			} else {
				return nil, false
			}
		}
		return found, true
	} else if funcCaseOr != nil {
		_, _, _, caseArgOr := argVarType.VariableTypeCases()
		if caseArgOr != nil {
			remainingTypesToMatch := []types.VariableType{}
			for _, argVarType := range caseArgOr.Elements {
				matched := false
				for _, element := range funcCaseOr.Elements {
					if types.VariableTypeEq(argVarType, element) {
						matched = true
						break
					}
				}
				if !matched {
					remainingTypesToMatch = append(remainingTypesToMatch, argVarType)
				}
			}
			argVarType = &types.OrVariableType{Elements: remainingTypesToMatch}
		}

		var found types.VariableType
		for _, element := range funcCaseOr.Elements {
			maybeInferred, ok := tryToInferGeneric(genericName, element, argVarType)
			if !ok {
				return nil, false
			}
			if maybeInferred != nil {
				if found == nil || types.VariableTypeContainedIn(found, maybeInferred) {
					found = maybeInferred
				} else {
					return nil, false
				}
			}
		}
		return found, true
	} else {
		return nil, true
	}
}

func expectTypeOfReferenceOrInvocation(expectedType types.VariableType, expression parser.ReferenceOrInvocation, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	overType, ok := binding.GetTypeByVariableName(scope, expression.Var.String)
	if !ok {
		return nil, type_error.PtrOnNodef(expression.Var.Node, "Not found in scope: %s", expression.Var.String)
	}

	if expression.Arguments != nil {
		arguments := []ast.Expression{}
		overFunction, ok := overType.(*types.Function)
		if !ok {
			return nil, type_error.PtrOnNodef(expression.Arguments.Node, "Can't invoke on not a function")
		}

		overFunction, generics, arguments, err := resolveFunctionGenerics(
			expression.Arguments.Node,
			overFunction,
			expression.Arguments.Generics,
			expression.Arguments.Arguments,
			&expectedType,
			file,
			scope,
		)
		if err != nil {
			return nil, err
		}

		if !types.VariableTypeContainedIn(overFunction.ReturnType, expectedType) {
			return nil, type_error.PtrOnNodef(expression.Var.Node, "expected type %s but found %s", types.PrintableName(expectedType), types.PrintableName(overFunction.ReturnType))
		}

		pkg, name := binding.GetPackageLevelAndUnaliasedNameOfVariable(scope, expression.Var)
		astExp := ast.Invocation{
			VariableType: overFunction.ReturnType,
			Over: ast.Reference{
				VariableType: overFunction,
				PackageName:  pkg,
				Name:         name,
			},
			Generics:  generics,
			Arguments: arguments,
		}

		return astExp, nil
	} else {
		if !types.VariableTypeContainedIn(overType, expectedType) {
			return nil, type_error.PtrOnNodef(expression.Var.Node, "expected type %s but found %s", types.PrintableName(expectedType), types.PrintableName(overType))
		}

		pkg, name := binding.GetPackageLevelAndUnaliasedNameOfVariable(scope, expression.Var)
		astExp := ast.Reference{
			VariableType: overType,
			PackageName:  pkg,
			Name:         name,
		}

		return astExp, nil
	}
}

func expectTypeOfLiteral(expectedType types.VariableType, expression parser.LiteralExpression, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	varType, err := typeOfExpression(expression, file, scope)
	if err != nil {
		return nil, err
	}
	if !types.VariableTypeContainedIn(varType, expectedType) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected type %s but found %s", types.PrintableName(expectedType), types.PrintableName(varType))
	}
	return ast.Literal{
		VariableType: varType,
		Literal:      expression.Literal,
	}, nil
}
