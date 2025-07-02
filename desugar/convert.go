package desugar

import (
	"github.com/xplosunn/tenecs/parser"
)

func ConvertFunctionType(parsed parser.FunctionType) FunctionType {
	return FunctionType{
		Generics:   convertSlice(parsed.Generics, convertName),
		Arguments:  convertSlice(parsed.Arguments, convertFunctionTypeArgument),
		ReturnType: convertTypeAnnotation(parsed.ReturnType),
	}
}

func convertSlice[In any, Out any](in []In, desugar func(In) Out) []Out {
	result := []Out{}
	for _, in := range in {
		result = append(result, desugar(in))
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func convertWhenNonNil[In any, Out any](ptr *In, desugar func(In) Out) *Out {
	if ptr == nil {
		return nil
	}
	result := desugar(*ptr)
	return &result
}

func convertPackage(parsed parser.Package) Package {
	return Package{
		Node:              parsed.Node,
		DotSeparatedNames: convertSlice(parsed.DotSeparatedNames, convertName),
	}
}

func convertName(parsed parser.Name) Name {
	return Name{
		Node:   parsed.Node,
		String: parsed.String,
	}
}

func convertImport(parsed parser.Import) Import {
	return Import{
		Node:             parsed.Node,
		DotSeparatedVars: convertSlice(parsed.DotSeparatedVars, convertName),
		As:               convertWhenNonNil(parsed.As, convertName),
	}
}

func convertTopLevelDeclaration(parsed parser.TopLevelDeclaration) TopLevelDeclaration {
	var result TopLevelDeclaration
	parser.TopLevelDeclarationExhaustiveSwitch(
		parsed,
		func(parsed parser.Declaration) {
			result = convertDeclaration(parsed)
		},
		func(parsed parser.Struct) {
			result = convertStruct(parsed)
		},
		func(parsed parser.TypeAlias) {
			result = convertTypeAlias(parsed)
		},
	)
	return result
}

func convertDeclaration(parsed parser.Declaration) Declaration {
	return Declaration{
		Name:           convertName(parsed.Name),
		TypeAnnotation: convertWhenNonNil(parsed.TypeAnnotation, convertTypeAnnotation),
		ShortCircuit:   convertWhenNonNil(parsed.ShortCircuit, convertDeclarationShortCircuit),
		ExpressionBox:  convertExpressionBox(parsed.ExpressionBox),
	}
}

func convertExpressionBox(parsed parser.ExpressionBox) ExpressionBox {
	return ExpressionBox{
		Node:                    parsed.Node,
		Expression:              convertExpression(parsed.Expression),
		AccessOrInvocationChain: convertSlice(parsed.AccessOrInvocationChain, convertAccessOrInvocation),
	}
}

func convertAccessOrInvocation(parsed parser.AccessOrInvocation) AccessOrInvocation {
	return AccessOrInvocation{
		Node:           parsed.Node,
		DotOrArrowName: convertWhenNonNil(parsed.DotOrArrowName, convertDotOrArrowName),
		Arguments:      convertWhenNonNil(parsed.Arguments, convertArgumentsList),
	}
}

func convertArgumentsList(parsed parser.ArgumentsList) ArgumentsList {
	return ArgumentsList{
		Node:      parsed.Node,
		Generics:  convertSlice(parsed.Generics, convertTypeAnnotation),
		Arguments: convertSlice(parsed.Arguments, convertNamedArgument),
	}
}

func convertNamedArgument(parsed parser.NamedArgument) NamedArgument {
	return NamedArgument{
		Node:     parsed.Node,
		Name:     convertWhenNonNil(parsed.Name, convertName),
		Argument: convertExpressionBox(parsed.Argument),
	}
}

func convertDotOrArrowName(parsed parser.DotOrArrowName) DotOrArrowName {
	return DotOrArrowName{
		Node:    parsed.Node,
		Dot:     parsed.Dot,
		Arrow:   parsed.Arrow,
		VarName: convertName(parsed.VarName),
	}
}

func convertExpression(parsed parser.Expression) Expression {
	var result Expression
	parser.ExpressionExhaustiveSwitch(
		parsed,
		func(parsed parser.LiteralExpression) {
			result = LiteralExpression{
				Node:    parsed.Node,
				Literal: parsed.Literal,
			}
		},
		func(parsed parser.ReferenceOrInvocation) {
			result = ReferenceOrInvocation{
				Var:       convertName(parsed.Var),
				Arguments: convertWhenNonNil(parsed.Arguments, convertArgumentsList),
			}
		},
		func(generics *parser.LambdaOrListGenerics, parsed parser.Lambda) {
			lambdaGenerics := []TypeAnnotation{}
			if generics != nil {
				lambdaGenerics = convertSlice(generics.Generics, convertTypeAnnotation)
			} else {
				lambdaGenerics = nil
			}
			result = Lambda{
				Node:      parsed.Node,
				Generics:  lambdaGenerics,
				Signature: convertLambdaSignature(parsed.Signature),
				Block:     convertSlice(parsed.Block, convertExpressionBox),
			}
		},
		func(parsed parser.Declaration) {
			result = Declaration{
				Name:           convertName(parsed.Name),
				TypeAnnotation: convertWhenNonNil(parsed.TypeAnnotation, convertTypeAnnotation),
				ShortCircuit:   convertWhenNonNil(parsed.ShortCircuit, convertDeclarationShortCircuit),
				ExpressionBox:  convertExpressionBox(parsed.ExpressionBox),
			}
		},
		func(parsed parser.If) {
			result = If{
				Node:      parsed.Node,
				Condition: convertExpressionBox(parsed.Condition),
				ThenBlock: convertSlice(parsed.ThenBlock, convertExpressionBox),
				ElseIfs:   convertSlice(parsed.ElseIfs, convertIfThen),
				ElseBlock: convertSlice(parsed.ElseBlock, convertExpressionBox),
			}
		},
		func(generics *parser.LambdaOrListGenerics, parsed parser.List) {
			listGenerics := []TypeAnnotation{}
			if generics != nil {
				listGenerics = convertSlice(generics.Generics, convertTypeAnnotation)
			} else {
				listGenerics = nil
			}
			result = List{
				Node:        parsed.Node,
				Generics:    listGenerics,
				Expressions: convertSlice(parsed.Expressions, convertExpressionBox),
			}
		},
		func(parsed parser.When) {
			result = When{
				Node:  parsed.Node,
				Over:  convertExpressionBox(parsed.Over),
				Is:    convertSlice(parsed.Is, convertWhenIs),
				Other: convertWhenNonNil(parsed.Other, convertWhenOther),
			}
		},
	)
	return result
}

func convertWhenOther(parsed parser.WhenOther) WhenOther {
	return WhenOther{
		Node:      parsed.Node,
		Name:      convertWhenNonNil(parsed.Name, convertName),
		ThenBlock: convertSlice(parsed.ThenBlock, convertExpressionBox),
	}
}

func convertWhenIs(parsed parser.WhenIs) WhenIs {
	return WhenIs{
		Node:      parsed.Node,
		Name:      convertWhenNonNil(parsed.Name, convertName),
		Type:      convertTypeAnnotation(parsed.Type),
		ThenBlock: convertSlice(parsed.ThenBlock, convertExpressionBox),
	}
}

func convertIfThen(parsed parser.IfThen) IfThen {
	return IfThen{
		Node:      parsed.Node,
		Condition: convertExpressionBox(parsed.Condition),
		ThenBlock: convertSlice(parsed.ThenBlock, convertExpressionBox),
	}
}

func convertLambdaSignature(parsed parser.LambdaSignature) LambdaSignature {
	return LambdaSignature{
		Node:       parsed.Node,
		Parameters: convertSlice(parsed.Parameters, convertParameter),
		ReturnType: convertWhenNonNil(parsed.ReturnType, convertTypeAnnotation),
	}
}

func convertParameter(parsed parser.Parameter) Parameter {
	return Parameter{
		Name: convertName(parsed.Name),
		Type: convertWhenNonNil(parsed.Type, convertTypeAnnotation),
	}
}

func convertDeclarationShortCircuit(parsed parser.DeclarationShortCircuit) DeclarationShortCircuit {
	return DeclarationShortCircuit{
		TypeAnnotation: convertWhenNonNil(parsed.TypeAnnotation, convertTypeAnnotation),
	}
}

func convertTypeAnnotation(parsed parser.TypeAnnotation) TypeAnnotation {
	return TypeAnnotation{
		Node:    parsed.Node,
		OrTypes: convertSlice(parsed.OrTypes, convertTypeAnnotationElement),
	}
}

func convertTypeAnnotationElement(parsed parser.TypeAnnotationElement) TypeAnnotationElement {
	var result TypeAnnotationElement
	parser.TypeAnnotationElementExhaustiveSwitch(
		parsed,
		func(parsed parser.SingleNameType) {
			result = SingleNameType{
				Node:     parsed.Node,
				TypeName: convertName(parsed.TypeName),
				Generics: convertSlice(parsed.Generics, convertTypeAnnotation),
			}
		},
		func(parsed parser.SingleNameType) {
			result = SingleNameType{
				Node:     parsed.Node,
				TypeName: convertName(parsed.TypeName),
				Generics: convertSlice(parsed.Generics, convertTypeAnnotation),
			}
		},
		func(parsed parser.FunctionType) {
			result = ConvertFunctionType(parsed)
		},
	)
	return result
}

func convertFunctionTypeArgument(parsed parser.FunctionTypeArgument) FunctionTypeArgument {
	return FunctionTypeArgument{
		Name: convertWhenNonNil(parsed.Name, convertName),
		Type: convertTypeAnnotation(parsed.Type),
	}
}

func convertStruct(parsed parser.Struct) Struct {
	return Struct{
		Name:      convertName(parsed.Name),
		Generics:  convertSlice(parsed.Generics, convertName),
		Variables: convertSlice(parsed.Variables, convertStructVariable),
	}
}

func convertStructVariable(parsed parser.StructVariable) StructVariable {
	return StructVariable{
		Name: convertName(parsed.Name),
		Type: convertTypeAnnotation(parsed.Type),
	}
}

func convertTypeAlias(parsed parser.TypeAlias) TypeAlias {
	return TypeAlias{
		Name:     convertName(parsed.Name),
		Generics: convertSlice(parsed.Generics, convertName),
		Type:     convertTypeAnnotation(parsed.Type),
	}
}
