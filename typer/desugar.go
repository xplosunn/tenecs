package typer

import (
	"errors"
	"github.com/xplosunn/tenecs/parser"
)

func desugarFileTopLevel(parsed parser.FileTopLevel) (parser.FileTopLevel, error) {
	var err error
	for i, topLevelDeclaration := range parsed.TopLevelDeclarations {
		parser.TopLevelDeclarationExhaustiveSwitch(
			topLevelDeclaration,
			func(topLevelDeclaration parser.Declaration) {
				if topLevelDeclaration.ShortCircuit != nil {
					err = errors.New("shortcircuit only allowed inside of functions")
					return
				}
				p, _, e := desugarExpressionBox(topLevelDeclaration.ExpressionBox, []parser.ExpressionBox{})
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

func desugarExpressionBox(parsed parser.ExpressionBox, restOfBlock []parser.ExpressionBox) (parser.ExpressionBox, []parser.ExpressionBox, error) {
	exp, restOfBlock, err := desugarExpression(parsed.Expression, restOfBlock)
	if err != nil {
		return parsed, restOfBlock, err
	}
	parsed.Expression = exp
	for i, accessOrInvocation := range parsed.AccessOrInvocationChain {
		if accessOrInvocation.Arguments != nil {
			for i2, argument := range accessOrInvocation.Arguments.Arguments {
				d, _, err := desugarExpressionBox(argument.Argument, []parser.ExpressionBox{})
				if err != nil {
					return parsed, restOfBlock, err
				}
				parsed.AccessOrInvocationChain[i].Arguments.Arguments[i2].Argument = d
			}
		}
	}
	return parsed, restOfBlock, nil
}

func desugarExpression(parsed parser.Expression, restOfBlock []parser.ExpressionBox) (parser.Expression, []parser.ExpressionBox, error) {
	var err error
	parser.ExpressionExhaustiveSwitch(
		parsed,
		func(expression parser.LiteralExpression) {

		},
		func(expression parser.ReferenceOrInvocation) {
			if expression.Arguments != nil {
				for i, argument := range expression.Arguments.Arguments {
					d, _, e := desugarExpressionBox(argument.Argument, []parser.ExpressionBox{})
					err = e
					if err != nil {
						return
					}
					expression.Arguments.Arguments[i].Argument = d
				}
			}
			parsed = expression
		},
		func(expression parser.Lambda) {
			d, e := desugarBlock(expression.Block)
			err = e
			if err != nil {
				return
			}
			expression.Block = d
			parsed = expression
		},
		func(expression parser.Declaration) {
			d, _, e := desugarExpressionBox(expression.ExpressionBox, []parser.ExpressionBox{})
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
					parsed = parser.When{
						Node: expression.Name.Node,
						Over: expression.ExpressionBox,
						Is: []parser.WhenIs{
							parser.WhenIs{
								Node: expression.ShortCircuit.TypeAnnotation.Node,
								Name: &expression.Name,
								Type: *expression.ShortCircuit.TypeAnnotation,
								ThenBlock: []parser.ExpressionBox{
									parser.ExpressionBox{
										Node: expression.ShortCircuit.TypeAnnotation.Node,
										Expression: parser.ReferenceOrInvocation{
											Var:       expression.Name,
											Arguments: nil,
										},
										AccessOrInvocationChain: []parser.AccessOrInvocation{},
									},
								},
							},
							parser.WhenIs{
								Node:      expression.TypeAnnotation.Node,
								Name:      &expression.Name,
								Type:      *expression.TypeAnnotation,
								ThenBlock: restOfBlock,
							},
						},
						Other: nil,
					}
				} else if expression.ShortCircuit.TypeAnnotation != nil {
					parsed = parser.When{
						Node: expression.Name.Node,
						Over: expression.ExpressionBox,
						Is: []parser.WhenIs{
							parser.WhenIs{
								Node: expression.ShortCircuit.TypeAnnotation.Node,
								Name: &expression.Name,
								Type: *expression.ShortCircuit.TypeAnnotation,
								ThenBlock: []parser.ExpressionBox{
									parser.ExpressionBox{
										Node: expression.ShortCircuit.TypeAnnotation.Node,
										Expression: parser.ReferenceOrInvocation{
											Var:       expression.Name,
											Arguments: nil,
										},
										AccessOrInvocationChain: []parser.AccessOrInvocation{},
									},
								},
							},
						},
						Other: &parser.WhenOther{
							Node:      expression.Name.Node,
							Name:      &expression.Name,
							ThenBlock: restOfBlock,
						},
					}
				} else {
					parsed = parser.When{
						Node: expression.Name.Node,
						Over: expression.ExpressionBox,
						Is: []parser.WhenIs{
							parser.WhenIs{
								Node:      expression.TypeAnnotation.Node,
								Name:      &expression.Name,
								Type:      *expression.TypeAnnotation,
								ThenBlock: restOfBlock,
							},
						},
						Other: &parser.WhenOther{
							Node: expression.TypeAnnotation.Node,
							Name: &expression.Name,
							ThenBlock: []parser.ExpressionBox{
								parser.ExpressionBox{
									Node: expression.TypeAnnotation.Node,
									Expression: parser.ReferenceOrInvocation{
										Var:       expression.Name,
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
			cond, _, e := desugarExpressionBox(expression.Condition, []parser.ExpressionBox{})
			err = e
			if err != nil {
				return
			}
			expression.Condition = cond

			then, e := desugarBlock(expression.ThenBlock)
			err = e
			if err != nil {
				return
			}
			expression.ThenBlock = then

			for i, elseIf := range expression.ElseIfs {
				cond, _, e := desugarExpressionBox(elseIf.Condition, []parser.ExpressionBox{})
				err = e
				if err != nil {
					return
				}
				elseIf.Condition = cond

				then, e := desugarBlock(elseIf.ThenBlock)
				err = e
				if err != nil {
					return
				}
				elseIf.ThenBlock = then
				expression.ElseIfs[i] = elseIf
			}

			elseThen, e := desugarBlock(expression.ElseBlock)
			err = e
			if err != nil {
				return
			}
			expression.ElseBlock = elseThen

			parsed = expression
		},
		func(expression parser.List) {
			for i, expressionBox := range expression.Expressions {
				d, _, e := desugarExpressionBox(expressionBox, []parser.ExpressionBox{})
				err = e
				if err != nil {
					return
				}
				expression.Expressions[i] = d
			}
			parsed = expression
		},
		func(expression parser.When) {
			over, _, e := desugarExpressionBox(expression.Over, []parser.ExpressionBox{})
			err = e
			if err != nil {
				return
			}
			expression.Over = over

			for i, is := range expression.Is {
				d, e := desugarBlock(is.ThenBlock)
				err = e
				if err != nil {
					return
				}
				expression.Is[i].ThenBlock = d
			}

			if expression.Other != nil {
				d, e := desugarBlock(expression.Other.ThenBlock)
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

func desugarBlock(block []parser.ExpressionBox) ([]parser.ExpressionBox, error) {
	for i := len(block) - 1; i >= 0; i-- {
		expressionBox := block[i]
		d, r, err := desugarExpressionBox(expressionBox, block[i+1:len(block)])
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
