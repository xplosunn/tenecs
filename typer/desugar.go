package typer

import (
	"errors"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/type_error"
)

func DesugarFileTopLevel(file string, parsed parser.FileTopLevel) (parser.FileTopLevel, error) {
	var err error
	for i, topLevelDeclaration := range parsed.TopLevelDeclarations {
		parser.TopLevelDeclarationExhaustiveSwitch(
			topLevelDeclaration,
			func(topLevelDeclaration parser.Declaration) {
				if topLevelDeclaration.ShortCircuit != nil {
					err = errors.New("shortcircuit only allowed inside of functions")
					return
				}
				p, _, e := desugarExpressionBox(file, topLevelDeclaration.ExpressionBox, []parser.ExpressionBox{})
				if e != nil {
					err = e
					return
				}
				topLevelDeclaration.ExpressionBox = p
				parsed.TopLevelDeclarations[i] = topLevelDeclaration
			},
			func(topLevelDeclaration parser.Struct) {},
			func(topLevelDeclaration parser.TypeAlias) {},
		)
	}
	return parsed, err
}

func desugarExpressionBox(file string, parsed parser.ExpressionBox, restOfBlock []parser.ExpressionBox) (parser.ExpressionBox, []parser.ExpressionBox, error) {
	exp, restOfBlock, err := desugarExpression(file, parsed.Expression, restOfBlock)
	if err != nil {
		return parsed, restOfBlock, err
	}
	parsed.Expression = exp
	for i, accessOrInvocation := range parsed.AccessOrInvocationChain {
		if accessOrInvocation.Arguments != nil {
			for i2, argument := range accessOrInvocation.Arguments.Arguments {
				d, _, err := desugarExpressionBox(file, argument.Argument, []parser.ExpressionBox{})
				if err != nil {
					return parsed, restOfBlock, err
				}
				parsed.AccessOrInvocationChain[i].Arguments.Arguments[i2].Argument = d
			}
		}
	}
	for i, accessOrInvocation := range parsed.AccessOrInvocationChain {
		if accessOrInvocation.DotOrArrowName != nil && accessOrInvocation.DotOrArrowName.Arrow {
			expressionBeforeThisArrow := parser.ExpressionBox{
				Node:                    parsed.Node,
				Expression:              parsed.Expression,
				AccessOrInvocationChain: []parser.AccessOrInvocation{},
			}
			if i > 0 {
				expressionBeforeThisArrow.AccessOrInvocationChain = parsed.AccessOrInvocationChain[0:i]
			}
			if accessOrInvocation.Arguments == nil {
				return parser.ExpressionBox{}, nil, type_error.PtrOnNodef(file, accessOrInvocation.DotOrArrowName.Node, "Arrow syntax requires parenthesis on the right-hand side")
			}
			newParsedExpression := parser.ReferenceOrInvocation{
				Var: accessOrInvocation.DotOrArrowName.VarName,
				Arguments: &parser.ArgumentsList{
					Node:     accessOrInvocation.Node,
					Generics: accessOrInvocation.Arguments.Generics,
					Arguments: append([]parser.NamedArgument{
						parser.NamedArgument{
							Node:     expressionBeforeThisArrow.Node,
							Argument: expressionBeforeThisArrow,
						},
					}, accessOrInvocation.Arguments.Arguments...),
				},
			}
			for i, argument := range newParsedExpression.Arguments.Arguments {
				desugared, restOfBlock, err := desugarExpressionBox(file, argument.Argument, nil)
				if err != nil {
					return desugared, nil, err
				}
				if len(restOfBlock) > 0 {
					panic("didn't expect rest of block when none is passed")
				}
				newParsedExpression.Arguments.Arguments[i].Argument = desugared
			}

			parsed.Expression = newParsedExpression
			if i < len(parsed.AccessOrInvocationChain) {
				parsed.AccessOrInvocationChain = parsed.AccessOrInvocationChain[i+1:]
			} else {
				parsed.AccessOrInvocationChain = []parser.AccessOrInvocation{}
			}
			return desugarExpressionBox(file, parsed, restOfBlock)
		}
	}
	return parsed, restOfBlock, nil
}

func desugarExpression(file string, parsed parser.Expression, restOfBlock []parser.ExpressionBox) (parser.Expression, []parser.ExpressionBox, error) {
	var err error
	parser.ExpressionExhaustiveSwitch(
		parsed,
		func(expression parser.LiteralExpression) {

		},
		func(expression parser.ReferenceOrInvocation) {
			if expression.Arguments != nil {
				for i, argument := range expression.Arguments.Arguments {
					d, _, e := desugarExpressionBox(file, argument.Argument, []parser.ExpressionBox{})
					err = e
					if err != nil {
						return
					}
					expression.Arguments.Arguments[i].Argument = d
				}
			}
			parsed = expression
		},
		func(generics *parser.LambdaOrListGenerics, expression parser.Lambda) {
			d, e := desugarBlock(file, expression.Block)
			err = e
			if err != nil {
				return
			}
			expression.Block = d
			parsedLambdaOrList := parser.LambdaOrList{
				Node:     expression.Node,
				Generics: generics,
				List:     nil,
				Lambda:   &expression,
			}
			if generics != nil {
				parsedLambdaOrList.Node = generics.Node
			}
			parsed = parsedLambdaOrList
		},
		func(expression parser.Declaration) {
			d, _, e := desugarExpressionBox(file, expression.ExpressionBox, []parser.ExpressionBox{})
			err = e
			if err != nil {
				return
			}
			expression.ExpressionBox = d
			parsed = expression
			if expression.ShortCircuit != nil {
				if expression.ShortCircuit.TypeAnnotation == nil && expression.TypeAnnotation == nil {
					err = errors.New("when shortciruiting one of the types needs to be annotated")
				} else if expression.ShortCircuit.TypeAnnotation != nil && expression.TypeAnnotation != nil {
					name := parser.Name{
						Node:   expression.Name.Node,
						String: expression.Name.String,
					}
					if name.String == "_" {
						name.String = "_unused_"
					}
					parsed = parser.When{
						Node: expression.Name.Node,
						Over: expression.ExpressionBox,
						Is: []parser.WhenIs{
							parser.WhenIs{
								Node: expression.ShortCircuit.TypeAnnotation.Node,
								Name: &name,
								Type: *expression.ShortCircuit.TypeAnnotation,
								ThenBlock: []parser.ExpressionBox{
									parser.ExpressionBox{
										Node: expression.ShortCircuit.TypeAnnotation.Node,
										Expression: parser.ReferenceOrInvocation{
											Var:       name,
											Arguments: nil,
										},
										AccessOrInvocationChain: []parser.AccessOrInvocation{},
									},
								},
							},
							parser.WhenIs{
								Node:      expression.TypeAnnotation.Node,
								Name:      &name,
								Type:      *expression.TypeAnnotation,
								ThenBlock: restOfBlock,
							},
						},
						Other: nil,
					}
				} else if expression.ShortCircuit.TypeAnnotation != nil {
					name := parser.Name{
						Node:   expression.Name.Node,
						String: expression.Name.String,
					}
					if name.String == "_" {
						name.String = "_unused_"
					}
					parsed = parser.When{
						Node: expression.Name.Node,
						Over: expression.ExpressionBox,
						Is: []parser.WhenIs{
							parser.WhenIs{
								Node: expression.ShortCircuit.TypeAnnotation.Node,
								Name: &name,
								Type: *expression.ShortCircuit.TypeAnnotation,
								ThenBlock: []parser.ExpressionBox{
									parser.ExpressionBox{
										Node: expression.ShortCircuit.TypeAnnotation.Node,
										Expression: parser.ReferenceOrInvocation{
											Var:       name,
											Arguments: nil,
										},
										AccessOrInvocationChain: []parser.AccessOrInvocation{},
									},
								},
							},
						},
						Other: &parser.WhenOther{
							Node:      expression.Name.Node,
							Name:      &name,
							ThenBlock: restOfBlock,
						},
					}
				} else {
					name := parser.Name{
						Node:   expression.Name.Node,
						String: expression.Name.String,
					}
					if name.String == "_" {
						name.String = "_unused_"
					}
					parsed = parser.When{
						Node: expression.Name.Node,
						Over: expression.ExpressionBox,
						Is: []parser.WhenIs{
							parser.WhenIs{
								Node:      expression.TypeAnnotation.Node,
								Name:      &name,
								Type:      *expression.TypeAnnotation,
								ThenBlock: restOfBlock,
							},
						},
						Other: &parser.WhenOther{
							Node: expression.TypeAnnotation.Node,
							Name: &name,
							ThenBlock: []parser.ExpressionBox{
								parser.ExpressionBox{
									Node: expression.TypeAnnotation.Node,
									Expression: parser.ReferenceOrInvocation{
										Var:       name,
										Arguments: nil,
									},
									AccessOrInvocationChain: []parser.AccessOrInvocation{},
								},
							},
						},
					}
				}
				restOfBlock = []parser.ExpressionBox{}
			}
		},
		func(expression parser.If) {
			cond, _, e := desugarExpressionBox(file, expression.Condition, []parser.ExpressionBox{})
			err = e
			if err != nil {
				return
			}
			expression.Condition = cond

			then, e := desugarBlock(file, expression.ThenBlock)
			err = e
			if err != nil {
				return
			}
			expression.ThenBlock = then

			for i, elseIf := range expression.ElseIfs {
				cond, _, e := desugarExpressionBox(file, elseIf.Condition, []parser.ExpressionBox{})
				err = e
				if err != nil {
					return
				}
				elseIf.Condition = cond

				then, e := desugarBlock(file, elseIf.ThenBlock)
				err = e
				if err != nil {
					return
				}
				elseIf.ThenBlock = then
				expression.ElseIfs[i] = elseIf
			}

			elseThen, e := desugarBlock(file, expression.ElseBlock)
			err = e
			if err != nil {
				return
			}
			expression.ElseBlock = elseThen

			parsed = expression
		},
		func(generics *parser.LambdaOrListGenerics, expression parser.List) {
			for i, expressionBox := range expression.Expressions {
				d, _, e := desugarExpressionBox(file, expressionBox, []parser.ExpressionBox{})
				err = e
				if err != nil {
					return
				}
				expression.Expressions[i] = d
			}
			parsedLambdaOrList := parser.LambdaOrList{
				Node:     expression.Node,
				Generics: generics,
				List:     &expression,
				Lambda:   nil,
			}
			if generics != nil {
				parsedLambdaOrList.Node = generics.Node
			}
			parsed = parsedLambdaOrList
		},
		func(expression parser.When) {
			over, _, e := desugarExpressionBox(file, expression.Over, []parser.ExpressionBox{})
			err = e
			if err != nil {
				return
			}
			expression.Over = over

			for i, is := range expression.Is {
				d, e := desugarBlock(file, is.ThenBlock)
				err = e
				if err != nil {
					return
				}
				expression.Is[i].ThenBlock = d
			}

			if expression.Other != nil {
				d, e := desugarBlock(file, expression.Other.ThenBlock)
				err = e
				if err != nil {
					return
				}
				expression.Other.ThenBlock = d
			}
			parsed = expression
		},
	)
	return parsed, restOfBlock, err
}

func desugarBlock(file string, block []parser.ExpressionBox) ([]parser.ExpressionBox, error) {
	for i := len(block) - 1; i >= 0; i-- {
		expressionBox := block[i]
		d, r, err := desugarExpressionBox(file, expressionBox, block[i+1:len(block)])
		if err != nil {
			return nil, err
		}
		if i > 0 {
			block = append(append(block[0:i], d), r...)
		} else {
			block = append([]parser.ExpressionBox{d}, r...)
		}
	}
	return block, nil
}
