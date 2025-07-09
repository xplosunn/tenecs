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
}

func ToIR(program ast.Program) Program {
	ctx := context{
		topLevelDeclarations: []Reference{},
	}
	for ref, _ := range program.Declarations {
		ctx.topLevelDeclarations = append(ctx.topLevelDeclarations, refToIR(ref))
	}
	for ref, _ := range program.NativeFunctions {
		ctx.topLevelDeclarations = append(ctx.topLevelDeclarations, refToIR(ref))
	}

	declarations := map[Reference]TopLevelFunction{}
	for ref, expression := range program.Declarations {
		declarations[refToIR(ref)] = topLevelDeclarationToIR(ctx, expression)
	}
	structFunctions := map[Reference]*types.Function{}
	for ref, function := range program.StructFunctions {
		structFunctions[refToIR(ref)] = function
	}
	nativeFunctions := map[NativeFunctionRef]*types.Function{}
	for ref, function := range program.NativeFunctions {
		nativeFunctions[NativeFunctionRef{
			Package: ref.Package,
			Name:    ref.Name,
		}] = function
	}

	return Program{
		Declarations:    declarations,
		StructFunctions: structFunctions,
		NativeFunctions: nativeFunctions,
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

func refToIR(ref ast.Ref) Reference {
	return Reference{
		Name: VariableName(&ref.Package, ref.Name),
	}
}

func topLevelDeclarationToIR(ctx context, expression ast.Expression) TopLevelFunction {
	topLevelFunction := TopLevelFunction{
		ParameterNames: []string{},
		Body: []Statement{
			Return{
				ReturnExpression: irStatementToExpression(ctx, expressionToIR(ctx, expression)),
			},
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
		panic("TODO irStatementToExpression If")
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

		isTopLevelReference := false
		for _, topLevelDeclaration := range ctx.topLevelDeclarations {
			if topLevelDeclaration.Name == reference.Name {
				isTopLevelReference = true
				break
			}
		}

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
			FieldName: caseAccess.Access,
		}
	} else if caseInvocation != nil {
		arguments := []Expression{}
		for _, argument := range caseInvocation.Arguments {
			arguments = append(arguments, irStatementToExpression(ctx, expressionToIR(ctx, argument)))
		}
		return Invocation{
			Over:      irStatementToExpression(ctx, expressionToIR(ctx, caseInvocation.Over)),
			Arguments: arguments,
		}
	} else if caseFunction != nil {
		parameterNames := []string{}
		for _, functionArgument := range caseFunction.VariableType.Arguments {
			parameterNames = append(parameterNames, functionArgument.Name)
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
			Name:       caseDeclaration.Name,
			Expression: irStatementToExpression(ctx, expressionToIR(ctx, caseDeclaration.Expression)),
		}
	} else if caseIf != nil {
		panic("TODO expressionToIR caseIf")
	} else if caseList != nil {
		panic("TODO expressionToIR caseList")
	} else if caseWhen != nil {
		panic("TODO expressionToIR caseWhen")
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}
