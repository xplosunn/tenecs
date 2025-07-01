package desugar

import (
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
			lambda := Lambda{
				Node:      desugarNode(parsed.Node),
				Signature: desugarLambdaSignature(parsed.Signature),
				Block:     desugarSlice(parsed.Block, desugarExpressionBox),
			}
			node := lambda.Node
			if generics != nil {
				node = desugarNode(generics.Node)
			}
			result = LambdaOrList{
				Node:     node,
				Generics: desugarWhenNonNil(generics, desugarLambdaOrListGenerics),
				List:     nil,
				Lambda:   &lambda,
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
			list := List{
				Node:        desugarNode(parsed.Node),
				Expressions: desugarSlice(parsed.Expressions, desugarExpressionBox),
			}
			node := list.Node
			if generics != nil {
				node = desugarNode(generics.Node)
			}
			result = LambdaOrList{
				Node:     node,
				Generics: desugarWhenNonNil(generics, desugarLambdaOrListGenerics),
				List:     &list,
				Lambda:   nil,
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

func desugarLambdaOrListGenerics(parsed parser.LambdaOrListGenerics) LambdaOrListGenerics {
	return LambdaOrListGenerics{
		Node:     desugarNode(parsed.Node),
		Generics: desugarSlice(parsed.Generics, desugarTypeAnnotation),
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
