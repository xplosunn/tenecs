package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
	"strings"
)

func Untypecheck(program *ast.Program) parser.FileTopLevel {
	pkg := []parser.Name{}
	for _, str := range strings.Split(program.Package, ".") {
		pkg = append(pkg, parser.Name{
			String: str,
		})
	}

	topLevelDeclarations := []parser.TopLevelDeclaration{}

	for constructorName, constructor := range program.StructFunctions {
		generics := []parser.Name{}
		for _, generic := range constructor.Generics {
			generics = append(generics, parser.Name{
				String: generic,
			})
		}
		if len(generics) == 0 {
			generics = nil
		}
		variables := []parser.StructVariable{}
		for _, arg := range constructor.Arguments {
			variables = append(variables, parser.StructVariable{
				Name: parser.Name{
					String: arg.Name,
				},
				Type: untypecheckTypeAnnotation(arg.VariableType),
			})
		}
		topLevelDeclarations = append(topLevelDeclarations, parser.Struct{
			Name: parser.Name{
				String: constructorName,
			},
			Generics:  generics,
			Variables: variables,
		})
	}

	for _, declaration := range program.Declarations {
		topLevelDeclarations = append(topLevelDeclarations, parser.Declaration{
			Name: parser.Name{
				String: declaration.Name,
			},
			TypeAnnotation: ptr(untypecheckTypeAnnotation(ast.VariableTypeOfExpression(declaration.Expression))),
			ShortCircuit:   nil,
			ExpressionBox:  untypecheckExpression(declaration.Expression),
		})
	}

	imports := []parser.Import{}
	for functionName, functionPkg := range program.NativeFunctionPackages {
		dotSeparatedVars := []parser.Name{}
		for _, pkgPart := range strings.Split(functionPkg, "_") {
			dotSeparatedVars = append(dotSeparatedVars, parser.Name{
				String: pkgPart,
			})
		}
		dotSeparatedVars = append(dotSeparatedVars, parser.Name{
			String: functionName,
		})
		imports = append(imports, parser.Import{
			DotSeparatedVars: dotSeparatedVars,
			As:               nil,
		})
	}

	return parser.FileTopLevel{
		Package: parser.Package{
			DotSeparatedNames: pkg,
		},
		Imports:              imports,
		TopLevelDeclarations: topLevelDeclarations,
	}
}

func ptr[T any](t T) *T {
	return &t
}

func untypecheckExpression(expression ast.Expression) parser.ExpressionBox {
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
		return parser.ExpressionBox{
			Expression: parser.LiteralExpression{
				Literal: caseLiteral.Literal,
			},
			AccessOrInvocationChain: nil,
		}
	} else if caseReference != nil {
		return parser.ExpressionBox{
			Expression: parser.ReferenceOrInvocation{
				Var: parser.Name{
					String: caseReference.Name,
				},
				Arguments: nil,
			},
			AccessOrInvocationChain: nil,
		}
	} else if caseAccess != nil {
		over := untypecheckExpression(caseAccess.Over)
		accessOrInvocationChain := []parser.AccessOrInvocation{}
		accessOrInvocationChain = append(accessOrInvocationChain, over.AccessOrInvocationChain...)
		accessOrInvocationChain = append(accessOrInvocationChain, parser.AccessOrInvocation{
			DotOrArrowName: &parser.DotOrArrowName{
				Dot:   true,
				Arrow: false,
				VarName: parser.Name{
					String: caseAccess.Access,
				},
			},
			Arguments: nil,
		})

		return parser.ExpressionBox{
			Expression:              over.Expression,
			AccessOrInvocationChain: accessOrInvocationChain,
		}
	} else if caseInvocation != nil {
		over := untypecheckExpression(caseInvocation.Over)
		accessOrInvocationChain := []parser.AccessOrInvocation{}
		accessOrInvocationChain = append(accessOrInvocationChain, over.AccessOrInvocationChain...)

		generics := []parser.TypeAnnotation{}
		for _, generic := range caseInvocation.Generics {
			generics = append(generics, untypecheckTypeAnnotation(generic))
		}

		arguments := []parser.NamedArgument{}
		for _, argument := range caseInvocation.Arguments {
			arguments = append(arguments, parser.NamedArgument{
				Name:     nil,
				Argument: untypecheckExpression(argument),
			})
		}
		accessOrInvocationChain = append(accessOrInvocationChain, parser.AccessOrInvocation{
			DotOrArrowName: nil,
			Arguments: &parser.ArgumentsList{
				Generics:  generics,
				Arguments: arguments,
			},
		})

		return parser.ExpressionBox{
			Expression:              over.Expression,
			AccessOrInvocationChain: accessOrInvocationChain,
		}
	} else if caseFunction != nil {
		generics := []parser.Name{}
		for _, generic := range caseFunction.VariableType.Generics {
			generics = append(generics, parser.Name{
				String: generic,
			})
		}
		if len(generics) == 0 {
			generics = nil
		}

		parameters := []parser.Parameter{}
		for _, argument := range caseFunction.VariableType.Arguments {
			parameters = append(parameters, parser.Parameter{
				Name: parser.Name{
					String: argument.Name,
				},
				Type: ptr(untypecheckTypeAnnotation(argument.VariableType)),
			})
		}

		block := []parser.ExpressionBox{}
		for _, expression := range caseFunction.Block {
			expBox := untypecheckExpression(expression)
			block = append(block, expBox)
		}

		return parser.ExpressionBox{
			Expression: parser.Lambda{
				Signature: parser.LambdaSignature{
					Generics:   generics,
					Parameters: parameters,
					ReturnType: ptr(untypecheckTypeAnnotation(caseFunction.VariableType.ReturnType)),
				},
				Block: block,
			},
			AccessOrInvocationChain: nil,
		}
	} else if caseDeclaration != nil {
		return parser.ExpressionBox{
			Expression: parser.Declaration{
				Name: parser.Name{
					String: caseDeclaration.Name,
				},
				TypeAnnotation: ptr(untypecheckTypeAnnotation(ast.VariableTypeOfExpression(caseDeclaration.Expression))),
				ShortCircuit:   nil,
				ExpressionBox:  untypecheckExpression(caseDeclaration.Expression),
			},
		}
	} else if caseIf != nil {
		thenBlock := []parser.ExpressionBox{}
		for _, expression := range caseIf.ThenBlock {
			expBox := untypecheckExpression(expression)
			thenBlock = append(thenBlock, expBox)
		}
		elseBlock := []parser.ExpressionBox{}
		for _, expression := range caseIf.ElseBlock {
			expBox := untypecheckExpression(expression)
			elseBlock = append(elseBlock, expBox)
		}
		return parser.ExpressionBox{
			Expression: parser.If{
				Condition: untypecheckExpression(caseIf.Condition),
				ThenBlock: thenBlock,
				ElseIfs:   nil,
				ElseBlock: elseBlock,
			},
			AccessOrInvocationChain: nil,
		}
	} else if caseList != nil {
		expressions := []parser.ExpressionBox{}
		for _, argument := range caseList.Arguments {
			expressions = append(expressions, untypecheckExpression(argument))
		}
		return parser.ExpressionBox{
			Expression: parser.List{
				Generic:     ptr(untypecheckTypeAnnotation(caseList.ContainedVariableType)),
				Expressions: expressions,
			},
			AccessOrInvocationChain: nil,
		}
	} else if caseWhen != nil {
		is := []parser.WhenIs{}
		for _, whenCase := range caseWhen.Cases {
			var name *parser.Name
			if whenCase.Name != nil {
				name = &parser.Name{
					String: *whenCase.Name,
				}
			}

			block := []parser.ExpressionBox{}
			for _, expression := range whenCase.Block {
				expBox := untypecheckExpression(expression)
				block = append(block, expBox)
			}

			is = append(is, parser.WhenIs{
				Name:      name,
				Type:      untypecheckTypeAnnotation(whenCase.VariableType),
				ThenBlock: block,
			})
		}

		var other *parser.WhenOther
		if len(caseWhen.OtherCase) > 0 {
			var name *parser.Name
			if caseWhen.OtherCaseName != nil {
				name = &parser.Name{
					String: *caseWhen.OtherCaseName,
				}
			}

			block := []parser.ExpressionBox{}
			for _, expression := range caseWhen.OtherCase {
				expBox := untypecheckExpression(expression)
				block = append(block, expBox)
			}

			other = &parser.WhenOther{
				Name:      name,
				ThenBlock: block,
			}
		}

		return parser.ExpressionBox{
			Expression: parser.When{
				Over:  untypecheckExpression(caseWhen.Over),
				Is:    is,
				Other: other,
			},
			AccessOrInvocationChain: nil,
		}
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func untypecheckTypeAnnotation(varType types.VariableType) parser.TypeAnnotation {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		return parser.TypeAnnotation{
			OrTypes: []parser.TypeAnnotationElement{
				parser.SingleNameType{
					TypeName: parser.Name{
						String: caseTypeArgument.Name,
					},
					Generics: nil,
				},
			},
		}
	} else if caseKnownType != nil {
		generics := []parser.TypeAnnotation{}
		for _, generic := range caseKnownType.Generics {
			generics = append(generics, untypecheckTypeAnnotation(generic))
		}
		if len(generics) == 0 {
			generics = nil
		}
		return parser.TypeAnnotation{
			OrTypes: []parser.TypeAnnotationElement{
				parser.SingleNameType{
					TypeName: parser.Name{
						String: caseKnownType.Name,
					},
					Generics: generics,
				},
			},
		}
	} else if caseFunction != nil {
		generics := []parser.Name{}
		for _, generic := range caseFunction.Generics {
			generics = append(generics, parser.Name{
				String: generic,
			})
		}
		if len(generics) == 0 {
			generics = nil
		}

		arguments := []parser.FunctionTypeArgument{}
		for _, argument := range caseFunction.Arguments {
			arguments = append(arguments, parser.FunctionTypeArgument{
				Name: &parser.Name{
					String: argument.Name,
				},
				Type: untypecheckTypeAnnotation(argument.VariableType),
			})
		}

		return parser.TypeAnnotation{
			OrTypes: []parser.TypeAnnotationElement{
				parser.FunctionType{
					Generics:   generics,
					Arguments:  arguments,
					ReturnType: untypecheckTypeAnnotation(caseFunction.ReturnType),
				},
			},
		}
	} else if caseOr != nil {
		result := parser.TypeAnnotation{
			OrTypes: []parser.TypeAnnotationElement{},
		}
		for _, orElem := range caseOr.Elements {
			result.OrTypes = append(result.OrTypes, untypecheckTypeAnnotation(orElem).OrTypes...)
		}
		return result
	} else {
		panic("cases on variableType")
	}
}
