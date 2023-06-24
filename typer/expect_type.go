package typer

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func expectTypeOfExpressionBox(expectedType types.VariableType, expressionBox parser.ExpressionBox, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	if len(expressionBox.AccessOrInvocationChain) == 0 {
		return expectTypeOfExpression(expectedType, expressionBox.Expression, universe)
	}

	varType, err := typeOfExpression(expressionBox.Expression, universe)
	if err != nil {
		return nil, err
	}

	astExp, err := expectTypeOfExpression(varType, expressionBox.Expression, universe)
	if err != nil {
		return nil, err
	}

	for _, accessOrInvocation := range expressionBox.AccessOrInvocationChain {
		astExp, err = determineTypeOfAccessOrInvocation(astExp, accessOrInvocation, universe)
		if err != nil {
			return nil, err
		}
	}

	return astExp, nil
}

func determineTypeOfAccessOrInvocation(over ast.Expression, accessOrInvocation parser.AccessOrInvocation, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	if accessOrInvocation.Arguments != nil {
		accessVarType, err := typeOfAccess(ast.VariableTypeOfExpression(over), accessOrInvocation.VarName)
		if err != nil {
			return nil, err
		}
		function, ok := accessVarType.(*types.Function)
		if !ok {
			return nil, type_error.PtrOnNodef(accessOrInvocation.Arguments.Node, "Should be a function in order to be invoked but is %s", printableName(accessVarType))
		}
		if len(function.Generics) != len(accessOrInvocation.Arguments.Generics) {
			return nil, type_error.PtrOnNodef(accessOrInvocation.Arguments.Node, "Invoked with wrong number of generics, expected %d but got %d", len(function.Generics), len(accessOrInvocation.Arguments.Generics))
		}
		if len(function.Arguments) != len(accessOrInvocation.Arguments.Arguments) {
			return nil, type_error.PtrOnNodef(accessOrInvocation.Arguments.Node, "Invoked with wrong number of arguments, expected %d but got %d", len(function.Arguments), len(accessOrInvocation.Arguments.Arguments))
		}

		resolvedGenericsFunction, generics, err := resolveFunctionGenerics(accessOrInvocation.Arguments.Node, function, accessOrInvocation.Arguments.Generics, universe)
		if err != nil {
			return nil, err
		}

		arguments := []ast.Expression{}
		for i, argument := range accessOrInvocation.Arguments.Arguments {
			expectedArgType := resolvedGenericsFunction.Arguments[i].VariableType
			astArg, err := expectTypeOfExpressionBox(expectedArgType, argument, universe)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, astArg)
		}

		astExp := ast.Invocation{
			VariableType: resolvedGenericsFunction.ReturnType,
			Over: ast.Access{
				VariableType: resolvedGenericsFunction,
				Over:         over,
				Access:       accessOrInvocation.VarName.String,
			},
			Generics:  generics,
			Arguments: arguments,
		}

		return astExp, nil
	} else {
		accessVarType, err := typeOfAccess(ast.VariableTypeOfExpression(over), accessOrInvocation.VarName)
		if err != nil {
			return nil, err
		}
		varType := accessVarType

		astExp := ast.Access{
			VariableType: varType,
			Over:         over,
			Access:       accessOrInvocation.VarName.String,
		}

		return astExp, nil
	}
}

func expectTypeOfExpression(expectedType types.VariableType, expression parser.Expression, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	var astExp ast.Expression
	var err *type_error.TypecheckError
	parser.ExpressionExhaustiveSwitch(
		expression,
		func(expression parser.Module) {
			astExp, err = expectTypeOfModule(expectedType, expression, universe)
		},
		func(expression parser.LiteralExpression) {
			astExp, err = expectTypeOfLiteral(expectedType, expression, universe)
		},
		func(expression parser.ReferenceOrInvocation) {
			astExp, err = expectTypeOfReferenceOrInvocation(expectedType, expression, universe)
		},
		func(expression parser.Lambda) {
			astExp, err = expectTypeOfLambda(expectedType, expression, universe)
		},
		func(expression parser.Declaration) {
			astExp, err = expectTypeOfDeclaration(expectedType, expression, universe)
		},
		func(expression parser.If) {
			astExp, err = expectTypeOfIf(expectedType, expression, universe)
		},
		func(expression parser.Array) {
			astExp, err = expectTypeOfArray(expectedType, expression, universe)
		},
		func(expression parser.When) {
			astExp, err = expectTypeOfWhen(expectedType, expression, universe)
		},
	)
	return astExp, err
}

func expectTypeOfWhen(expectedType types.VariableType, expression parser.When, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	typeOfOver, err := typeOfExpressionBox(expression.Over, universe)
	if err != nil {
		return nil, err
	}
	typeOverOr, ok := typeOfOver.(*types.OrVariableType)
	if !ok {
		return nil, type_error.PtrOnNodef(expression.Node, "use when only on an or type, not %s", printableName(typeOfOver))
	}

	missingCases := map[types.VariableType]bool{}
	for _, varType := range typeOverOr.Elements {
		missingCases[varType] = true
	}

	astOver, err := expectTypeOfExpressionBox(typeOfOver, expression.Over, universe)
	if err != nil {
		return nil, err
	}
	overRefName := ""
	if overRef, ok := astOver.(ast.Reference); ok {
		overRefName = overRef.Name
	}

	cases := map[types.VariableType][]ast.Expression{}

	for _, whenIs := range expression.Is {
		varType, err := validateTypeAnnotationInUniverse(whenIs.Is, universe)
		if err != nil {
			return nil, err
		}
		if missingCases[varType] {
			delete(missingCases, varType)
			localUniverse := universe
			if overRefName != "" {
				localUniverse, err = binding.CopyOverridingVariableType(universe, overRefName, varType)
				if err != nil {
					return nil, err
				}
			}
			astThen, err := expectTypeOfBlock(expectedType, whenIs.Node, whenIs.ThenBlock, localUniverse)
			if err != nil {
				return nil, err
			}
			cases[varType] = astThen
		}
	}

	if expression.Other != nil {
		orCases := []types.VariableType{}
		for variableType, _ := range missingCases {
			orCases = append(orCases, variableType)
		}
		missingCases = nil
		varType := &types.OrVariableType{Elements: orCases}
		localUniverse := universe
		if overRefName != "" {
			localUniverse, err = binding.CopyOverridingVariableType(universe, overRefName, varType)
			if err != nil {
				return nil, err
			}
		}
		astThen, err := expectTypeOfBlock(expectedType, expression.Other.Node, expression.Other.ThenBlock, localUniverse)
		if err != nil {
			return nil, err
		}
		cases[varType] = astThen
	}

	if len(missingCases) > 0 {
		varTypeNames := ""
		for varType, _ := range missingCases {
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
	}, nil
}

func expectTypeOfArray(expectedType types.VariableType, expression parser.Array, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	expectedArray, ok := expectedType.(*types.Array)
	if !ok {
		return nil, type_error.PtrOnNodef(expression.Node, "Expected %s but got array", printableName(expectedType))
	}
	expectedArrayOf := types.VariableTypeFromStructFieldVariableType(expectedArray.OfType)

	if expression.Generic != nil {
		varType, err := validateTypeAnnotationInUniverse(*expression.Generic, universe)
		if err != nil {
			return nil, err
		}
		_, ok := types.StructFieldVariableTypeFromVariableType(varType)
		if !ok {
			return nil, type_error.PtrOnNodef(expression.Node, "not a valid generic: %s", printableName(varType))
		}
		if !variableTypeContainedIn(varType, expectedArrayOf) {
			return nil, type_error.PtrOnNodef(expression.Node, "expected array of %s but got array of %s", printableName(expectedArrayOf), printableName(varType))
		}
	}

	astArguments := []ast.Expression{}
	for _, expressionBox := range expression.Expressions {
		astExp, err := expectTypeOfExpressionBox(expectedArrayOf, expressionBox, universe)
		if err != nil {
			return nil, err
		}
		astArguments = append(astArguments, astExp)
	}

	return ast.Array{
		ContainedVariableType: expectedArray.OfType,
		Arguments:             astArguments,
	}, nil
}

func expectTypeOfIf(expectedType types.VariableType, expression parser.If, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	astCondition, err := expectTypeOfExpressionBox(&types.BasicType{Type: "Boolean"}, expression.Condition, universe)
	if err != nil {
		return nil, err
	}

	thenBlock, err := expectTypeOfBlock(expectedType, expression.Node, expression.ThenBlock, universe)
	if err != nil {
		return nil, err
	}
	var elseBlock []ast.Expression = nil
	if len(expression.ElseBlock) != 0 {
		block, err := expectTypeOfBlock(expectedType, expression.Node, expression.ElseBlock, universe)
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

func expectTypeOfDeclaration(expectedType types.VariableType, expression parser.Declaration, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	_, ok := expectedType.(*types.Void)
	if !ok {
		return nil, type_error.PtrOnNodef(expression.Name.Node, "Expected type %s but got void", printableName(expectedType))
	}
	expectedType, err := typeOfExpressionBox(expression.ExpressionBox, universe)
	if err != nil {
		return nil, err
	}
	astExp, err := expectTypeOfExpressionBox(expectedType, expression.ExpressionBox, universe)
	if err != nil {
		return nil, err
	}
	return ast.Declaration{
		Name:       expression.Name.String,
		Expression: astExp,
	}, nil
}

func expectTypeOfLambda(expectedType types.VariableType, expression parser.Lambda, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
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
		localUniverse, err = binding.CopyAddingType(localUniverse, generic, &types.TypeArgument{Name: generic.String})
		if err != nil {
			return nil, err
		}
	}

	if len(expression.Parameters) != len(expectedFunction.Arguments) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected %d params but got %d", len(expectedFunction.Arguments), len(expression.Parameters))
	}
	for i, parameter := range expression.Parameters {
		if parameter.Type != nil {
			paramType, err := validateTypeAnnotationInUniverse(*parameter.Type, localUniverse)
			if err != nil {
				return nil, err
			}
			if !variableTypeContainedIn(expectedFunction.Arguments[i].VariableType, paramType) {
				return nil, type_error.PtrOnNodef(expression.Node, "in parameter position %d expected type %s but you have annotated %s", i, printableName(expectedFunction.Arguments[i].VariableType), printableName(paramType))
			}
		}
		localUniverse, err = binding.CopyAddingVariable(localUniverse, parameter.Name, expectedFunction.Arguments[i].VariableType)
		if err != nil {
			return nil, err
		}
	}

	if expression.ReturnType != nil {
		returnType, err := validateTypeAnnotationInUniverse(*expression.ReturnType, localUniverse)
		if err != nil {
			return nil, err
		}
		if !variableTypeContainedIn(returnType, expectedFunction.ReturnType) {
			return nil, type_error.PtrOnNodef(expression.Node, "in return type expected type %s but you have annotated %s", printableName(expectedFunction.ReturnType), printableName(returnType))
		}
	}

	astBlock, err := expectTypeOfBlock(expectedFunction.ReturnType, expression.Node, expression.Block, localUniverse)
	if err != nil {
		return nil, err
	}

	return &ast.Function{
		VariableType: expectedFunction,
		Block:        astBlock,
	}, nil
}

func expectTypeOfBlock(expectedType types.VariableType, node parser.Node, block []parser.ExpressionBox, universe binding.Universe) ([]ast.Expression, *type_error.TypecheckError) {
	result := []ast.Expression{}

	if len(block) == 0 {
		if variableTypeContainedIn(&types.Void{}, expectedType) {
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
			expectedTypeOfExpressionBox, err = typeOfExpressionBox(expressionBox, localUniverse)
			if err != nil {
				return nil, err
			}
		}
		astExp, err := expectTypeOfExpressionBox(expectedTypeOfExpressionBox, expressionBox, localUniverse)
		if err != nil {
			return nil, err
		}
		result = append(result, astExp)
		astDec, isDec := astExp.(ast.Declaration)
		if isDec {
			localUniverse, err = binding.CopyAddingVariable(localUniverse, parser.Name{
				String: astDec.Name,
			}, ast.VariableTypeOfExpression(astDec.Expression))
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func resolveFunctionGenerics(node parser.Node, function *types.Function, genericsPassed []parser.TypeAnnotation, universe binding.Universe) (*types.Function, []types.StructFieldVariableType, *type_error.TypecheckError) {
	generics := []types.StructFieldVariableType{}
	genericsMap := map[string]types.StructFieldVariableType{}
	if len(genericsPassed) != len(function.Generics) {
		return nil, nil, type_error.PtrOnNodef(node, "expected %d generics but got %d", len(function.Generics), len(genericsPassed))
	}
	for i, generic := range genericsPassed {
		varType, err := validateTypeAnnotationInUniverse(generic, universe)
		if err != nil {
			return nil, nil, err
		}
		structVarType, ok := types.StructFieldVariableTypeFromVariableType(varType)
		if !ok {
			return nil, nil, type_error.PtrOnNodef(generic.Node, "invalid generic type %s", printableName(varType))
		}
		genericsMap[function.Generics[i]] = structVarType
		generics = append(generics, structVarType)
	}

	arguments := []types.FunctionArgument{}
	for _, argument := range function.Arguments {
		arguments = append(arguments, argument)
	}
	for i := 0; i < len(arguments); i++ {
		for genericName, resolveTo := range genericsMap {
			newVarType, err := resolveGeneric(arguments[i].VariableType, genericName, resolveTo)
			if err != nil {
				return nil, nil, err
			}
			arguments[i].VariableType = newVarType
		}
	}

	returnType := function.ReturnType
	for genericName, resolveTo := range genericsMap {
		newVarType, err := resolveGeneric(returnType, genericName, resolveTo)
		if err != nil {
			return nil, nil, err
		}
		returnType = newVarType
	}

	return &types.Function{
		Generics:   nil,
		Arguments:  arguments,
		ReturnType: returnType,
	}, generics, nil
}

func expectTypeOfReferenceOrInvocation(expectedType types.VariableType, expression parser.ReferenceOrInvocation, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
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
		overFunction, generics, err := resolveFunctionGenerics(expression.Arguments.Node, overFunction, expression.Arguments.Generics, universe)
		if err != nil {
			return nil, err
		}
		for i, argument := range expression.Arguments.Arguments {
			expectedArgType := overFunction.Arguments[i].VariableType
			astArg, err := expectTypeOfExpressionBox(expectedArgType, argument, universe)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, astArg)
		}

		astExp := ast.Invocation{
			VariableType: overFunction.ReturnType,
			Over: ast.Reference{
				VariableType: overFunction,
				Name:         expression.Var.String,
			},
			Generics:  generics,
			Arguments: arguments,
		}

		return astExp, nil
	} else {
		if !variableTypeContainedIn(overType, expectedType) {
			return nil, type_error.PtrOnNodef(expression.Var.Node, "expected type %s but found %s", printableName(expectedType), printableName(overType))
		}

		astExp := ast.Reference{
			VariableType: overType,
			Name:         expression.Var.String,
		}

		return astExp, nil
	}
}

func expectTypeOfLiteral(expectedType types.VariableType, expression parser.LiteralExpression, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	varType, err := typeOfExpression(expression, universe)
	if err != nil {
		return nil, err
	}
	if !variableTypeContainedIn(varType, expectedType) {
		return nil, type_error.PtrOnNodef(expression.Node, "expected type %s but found %s", printableName(expectedType), printableName(varType))
	}
	return ast.Literal{
		VariableType: varType.(*types.BasicType),
		Literal:      expression.Literal,
	}, nil
}

func expectTypeOfModule(expectedType types.VariableType, expression parser.Module, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	_, resolutionErr := binding.GetTypeByTypeName(universe, expression.Implementing.String, []types.StructFieldVariableType{})
	if resolutionErr != nil {
		return nil, TypecheckErrorFromResolutionError(expression.Node, resolutionErr)
	}
	expectedInterface, ok := expectedType.(*types.Interface)
	if !ok {
		return nil, type_error.PtrOnNodef(expression.Node, "Expected %s but got %s", printableName(expectedType), expression.Implementing.String)
	}

	declarations := []parser.Declaration{}
	for _, moduleDeclaration := range expression.Declarations {
		if moduleDeclaration.Public {
			if expectedInterface.Variables[moduleDeclaration.Name.String] == nil {
				return nil, type_error.PtrOnNodef(expression.Node, "variable %s should not be public", moduleDeclaration.Name.String)
			}
		} else {
			if expectedInterface.Variables[moduleDeclaration.Name.String] != nil {
				return nil, type_error.PtrOnNodef(expression.Node, "variable %s should be public", moduleDeclaration.Name.String)
			}
		}
		declarations = append(declarations, parser.Declaration{
			Name: moduleDeclaration.Name,
			ExpressionBox: parser.ExpressionBox{
				Expression: moduleDeclaration.Expression,
			},
		})
	}

	astExpMap, err := TypecheckDeclarations(&expectedInterface.Variables, expression.Node, declarations, universe)
	if err != nil {
		return nil, err
	}
	return ast.Module{
		Implements: expectedInterface,
		Variables:  astExpMap,
	}, nil
}
