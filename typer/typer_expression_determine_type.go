package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func determineTypeOfExpressionBox(validateFunctionBlock bool, expressionBox parser.ExpressionBox, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	expression, accessOrInvocations := parser.ExpressionBoxFields(expressionBox)
	universe, astExp, err := determineTypeOfExpression(validateFunctionBlock, expression, universe)
	if err != nil {
		return nil, nil, err
	}
	if accessOrInvocations == nil || len(accessOrInvocations) == 0 {
		return universe, astExp, nil
	}

	invocationOverAstExp := astExp
	accessChain := []ast.AccessAndMaybeInvocation{}

	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid, caseArray := ast.VariableTypeOfExpression(astExp).VariableTypeCases()
	_ = caseTypeArgument
	_ = caseFunction
	_ = caseBasicType
	_ = caseVoid
	_ = caseArray
	currentUniverse := universe
	if caseStruct != nil {
		currentUniverse, err = binding.NewFromStructVariables(parser.GetExpressionNode(expression), caseStruct.Fields, universe)
		if err != nil {
			return nil, nil, err
		}
	} else if caseInterface != nil {
		currentUniverse, err = binding.NewFromInterfaceVariables(parser.GetExpressionNode(expression), caseInterface.Variables, universe)
		if err != nil {
			return nil, nil, err
		}
	} else {
		return nil, nil, type_error.PtrOnNodef(parser.GetExpressionNode(expression), "should be an interface or struct to continue chained calls but found %s", printableName(ast.VariableTypeOfExpression(astExp)))
	}
	for i, accessOrInvocation := range accessOrInvocations {
		varType, ok := binding.GetTypeByVariableName(currentUniverse, accessOrInvocation.VarName.String)
		if !ok {
			return nil, nil, type_error.PtrOnNodef(accessOrInvocation.VarName.Node, "not found in scope: "+accessOrInvocation.VarName.String)
		}

		if accessOrInvocation.Arguments == nil {
			accessChain = append(accessChain, ast.AccessAndMaybeInvocation{
				VariableType:  varType,
				Access:        accessOrInvocation.VarName.String,
				ArgumentsList: nil,
			})
		} else {
			argumentsList := *accessOrInvocation.Arguments
			caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid, caseArray := varType.VariableTypeCases()
			_ = caseTypeArgument
			_ = caseStruct
			_ = caseInterface
			_ = caseBasicType
			_ = caseVoid
			_ = caseArray

			if caseFunction == nil {
				return nil, nil, type_error.PtrOnNodef(accessOrInvocation.VarName.Node, "%s should be a function for invocation but found %s", accessOrInvocation.VarName.String, printableName(varType))
			}
			returnType, astArgumentsList, err := determineTypeReturnedFromFunctionInvocation(validateFunctionBlock, argumentsList, *caseFunction, universe)
			if err != nil {
				return nil, nil, err
			}
			varType = returnType
			accessChain = append(accessChain, ast.AccessAndMaybeInvocation{
				VariableType:  varType,
				Access:        accessOrInvocation.VarName.String,
				ArgumentsList: astArgumentsList,
			})
		}

		if i < len(accessOrInvocations)-1 {
			caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid, caseArray := varType.VariableTypeCases()
			_ = caseTypeArgument
			_ = caseStruct
			_ = caseFunction
			_ = caseBasicType
			_ = caseVoid
			_ = caseArray
			if caseInterface == nil {
				return nil, nil, type_error.PtrOnNodef(accessOrInvocation.VarName.Node, "%s should be an interface to continue chained calls but found %s", accessOrInvocation.VarName.String, printableName(varType))
			}
			currentUniverse, _ = binding.NewFromInterfaceVariables(accessOrInvocation.VarName.Node, caseInterface.Variables, currentUniverse)
		} else {
			return universe, ast.WithAccessAndMaybeInvocation{
				VariableType: varType,
				Over:         invocationOverAstExp,
				AccessChain:  accessChain,
			}, nil
		}
	}

	panic("should have returned")
}

func determineTypeOfExpression(validateFunctionBlock bool, expression parser.Expression, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	resultUniverse := universe
	var resultExpression ast.Expression
	var err *type_error.TypecheckError
	parser.ExpressionExhaustiveSwitch(
		expression,
		func(expression parser.Module) {
			resultUniverse, resultExpression, err = determineTypeOfModule(validateFunctionBlock, expression, universe)
		},
		func(expression parser.LiteralExpression) {
			resultExpression = determineTypeOfLiteral(expression.Literal)
		},
		func(expression parser.ReferenceOrInvocation) {
			resultExpression, err = determineTypeOfReferenceOrInvocation(validateFunctionBlock, expression, universe)
		},
		func(expression parser.Lambda) {
			resultUniverse, resultExpression, err = determineTypeOfLambda(validateFunctionBlock, expression, universe)
		},
		func(expression parser.Declaration) {
			resultUniverse, resultExpression, err = determineTypeOfDeclaration(validateFunctionBlock, expression, universe)
		},
		func(expression parser.If) {
			resultUniverse, resultExpression, err = determineTypeOfIf(validateFunctionBlock, expression, universe)
		},
		func(expression parser.Array) {
			resultExpression, err = determineTypeOfArray(validateFunctionBlock, expression, universe)
		},
	)
	return resultUniverse, resultExpression, err
}

func determineTypeOfModule(validateFunctionBlock bool, module parser.Module, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	implementing, declarations := parser.ModuleFields(module)
	implementingVarType, ok := binding.GetTypeByTypeName(universe, implementing.String)
	if !ok {
		return nil, nil, type_error.PtrOnNodef(implementing.Node, "No interface %s found", implementing.String)
	}
	_, _, caseInterface, _, _, _, _ := implementingVarType.VariableTypeCases()
	if caseInterface == nil {
		return nil, nil, type_error.PtrOnNodef(implementing.Node, "Expected %s to be an interface but it's %s", implementing.String, printableName(implementingVarType))
	}
	for interfVarName, _ := range caseInterface.Variables {
		found := false
		for _, declaration := range declarations {
			if declaration.Name.String == interfVarName {
				found = true
				break
			}
		}
		if !found {
			return nil, nil, type_error.PtrOnNodef(implementing.Node, "interface %s has variable '%s' that needs to be implemented", implementing.String, interfVarName)
		}
	}
	astModule := ast.Module{
		Implements: caseInterface,
		Variables:  map[string]ast.Expression{},
	}
	typeOfInterfaceVarWithName := map[string]types.VariableType{}
	for interfVarName, interfVarType := range caseInterface.Variables {
		typeOfInterfaceVarWithName[interfVarName] = interfVarType
	}
	localUniverse := universe
	for _, declaration := range declarations {
		typeOfInterfaceVarWithSameName := typeOfInterfaceVarWithName[declaration.Name.String]
		if typeOfInterfaceVarWithSameName != nil && !declaration.Public {
			return nil, nil, type_error.PtrOnNodef(declaration.Name.Node, "variable %s should be public", declaration.Name.String)
		}
		if typeOfInterfaceVarWithSameName == nil && declaration.Public {
			return nil, nil, type_error.PtrOnNodef(declaration.Name.Node, "variable %s should not be public", declaration.Name.String)
		}
		var exp ast.Expression
		var err *type_error.TypecheckError
		if typeOfInterfaceVarWithSameName != nil {
			_, exp, err = expectTypeOfExpression(false, declaration.Expression, typeOfInterfaceVarWithSameName, localUniverse)
		} else {
			_, exp, err = determineTypeOfExpression(false, declaration.Expression, localUniverse)
		}
		if err != nil {
			return nil, nil, err
		}
		astModule.Variables[declaration.Name.String] = exp
		localUniverse, err = binding.CopyAddingVariable(localUniverse, declaration.Name, ast.VariableTypeOfExpression(exp))
		if err != nil {
			return nil, nil, err
		}
	}
	if validateFunctionBlock {
		for _, declaration := range declarations {
			_, ok := declaration.Expression.(parser.Lambda)
			if !ok {
				continue
			}
			typeOfInterfaceVarWithSameName := typeOfInterfaceVarWithName[declaration.Name.String]
			var exp ast.Expression
			var err *type_error.TypecheckError
			if typeOfInterfaceVarWithSameName != nil {
				_, exp, err = expectTypeOfExpression(true, declaration.Expression, typeOfInterfaceVarWithSameName, localUniverse)
			} else {
				_, exp, err = determineTypeOfExpression(true, declaration.Expression, localUniverse)
			}
			if err != nil {
				return nil, nil, err
			}
			astModule.Variables[declaration.Name.String] = exp
		}
	}
	return universe, astModule, nil
}

func determineTypeOfDeclaration(validateFunctionBlock bool, expression parser.Declaration, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	fieldName, fieldExpression := parser.DeclarationFields(expression)
	_, programExp, err := determineTypeOfExpressionBox(validateFunctionBlock, fieldExpression, universe)
	if err != nil {
		return nil, nil, err
	}
	varType := ast.VariableTypeOfExpression(programExp)
	updatedUniverse, err := binding.CopyAddingVariable(universe, fieldName, varType)
	if err != nil {
		return nil, nil, err
	}
	declarationProgramExp := ast.Declaration{
		VariableType: &standard_library.Void,
		Name:         fieldName.String,
		Expression:   programExp,
	}
	return updatedUniverse, declarationProgramExp, nil
}

func determineTypeOfLambda(validateFunctionBlock bool, expression parser.Lambda, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	localUniverse := universe
	generics, parameters, annotatedReturnType, block := parser.LambdaFields(expression)
	_ = block
	genericsStrings := []string{}
	for _, generic := range generics {
		genericsStrings = append(genericsStrings, generic.String)
	}
	if generics == nil {
		genericsStrings = nil
	}
	function := &types.Function{
		Generics:   genericsStrings,
		Arguments:  []types.FunctionArgument{},
		ReturnType: nil,
	}
	for _, generic := range generics {
		u, err := binding.CopyAddingType(localUniverse, generic, &types.TypeArgument{Name: generic.String})
		if err != nil {
			return nil, nil, err
		}
		localUniverse = u
	}
	for _, parameter := range parameters {
		if parameter.Type == nil {
			return nil, nil, type_error.PtrOnNodef(parameter.Name.Node, "parameter '%s' needs to be type annotated as the variable is not public", parameter.Name.String)
		}

		varType, err := validateTypeAnnotationInUniverse(*parameter.Type, localUniverse)
		if err != nil {
			return nil, nil, err
		}
		function.Arguments = append(function.Arguments, types.FunctionArgument{
			Name:         parameter.Name.String,
			VariableType: varType,
		})
	}
	if annotatedReturnType == nil {
		return nil, nil, type_error.PtrOnNodef(expression.Node, "return type needs to be type annotated as the variable is not public")
	}
	varType, err := validateTypeAnnotationInUniverse(*annotatedReturnType, localUniverse)
	if err != nil {
		return nil, nil, err
	}
	function.ReturnType = varType

	functionArgumentNames := []parser.Name{}
	for _, parameter := range expression.Parameters {
		functionArgumentNames = append(functionArgumentNames, parameter.Name)
	}
	functionArgumentVariableTypes := []types.VariableType{}
	for _, argument := range function.Arguments {
		functionArgumentVariableTypes = append(functionArgumentVariableTypes, argument.VariableType)
	}
	localUniverse, err = binding.CopyAddingFunctionArguments(localUniverse, functionArgumentNames, functionArgumentVariableTypes)
	if err != nil {
		return nil, nil, err
	}

	functionBlock := []ast.Expression{}
	if validateFunctionBlock {
		_, hasVoidReturnType := function.ReturnType.(*types.Void)
		if !hasVoidReturnType && len(block) == 0 {
			return nil, nil, type_error.PtrOnNodef(expression.Node, "Function has return type of %s but has empty body", printableName(function.ReturnType))
		}
		for i, blockExp := range block {
			if i < len(block)-1 {
				u, astExp, err := determineTypeOfExpressionBox(true, blockExp, localUniverse)
				if err != nil {
					return nil, nil, err
				}
				functionBlock = append(functionBlock, astExp)
				localUniverse = u
			} else {
				_, astExp, err := expectTypeOfExpressionBox(true, blockExp, varType, localUniverse)
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
}

func determineTypeOfLiteral(literal parser.Literal) ast.Expression {
	var varType *types.BasicType
	parser.LiteralExhaustiveSwitch(
		literal,
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
	return ast.Literal{
		VariableType: varType,
		Literal:      literal,
	}
}

func determineTypeOfArray(validateFunctionBlock bool, array parser.Array, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	if len(array.Expressions) == 0 && array.Generic == nil {
		return nil, type_error.PtrOnNodef(array.Node, "empty array requires type annotation")
	}
	if array.Generic == nil {
		return nil, type_error.PtrOnNodef(array.Node, "array type inference not yet implemented")
	}
	varType, err := validateTypeAnnotationInUniverse(array.Generic, universe)
	if err != nil {
		return nil, err
	}
	variableType, ok := types.StructFieldVariableTypeFromVariableType(varType)
	if !ok {
		return nil, type_error.PtrOnNodef(array.Node, "invalid array type %s", printableName(varType))
	}
	arguments := []ast.Expression{}
	for _, expressionBox := range array.Expressions {
		_, astExp, err := expectTypeOfExpressionBox(validateFunctionBlock, expressionBox, varType, universe)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, astExp)
	}
	return ast.Array{
		ContainedVariableType: variableType,
		Arguments:             arguments,
	}, nil
}

func determineTypeOfIf(validateFunctionBlock bool, caseIf parser.If, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	u2, conditionProgramExp, err := expectTypeOfExpressionBox(validateFunctionBlock, caseIf.Condition, &standard_library.BasicTypeBoolean, universe)
	if err != nil {
		return nil, nil, err
	}
	universe = u2

	varTypeOfBlock := func(expressionBoxess []parser.ExpressionBox, universe binding.Universe) (binding.Universe, []ast.Expression, types.VariableType, *type_error.TypecheckError) {
		if len(expressionBoxess) == 0 {
			return universe, []ast.Expression{}, &standard_library.Void, nil
		}
		localUniverse := universe
		programExpressions := []ast.Expression{}
		for i, exp := range expressionBoxess {
			u, programExp, err := determineTypeOfExpressionBox(validateFunctionBlock, exp, localUniverse)
			if err != nil {
				return nil, nil, nil, err
			}
			localUniverse = u
			varType := ast.VariableTypeOfExpression(programExp)
			programExpressions = append(programExpressions, programExp)
			if i == len(expressionBoxess)-1 {
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
		u2, elseProgramExpressions, elseVarType, err := varTypeOfBlock(caseIf.ElseBlock, universe)
		if err != nil {
			return nil, nil, err
		}
		universe = u2
		if !variableTypeEq(thenVarType, elseVarType) {
			return nil, nil, type_error.PtrOnNodef(caseIf.Node, "if and else blocks should yield the same type, but if is %s and then is %s", printableName(thenVarType), printableName(elseVarType))
		}
		return universe, ast.If{
			VariableType: thenVarType,
			Condition:    conditionProgramExp,
			ThenBlock:    thenProgramExpressions,
			ElseBlock:    elseProgramExpressions,
		}, nil
	} else {
		return universe, ast.If{
			VariableType: &standard_library.Void,
			Condition:    conditionProgramExp,
			ThenBlock:    thenProgramExpressions,
			ElseBlock:    []ast.Expression{},
		}, nil
	}
}

func determineTypeOfReferenceOrInvocation(validateFunctionBlock bool, referenceOrInvocation parser.ReferenceOrInvocation, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	refName, argumentsPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)

	varType, ok := binding.GetTypeByVariableName(universe, refName.String)
	if !ok {
		return nil, type_error.PtrOnNodef(refName.Node, "not found in scope: "+refName.String)
	}

	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid, caseArray := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName.String,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrOnNodef(refName.Node, "%s should be a function for invocation but found %s", refName.String, printableName(varType))
		}
	} else if caseStruct != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName.String,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrOnNodef(refName.Node, "%s should be a function for invocation but found %s", refName.String, printableName(varType))
		}
	} else if caseInterface != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName.String,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrOnNodef(refName.Node, "%s should be a function for invocation but found %s", refName.String, printableName(varType))
		}
	} else if caseFunction != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName.String,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			argumentsList := *argumentsPtr
			returnType, astArgumentsList, err := determineTypeReturnedFromFunctionInvocation(validateFunctionBlock, argumentsList, *caseFunction, universe)
			if err != nil {
				return nil, err
			}
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  returnType,
				Name:          refName.String,
				ArgumentsList: astArgumentsList,
			}
			return programExp, nil
		}
	} else if caseBasicType != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName.String,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrOnNodef(refName.Node, "%s should be a function for invocation but found %s", refName.String, printableName(varType))
		}
	} else if caseVoid != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName.String,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrOnNodef(refName.Node, "%s should be a function for invocation but found %s", refName.String, printableName(varType))
		}
	} else if caseArray != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName.String,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrOnNodef(refName.Node, "%s should be a function for invocation but found %s", refName.String, printableName(varType))
		}
	} else {
		panic(fmt.Errorf("code on %v", varType))
	}
}

func determineTypeReturnedFromFunctionInvocation(validateFunctionBlock bool, argumentsList parser.ArgumentsList, caseFunction types.Function, universe binding.Universe) (types.VariableType, *ast.ArgumentsList, *type_error.TypecheckError) {
	if len(argumentsList.Arguments) != len(caseFunction.Arguments) {
		return nil, nil, type_error.PtrOnNodef(argumentsList.Node, "Expected %d arguments but got %d", len(caseFunction.Arguments), len(argumentsList.Arguments))
	}
	if len(argumentsList.Generics) != len(caseFunction.Generics) {
		return nil, nil, type_error.PtrOnNodef(argumentsList.Node, "Expected %d generics annotated but got %d", len(caseFunction.Generics), len(argumentsList.Generics))
	}
	genericArgumentTypes := []types.StructFieldVariableType{}
	for _, generic := range argumentsList.Generics {
		genericVarType, err := validateTypeAnnotationInUniverse(generic, universe)
		if err != nil {
			return nil, nil, type_error.PtrOnNodef(generic.Node, "not found annotated generic type %s", generic.TypeName)
		}
		genericType, ok := types.StructFieldVariableTypeFromVariableType(genericVarType)
		if !ok {
			return nil, nil, type_error.PtrOnNodef(generic.Node, "not a valid annotated generic %s", generic.TypeName)
		}
		genericArgumentTypes = append(genericArgumentTypes, genericType)
	}
	argumentProgramExpressions := []ast.Expression{}
	for i2, argument := range argumentsList.Arguments {
		expectedType := caseFunction.Arguments[i2].VariableType
		expectedTypeArg, isGeneric := expectedType.(*types.TypeArgument)
		if isGeneric {
			caseFunctionGenericIndex := -1
			for index, functionGeneric := range caseFunction.Generics {
				if functionGeneric == expectedTypeArg.Name {
					caseFunctionGenericIndex = index
					break
				}
			}
			if caseFunctionGenericIndex == -1 {
				return nil, nil, type_error.PtrOnNodef(parser.GetExpressionNode(argument.Expression), "unexpected error not found generic %s", expectedTypeArg.Name)
			}
			expectedType = types.VariableTypeFromStructFieldVariableType(genericArgumentTypes[caseFunctionGenericIndex])
		}
		_, programExp, err := expectTypeOfExpressionBox(validateFunctionBlock, argument, expectedType, universe)
		if err != nil {
			return nil, nil, err
		}
		argumentProgramExpressions = append(argumentProgramExpressions, programExp)
	}
	returnType := caseFunction.ReturnType
	returnTypeArg, isGeneric := returnType.(*types.TypeArgument)
	if isGeneric {
		caseFunctionGenericIndex := -1
		for index, functionGeneric := range caseFunction.Generics {
			if functionGeneric == returnTypeArg.Name {
				caseFunctionGenericIndex = index
				break
			}
		}
		if caseFunctionGenericIndex == -1 {
			return nil, nil, type_error.PtrOnNodef(argumentsList.Node, "unexpected error not found return generic %s", returnTypeArg.Name)
		}
		invocationGeneric := argumentsList.Generics[caseFunctionGenericIndex]
		newReturnType, err := validateTypeAnnotationInUniverse(invocationGeneric, universe)
		if err != nil {
			return nil, nil, type_error.PtrOnNodef(invocationGeneric.Node, "not found return generic type %s", invocationGeneric.TypeName)
		}
		returnType = newReturnType
	}
	returnTypeStruct, isStruct := returnType.(*types.Struct)
	if isStruct && len(caseFunction.Generics) > 0 {
		structToReturn := &types.Struct{
			Package: returnTypeStruct.Package,
			Name:    returnTypeStruct.Name,
			Fields:  map[string]types.StructFieldVariableType{},
		}
		for i, generic := range argumentsList.Generics {
			genericVarType, err := validateTypeAnnotationInUniverse(generic, universe)
			if err != nil {
				return nil, nil, type_error.PtrOnNodef(generic.Node, "not found annotated generic type %s", generic.TypeName)
			}
			structFieldVarType, ok := types.StructFieldVariableTypeFromVariableType(genericVarType)
			if !ok {
				return nil, nil, type_error.PtrOnNodef(generic.Node, "not a valid annotated generic type %s", generic.TypeName)
			}
			for fieldName, fieldVariableType := range returnTypeStruct.Fields {
				resolvedVarType, err := resolveGeneric(fieldVariableType, caseFunction.Generics[i], structFieldVarType)
				if err != nil {
					return nil, nil, err
				}
				structToReturn.Fields[fieldName] = resolvedVarType
			}

		}
		returnType = structToReturn
	}

	return returnType, &ast.ArgumentsList{
		Generics:  genericArgumentTypes,
		Arguments: argumentProgramExpressions,
	}, nil
}

func resolveGeneric(over types.StructFieldVariableType, genericName string, resolveWith types.StructFieldVariableType) (types.StructFieldVariableType, *type_error.TypecheckError) {
	caseTypeArgument, caseStruct, caseBasicType, caseVoid, caseArray := over.StructFieldVariableTypeCases()
	if caseTypeArgument != nil {
		if caseTypeArgument.Name == genericName {
			return resolveWith, nil
		}
		return caseTypeArgument, nil
	} else if caseStruct != nil {
		panic("todo resolveGeneric caseStruct")
	} else if caseBasicType != nil {
		return caseBasicType, nil
	} else if caseVoid != nil {
		return caseVoid, nil
	} else if caseArray != nil {
		panic("todo resolveGeneric caseArray")
	} else {
		panic(fmt.Errorf("cases on %v", over))
	}
}
