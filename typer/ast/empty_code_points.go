package ast

func EmptyCodePoints(program Program) Program {
	declarations := map[Ref]Expression{}
	for ref, expression := range program.Declarations {
		declarations[ref] = emptyCodePointsInExpression(expression)
	}
	return Program{
		Declarations:    declarations,
		TypeAliases:     program.TypeAliases,
		StructFunctions: program.StructFunctions,
		NativeFunctions: program.NativeFunctions,
		FieldsByType:    program.FieldsByType,
	}
}

func emptyCodePointsInExpressions(expressions []Expression) []Expression {
	result := []Expression{}
	for _, expression := range expressions {
		result = append(result, emptyCodePointsInExpression(expression))
	}
	return result
}

func emptyCodePointsInExpression(expression Expression) Expression {
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
		return Literal{
			CodePoint:    CodePoint{},
			VariableType: caseLiteral.VariableType,
			Literal:      caseLiteral.Literal,
		}
	} else if caseReference != nil {
		return Reference{
			CodePoint:    CodePoint{},
			VariableType: caseReference.VariableType,
			PackageName:  caseReference.PackageName,
			Name:         caseReference.Name,
		}
	} else if caseAccess != nil {
		return Access{
			CodePoint:    CodePoint{},
			VariableType: caseAccess.VariableType,
			Over:         emptyCodePointsInExpression(caseAccess.Over),
			Access:       caseAccess.Access,
		}
	} else if caseInvocation != nil {
		return Invocation{
			CodePoint:    CodePoint{},
			VariableType: caseInvocation.VariableType,
			Over:         emptyCodePointsInExpression(caseInvocation.Over),
			Generics:     caseInvocation.Generics,
			Arguments:    emptyCodePointsInExpressions(caseInvocation.Arguments),
		}
	} else if caseFunction != nil {
		return &Function{
			CodePoint:    CodePoint{},
			VariableType: caseFunction.VariableType,
			Block:        emptyCodePointsInExpressions(caseFunction.Block),
		}
	} else if caseDeclaration != nil {
		return Declaration{
			CodePoint:  CodePoint{},
			Name:       caseDeclaration.Name,
			Expression: emptyCodePointsInExpression(caseDeclaration.Expression),
		}
	} else if caseIf != nil {
		return If{
			CodePoint:    CodePoint{},
			VariableType: caseIf.VariableType,
			Condition:    emptyCodePointsInExpression(caseIf.Condition),
			ThenBlock:    emptyCodePointsInExpressions(caseIf.ThenBlock),
			ElseBlock:    emptyCodePointsInExpressions(caseIf.ElseBlock),
		}
	} else if caseList != nil {
		return List{
			CodePoint:             CodePoint{},
			ContainedVariableType: caseList.ContainedVariableType,
			Arguments:             emptyCodePointsInExpressions(caseList.Arguments),
		}
	} else if caseWhen != nil {
		cases := []WhenCase{}
		for _, whenCase := range caseWhen.Cases {
			cases = append(cases, WhenCase{
				Name:         whenCase.Name,
				VariableType: whenCase.VariableType,
				Block:        emptyCodePointsInExpressions(whenCase.Block),
			})
		}

		return When{
			CodePoint:     CodePoint{},
			VariableType:  caseWhen.VariableType,
			Over:          emptyCodePointsInExpression(caseWhen.Over),
			Cases:         cases,
			OtherCase:     emptyCodePointsInExpressions(caseWhen.OtherCase),
			OtherCaseName: caseWhen.OtherCaseName,
		}
	} else {
		panic("ExpressionCases")
	}
}
