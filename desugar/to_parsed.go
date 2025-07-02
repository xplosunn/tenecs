package desugar

import "github.com/xplosunn/tenecs/parser"

func ToParsed(desugared FileTopLevel) parser.FileTopLevel {
	return parser.FileTopLevel{
		Tokens:               desugared.Tokens,
		Package:              toParsedPackage(desugared.Package),
		Imports:              toParsedSlice(desugared.Imports, toParsedImport),
		TopLevelDeclarations: toParsedSlice(desugared.TopLevelDeclarations, toParsedTopLevelDeclaration),
	}
}

func toParsedFunctionType(desugared FunctionType) parser.FunctionType {
	return parser.FunctionType{
		Generics:   toParsedSlice(desugared.Generics, toParsedName),
		Arguments:  toParsedSlice(desugared.Arguments, toParsedFunctionTypeArgument),
		ReturnType: toParsedTypeAnnotation(desugared.ReturnType),
	}
}

func toParsedSlice[In any, Out any](in []In, toParsed func(In) Out) []Out {
	result := []Out{}
	for _, in := range in {
		result = append(result, toParsed(in))
	}
	return result
}

func toParsedWhenNonNil[In any, Out any](ptr *In, toParsed func(In) Out) *Out {
	if ptr == nil {
		return nil
	}
	result := toParsed(*ptr)
	return &result
}

func toParsedPackage(desugared Package) parser.Package {
	return parser.Package{
		Node:              desugared.Node,
		DotSeparatedNames: toParsedSlice(desugared.DotSeparatedNames, toParsedName),
	}
}

func toParsedName(desugared Name) parser.Name {
	return parser.Name{
		Node:   desugared.Node,
		String: desugared.String,
	}
}

func toParsedImport(desugared Import) parser.Import {
	return parser.Import{
		Node:             desugared.Node,
		DotSeparatedVars: toParsedSlice(desugared.DotSeparatedVars, toParsedName),
		As:               toParsedWhenNonNil(desugared.As, toParsedName),
	}
}

func toParsedTopLevelDeclaration(desugared TopLevelDeclaration) parser.TopLevelDeclaration {
	var result parser.TopLevelDeclaration
	TopLevelDeclarationExhaustiveSwitch(
		desugared,
		func(desugared Declaration) {
			result = toParsedDeclaration(desugared)
		},
		func(desugared Struct) {
			result = toParsedStruct(desugared)
		},
		func(desugared TypeAlias) {
			result = toParsedTypeAlias(desugared)
		},
	)
	return result
}

func toParsedDeclaration(desugared Declaration) parser.Declaration {
	return parser.Declaration{
		Name:           toParsedName(desugared.Name),
		TypeAnnotation: toParsedWhenNonNil(desugared.TypeAnnotation, toParsedTypeAnnotation),
		ShortCircuit:   toParsedWhenNonNil(desugared.ShortCircuit, toParsedDeclarationShortCircuit),
		ExpressionBox:  toParsedExpressionBox(desugared.ExpressionBox),
	}
}

func toParsedExpressionBox(desugared ExpressionBox) parser.ExpressionBox {
	return parser.ExpressionBox{
		Node:                    desugared.Node,
		Expression:              toParsedExpression(desugared.Expression),
		AccessOrInvocationChain: toParsedSlice(desugared.AccessOrInvocationChain, toParsedAccessOrInvocation),
	}
}

func toParsedAccessOrInvocation(desugared AccessOrInvocation) parser.AccessOrInvocation {
	return parser.AccessOrInvocation{
		Node:           desugared.Node,
		DotOrArrowName: toParsedWhenNonNil(desugared.DotOrArrowName, toParsedDotOrArrowName),
		Arguments:      toParsedWhenNonNil(desugared.Arguments, toParsedArgumentsList),
	}
}

func toParsedArgumentsList(desugared ArgumentsList) parser.ArgumentsList {
	return parser.ArgumentsList{
		Node:      desugared.Node,
		Generics:  toParsedSlice(desugared.Generics, toParsedTypeAnnotation),
		Arguments: toParsedSlice(desugared.Arguments, toParsedNamedArgument),
	}
}

func toParsedNamedArgument(desugared NamedArgument) parser.NamedArgument {
	return parser.NamedArgument{
		Node:     desugared.Node,
		Name:     toParsedWhenNonNil(desugared.Name, toParsedName),
		Argument: toParsedExpressionBox(desugared.Argument),
	}
}

func toParsedDotOrArrowName(desugared DotOrArrowName) parser.DotOrArrowName {
	return parser.DotOrArrowName{
		Node:    desugared.Node,
		Dot:     desugared.Dot,
		Arrow:   desugared.Arrow,
		VarName: toParsedName(desugared.VarName),
	}
}

func toParsedExpression(desugared Expression) parser.Expression {
	var result parser.Expression
	ExpressionExhaustiveSwitch(
		desugared,
		func(desugared LiteralExpression) {
			result = parser.LiteralExpression{
				Node:    desugared.Node,
				Literal: desugared.Literal,
			}
		},
		func(desugared ReferenceOrInvocation) {
			result = parser.ReferenceOrInvocation{
				Var:       toParsedName(desugared.Var),
				Arguments: toParsedWhenNonNil(desugared.Arguments, toParsedArgumentsList),
			}
		},
		func(desugared Lambda) {
			lambda := parser.Lambda{
				Node:      desugared.Node,
				Signature: toParsedLambdaSignature(desugared.Signature),
				Block:     toParsedSlice(desugared.Block, toParsedExpressionBox),
			}
			genericTypeAnnotations := toParsedSlice(desugared.Generics, toParsedTypeAnnotation)
			var generics *parser.LambdaOrListGenerics
			if len(genericTypeAnnotations) > 0 {
				generics = &parser.LambdaOrListGenerics{
					Node:     lambda.Node,
					Generics: genericTypeAnnotations,
				}
			}
			result = parser.LambdaOrList{
				Node:     lambda.Node,
				Generics: generics,
				List:     nil,
				Lambda:   &lambda,
			}
		},
		func(desugared Declaration) {
			result = parser.Declaration{
				Name:           toParsedName(desugared.Name),
				TypeAnnotation: toParsedWhenNonNil(desugared.TypeAnnotation, toParsedTypeAnnotation),
				ShortCircuit:   toParsedWhenNonNil(desugared.ShortCircuit, toParsedDeclarationShortCircuit),
				ExpressionBox:  toParsedExpressionBox(desugared.ExpressionBox),
			}
		},
		func(desugared If) {
			result = parser.If{
				Node:      desugared.Node,
				Condition: toParsedExpressionBox(desugared.Condition),
				ThenBlock: toParsedSlice(desugared.ThenBlock, toParsedExpressionBox),
				ElseIfs:   toParsedSlice(desugared.ElseIfs, toParsedIfThen),
				ElseBlock: toParsedSlice(desugared.ElseBlock, toParsedExpressionBox),
			}
		},
		func(desugared List) {
			list := parser.List{
				Node:        desugared.Node,
				Expressions: toParsedSlice(desugared.Expressions, toParsedExpressionBox),
			}
			genericTypeAnnotations := toParsedSlice(desugared.Generics, toParsedTypeAnnotation)
			var generics *parser.LambdaOrListGenerics
			if len(genericTypeAnnotations) > 0 {
				generics = &parser.LambdaOrListGenerics{
					Node:     list.Node,
					Generics: genericTypeAnnotations,
				}
			}
			result = parser.LambdaOrList{
				Node:     list.Node,
				Generics: generics,
				List:     &list,
				Lambda:   nil,
			}
		},
		func(desugared When) {
			result = parser.When{
				Node:  desugared.Node,
				Over:  toParsedExpressionBox(desugared.Over),
				Is:    toParsedSlice(desugared.Is, toParsedWhenIs),
				Other: toParsedWhenNonNil(desugared.Other, toParsedWhenOther),
			}
		},
	)
	return result
}

func toParsedWhenOther(desugared WhenOther) parser.WhenOther {
	return parser.WhenOther{
		Node:      desugared.Node,
		Name:      toParsedWhenNonNil(desugared.Name, toParsedName),
		ThenBlock: toParsedSlice(desugared.ThenBlock, toParsedExpressionBox),
	}
}

func toParsedWhenIs(desugared WhenIs) parser.WhenIs {
	return parser.WhenIs{
		Node:      desugared.Node,
		Name:      toParsedWhenNonNil(desugared.Name, toParsedName),
		Type:      toParsedTypeAnnotation(desugared.Type),
		ThenBlock: toParsedSlice(desugared.ThenBlock, toParsedExpressionBox),
	}
}

func toParsedIfThen(desugared IfThen) parser.IfThen {
	return parser.IfThen{
		Node:      desugared.Node,
		Condition: toParsedExpressionBox(desugared.Condition),
		ThenBlock: toParsedSlice(desugared.ThenBlock, toParsedExpressionBox),
	}
}

func toParsedLambdaSignature(desugared LambdaSignature) parser.LambdaSignature {
	return parser.LambdaSignature{
		Node:       desugared.Node,
		Parameters: toParsedSlice(desugared.Parameters, toParsedParameter),
		ReturnType: toParsedWhenNonNil(desugared.ReturnType, toParsedTypeAnnotation),
	}
}

func toParsedParameter(desugared Parameter) parser.Parameter {
	return parser.Parameter{
		Name: toParsedName(desugared.Name),
		Type: toParsedWhenNonNil(desugared.Type, toParsedTypeAnnotation),
	}
}

func toParsedDeclarationShortCircuit(desugared DeclarationShortCircuit) parser.DeclarationShortCircuit {
	return parser.DeclarationShortCircuit{
		TypeAnnotation: toParsedWhenNonNil(desugared.TypeAnnotation, toParsedTypeAnnotation),
	}
}

func toParsedTypeAnnotation(desugared TypeAnnotation) parser.TypeAnnotation {
	return parser.TypeAnnotation{
		Node:    desugared.Node,
		OrTypes: toParsedSlice(desugared.OrTypes, toParsedTypeAnnotationElement),
	}
}

func toParsedTypeAnnotationElement(desugared TypeAnnotationElement) parser.TypeAnnotationElement {
	var result parser.TypeAnnotationElement
	TypeAnnotationElementExhaustiveSwitch(
		desugared,
		func(desugared SingleNameType) {
			result = parser.SingleNameType{
				Node:     desugared.Node,
				TypeName: toParsedName(desugared.TypeName),
				Generics: toParsedSlice(desugared.Generics, toParsedTypeAnnotation),
			}
		},
		func(desugared SingleNameType) {
			result = parser.SingleNameType{
				Node:     desugared.Node,
				TypeName: toParsedName(desugared.TypeName),
				Generics: toParsedSlice(desugared.Generics, toParsedTypeAnnotation),
			}
		},
		func(desugared FunctionType) {
			result = toParsedFunctionType(desugared)
		},
	)
	return result
}

func toParsedFunctionTypeArgument(desugared FunctionTypeArgument) parser.FunctionTypeArgument {
	return parser.FunctionTypeArgument{
		Name: toParsedWhenNonNil(desugared.Name, toParsedName),
		Type: toParsedTypeAnnotation(desugared.Type),
	}
}

func toParsedStruct(desugared Struct) parser.Struct {
	return parser.Struct{
		Name:      toParsedName(desugared.Name),
		Generics:  toParsedSlice(desugared.Generics, toParsedName),
		Variables: toParsedSlice(desugared.Variables, toParsedStructVariable),
	}
}

func toParsedStructVariable(desugared StructVariable) parser.StructVariable {
	return parser.StructVariable{
		Name: toParsedName(desugared.Name),
		Type: toParsedTypeAnnotation(desugared.Type),
	}
}

func toParsedTypeAlias(desugared TypeAlias) parser.TypeAlias {
	return parser.TypeAlias{
		Name:     toParsedName(desugared.Name),
		Generics: toParsedSlice(desugared.Generics, toParsedName),
		Type:     toParsedTypeAnnotation(desugared.Type),
	}
}
