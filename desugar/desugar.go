package desugar

import (
	"errors"

	"github.com/xplosunn/tenecs/parser"
)

func Desugar(parsed parser.FileTopLevel) FileTopLevel {
	return FileTopLevel{
		Tokens:               parsed.Tokens,
		Package:              desugarPackage(parsed.Package),
		Imports:              desugarSlice(parsed.Imports, desugarImport),
		TopLevelDeclarations: desugarSlice(parsed.TopLevelDeclarations, desugarTopLevelDeclaration),
	}
}

func DesugarFunctionType(parsed parser.FunctionType) FunctionType {
	return FunctionType{
		Generics:   desugarSlice(parsed.Generics, desugarName),
		Arguments:  desugarSlice(parsed.Arguments, desugarFunctionTypeArgument),
		ReturnType: desugarTypeAnnotation(parsed.ReturnType),
	}
}

func desugarSlice[In any, Out any](in []In, desugar func(In) Out) []Out {
	result := []Out{}
	for _, in := range in {
		result = append(result, desugar(in))
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func desugarWhenNonNil[In any, Out any](ptr *In, desugar func(In) Out) *Out {
	if ptr == nil {
		return nil
	}
	result := desugar(*ptr)
	return &result
}

func desugarPackage(parsed parser.Package) Package {
	return Package{
		Node:              desugarNode(parsed.Node),
		DotSeparatedNames: desugarSlice(parsed.DotSeparatedNames, desugarName),
	}
}

func desugarName(parsed parser.Name) Name {
	return Name{
		Node:   desugarNode(parsed.Node),
		String: parsed.String,
	}
}

func desugarNode(parsed parser.Node) Node {
	return Node{
		Pos:    parsed.Pos,
		EndPos: parsed.EndPos,
	}
}

func desugarImport(parsed parser.Import) Import {
	return Import{
		Node:             desugarNode(parsed.Node),
		DotSeparatedVars: desugarSlice(parsed.DotSeparatedVars, desugarName),
		As:               desugarWhenNonNil(parsed.As, desugarName),
	}
}

func desugarTopLevelDeclaration(parsed parser.TopLevelDeclaration) TopLevelDeclaration {
	var result TopLevelDeclaration
	parser.TopLevelDeclarationExhaustiveSwitch(
		parsed,
		func(parsed parser.Declaration) {
			result = desugarDeclaration(parsed)
		},
		func(parsed parser.Struct) {
			result = desugarStruct(parsed)
		},
		func(parsed parser.TypeAlias) {
			result = desugarTypeAlias(parsed)
		},
	)
	return result
}

func desugarDeclaration(parsed parser.Declaration) Declaration {
	return Declaration{
		Name:           desugarName(parsed.Name),
		TypeAnnotation: desugarWhenNonNil(parsed.TypeAnnotation, desugarTypeAnnotation),
		ShortCircuit:   desugarWhenNonNil(parsed.ShortCircuit, desugarDeclarationShortCircuit),
		ExpressionBox:  desugarExpressionBox(parsed.ExpressionBox),
	}
}

func desugarExpressionBox(parsed parser.ExpressionBox) ExpressionBox {
	return ExpressionBox{
		Node:                    desugarNode(parsed.Node),
		Expression:              desugarExpression(parsed.Expression),
		AccessOrInvocationChain: desugarSlice(parsed.AccessOrInvocationChain, desugarAccessOrInvocation),
	}
}

func desugarAccessOrInvocation(parsed parser.AccessOrInvocation) AccessOrInvocation {
	return AccessOrInvocation{
		Node:           desugarNode(parsed.Node),
		DotOrArrowName: desugarWhenNonNil(parsed.DotOrArrowName, desugarDotOrArrowName),
		Arguments:      desugarWhenNonNil(parsed.Arguments, desugarArgumentsList),
	}
}

func desugarArgumentsList(parsed parser.ArgumentsList) ArgumentsList {
	return ArgumentsList{
		Node:      desugarNode(parsed.Node),
		Generics:  desugarSlice(parsed.Generics, desugarTypeAnnotation),
		Arguments: desugarSlice(parsed.Arguments, desugarNamedArgument),
	}
}

func desugarNamedArgument(parsed parser.NamedArgument) NamedArgument {
	return NamedArgument{
		Node:     desugarNode(parsed.Node),
		Name:     desugarWhenNonNil(parsed.Name, desugarName),
		Argument: desugarExpressionBox(parsed.Argument),
	}
}

func desugarDotOrArrowName(parsed parser.DotOrArrowName) DotOrArrowName {
	return DotOrArrowName{
		Node:    desugarNode(parsed.Node),
		Dot:     parsed.Dot,
		Arrow:   parsed.Arrow,
		VarName: desugarName(parsed.VarName),
	}
}

func desugarExpression(parsed parser.Expression) Expression {
	var result Expression
	parser.ExpressionExhaustiveSwitch(
		parsed,
		func(parsed parser.LiteralExpression) {
			result = LiteralExpression{
				Node:    desugarNode(parsed.Node),
				Literal: parsed.Literal,
			}
		},
		func(parsed parser.ReferenceOrInvocation) {
			result = ReferenceOrInvocation{
				Var:       desugarName(parsed.Var),
				Arguments: desugarWhenNonNil(parsed.Arguments, desugarArgumentsList),
			}
		},
		func(generics *parser.LambdaOrListGenerics, parsed parser.Lambda) {
			lambdaGenerics := []TypeAnnotation{}
			if generics != nil {
				lambdaGenerics = desugarSlice(generics.Generics, desugarTypeAnnotation)
			} else {
				lambdaGenerics = nil
			}
			result = Lambda{
				Node:      desugarNode(parsed.Node),
				Generics:  lambdaGenerics,
				Signature: desugarLambdaSignature(parsed.Signature),
				Block:     desugarSlice(parsed.Block, desugarExpressionBox),
			}
		},
		func(parsed parser.Declaration) {
			result = Declaration{
				Name:           desugarName(parsed.Name),
				TypeAnnotation: desugarWhenNonNil(parsed.TypeAnnotation, desugarTypeAnnotation),
				ShortCircuit:   desugarWhenNonNil(parsed.ShortCircuit, desugarDeclarationShortCircuit),
				ExpressionBox:  desugarExpressionBox(parsed.ExpressionBox),
			}
		},
		func(parsed parser.If) {
			result = If{
				Node:      desugarNode(parsed.Node),
				Condition: desugarExpressionBox(parsed.Condition),
				ThenBlock: desugarSlice(parsed.ThenBlock, desugarExpressionBox),
				ElseIfs:   desugarSlice(parsed.ElseIfs, desugarIfThen),
				ElseBlock: desugarSlice(parsed.ElseBlock, desugarExpressionBox),
			}
		},
		func(generics *parser.LambdaOrListGenerics, parsed parser.List) {
			listGenerics := []TypeAnnotation{}
			if generics != nil {
				listGenerics = desugarSlice(generics.Generics, desugarTypeAnnotation)
			} else {
				listGenerics = nil
			}
			result = List{
				Node:        desugarNode(parsed.Node),
				Generics:    listGenerics,
				Expressions: desugarSlice(parsed.Expressions, desugarExpressionBox),
			}
		},
		func(parsed parser.When) {
			result = When{
				Node:  desugarNode(parsed.Node),
				Over:  desugarExpressionBox(parsed.Over),
				Is:    desugarSlice(parsed.Is, desugarWhenIs),
				Other: desugarWhenNonNil(parsed.Other, desugarWhenOther),
			}
		},
	)
	return result
}

func desugarWhenOther(parsed parser.WhenOther) WhenOther {
	return WhenOther{
		Node:      desugarNode(parsed.Node),
		Name:      desugarWhenNonNil(parsed.Name, desugarName),
		ThenBlock: desugarSlice(parsed.ThenBlock, desugarExpressionBox),
	}
}

func desugarWhenIs(parsed parser.WhenIs) WhenIs {
	return WhenIs{
		Node:      desugarNode(parsed.Node),
		Name:      desugarWhenNonNil(parsed.Name, desugarName),
		Type:      desugarTypeAnnotation(parsed.Type),
		ThenBlock: desugarSlice(parsed.ThenBlock, desugarExpressionBox),
	}
}

func desugarIfThen(parsed parser.IfThen) IfThen {
	return IfThen{
		Node:      desugarNode(parsed.Node),
		Condition: desugarExpressionBox(parsed.Condition),
		ThenBlock: desugarSlice(parsed.ThenBlock, desugarExpressionBox),
	}
}

func desugarLambdaSignature(parsed parser.LambdaSignature) LambdaSignature {
	return LambdaSignature{
		Node:       desugarNode(parsed.Node),
		Parameters: desugarSlice(parsed.Parameters, desugarParameter),
		ReturnType: desugarWhenNonNil(parsed.ReturnType, desugarTypeAnnotation),
	}
}

func desugarParameter(parsed parser.Parameter) Parameter {
	return Parameter{
		Name: desugarName(parsed.Name),
		Type: desugarWhenNonNil(parsed.Type, desugarTypeAnnotation),
	}
}

func desugarDeclarationShortCircuit(parsed parser.DeclarationShortCircuit) DeclarationShortCircuit {
	return DeclarationShortCircuit{
		TypeAnnotation: desugarWhenNonNil(parsed.TypeAnnotation, desugarTypeAnnotation),
	}
}

func desugarTypeAnnotation(parsed parser.TypeAnnotation) TypeAnnotation {
	return TypeAnnotation{
		Node:    desugarNode(parsed.Node),
		OrTypes: desugarSlice(parsed.OrTypes, desugarTypeAnnotationElement),
	}
}

func desugarTypeAnnotationElement(parsed parser.TypeAnnotationElement) TypeAnnotationElement {
	var result TypeAnnotationElement
	parser.TypeAnnotationElementExhaustiveSwitch(
		parsed,
		func(parsed parser.SingleNameType) {
			result = SingleNameType{
				Node:     desugarNode(parsed.Node),
				TypeName: desugarName(parsed.TypeName),
				Generics: desugarSlice(parsed.Generics, desugarTypeAnnotation),
			}
		},
		func(parsed parser.SingleNameType) {
			result = SingleNameType{
				Node:     desugarNode(parsed.Node),
				TypeName: desugarName(parsed.TypeName),
				Generics: desugarSlice(parsed.Generics, desugarTypeAnnotation),
			}
		},
		func(parsed parser.FunctionType) {
			result = DesugarFunctionType(parsed)
		},
	)
	return result
}

func desugarFunctionTypeArgument(parsed parser.FunctionTypeArgument) FunctionTypeArgument {
	return FunctionTypeArgument{
		Name: desugarWhenNonNil(parsed.Name, desugarName),
		Type: desugarTypeAnnotation(parsed.Type),
	}
}

func desugarStruct(parsed parser.Struct) Struct {
	return Struct{
		Name:      desugarName(parsed.Name),
		Generics:  desugarSlice(parsed.Generics, desugarName),
		Variables: desugarSlice(parsed.Variables, desugarStructVariable),
	}
}

func desugarStructVariable(parsed parser.StructVariable) StructVariable {
	return StructVariable{
		Name: desugarName(parsed.Name),
		Type: desugarTypeAnnotation(parsed.Type),
	}
}

func desugarTypeAlias(parsed parser.TypeAlias) TypeAlias {
	return TypeAlias{
		Name:     desugarName(parsed.Name),
		Generics: desugarSlice(parsed.Generics, desugarName),
		Type:     desugarTypeAnnotation(parsed.Type),
	}
}

func DesugarFileTopLevel(file string, parsed FileTopLevel) (FileTopLevel, error) {
	var err error
	for i, topLevelDeclaration := range parsed.TopLevelDeclarations {
		TopLevelDeclarationExhaustiveSwitch(
			topLevelDeclaration,
			func(topLevelDeclaration Declaration) {
				if topLevelDeclaration.ShortCircuit != nil {
					err = errors.New("shortcircuit only allowed inside of functions")
					return
				}
				p, _, e := desugarExpressionBoxComplex(file, topLevelDeclaration.ExpressionBox, []ExpressionBox{})
				if e != nil {
					err = e
					return
				}
				topLevelDeclaration.ExpressionBox = p
				parsed.TopLevelDeclarations[i] = topLevelDeclaration
			},
			func(topLevelDeclaration Struct) {},
			func(topLevelDeclaration TypeAlias) {},
		)
	}
	return parsed, err
}

func desugarExpressionBoxComplex(file string, parsed ExpressionBox, restOfBlock []ExpressionBox) (ExpressionBox, []ExpressionBox, error) {
	exp, restOfBlock, err := desugarExpressionComplex(file, parsed.Expression, restOfBlock)
	if err != nil {
		return parsed, restOfBlock, err
	}
	parsed.Expression = exp
	for i, accessOrInvocation := range parsed.AccessOrInvocationChain {
		if accessOrInvocation.Arguments != nil {
			for i2, argument := range accessOrInvocation.Arguments.Arguments {
				d, _, err := desugarExpressionBoxComplex(file, argument.Argument, []ExpressionBox{})
				if err != nil {
					return parsed, restOfBlock, err
				}
				parsed.AccessOrInvocationChain[i].Arguments.Arguments[i2].Argument = d
			}
		}
	}
	for i, accessOrInvocation := range parsed.AccessOrInvocationChain {
		if accessOrInvocation.DotOrArrowName != nil && accessOrInvocation.DotOrArrowName.Arrow {
			expressionBeforeThisArrow := ExpressionBox{
				Node:                    parsed.Node,
				Expression:              parsed.Expression,
				AccessOrInvocationChain: []AccessOrInvocation{},
			}
			if i > 0 {
				expressionBeforeThisArrow.AccessOrInvocationChain = parsed.AccessOrInvocationChain[0:i]
			}
			if accessOrInvocation.Arguments == nil {
				return ExpressionBox{}, nil, errors.New("Arrow syntax requires parenthesis on the right-hand side")
			}
			newParsedExpression := ReferenceOrInvocation{
				Var: accessOrInvocation.DotOrArrowName.VarName,
				Arguments: &ArgumentsList{
					Node:     accessOrInvocation.Node,
					Generics: accessOrInvocation.Arguments.Generics,
					Arguments: append([]NamedArgument{
						NamedArgument{
							Node:     expressionBeforeThisArrow.Node,
							Argument: expressionBeforeThisArrow,
						},
					}, accessOrInvocation.Arguments.Arguments...),
				},
			}
			for i, argument := range newParsedExpression.Arguments.Arguments {
				desugared, restOfBlock, err := desugarExpressionBoxComplex(file, argument.Argument, nil)
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
				parsed.AccessOrInvocationChain = []AccessOrInvocation{}
			}
			return desugarExpressionBoxComplex(file, parsed, restOfBlock)
		}
	}
	return parsed, restOfBlock, nil
}

func desugarExpressionComplex(file string, parsed Expression, restOfBlock []ExpressionBox) (Expression, []ExpressionBox, error) {
	var err error
	ExpressionExhaustiveSwitch(
		parsed,
		func(expression LiteralExpression) {

		},
		func(expression ReferenceOrInvocation) {
			if expression.Arguments != nil {
				for i, argument := range expression.Arguments.Arguments {
					d, _, e := desugarExpressionBoxComplex(file, argument.Argument, []ExpressionBox{})
					err = e
					if err != nil {
						return
					}
					expression.Arguments.Arguments[i].Argument = d
				}
			}
			parsed = expression
		},
		func(expression Lambda) {
			d, e := desugarBlock(file, expression.Block)
			err = e
			if err != nil {
				return
			}
			expression.Block = d
			parsed = expression
		},
		func(expression Declaration) {
			d, _, e := desugarExpressionBoxComplex(file, expression.ExpressionBox, []ExpressionBox{})
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
					name := Name{
						Node:   expression.Name.Node,
						String: expression.Name.String,
					}
					if name.String == "_" {
						name.String = "_unused_"
					}
					parsed = When{
						Node: expression.Name.Node,
						Over: expression.ExpressionBox,
						Is: []WhenIs{
							WhenIs{
								Node: expression.ShortCircuit.TypeAnnotation.Node,
								Name: &name,
								Type: *expression.ShortCircuit.TypeAnnotation,
								ThenBlock: []ExpressionBox{
									ExpressionBox{
										Node: expression.ShortCircuit.TypeAnnotation.Node,
										Expression: ReferenceOrInvocation{
											Var:       name,
											Arguments: nil,
										},
										AccessOrInvocationChain: []AccessOrInvocation{},
									},
								},
							},
							WhenIs{
								Node:      expression.TypeAnnotation.Node,
								Name:      &name,
								Type:      *expression.TypeAnnotation,
								ThenBlock: restOfBlock,
							},
						},
						Other: nil,
					}
				} else if expression.ShortCircuit.TypeAnnotation != nil {
					name := Name{
						Node:   expression.Name.Node,
						String: expression.Name.String,
					}
					if name.String == "_" {
						name.String = "_unused_"
					}
					parsed = When{
						Node: expression.Name.Node,
						Over: expression.ExpressionBox,
						Is: []WhenIs{
							WhenIs{
								Node: expression.ShortCircuit.TypeAnnotation.Node,
								Name: &name,
								Type: *expression.ShortCircuit.TypeAnnotation,
								ThenBlock: []ExpressionBox{
									ExpressionBox{
										Node: expression.ShortCircuit.TypeAnnotation.Node,
										Expression: ReferenceOrInvocation{
											Var:       name,
											Arguments: nil,
										},
										AccessOrInvocationChain: []AccessOrInvocation{},
									},
								},
							},
						},
						Other: &WhenOther{
							Node:      expression.Name.Node,
							Name:      &name,
							ThenBlock: restOfBlock,
						},
					}
				} else {
					name := Name{
						Node:   expression.Name.Node,
						String: expression.Name.String,
					}
					if name.String == "_" {
						name.String = "_unused_"
					}
					parsed = When{
						Node: expression.Name.Node,
						Over: expression.ExpressionBox,
						Is: []WhenIs{
							WhenIs{
								Node:      expression.TypeAnnotation.Node,
								Name:      &name,
								Type:      *expression.TypeAnnotation,
								ThenBlock: restOfBlock,
							},
						},
						Other: &WhenOther{
							Node: expression.TypeAnnotation.Node,
							Name: &name,
							ThenBlock: []ExpressionBox{
								ExpressionBox{
									Node: expression.TypeAnnotation.Node,
									Expression: ReferenceOrInvocation{
										Var:       name,
										Arguments: nil,
									},
									AccessOrInvocationChain: []AccessOrInvocation{},
								},
							},
						},
					}
				}
				restOfBlock = []ExpressionBox{}
			}
		},
		func(expression If) {
			cond, _, e := desugarExpressionBoxComplex(file, expression.Condition, []ExpressionBox{})
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
				cond, _, e := desugarExpressionBoxComplex(file, elseIf.Condition, []ExpressionBox{})
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
		func(expression List) {
			for i, expressionBox := range expression.Expressions {
				d, _, e := desugarExpressionBoxComplex(file, expressionBox, []ExpressionBox{})
				err = e
				if err != nil {
					return
				}
				expression.Expressions[i] = d
			}
			parsed = expression
		},
		func(expression When) {
			over, _, e := desugarExpressionBoxComplex(file, expression.Over, []ExpressionBox{})
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

func desugarBlock(file string, block []ExpressionBox) ([]ExpressionBox, error) {
	for i := len(block) - 1; i >= 0; i-- {
		expressionBox := block[i]
		d, r, err := desugarExpressionBoxComplex(file, expressionBox, block[i+1:len(block)])
		if err != nil {
			return nil, err
		}
		if i > 0 {
			block = append(append(block[0:i], d), r...)
		} else {
			block = append([]ExpressionBox{d}, r...)
		}
	}
	return block, nil
}
