package ir

import (
	"fmt"
	"strings"

	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
)

type context struct {
	topLevelDeclarations []Reference
	generics             map[string]ArrayIndexAccess
	nextGenericId        int
}

func newContext() context {
	return context{
		topLevelDeclarations: []Reference{},
		generics:             map[string]ArrayIndexAccess{},
		nextGenericId:        0,
	}
}

func (ctx context) isTopLevelDeclaration(ref Reference) bool {
	for _, topLevelDeclaration := range ctx.topLevelDeclarations {
		if topLevelDeclaration.Name == ref.Name {
			return true
		}
	}
	return false
}

func (ctx context) resolvePrintableNameOrReference(variableType types.VariableType) Expression {
	caseTypeArgument, caseList, caseKnownType, _, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		resolved, ok := ctx.generics[caseTypeArgument.Name]
		if !ok {
			panic(fmt.Errorf("generic %s not found in context", caseTypeArgument.Name))
		}
		return resolved
	} else if caseList != nil {
		return NonFunctionTypeName{
			Prefix: "List<",
			CommaSpaceSeparated: []TypeNamesOrConcat{
				TypeNamesOrConcat{
					Expressions: []Expression{
						ctx.resolvePrintableNameOrReference(caseList.Generic),
					},
				},
			},
			Suffix: ">",
		}
	} else if caseKnownType != nil {
		resolvedGenerics := []TypeNamesOrConcat{}
		for _, generic := range caseKnownType.Generics {
			resolvedGenerics = append(resolvedGenerics, TypeNamesOrConcat{
				Expressions: []Expression{
					ctx.resolvePrintableNameOrReference(generic),
				},
			})
		}
		prefix := types.PrintableNameKnownTypeWithoutGenerics(caseKnownType)
		suffix := ""
		if len(resolvedGenerics) > 0 {
			prefix += "<"
			suffix += ">"
		}
		return NonFunctionTypeName{
			Prefix:              prefix,
			CommaSpaceSeparated: resolvedGenerics,
			Suffix:              suffix,
		}
	} else if caseOr != nil {
		expressions := []Expression{}
		for _, element := range caseOr.Elements {
			expressions = append(expressions, ctx.resolvePrintableNameOrReference(element))
		}
		return TypeNamesOrConcat{
			Expressions: expressions,
		}
	} else {
		//TODO FIXME there's other cases with unresolved generics probably besides the ones done above
		return Literal{
			Value: parser.LiteralString{
				Value: `"` + types.PrintableNameSorted(variableType) + `"`,
			},
		}
	}
}

func (ctx context) CopyAddingFunctionGenerics(functionGenericNames []string) (context, string) {
	nextGenericId := ctx.nextGenericId + 1
	newGenerics := map[string]ArrayIndexAccess{}
	for name, generic := range ctx.generics {
		newGenerics[name] = generic
	}
	genericsVarName := fmt.Sprintf("generics%d", nextGenericId)
	for i, genericName := range functionGenericNames {
		newGenerics[genericName] = ArrayIndexAccess{
			Over: Reference{
				Name: genericsVarName,
			},
			Index: i,
		}
	}
	return context{
		topLevelDeclarations: ctx.topLevelDeclarations,
		generics:             newGenerics,
		nextGenericId:        nextGenericId,
	}, genericsVarName
}

func ToIR(program ast.Program) Program {
	ctx := newContext()
	for ref, _ := range program.Declarations {
		ctx.topLevelDeclarations = append(ctx.topLevelDeclarations, RefToIR(ref))
	}
	for ref, _ := range program.NativeFunctions {
		ctx.topLevelDeclarations = append(ctx.topLevelDeclarations, RefToIR(ref))
	}
	for ref, _ := range program.StdLibStructFunctions {
		ctx.topLevelDeclarations = append(ctx.topLevelDeclarations, RefToIR(ref))
	}

	declarations := map[Reference]TopLevelFunction{}
	for ref, expression := range program.Declarations {
		declarations[RefToIR(ref)] = topLevelDeclarationToIR(ctx, expression)
	}
	structFunctions := map[Reference]*StructFunction{}
	for ref, function := range program.StructFunctions {
		argumentNames := []string{}
		for _, argument := range function.Arguments {
			argumentNames = append(argumentNames, VariableName(nil, argument.Name))
		}
		_, _, caseKnownType, _, _ := function.ReturnType.VariableTypeCases()
		if caseKnownType == nil {
			panic("struct function should return struct")
		}
		genericsArgName := "generics"
		returnTypePrefix := types.PrintableNameKnownTypeWithoutGenerics(caseKnownType)
		returnTypeCommaSpaceSeparated := []TypeNamesOrConcat{}
		returnTypeSuffix := ""
		if len(caseKnownType.Generics) > 0 {
			returnTypePrefix += "<"
			for i, _ := range caseKnownType.Generics {
				returnTypeCommaSpaceSeparated = append(returnTypeCommaSpaceSeparated, TypeNamesOrConcat{
					Expressions: []Expression{
						ArrayIndexAccess{
							Over: Reference{
								Name: genericsArgName,
							},
							Index: i,
						},
					},
				})
			}
			returnTypeSuffix += ">"
		}
		updatedFunction := StructFunction{
			GenericsArgumentName: genericsArgName,
			OtherArgumentNames:   argumentNames,
			ReturnType: NonFunctionTypeName{
				Prefix:              returnTypePrefix,
				CommaSpaceSeparated: returnTypeCommaSpaceSeparated,
				Suffix:              returnTypeSuffix,
			},
		}
		structFunctions[RefToIR(ref)] = &updatedFunction
	}
	stdLibStructFunctions := map[Reference]*StructFunction{}
	for ref, function := range program.StdLibStructFunctions {
		argumentNames := []string{}
		for _, argument := range function.Arguments {
			argumentNames = append(argumentNames, VariableName(nil, argument.Name))
		}
		_, _, caseKnownType, _, _ := function.ReturnType.VariableTypeCases()
		if caseKnownType == nil {
			panic("struct function should return struct")
		}
		genericsArgName := "generics"
		returnTypePrefix := types.PrintableNameKnownTypeWithoutGenerics(caseKnownType)
		returnTypeCommaSpaceSeparated := []TypeNamesOrConcat{}
		returnTypeSuffix := ""
		if len(caseKnownType.Generics) > 0 {
			returnTypePrefix += "<"
			for i, _ := range caseKnownType.Generics {
				returnTypeCommaSpaceSeparated = append(returnTypeCommaSpaceSeparated, TypeNamesOrConcat{
					Expressions: []Expression{
						ArrayIndexAccess{
							Over: Reference{
								Name: genericsArgName,
							},
							Index: i,
						},
					},
				})
			}
			returnTypeSuffix += ">"
		}
		updatedFunction := StructFunction{
			GenericsArgumentName: genericsArgName,
			OtherArgumentNames:   argumentNames,
			ReturnType: NonFunctionTypeName{
				Prefix:              returnTypePrefix,
				CommaSpaceSeparated: returnTypeCommaSpaceSeparated,
				Suffix:              returnTypeSuffix,
			},
		}
		stdLibStructFunctions[RefToIR(ref)] = &updatedFunction
	}
	nativeFunctions := []NativeFunctionRef{}
	for ref, _ := range program.NativeFunctions {
		nativeFunctions = append(nativeFunctions, NativeFunctionRef{
			Package: ref.Package,
			Name:    ref.Name,
		})
	}

	return Program{
		Declarations:          declarations,
		StructFunctions:       structFunctions,
		NativeFunctions:       nativeFunctions,
		StdLibStructFunctions: stdLibStructFunctions,
	}
}

func VariableName(packageName *string, name string) string {
	if packageName == nil {
		return "_" + name
	} else if *packageName == "" {
		panic("package should not have empty name")
	} else {
		return strings.ReplaceAll(*packageName, ".", "_") + "__" + name
	}
}

func RefToIR(ref ast.Ref) Reference {
	return Reference{
		Name: VariableName(&ref.Package, ref.Name),
	}
}

func topLevelDeclarationToIR(ctx context, expression ast.Expression) TopLevelFunction {
	topLevelFunction := TopLevelFunction{
		Body: Return{
			ReturnExpression: irStatementToExpression(ctx, expressionToIR(ctx, expression)),
		},
	}
	return topLevelFunction
}

func irStatementToExpression(ctx context, statement Statement) Expression {
	switch s := statement.(type) {
	case LocalFunction:
		return s
	case Return:
		return s.ReturnExpression
	case Invocation:
		return s
	case InvocationOverTopLevelFunction:
		return s
	case FieldAccess:
		return s
	case ObjectInstantiation:
		return s
	case Literal:
		return s
	case Reference:
		return s
	case If:
		return LocalFunction{
			ParameterNames: []string{},
			Block: []Statement{
				s,
				Return{
					ReturnExpression: Literal{
						Value: parser.LiteralNull{},
					},
				},
			},
		}
	case IfForTypeMatching:
		return LocalFunction{
			ParameterNames: []string{},
			Block: []Statement{
				s,
				Return{
					ReturnExpression: Literal{
						Value: parser.LiteralNull{},
					},
				},
			},
		}
	case ArrayIndexAccess:
		return s
	case VariableDeclaration:
		return LocalFunction{
			ParameterNames: []string{},
			Block: []Statement{
				s,
				Return{
					ReturnExpression: Literal{
						Value: parser.LiteralNull{},
					},
				},
			},
		}
	default:
		panic(fmt.Sprintf("unsupported statement type for conversion to expression: %T", statement))
	}
}

func expressionToIR(ctx context, expression ast.Expression) Statement {
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
		literalType := ""
		parser.LiteralExhaustiveSwitch(
			caseLiteral.Literal,
			func(literal float64) {
				literalType = "Float"
			}, func(literal int) {
				literalType = "Int"
			}, func(literal string) {
				literalType = "String"
			}, func(literal bool) {
				literalType = "Boolean"
			}, func() {
				literalType = "Void"
			},
		)
		return ObjectInstantiation{
			Fields: map[string]Expression{
				"$type": Literal{
					Value: parser.LiteralString{
						Value: `"` + literalType + `"`,
					},
				},
				"value": Literal{
					Value: caseLiteral.Literal,
				},
			},
		}
	} else if caseReference != nil {
		reference := Reference{
			Name: VariableName(caseReference.PackageName, caseReference.Name),
		}

		isTopLevelReference := ctx.isTopLevelDeclaration(reference)

		if isTopLevelReference {
			return InvocationOverTopLevelFunction{
				Over: reference,
			}
		} else {
			return reference
		}
	} else if caseAccess != nil {
		return FieldAccess{
			Over:      irStatementToExpression(ctx, expressionToIR(ctx, caseAccess.Over)),
			FieldName: VariableName(nil, caseAccess.Access),
		}
	} else if caseInvocation != nil {
		genericsPassed := []Expression{}
		for _, generic := range caseInvocation.Generics {
			generic = types.FlattenOr(generic)
			_, _, _, _, caseOr := generic.VariableTypeCases()
			if caseOr == nil {
				caseOr = &types.OrVariableType{
					Elements: []types.VariableType{generic},
				}
			}
			genericArray := ArrayInstantiation{Arguments: []Expression{}}
			for _, element := range caseOr.Elements {
				genericArray.Arguments = append(genericArray.Arguments, ctx.resolvePrintableNameOrReference(element))
			}
			genericsPassed = append(genericsPassed, genericArray)
		}
		arguments := []Expression{
			ArrayInstantiation{
				Arguments: genericsPassed,
			},
		}
		for _, argument := range caseInvocation.Arguments {
			arguments = append(arguments, irStatementToExpression(ctx, expressionToIR(ctx, argument)))
		}
		return Invocation{
			CodePointAsFirstArgument: codePointAsFirstArgumentOfInvocation(*caseInvocation),
			Over:                     irStatementToExpression(ctx, expressionToIR(ctx, caseInvocation.Over)),
			Arguments:                arguments,
		}
	} else if caseFunction != nil {
		ctx, genericsVarName := ctx.CopyAddingFunctionGenerics(caseFunction.VariableType.Generics)
		parameterNames := []string{genericsVarName}
		for _, functionArgument := range caseFunction.VariableType.Arguments {
			parameterNames = append(parameterNames, VariableName(nil, functionArgument.Name))
		}
		block := []Statement{}
		for i, exp := range caseFunction.Block {
			newExp := expressionToIR(ctx, exp)
			if i < len(caseFunction.Block)-1 {
				block = append(block, newExp)
			} else {
				block = append(block, Return{
					ReturnExpression: irStatementToExpression(ctx, newExp),
				})
			}
		}
		return LocalFunction{
			ParameterNames: parameterNames,
			Block:          block,
		}
	} else if caseDeclaration != nil {
		return VariableDeclaration{
			Name:       VariableName(nil, caseDeclaration.Name),
			Expression: irStatementToExpression(ctx, expressionToIR(ctx, caseDeclaration.Expression)),
		}
	} else if caseIf != nil {
		condition := FieldAccess{
			Over:      irStatementToExpression(ctx, expressionToIR(ctx, caseIf.Condition)),
			FieldName: "value",
		}

		thenBlock := []Statement{}
		for i, exp := range caseIf.ThenBlock {
			newExp := expressionToIR(ctx, exp)
			if i < len(caseIf.ThenBlock)-1 {
				thenBlock = append(thenBlock, newExp)
			} else {
				thenBlock = append(thenBlock, Return{
					ReturnExpression: irStatementToExpression(ctx, newExp),
				})
			}
		}

		elseBlock := []Statement{}
		for i, exp := range caseIf.ElseBlock {
			newExp := expressionToIR(ctx, exp)
			if i < len(caseIf.ElseBlock)-1 {
				elseBlock = append(elseBlock, newExp)
			} else {
				elseBlock = append(elseBlock, Return{
					ReturnExpression: irStatementToExpression(ctx, newExp),
				})
			}
		}

		if len(elseBlock) == 0 {
			elseBlock = append(elseBlock, Return{
				ReturnExpression: Literal{
					Value: parser.LiteralNull{},
				},
			})
		}

		return Invocation{
			CodePointAsFirstArgument: nil,
			Over: LocalFunction{
				ParameterNames: []string{},
				Block: []Statement{
					If{
						Condition: condition,
						ThenBlock: thenBlock,
						ElseBlock: elseBlock,
					},
				},
			},
			Arguments: []Expression{},
		}
	} else if caseList != nil {
		arguments := []Expression{}
		for _, argument := range caseList.Arguments {
			arguments = append(arguments, irStatementToExpression(ctx, expressionToIR(ctx, argument)))
		}
		return ObjectInstantiation{
			Fields: map[string]Expression{
				"$type": NonFunctionTypeName{
					Prefix: "List<",
					CommaSpaceSeparated: []TypeNamesOrConcat{
						TypeNamesOrConcat{
							Expressions: []Expression{
								ctx.resolvePrintableNameOrReference(caseList.ContainedVariableType),
							},
						},
					},
					Suffix: ">",
				},
				"value": ArrayInstantiation{
					Arguments: arguments,
				},
			},
		}
	} else if caseWhen != nil {
		overExpression := irStatementToExpression(ctx, expressionToIR(ctx, caseWhen.Over))

		overVarDecl := VariableDeclaration{
			Name:       "__over",
			Expression: overExpression,
		}

		whenCases := []IfForTypeMatching{}
		for _, whenCase := range caseWhen.Cases {
			block := []Statement{}
			if whenCase.Name != nil {
				block = append(block, VariableDeclaration{
					Name: VariableName(nil, *whenCase.Name),
					Expression: Reference{
						Name: "__over",
					},
				})
			}
			for i, exp := range whenCase.Block {
				newExp := expressionToIR(ctx, exp)
				if i < len(whenCase.Block)-1 {
					block = append(block, newExp)
				} else {
					block = append(block, Return{
						ReturnExpression: irStatementToExpression(ctx, newExp),
					})
				}
			}
			whenCases = append(whenCases, IfForTypeMatching{
				Condition: IfArrayIncludesElseEqual{
					Left: ctx.resolvePrintableNameOrReference(whenCase.VariableType),
					Right: FieldAccess{
						Over: Reference{
							Name: "__over",
						},
						FieldName: "$type",
					},
				},
				ThenBlock: block,
				ElseBlock: []Statement{},
			})
		}
		otherCaseBlock := []Statement{}
		if caseWhen.OtherCaseName != nil {
			otherCaseBlock = append(otherCaseBlock, VariableDeclaration{
				Name: VariableName(nil, *caseWhen.OtherCaseName),
				Expression: Reference{
					Name: "__over",
				},
			})
		}
		for i, exp := range caseWhen.OtherCase {
			newExp := expressionToIR(ctx, exp)
			if i < len(caseWhen.OtherCase)-1 {
				otherCaseBlock = append(otherCaseBlock, newExp)
			} else {
				otherCaseBlock = append(otherCaseBlock, Return{
					ReturnExpression: irStatementToExpression(ctx, newExp),
				})
			}
		}
		if len(otherCaseBlock) == 0 {
			otherCaseBlock = append(otherCaseBlock, Return{
				ReturnExpression: Literal{
					Value: parser.LiteralNull{},
				},
			})
		}

		if len(whenCases) == 0 {
			panic("whenCases is empty")
		}

		allIfElseChainedTogether := IfForTypeMatching{}
		for i := len(whenCases) - 1; i >= 0; i-- {
			if i == len(whenCases)-1 {
				allIfElseChainedTogether = whenCases[i]
				allIfElseChainedTogether.ElseBlock = otherCaseBlock
			} else {
				whenCases[i].ElseBlock = []Statement{allIfElseChainedTogether}
				allIfElseChainedTogether = whenCases[i]
			}
		}

		return Invocation{
			CodePointAsFirstArgument: nil,
			Over: LocalFunction{
				ParameterNames: []string{},
				Block: []Statement{
					overVarDecl,
					allIfElseChainedTogether,
				},
			},
			Arguments: []Expression{},
		}
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func codePointAsFirstArgumentOfInvocation(invocation ast.Invocation) *CodePoint {
	_, _, _, caseFunction, _ := ast.VariableTypeOfExpression(invocation.Over).VariableTypeCases()
	if caseFunction == nil {
		panic(fmt.Sprintf("codePointAsFirstArgumentOfInvocation over %+v", invocation.Over))
	}
	if caseFunction.CodePointAsFirstArgument {
		return &CodePoint{
			FileName: invocation.CodePoint.FileName,
			Line:     invocation.CodePoint.Line,
		}
	} else {
		return nil
	}
}
