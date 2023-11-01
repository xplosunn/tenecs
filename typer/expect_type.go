package typer

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func expectTypeOfExpressionBox(expectedType types.VariableType, expressionBox parser.ExpressionBox, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	if len(expressionBox.AccessOrInvocationChain) == 0 {
		return expectTypeOfExpression(expectedType, expressionBox.Expression, file, universe)
	}

	varType, err := typeOfExpression(expressionBox.Expression, file, universe)
	if err != nil {
		return nil, err
	}

	astExp, err := expectTypeOfExpression(varType, expressionBox.Expression, file, universe)
	if err != nil {
		return nil, err
	}

	for i, accessOrInvocation := range expressionBox.AccessOrInvocationChain {
		astExp, err = determineTypeOfAccessOrInvocation(astExp, accessOrInvocation, file, universe)
		if err != nil {
			return nil, err
		}
		if i == len(expressionBox.AccessOrInvocationChain)-1 {
			gotVarType := ast.VariableTypeOfExpression(astExp)
			if !types.VariableTypeContainedIn(gotVarType, expectedType) {
				return nil, type_error.PtrOnNodef(accessOrInvocation.VarName.Node, "Expected %s but got %s", printableName(expectedType), printableName(gotVarType))
			}
		}
	}

	return astExp, nil
}

func determineTypeOfAccessOrInvocation(over ast.Expression, accessOrInvocation parser.AccessOrInvocation, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	lhsVarType := ast.VariableTypeOfExpression(over)
	astExp := over
	var err *type_error.TypecheckError
	if accessOrInvocation.VarName != nil {
		lhsVarType, err = typeOfAccess(lhsVarType, *accessOrInvocation.VarName, universe)
		if err != nil {
			return nil, err
		}

		astExp = ast.Access{
			VariableType: lhsVarType,
			Over:         over,
			Access:       accessOrInvocation.VarName.String,
		}
	}
	if accessOrInvocation.Arguments != nil {
		function, ok := lhsVarType.(*types.Function)
		if !ok {
			return nil, type_error.PtrOnNodef(accessOrInvocation.Arguments.Node, "Should be a function in order to be invoked but is %s", printableName(lhsVarType))
		}
		if len(function.Arguments) != len(accessOrInvocation.Arguments.Arguments) {
			return nil, type_error.PtrOnNodef(accessOrInvocation.Arguments.Node, "Invoked with wrong number of arguments, expected %d but got %d", len(function.Arguments), len(accessOrInvocation.Arguments.Arguments))
		}

		resolvedGenericsFunction, generics, arguments, err := resolveFunctionGenerics(
			accessOrInvocation.Arguments.Node,
			function,
			accessOrInvocation.Arguments.Generics,
			accessOrInvocation.Arguments.Arguments,
			file,
			universe,
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

func expectTypeOfExpression(expectedType types.VariableType, expression parser.Expression, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	var astExp ast.Expression
	var err *type_error.TypecheckError
	parser.ExpressionExhaustiveSwitch(
		expression,
		func(expression parser.Implementation) {
			astExp, err = expectTypeOfImplementation(expectedType, expression, file, universe)
		},
		func(expression parser.LiteralExpression) {
			astExp, err = expectTypeOfLiteral(expectedType, expression, file, universe)
		},
		func(expression parser.ReferenceOrInvocation) {
			astExp, err = expectTypeOfReferenceOrInvocation(expectedType, expression, file, universe)
		},
		func(expression parser.Lambda) {
			astExp, err = expectTypeOfLambda(expectedType, expression, file, universe)
		},
		func(expression parser.Declaration) {
			astExp, err = expectTypeOfDeclaration(expectedType, expression, file, universe)
		},
		func(expression parser.If) {
			astExp, err = expectTypeOfIf(expectedType, expression, file, universe)
		},
		func(expression parser.Array) {
			astExp, err = expectTypeOfArray(expectedType, expression, file, universe)
		},
		func(expression parser.When) {
			astExp, err = expectTypeOfWhen(expectedType, expression, file, universe)
		},
	)
	return astExp, err
}

func expectTypeOfWhen(expectedType types.VariableType, expression parser.When, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	typeOfOver, err := typeOfExpressionBox(expression.Over, file, universe)
	if err != nil {
		return nil, err
	}
	typeOverOr, ok := typeOfOver.(*types.OrVariableType)
	if !ok {
		return nil, type_error.PtrOnNodef(expression.Node, "use when only on an or type, not %s", printableName(typeOfOver))
	}

	missingCases := map[string]types.VariableType{}
	for _, varType := range typeOverOr.Elements {
		missingCases[printableName(varType)] = varType
	}

	astOver, err := expectTypeOfExpressionBox(typeOfOver, expression.Over, file, universe)
	if err != nil {
		return nil, err
	}

	cases := map[types.VariableType][]ast.Expression{}
	caseNames := map[types.VariableType]*string{}

	for _, whenIs := range expression.Is {
		varType, err := validateTypeAnnotationInUniverse(whenIs.Type, file, universe)
		if err != nil {
			return nil, err
		}
		if missingCases[printableName(varType)] != nil {
			delete(missingCases, printableName(varType))
			localUniverse := universe
			if whenIs.Name != nil {
				localUniverse, err = binding.CopyAddingLocalVariable(localUniverse, *whenIs.Name, varType)
				if err != nil {
					return nil, err
				}
			}
			astThen, err := expectTypeOfBlock(expectedType, whenIs.Node, whenIs.ThenBlock, file, localUniverse)
			if err != nil {
				return nil, err
			}
			cases[varType] = astThen
			if whenIs.Name != nil {
				caseNames[varType] = &whenIs.Name.String
			}
		} else {
			return nil, type_error.PtrOnNodef(whenIs.Node, "no matching for %s in %s", printableName(varType), printableName(typeOfOver))
		}
	}

	if expression.Other != nil {
		orCases := []types.VariableType{}
		for _, variableType := range missingCases {
			orCases = append(orCases, variableType)
		}
		missingCases = nil
		varType := &types.OrVariableType{Elements: orCases}
		localUniverse := universe
		if expression.Other.Name != nil {
			localUniverse, err = binding.CopyAddingLocalVariable(universe, *expression.Other.Name, varType)
			if err != nil {
				return nil, err
			}
		}
		astThen, err := expectTypeOfBlock(expectedType, expression.Other.Node, expression.Other.ThenBlock, file, localUniverse)
		if err != nil {
			return nil, err
		}
		cases[varType] = astThen
		if expression.Other.Name != nil {
			caseNames[varType] = &expression.Other.Name.String
		}
	}

	if len(missingCases) > 0 {
		varTypeNames := ""
		for _, varType := range missingCases {
			if varTypeNames != "" {
				varTypeNames += ", "
			}
			varTypeNames += printableName(varType)
		}
		return nil, type_error.PtrOnNodef(expression.Node, "missing cases for %s", varTypeNames)
	}

	return ast.When{
		VariableType: expectedType,
		Over:         astOver,
		Cases:        cases,
		CaseNames:    caseNames,
	}, nil
}

func expectTypeOfArray(expectedType types.VariableType, expression parser.Array, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	var expectedArrayOf types.VariableType

	if expression.Generic != nil {
		varType, err := validateTypeAnnotationInUniverse(*expression.Generic, file, universe)
		if err != nil {
			return nil, err
		}
		expectedArrayOf = varType
	} else {
		or := &types.OrVariableType{
			Elements: []types.VariableType{},
		}
		for _, expressionBox := range expression.Expressions {
			varType, err := typeOfExpressionBox(expressionBox, file, universe)
			if err != nil {
				return nil, err
			}
			types.VariableTypeAddToOr(varType, or)
		}
		if len(or.Elements) == 0 {
			panic("TODO expectTypeOfArray invalid")
		} else if len(or.Elements) == 1 {
			expectedArrayOf = or.Elements[0]
		} else {
			expectedArrayOf = or
		}
	}

	expectedArray, ok := types.Array(expectedArrayOf)
	if !ok {
		return nil, type_error.PtrOnNodef(expression.Node, "not a valid generic: %s", printableName(expectedArrayOf))
	}
	if !types.VariableTypeContainedIn(expectedArray, expectedType) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected %s but got %s", printableName(expectedType), printableName(expectedArray))
	}

	astArguments := []ast.Expression{}
	for _, expressionBox := range expression.Expressions {
		astExp, err := expectTypeOfExpressionBox(expectedArrayOf, expressionBox, file, universe)
		if err != nil {
			return nil, err
		}
		astArguments = append(astArguments, astExp)
	}

	return ast.Array{
		ContainedVariableType: expectedArrayOf,
		Arguments:             astArguments,
	}, nil
}

func expectTypeOfIf(expectedType types.VariableType, expression parser.If, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	astCondition, err := expectTypeOfExpressionBox(types.Boolean(), expression.Condition, file, universe)
	if err != nil {
		return nil, err
	}

	thenBlock, err := expectTypeOfBlock(expectedType, expression.Node, expression.ThenBlock, file, universe)
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
		block, err := expectTypeOfBlock(expectedType, expression.Node, expression.ElseBlock, file, universe)
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

func expectTypeOfDeclaration(expectedDeclarationType types.VariableType, expression parser.Declaration, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	if !types.VariableTypeEq(expectedDeclarationType, types.Void()) {
		return nil, type_error.PtrOnNodef(expression.Name.Node, "Expected type %s but got void", printableName(expectedDeclarationType))
	}

	var expectedType types.VariableType
	var err *type_error.TypecheckError
	if expression.TypeAnnotation != nil {
		expectedType, err = validateTypeAnnotationInUniverse(*expression.TypeAnnotation, file, universe)
	} else {
		expectedType, err = typeOfExpressionBox(expression.ExpressionBox, file, universe)
	}
	if err != nil {
		return nil, err
	}
	astExp, err := expectTypeOfExpressionBox(expectedType, expression.ExpressionBox, file, universe)
	if err != nil {
		return nil, err
	}
	return ast.Declaration{
		Name:       expression.Name.String,
		Expression: astExp,
	}, nil
}

func expectTypeOfLambda(expectedType types.VariableType, expression parser.Lambda, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	expectedFunction, ok := expectedType.(*types.Function)
	if !ok {
		return nil, type_error.PtrOnNodef(expression.Node, "Expected %s but got a function", printableName(expectedType))
	}

	if len(expression.Generics) != len(expectedFunction.Generics) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected %d generics but got %d", len(expectedFunction.Generics), len(expression.Generics))
	}

	localUniverse := universe
	var err *type_error.TypecheckError
	for _, generic := range expression.Generics {
		localUniverse, err = binding.CopyAddingTypeToAllFiles(localUniverse, generic, &types.TypeArgument{Name: generic.String})
		if err != nil {
			return nil, err
		}
	}

	if len(expression.Parameters) != len(expectedFunction.Arguments) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected %d params but got %d", len(expectedFunction.Arguments), len(expression.Parameters))
	}
	for i, parameter := range expression.Parameters {
		if parameter.Type != nil {
			paramType, err := validateTypeAnnotationInUniverse(*parameter.Type, file, localUniverse)
			if err != nil {
				return nil, err
			}
			if !types.VariableTypeContainedIn(expectedFunction.Arguments[i].VariableType, paramType) {
				return nil, type_error.PtrOnNodef(expression.Node, "in parameter position %d expected type %s but you have annotated %s", i, printableName(expectedFunction.Arguments[i].VariableType), printableName(paramType))
			}
		}
		localUniverse, err = binding.CopyAddingLocalVariable(localUniverse, parameter.Name, expectedFunction.Arguments[i].VariableType)
		if err != nil {
			return nil, err
		}
	}

	if expression.ReturnType != nil {
		returnType, err := validateTypeAnnotationInUniverse(*expression.ReturnType, file, localUniverse)
		if err != nil {
			return nil, err
		}
		if !types.VariableTypeContainedIn(returnType, expectedFunction.ReturnType) {
			return nil, type_error.PtrOnNodef(expression.Node, "in return type expected type %s but you have annotated %s", printableName(expectedFunction.ReturnType), printableName(returnType))
		}
	}

	astBlock, err := expectTypeOfBlock(expectedFunction.ReturnType, expression.Node, expression.Block, file, localUniverse)
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
			Name:         expression.Parameters[i].Name.String,
			VariableType: arg.VariableType,
		})
	}

	return &ast.Function{
		VariableType: varType,
		Block:        astBlock,
	}, nil
}

func expectTypeOfBlock(expectedType types.VariableType, node parser.Node, block []parser.ExpressionBox, file string, universe binding.Universe) ([]ast.Expression, *type_error.TypecheckError) {
	result := []ast.Expression{}

	if len(block) == 0 {
		if types.VariableTypeContainedIn(types.Void(), expectedType) {
			return result, nil
		} else {
			return nil, type_error.PtrOnNodef(node, "empty block only allowed for Void type")
		}
	}

	localUniverse := universe
	for i, expressionBox := range block {
		var expectedTypeOfExpressionBox = expectedType
		var err *type_error.TypecheckError
		if i < len(block)-1 {
			expectedTypeOfExpressionBox, err = typeOfExpressionBox(expressionBox, file, localUniverse)
			if err != nil {
				return nil, err
			}
		}
		astExp, err := expectTypeOfExpressionBox(expectedTypeOfExpressionBox, expressionBox, file, localUniverse)
		if err != nil {
			return nil, err
		}
		result = append(result, astExp)
		astDec, isDec := astExp.(ast.Declaration)
		if isDec {
			localUniverse, err = binding.CopyAddingLocalVariable(localUniverse, parser.Name{
				String: astDec.Name,
			}, ast.VariableTypeOfExpression(astDec.Expression))
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func resolveFunctionGenerics(node parser.Node, function *types.Function, genericsPassed []parser.TypeAnnotation, argumentsPassed []parser.ExpressionBox, file string, universe binding.Universe) (*types.Function, []types.VariableType, []ast.Expression, *type_error.TypecheckError) {
	generics := []types.VariableType{}

	if len(genericsPassed) == 0 && len(function.Generics) > 0 {
		inferredGenerics, err := attemptGenericInference(node, function, argumentsPassed, file, universe)
		if err != nil {
			return nil, nil, nil, err
		}
		generics = inferredGenerics
	} else {
		if len(genericsPassed) != len(function.Generics) {
			return nil, nil, nil, type_error.PtrOnNodef(node, "expected %d generics but got %d", len(function.Generics), len(genericsPassed))
		}
		for _, generic := range genericsPassed {
			varType, err := validateTypeAnnotationInUniverse(generic, file, universe)
			if err != nil {
				return nil, nil, nil, err
			}
			if !varType.CanBeStructField() {
				return nil, nil, nil, type_error.PtrOnNodef(generic.Node, "invalid generic type %s", printableName(varType))
			}
			generics = append(generics, varType)
		}
	}
	genericsMap := map[string]types.VariableType{}
	for i, varType := range generics {
		genericsMap[function.Generics[i]] = varType
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
		expectedArgType := arguments[i].VariableType
		astArg, err := expectTypeOfExpressionBox(expectedArgType, argument, file, universe)
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

func attemptGenericInference(node parser.Node, function *types.Function, argumentsPassed []parser.ExpressionBox, file string, universe binding.Universe) ([]types.VariableType, *type_error.TypecheckError) {
	resolvedGenerics := []types.VariableType{}
	for _, functionGenericName := range function.Generics {
		var found types.VariableType
		for i, arg := range argumentsPassed {
			typeOfArg, err := typeOfExpressionBox(arg, file, universe)
			if err != nil {
				return nil, err
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
		if found == nil {
			return nil, type_error.PtrOnNodef(node, "Could not infer generics, please annotate them")
		}
		resolvedGenerics = append(resolvedGenerics, found)
	}
	return resolvedGenerics, nil
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
		panic("TODO tryToInferGeneric Or")
	} else {
		return nil, true
	}
}

func expectTypeOfReferenceOrInvocation(expectedType types.VariableType, expression parser.ReferenceOrInvocation, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	overType, ok := binding.GetTypeByVariableName(universe, expression.Var.String)
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
			file,
			universe,
		)
		if err != nil {
			return nil, err
		}

		if !types.VariableTypeContainedIn(overFunction.ReturnType, expectedType) {
			return nil, type_error.PtrOnNodef(expression.Var.Node, "expected type %s but found %s", printableName(expectedType), printableName(overFunction.ReturnType))
		}

		pkg, name := binding.GetPackageLevelAndUnaliasedNameOfVariable(universe, expression.Var)
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
			return nil, type_error.PtrOnNodef(expression.Var.Node, "expected type %s but found %s", printableName(expectedType), printableName(overType))
		}

		pkg, name := binding.GetPackageLevelAndUnaliasedNameOfVariable(universe, expression.Var)
		astExp := ast.Reference{
			VariableType: overType,
			PackageName:  pkg,
			Name:         name,
		}

		return astExp, nil
	}
}

func expectTypeOfLiteral(expectedType types.VariableType, expression parser.LiteralExpression, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	varType, err := typeOfExpression(expression, file, universe)
	if err != nil {
		return nil, err
	}
	if !types.VariableTypeContainedIn(varType, expectedType) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected type %s but found %s", printableName(expectedType), printableName(varType))
	}
	return ast.Literal{
		VariableType: varType,
		Literal:      expression.Literal,
	}, nil
}

func expectTypeOfImplementation(expectedType types.VariableType, expression parser.Implementation, file string, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	generics := []types.VariableType{}
	for _, generic := range expression.Generics {
		varType, err := validateTypeAnnotationInUniverse(generic, file, universe)
		if err != nil {
			return nil, err
		}
		generics = append(generics, varType)
	}
	_, resolutionErr := binding.GetTypeByTypeName(universe, file, expression.Implementing.String, generics)
	if resolutionErr != nil {
		return nil, TypecheckErrorFromResolutionError(expression.Node, resolutionErr)
	}
	expectedInterface, ok := expectedType.(*types.KnownType)
	if !ok {
		return nil, type_error.PtrOnNodef(expression.Node, "Expected %s but got %s", printableName(expectedType), expression.Implementing.String)
	}

	expectedInterfaceFields, resolutionErr := binding.GetFields(universe, expectedInterface)
	if resolutionErr != nil {
		return nil, TypecheckErrorFromResolutionError(expression.Node, resolutionErr)
	}

	declarations := []parser.Declaration{}
	for _, implementationDeclaration := range expression.Declarations {
		if implementationDeclaration.Public {
			if expectedInterfaceFields[implementationDeclaration.Name.String] == nil {
				return nil, type_error.PtrOnNodef(expression.Node, "variable %s should not be public", implementationDeclaration.Name.String)
			}
		} else {
			if expectedInterfaceFields[implementationDeclaration.Name.String] != nil {
				return nil, type_error.PtrOnNodef(expression.Node, "variable %s should be public", implementationDeclaration.Name.String)
			}
		}
		declarations = append(declarations, parser.Declaration{
			Name:           implementationDeclaration.Name,
			TypeAnnotation: implementationDeclaration.TypeAnnotation,
			ExpressionBox: parser.ExpressionBox{
				Expression: implementationDeclaration.Expression,
			},
		})
	}

	astExpMap, err := TypecheckDeclarations(&expectedInterfaceFields, nil, expression.Node, map[string][]parser.Declaration{file: declarations}, universe)
	if err != nil {
		return nil, err
	}
	return ast.Implementation{
		Implements: expectedInterface,
		Variables:  astExpMap,
	}, nil
}
