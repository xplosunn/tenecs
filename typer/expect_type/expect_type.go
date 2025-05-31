package expect_type

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/type_of"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/exp/slices"
)

var ForbiddenVariableNames = []string{
	"true",
	"false",
}

func ExpectTypeOfExpressionBox(expectedType types.VariableType, expressionBox parser.ExpressionBox, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	if len(expressionBox.AccessOrInvocationChain) == 0 {
		return expectTypeOfExpression(expectedType, expressionBox.Expression, file, scope)
	}

	varType, err := type_of.TypeOfExpression(expressionBox.Expression, file, scope)
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
				return nil, type_error.PtrOnNodef(file, accessOrInvocation.Node, "Expected %s but got %s", types.PrintableName(expectedType), types.PrintableName(gotVarType))
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
		lhsVarType, err = type_of.TypeOfAccess(lhsVarType, accessOrInvocation.DotOrArrowName.VarName, file, scope)
		if err != nil {
			return nil, err
		}

		astExp = ast.Access{
			CodePoint:    over.SourceCodePoint(),
			VariableType: lhsVarType,
			Over:         over,
			Access:       accessOrInvocation.DotOrArrowName.VarName.String,
		}
	}
	if accessOrInvocation.Arguments != nil {
		function, ok := lhsVarType.(*types.Function)
		if !ok {
			return nil, type_error.PtrOnNodef(file, accessOrInvocation.Arguments.Node, "Should be a function in order to be invoked but is %s", types.PrintableName(lhsVarType))
		}
		if len(function.Arguments) != len(accessOrInvocation.Arguments.Arguments) {
			return nil, type_error.PtrOnNodef(file, accessOrInvocation.Arguments.Node, "Invoked with wrong number of arguments, expected %d but got %d", len(function.Arguments), len(accessOrInvocation.Arguments.Arguments))
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
			CodePoint:    over.SourceCodePoint(),
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
		func(generics *parser.LambdaOrListGenerics, expression parser.Lambda) {
			astExp, err = expectTypeOfLambda(expectedType, generics, expression, file, scope)
		},
		func(expression parser.Declaration) {
			astExp, err = expectTypeOfDeclaration(expectedType, expression, file, scope)
		},
		func(expression parser.If) {
			astExp, err = expectTypeOfIf(expectedType, expression, file, scope)
		},
		func(generics *parser.LambdaOrListGenerics, expression parser.List) {
			astExp, err = expectTypeOfList(expectedType, generics, expression, file, scope)
		},
		func(expression parser.When) {
			astExp, err = expectTypeOfWhen(expectedType, expression, file, scope)
		},
	)
	return astExp, err
}

func expectTypeOfWhen(expectedType types.VariableType, expression parser.When, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	typeOfOver, err := type_of.TypeOfExpressionBox(expression.Over, file, scope)
	if err != nil {
		return nil, err
	}
	typeOfOver = types.FlattenOr(typeOfOver)
	_, _, _, _, typeOverOr := typeOfOver.VariableTypeCases()
	if typeOverOr == nil {
		typeOverOr = &types.OrVariableType{
			Elements: []types.VariableType{typeOfOver},
		}
	}

	missingCases := []types.VariableType{}
	for _, varType := range typeOverOr.Elements {
		missingCases = append(missingCases, varType)
	}

	astOver, err := ExpectTypeOfExpressionBox(typeOfOver, expression.Over, file, scope)
	if err != nil {
		return nil, err
	}

	cases := []ast.WhenCase{}

	for _, whenIs := range expression.Is {
		varType, err := scopecheck.ValidateTypeAnnotationInScope(whenIs.Type, file, scope)
		if err != nil {
			return nil, type_error.FromScopeCheckError(file, err)
		}

		{
			_, matchableErr := AsMatchable(varType, binding.GetAllFieldsWithRef(scope))
			if matchableErr != nil {
				return nil, type_error.PtrOnNodef(file, whenIs.Type.Node, matchableErr.Error())
			}
		}

		isMissingCase := types.VariableTypeContainedIn(varType, &types.OrVariableType{
			Elements: missingCases,
		})
		if isMissingCase {
			_, _, _, _, missingOrToDeleteOr := varType.VariableTypeCases()
			if missingOrToDeleteOr == nil {
				missingOrToDeleteOr = &types.OrVariableType{
					Elements: []types.VariableType{
						varType,
					},
				}
			}
			for _, elementToDelete := range missingOrToDeleteOr.Elements {
				for i, missingCase := range missingCases {
					if types.VariableTypeContainedIn(elementToDelete, missingCase) {
						missingCases = append(missingCases[:i], missingCases[i+1:]...)
						break
					}
				}

			}
			localScope := scope
			if whenIs.Name != nil {
				var err *binding.ResolutionError
				localScope, err = binding.CopyAddingLocalVariable(localScope, *whenIs.Name, varType)
				if err != nil {
					return nil, type_error.FromResolutionError(file, whenIs.Name.Node, err)
				}
			}
			astThen, err := expectTypeOfBlock(expectedType, whenIs.Node, whenIs.ThenBlock, file, localScope)
			if err != nil {
				return nil, err
			}
			whenCase := ast.WhenCase{
				Name:         nil,
				VariableType: varType,
				Block:        astThen,
			}
			if whenIs.Name != nil {
				whenCase.Name = &whenIs.Name.String
			}
			cases = append(cases, whenCase)
		} else {
			return nil, type_error.PtrOnNodef(file, whenIs.Node, "no matching for %s in %s", types.PrintableName(varType), types.PrintableName(typeOfOver))
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
				return nil, type_error.FromResolutionError(file, expression.Other.Name.Node, err)
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
		return nil, type_error.PtrOnNodef(file, expression.Node, "missing cases for %s", varTypeNames)
	}

	return ast.When{
		CodePoint:     codePoint(file, expression.Node),
		VariableType:  expectedType,
		Over:          astOver,
		Cases:         cases,
		OtherCase:     otherCase,
		OtherCaseName: otherCaseName,
	}, nil
}

func expectTypeOfList(expectedType types.VariableType, generics *parser.LambdaOrListGenerics, expression parser.List, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	var expectedListOf types.VariableType

	if generics != nil {
		if len(generics.Generics) != 1 {
			return nil, type_error.PtrOnNodef(file, generics.Node, "Expected 1 generic")
		}
		varType, err := scopecheck.ValidateTypeAnnotationInScope(generics.Generics[0], file, scope)
		if err != nil {
			return nil, type_error.FromScopeCheckError(file, err)
		}
		expectedListOf = varType
	} else if len(expression.Expressions) == 0 {
		_, caseList, _, _, _ := expectedType.VariableTypeCases()
		if caseList != nil {
			return ast.List{
				CodePoint:             codePoint(file, expression.Node),
				ContainedVariableType: caseList.Generic,
				Arguments:             []ast.Expression{},
			}, nil
		} else {
			return nil, type_error.PtrOnNodef(file, expression.Node, "Could not infer list generic, please annotate it")
		}
	} else {
		or := &types.OrVariableType{
			Elements: []types.VariableType{},
		}
		for _, expressionBox := range expression.Expressions {
			varType, err := type_of.TypeOfExpressionBox(expressionBox, file, scope)
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

	expectedList := &types.List{
		Generic: expectedListOf,
	}
	if !types.VariableTypeContainedIn(expectedList, expectedType) {
		return nil, type_error.PtrOnNodef(file, expression.Node, "expected %s but got %s", types.PrintableName(expectedType), types.PrintableName(expectedList))
	}

	astArguments := []ast.Expression{}
	for _, expressionBox := range expression.Expressions {
		astExp, err := ExpectTypeOfExpressionBox(expectedListOf, expressionBox, file, scope)
		if err != nil {
			return nil, err
		}
		astArguments = append(astArguments, astExp)
	}

	return ast.List{
		CodePoint:             codePoint(file, expression.Node),
		ContainedVariableType: expectedListOf,
		Arguments:             astArguments,
	}, nil
}

func expectTypeOfIf(expectedType types.VariableType, expression parser.If, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	astCondition, err := ExpectTypeOfExpressionBox(types.Boolean(), expression.Condition, file, scope)
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
		CodePoint:    codePoint(file, expression.Node),
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
	if slices.Contains(ForbiddenVariableNames, expression.Name.String) {
		return nil, type_error.PtrOnNodef(file, expression.Name.Node, "Variable can't be named '%s'", expression.Name.String)
	}
	if !types.VariableTypeEq(expectedDeclarationType, types.Void()) {
		return nil, type_error.PtrOnNodef(file, expression.Name.Node, "Expected type %s but got void", types.PrintableName(expectedDeclarationType))
	}

	var expectedType types.VariableType
	var err *type_error.TypecheckError
	if expression.TypeAnnotation != nil {
		var err2 scopecheck.ScopeCheckError
		expectedType, err2 = scopecheck.ValidateTypeAnnotationInScope(*expression.TypeAnnotation, file, scope)
		err = type_error.FromScopeCheckError(file, err2)
	} else {
		expectedType, err = type_of.TypeOfExpressionBox(expression.ExpressionBox, file, scope)
	}
	if err != nil {
		return nil, err
	}
	scope, resolutionErr := binding.CopyAddingLocalVariable(scope, expression.Name, expectedType)
	if resolutionErr != nil {
		return nil, type_error.FromResolutionError(file, expression.Name.Node, resolutionErr)
	}
	astExp, err := ExpectTypeOfExpressionBox(expectedType, expression.ExpressionBox, file, scope)
	if err != nil {
		return nil, err
	}
	return ast.Declaration{
		CodePoint:  codePoint(file, expression.Name.Node),
		Name:       expression.Name.String,
		Expression: astExp,
	}, nil
}

func expectTypeOfLambda(expectedType types.VariableType, generics *parser.LambdaOrListGenerics, expression parser.Lambda, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	_, _, _, expectedFunction, expectedOr := expectedType.VariableTypeCases()
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
		return nil, type_error.PtrOnNodef(file, expression.Node, "Expected %s but got a function", types.PrintableName(expectedType))
	}

	signatureGenerics := []parser.TypeAnnotation{}
	if generics != nil {
		signatureGenerics = generics.Generics
	}
	if len(signatureGenerics) != len(expectedFunction.Generics) {
		return nil, type_error.PtrOnNodef(file, expression.Node, "expected %d generics but got %d", len(expectedFunction.Generics), len(signatureGenerics))
	}

	localScope := scope
	for _, genericTypeAnnotation := range signatureGenerics {
		generic, singleTypeNameErr := type_of.ExpectSingleTypeName(genericTypeAnnotation, file)
		if singleTypeNameErr != nil {
			return nil, singleTypeNameErr
		}
		var err *binding.ResolutionError
		localScope, err = binding.CopyAddingTypeToAllFiles(localScope, generic, &types.TypeArgument{Name: generic.String})
		if err != nil {
			return nil, type_error.FromResolutionError(file, generic.Node, err)
		}
	}

	if len(expression.Signature.Parameters) != len(expectedFunction.Arguments) {
		return nil, type_error.PtrOnNodef(file, expression.Node, "expected %d params but got %d", len(expectedFunction.Arguments), len(expression.Signature.Parameters))
	}
	for i, parameter := range expression.Signature.Parameters {
		if parameter.Type != nil {
			var err scopecheck.ScopeCheckError
			paramType, err := scopecheck.ValidateTypeAnnotationInScope(*parameter.Type, file, localScope)
			if err != nil {
				return nil, type_error.FromScopeCheckError(file, err)
			}
			if !types.VariableTypeContainedIn(expectedFunction.Arguments[i].VariableType, paramType) {
				return nil, type_error.PtrOnNodef(file, expression.Node, "in parameter position %d expected type %s but you have annotated %s", i, types.PrintableName(expectedFunction.Arguments[i].VariableType), types.PrintableName(paramType))
			}
		}
		var err *binding.ResolutionError
		localScope, err = binding.CopyAddingLocalVariable(localScope, parameter.Name, expectedFunction.Arguments[i].VariableType)
		if err != nil {
			return nil, type_error.FromResolutionError(file, parameter.Name.Node, err)
		}
	}

	expectedTypeOfBlock := expectedFunction.ReturnType
	if expression.Signature.ReturnType != nil {
		returnType, err := scopecheck.ValidateTypeAnnotationInScope(*expression.Signature.ReturnType, file, localScope)
		if err != nil {
			return nil, type_error.FromScopeCheckError(file, err)
		}
		if !types.VariableTypeContainedIn(returnType, expectedFunction.ReturnType) {
			return nil, type_error.PtrOnNodef(file, expression.Node, "in return type expected type %s but you have annotated %s", types.PrintableName(expectedFunction.ReturnType), types.PrintableName(returnType))
		}
		expectedTypeOfBlock = returnType
	}

	astBlock, err := expectTypeOfBlock(expectedTypeOfBlock, expression.Node, expression.Block, file, localScope)
	if err != nil {
		return nil, err
	}

	varType := &types.Function{
		CodePointAsFirstArgument: expectedFunction.CodePointAsFirstArgument,
		Generics:                 expectedFunction.Generics,
		Arguments:                []types.FunctionArgument{},
		ReturnType:               expectedFunction.ReturnType,
	}
	for i, arg := range expectedFunction.Arguments {
		varType.Arguments = append(varType.Arguments, types.FunctionArgument{
			Name:         expression.Signature.Parameters[i].Name.String,
			VariableType: arg.VariableType,
		})
	}

	return &ast.Function{
		CodePoint:    codePoint(file, expression.Node),
		VariableType: varType,
		Block:        astBlock,
	}, nil
}

func expectTypeOfBlock(expectedType types.VariableType, node parser.Node, block []parser.ExpressionBox, file string, scope binding.Scope) ([]ast.Expression, *type_error.TypecheckError) {
	result := []ast.Expression{}

	if len(block) == 0 {
		return nil, type_error.PtrOnNodef(file, node, "empty function block not allowed (maybe you want to return null?)")
	}

	localScope := scope
	for i, expressionBox := range block {
		var expectedTypeOfExpressionBox = expectedType
		var err *type_error.TypecheckError
		if i < len(block)-1 {
			expectedTypeOfExpressionBox, err = type_of.TypeOfExpressionBox(expressionBox, file, localScope)
			if err != nil {
				return nil, err
			}
		}
		astExp, err := ExpectTypeOfExpressionBox(expectedTypeOfExpressionBox, expressionBox, file, localScope)
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
				return nil, type_error.FromResolutionError(file, expressionBox.Node, err)
			}
		}
	}

	return result, nil
}

func resolveFunctionGenerics(node parser.Node, function *types.Function, genericsPassed []parser.TypeAnnotation, argumentsPassed []parser.NamedArgument, expectedReturnType *types.VariableType, file string, scope binding.Scope) (*types.Function, []types.VariableType, []ast.Expression, *type_error.TypecheckError) {
	generics := []types.VariableType{}

	genericsPassedContainsUnderscore := false
	for _, passed := range genericsPassed {
		for _, element := range passed.OrTypes {
			var err *type_error.TypecheckError
			parser.TypeAnnotationElementExhaustiveSwitch(
				element,
				func(underscoreTypeAnnotation parser.SingleNameType) {
					if len(passed.OrTypes) > 1 {
						err = type_error.PtrOnNodef(file, underscoreTypeAnnotation.Node, "Cannot infer part of an or type")
						return
					}
					genericsPassedContainsUnderscore = true
				},
				func(typeAnnotation parser.SingleNameType) {},
				func(typeAnnotation parser.FunctionType) {
					err = type_error.PtrOnNodef(file, node, "Can't pass a function as a generic")
				},
			)
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}

	if genericsPassedContainsUnderscore || (len(genericsPassed) == 0 && len(function.Generics) > 0) {
		inferredGenerics, err := type_of.AttemptGenericInference(node, function, argumentsPassed, genericsPassed, expectedReturnType, file, scope)
		if err != nil {
			return nil, nil, nil, err
		}
		generics = inferredGenerics
	} else {
		if len(genericsPassed) != len(function.Generics) {
			return nil, nil, nil, type_error.PtrOnNodef(file, node, "expected %d generics but got %d", len(function.Generics), len(genericsPassed))
		}
		for _, generic := range genericsPassed {
			varType, err := scopecheck.ValidateTypeAnnotationInScope(generic, file, scope)
			if err != nil {
				return nil, nil, nil, type_error.FromScopeCheckError(file, err)
			}
			generics = append(generics, varType)
		}
	}
	genericsMap := map[string]types.VariableType{}
	for i, varType := range generics {
		genericsMap[function.Generics[i]] = varType
	}

	if len(argumentsPassed) != len(function.Arguments) {
		return nil, nil, nil, type_error.PtrOnNodef(file, node, "expected %d arguments but got %d", len(function.Arguments), len(argumentsPassed))
	}
	arguments := []types.FunctionArgument{}
	for _, argument := range function.Arguments {
		arguments = append(arguments, argument)
	}
	for i := 0; i < len(arguments); i++ {
		newVarType, err := binding.ResolveGeneric(arguments[i].VariableType, genericsMap)
		if err != nil {
			return nil, nil, nil, type_error.FromResolutionError(file, node, err)
		}
		arguments[i].VariableType = newVarType
	}
	astArguments := []ast.Expression{}
	for i, argument := range argumentsPassed {
		if argument.Name != nil && argument.Name.String != arguments[i].Name {
			return nil, nil, nil, type_error.PtrOnNodef(file, argument.Name.Node, "name of argument should be '%s'", arguments[i].Name)
		}
		expectedArgType := arguments[i].VariableType
		astArg, err := ExpectTypeOfExpressionBox(expectedArgType, argument.Argument, file, scope)
		if err != nil {
			return nil, nil, nil, err
		}
		astArguments = append(astArguments, astArg)
	}

	returnType := function.ReturnType
	newVarType, err2 := binding.ResolveGeneric(returnType, genericsMap)
	if err2 != nil {
		return nil, nil, nil, type_error.FromResolutionError(file, node, err2)
	}
	returnType = newVarType

	return &types.Function{
		Generics:   nil,
		Arguments:  arguments,
		ReturnType: returnType,
	}, generics, astArguments, nil
}

func expectTypeOfReferenceOrInvocation(expectedType types.VariableType, expression parser.ReferenceOrInvocation, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	overType, ok := binding.GetTypeByVariableName(scope, file, expression.Var.String)
	if !ok {
		return nil, type_error.PtrOnNodef(file, expression.Var.Node, "Not found in scope: %s", expression.Var.String)
	}

	if expression.Arguments != nil {
		arguments := []ast.Expression{}
		overFunction, ok := overType.(*types.Function)
		if !ok {
			return nil, type_error.PtrOnNodef(file, expression.Arguments.Node, "Can't invoke on not a function")
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
			return nil, type_error.PtrOnNodef(file, expression.Var.Node, "expected type %s but found %s", types.PrintableName(expectedType), types.PrintableName(overFunction.ReturnType))
		}

		pkg, name := binding.GetPackageLevelAndUnaliasedNameOfVariable(scope, file, expression.Var)
		astExp := ast.Invocation{
			CodePoint:    codePoint(file, expression.Var.Node),
			VariableType: overFunction.ReturnType,
			Over: ast.Reference{
				CodePoint:    codePoint(file, expression.Var.Node),
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
			return nil, type_error.PtrOnNodef(file, expression.Var.Node, "expected type %s but found %s", types.PrintableName(expectedType), types.PrintableName(overType))
		}

		pkg, name := binding.GetPackageLevelAndUnaliasedNameOfVariable(scope, file, expression.Var)
		astExp := ast.Reference{
			CodePoint:    codePoint(file, expression.Var.Node),
			VariableType: overType,
			PackageName:  pkg,
			Name:         name,
		}

		return astExp, nil
	}
}

func expectTypeOfLiteral(expectedType types.VariableType, expression parser.LiteralExpression, file string, scope binding.Scope) (ast.Expression, *type_error.TypecheckError) {
	varType, err := type_of.TypeOfExpression(expression, file, scope)
	if err != nil {
		return nil, err
	}
	if !types.VariableTypeContainedIn(varType, expectedType) {
		return nil, type_error.PtrOnNodef(file, expression.Node, "expected type %s but found %s", types.PrintableName(expectedType), types.PrintableName(varType))
	}
	return ast.Literal{
		CodePoint:    codePoint(file, expression.Node),
		VariableType: varType,
		Literal:      expression.Literal,
	}, nil
}

func codePoint(fileName string, node parser.Node) ast.CodePoint {
	var emptyNode parser.Node
	if node == emptyNode {
		panic("missing node")
	}
	return ast.CodePoint{
		FileName: fileName,
		Line:     node.Pos.Line,
	}
}
