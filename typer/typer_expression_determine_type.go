package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
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

	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := ast.VariableTypeOfExpression(astExp).VariableTypeCases()
	_ = caseTypeArgument
	_ = caseFunction
	_ = caseBasicType
	_ = caseVoid
	currentUniverse := universe
	if caseStruct != nil {
		structVariables, err := binding.GetGlobalStructVariables(universe, *caseStruct)
		if err != nil {
			return nil, nil, err
		}
		currentUniverse = binding.NewFromStructVariables(structVariables, universe)
	} else if caseInterface != nil {
		interfaceVariables, err := binding.GetGlobalInterfaceVariables(universe, *caseInterface)
		if err != nil {
			return nil, nil, err
		}
		currentUniverse = binding.NewFromInterfaceVariables(interfaceVariables, universe)
	} else {
		return nil, nil, type_error.PtrTypeCheckErrorf("should be an interface or struct to continue chained calls but found %s", printableName(ast.VariableTypeOfExpression(astExp)))
	}
	for i, accessOrInvocation := range accessOrInvocations {
		varType, ok := binding.GetTypeByVariableName(currentUniverse, accessOrInvocation.VarName)
		if !ok {
			return nil, nil, &type_error.TypecheckError{Message: "not found in scope: " + accessOrInvocation.VarName}
		}

		if accessOrInvocation.Arguments == nil {
			accessChain = append(accessChain, ast.AccessAndMaybeInvocation{
				VariableType:  varType,
				Access:        accessOrInvocation.VarName,
				ArgumentsList: nil,
			})
		} else {
			argumentsList := *accessOrInvocation.Arguments
			caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.VariableTypeCases()
			_ = caseTypeArgument
			_ = caseStruct
			_ = caseInterface
			_ = caseBasicType
			_ = caseVoid

			if caseFunction == nil {
				return nil, nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", accessOrInvocation.VarName, printableName(varType))
			}
			returnType, astArgumentsList, err := determineTypeReturnedFromFunctionInvocation(validateFunctionBlock, argumentsList, *caseFunction, universe)
			if err != nil {
				return nil, nil, err
			}
			varType = returnType
			accessChain = append(accessChain, ast.AccessAndMaybeInvocation{
				VariableType:  varType,
				Access:        accessOrInvocation.VarName,
				ArgumentsList: astArgumentsList,
			})
		}

		if i < len(accessOrInvocations)-1 {
			caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.VariableTypeCases()
			_ = caseTypeArgument
			_ = caseStruct
			_ = caseFunction
			_ = caseBasicType
			_ = caseVoid
			if caseInterface == nil {
				return nil, nil, type_error.PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", accessOrInvocation.VarName, printableName(varType))
			}
			interfaceVariables, err := binding.GetGlobalInterfaceVariables(currentUniverse, *caseInterface)
			if err != nil {
				return nil, nil, err
			}
			currentUniverse = binding.NewFromInterfaceVariables(interfaceVariables, currentUniverse)
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
	caseModule, caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseModule != nil {
		return determineTypeOfModule(validateFunctionBlock, *caseModule, universe)
	} else if caseLiteralExp != nil {
		return universe, determineTypeOfLiteral(caseLiteralExp.Literal), nil
	} else if caseReferenceOrInvocation != nil {
		varType, err := determineTypeOfReferenceOrInvocation(validateFunctionBlock, *caseReferenceOrInvocation, universe)
		return universe, varType, err
	} else if caseLambda != nil {
		return determineTypeOfLambda(validateFunctionBlock, *caseLambda, universe)
	} else if caseDeclaration != nil {
		return determineTypeOfDeclaration(validateFunctionBlock, *caseDeclaration, universe)
	} else if caseIf != nil {
		return determineTypeOfIf(validateFunctionBlock, *caseIf, universe)
	} else {
		panic(fmt.Errorf("code on %v", expression))
	}
}

func determineTypeOfModule(validateFunctionBlock bool, module parser.Module, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	implementing, declarations := parser.ModuleFields(module)
	implementingVarType, ok := binding.GetTypeByTypeName(universe, implementing)
	if !ok {
		return nil, nil, type_error.PtrTypeCheckErrorf("No interface %s found", implementing)
	}
	_, _, caseInterface, _, _, _ := implementingVarType.VariableTypeCases()
	if caseInterface == nil {
		return nil, nil, type_error.PtrTypeCheckErrorf("Expected %s to be an interface but it's %s", implementing, printableName(implementingVarType))
	}
	interf := *caseInterface
	interfaceVariables, err := binding.GetGlobalInterfaceVariables(universe, interf)
	if err != nil {
		return nil, nil, err
	}
	for interfVarName, _ := range interfaceVariables {
		found := false
		for _, declaration := range declarations {
			if declaration.Name == interfVarName {
				found = true
				break
			}
		}
		if !found {
			return nil, nil, type_error.PtrTypeCheckErrorf("interface %s has variable '%s' that needs to be implemented", implementing, interfVarName)
		}
	}
	astModule := ast.Module{
		Implements: interf,
		Variables:  map[string]ast.Expression{},
	}
	typeOfInterfaceVarWithName := map[string]types.VariableType{}
	for interfVarName, interfVarType := range interfaceVariables {
		typeOfInterfaceVarWithName[interfVarName] = interfVarType
	}
	localUniverse := universe
	for _, declaration := range declarations {
		typeOfInterfaceVarWithSameName := typeOfInterfaceVarWithName[declaration.Name]
		if typeOfInterfaceVarWithSameName != nil && !declaration.Public {
			return nil, nil, type_error.PtrTypeCheckErrorf("variable %s should be public", declaration.Name)
		}
		if typeOfInterfaceVarWithSameName == nil && declaration.Public {
			return nil, nil, type_error.PtrTypeCheckErrorf("variable %s should not be public", declaration.Name)
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
		astModule.Variables[declaration.Name] = exp
		localUniverse, err = binding.CopyAddingVariable(localUniverse, declaration.Name, ast.VariableTypeOfExpression(exp))
		if err != nil {
			return nil, nil, err
		}
	}
	if validateFunctionBlock {
		for _, declaration := range declarations {
			_, _, _, caseLambda, _, _ := declaration.Expression.ExpressionCases()
			if caseLambda == nil {
				continue
			}
			typeOfInterfaceVarWithSameName := typeOfInterfaceVarWithName[declaration.Name]
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
			astModule.Variables[declaration.Name] = exp
		}
	}
	return universe, astModule, nil
}

func determineTypeOfDeclaration(validateFunctionBlock bool, expression parser.Declaration, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	fieldName, fieldExpression := parser.DeclarationFields(expression)
	updatedUniverse, programExp, err := determineTypeOfExpressionBox(validateFunctionBlock, fieldExpression, universe)
	if err != nil {
		return nil, nil, err
	}
	varType := ast.VariableTypeOfExpression(programExp)
	updatedUniverse, err = binding.CopyAddingVariable(updatedUniverse, fieldName, varType)
	if err != nil {
		return nil, nil, err
	}
	declarationProgramExp := ast.Declaration{
		VariableType: void,
		Name:         fieldName,
		Expression:   programExp,
	}
	return updatedUniverse, declarationProgramExp, nil
}

func determineTypeOfLambda(validateFunctionBlock bool, expression parser.Lambda, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	localUniverse := universe
	generics, parameters, annotatedReturnType, block := parser.LambdaFields(expression)
	_ = block
	function := types.Function{
		Generics:   generics,
		Arguments:  []types.FunctionArgument{},
		ReturnType: nil,
	}
	for _, generic := range generics {
		u, err := binding.CopyAddingType(localUniverse, generic, types.TypeArgument{Name: generic})
		if err != nil {
			return nil, nil, err
		}
		localUniverse = u
	}
	for _, parameter := range parameters {
		if parameter.Type == nil {
			return nil, nil, type_error.PtrTypeCheckErrorf("parameter '%s' needs to be type annotated as the variable is not public", parameter.Name)
		}

		varType, err := validateTypeAnnotationInUniverse(*parameter.Type, localUniverse)
		if err != nil {
			return nil, nil, err
		}
		function.Arguments = append(function.Arguments, types.FunctionArgument{
			Name:         parameter.Name,
			VariableType: varType,
		})
	}
	if annotatedReturnType == nil {
		return nil, nil, type_error.PtrTypeCheckErrorf("return type needs to be type annotated as the variable is not public")
	}
	varType, err := validateTypeAnnotationInUniverse(*annotatedReturnType, localUniverse)
	if err != nil {
		return nil, nil, err
	}
	function.ReturnType = varType

	localUniverse, err = binding.CopyAddingFunctionArguments(localUniverse, function.Arguments)
	if err != nil {
		return nil, nil, err
	}

	functionBlock := []ast.Expression{}
	if validateFunctionBlock {
		if function.ReturnType != void && len(block) == 0 {
			return nil, nil, type_error.PtrTypeCheckErrorf("Function has return type of %s but has empty body", printableName(function.ReturnType))
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
	return ast.Literal{
		VariableType: varType,
		Literal:      literal,
	}
}

func determineTypeOfIf(validateFunctionBlock bool, caseIf parser.If, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	u2, conditionProgramExp, err := expectTypeOfExpressionBox(validateFunctionBlock, caseIf.Condition, basicTypeBoolean, universe)
	if err != nil {
		return nil, nil, err
	}
	universe = u2

	varTypeOfBlock := func(expressionBoxess []parser.ExpressionBox, universe binding.Universe) (binding.Universe, []ast.Expression, types.VariableType, *type_error.TypecheckError) {
		if len(expressionBoxess) == 0 {
			return universe, []ast.Expression{}, void, nil
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
			return nil, nil, type_error.PtrTypeCheckErrorf("if and else blocks should yield the same type, but if is %s and then is %s", printableName(thenVarType), printableName(elseVarType))
		}
		return universe, ast.If{
			VariableType: thenVarType,
			Condition:    conditionProgramExp,
			ThenBlock:    thenProgramExpressions,
			ElseBlock:    elseProgramExpressions,
		}, nil
	} else {
		return universe, ast.If{
			VariableType: void,
			Condition:    conditionProgramExp,
			ThenBlock:    thenProgramExpressions,
			ElseBlock:    []ast.Expression{},
		}, nil
	}
}

func determineTypeOfReferenceOrInvocation(validateFunctionBlock bool, referenceOrInvocation parser.ReferenceOrInvocation, universe binding.Universe) (ast.Expression, *type_error.TypecheckError) {
	refName, argumentsPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)

	varType, ok := binding.GetTypeByVariableName(universe, refName)
	if !ok {
		return nil, &type_error.TypecheckError{Message: "not found in scope: " + refName}
	}

	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", refName, printableName(varType))
		}
	} else if caseStruct != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", refName, printableName(varType))
		}
	} else if caseInterface != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", refName, printableName(varType))
		}
	} else if caseFunction != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName,
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
				Name:          refName,
				ArgumentsList: astArgumentsList,
			}
			return programExp, nil
		}
	} else if caseBasicType != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", refName, printableName(varType))
		}
	} else if caseVoid != nil {
		if argumentsPtr == nil {
			programExp := ast.ReferenceAndMaybeInvocation{
				VariableType:  varType,
				Name:          refName,
				ArgumentsList: nil,
			}
			return programExp, nil
		} else {
			return nil, type_error.PtrTypeCheckErrorf("%s should be a function for invocation but found %s", refName, printableName(varType))
		}
	} else {
		panic(fmt.Errorf("code on %v", varType))
	}
}

func determineTypeReturnedFromFunctionInvocation(validateFunctionBlock bool, argumentsList parser.ArgumentsList, caseFunction types.Function, universe binding.Universe) (types.VariableType, *ast.ArgumentsList, *type_error.TypecheckError) {
	if len(argumentsList.Arguments) != len(caseFunction.Arguments) {
		return nil, nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(caseFunction.Arguments), len(argumentsList.Arguments))}
	}
	if len(argumentsList.Generics) != len(caseFunction.Generics) {
		return nil, nil, &type_error.TypecheckError{Message: fmt.Sprintf("Expected %d generics annotated but got %d", len(caseFunction.Generics), len(argumentsList.Generics))}
	}
	argumentProgramExpressions := []ast.Expression{}
	for i2, argument := range argumentsList.Arguments {
		expectedType := caseFunction.Arguments[i2].VariableType
		expectedTypeArg, isGeneric := expectedType.(types.TypeArgument)
		if isGeneric {
			caseFunctionGenericIndex := -1
			for index, functionGeneric := range caseFunction.Generics {
				if functionGeneric == expectedTypeArg.Name {
					caseFunctionGenericIndex = index
					break
				}
			}
			if caseFunctionGenericIndex == -1 {
				return nil, nil, &type_error.TypecheckError{Message: fmt.Sprintf("unexpected error not found generic %s", expectedTypeArg.Name)}
			}
			invocationGeneric := argumentsList.Generics[caseFunctionGenericIndex]
			newExpectedType, err := validateTypeAnnotationInUniverse(parser.SingleNameType{TypeName: invocationGeneric}, universe)
			if err != nil {
				return nil, nil, &type_error.TypecheckError{Message: fmt.Sprintf("not found annotated generic type %s", invocationGeneric)}
			}
			expectedType = newExpectedType
		}
		_, programExp, err := expectTypeOfExpressionBox(validateFunctionBlock, argument, expectedType, universe)
		if err != nil {
			return nil, nil, err
		}
		argumentProgramExpressions = append(argumentProgramExpressions, programExp)
	}
	returnType := caseFunction.ReturnType
	returnTypeArg, isGeneric := returnType.(types.TypeArgument)
	if isGeneric {
		caseFunctionGenericIndex := -1
		for index, functionGeneric := range caseFunction.Generics {
			if functionGeneric == returnTypeArg.Name {
				caseFunctionGenericIndex = index
				break
			}
		}
		if caseFunctionGenericIndex == -1 {
			return nil, nil, &type_error.TypecheckError{Message: fmt.Sprintf("unexpected error not found return generic %s", returnTypeArg.Name)}
		}
		invocationGeneric := argumentsList.Generics[caseFunctionGenericIndex]
		newReturnType, err := validateTypeAnnotationInUniverse(parser.SingleNameType{TypeName: invocationGeneric}, universe)
		if err != nil {
			return nil, nil, &type_error.TypecheckError{Message: fmt.Sprintf("not found return generic type %s", invocationGeneric)}
		}
		returnType = newReturnType
	}
	returnTypeStruct, isStruct := returnType.(types.Struct)
	if isStruct && len(caseFunction.Generics) > 0 {
		if returnTypeStruct.ResolvedTypeArguments == nil {
			returnTypeStruct.ResolvedTypeArguments = []types.ResolvedTypeArgument{}
		}
		for i, generic := range argumentsList.Generics {
			genericVarType, err := validateTypeAnnotationInUniverse(parser.SingleNameType{TypeName: generic}, universe)
			if err != nil {
				return nil, nil, &type_error.TypecheckError{Message: fmt.Sprintf("not found annotated generic type %s", generic)}
			}
			structVarType, ok := types.StructVariableTypeFromVariableType(genericVarType)
			if !ok {
				return nil, nil, &type_error.TypecheckError{Message: fmt.Sprintf("not a valid annotated generic type %s", generic)}
			}
			returnTypeStruct.ResolvedTypeArguments = append(returnTypeStruct.ResolvedTypeArguments, types.ResolvedTypeArgument{
				Name:               caseFunction.Generics[i],
				StructVariableType: structVarType,
			})
		}
		returnType = returnTypeStruct
	}

	return returnType, &ast.ArgumentsList{Arguments: argumentProgramExpressions}, nil
}
