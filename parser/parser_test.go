package parser_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestParseString(t *testing.T) {
	testCases := testcode.GetAll()
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := parser.ParseString(testCase.Content)
			assert.NoError(t, err)
		})
	}
}

func TestParseSignatureString(t *testing.T) {
	testCases := []testcode.TestCode{
		{
			Name:    "Boolean and (unnamed)",
			Content: "(Boolean, () ~> Boolean) ~> Boolean",
		},
		{
			Name:    "Boolean and",
			Content: "(a: Boolean, b: () ~> Boolean) ~> Boolean",
		},
		{
			Name:    "Boolean not",
			Content: "(b: Boolean) ~> Boolean",
		},
		{
			Name:    "eq",
			Content: "<T>(first: T, second: T) ~> Boolean",
		},
		{
			Name:    "List mapNotNull",
			Content: "<A, B>(list: List<A>, f: (A) ~> B | Void) ~> List<B>",
		},
	}
	for _, _testCase := range testCases {
		testCase := _testCase
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := parser.ParseFunctionTypeString(testCase.Content)
			assert.NoError(t, err)
		})
	}
}

func TestParserGrammar(t *testing.T) {
	expected := `FileTopLevel = Package Import* TopLevelDeclaration* .
Package = "package" (Name ("." Name)*)? .
Name = <ident> .
Import = "import" (Name ("." Name)*)? ("as" Name)? .
TopLevelDeclaration = Struct | TypeAlias | Declaration .
Struct = "struct" Name ("<" (Name ("," Name)*)? ">")? "(" (StructVariable ("," StructVariable)*)? ")" .
StructVariable = Name ":" TypeAnnotation .
TypeAnnotation = TypeAnnotationElement ("|" TypeAnnotationElement)* .
TypeAnnotationElement = SingleNameType | FunctionType .
SingleNameType = Name ("<" TypeAnnotation ("," TypeAnnotation)* ">")? .
FunctionType = ("<" Name ("," Name)* ">")? "(" (FunctionTypeArgument ("," FunctionTypeArgument)*)? ")" "~" ">" TypeAnnotation .
FunctionTypeArgument = (Name ":")? TypeAnnotation .
TypeAlias = "typealias" Name ("<" (Name ("," Name)*)? ">")? "=" TypeAnnotation .
Declaration = Name ":" TypeAnnotation? DeclarationShortCircuit? "=" ExpressionBox .
DeclarationShortCircuit = "?" TypeAnnotation? .
ExpressionBox = Expression AccessOrInvocation* .
Expression = When | If | Declaration | LiteralExpression | ReferenceOrInvocation | LambdaOrList .
When = "when" ExpressionBox "{" WhenIs* WhenOther? "}" .
WhenIs = "is" (Name ":")? TypeAnnotation "=" ">" "{" ExpressionBox* "}" .
WhenOther = "other" Name? "=" ">" "{" ExpressionBox* "}" .
If = "if" ExpressionBox "{" ExpressionBox* "}" ("else" IfThen)* ("else" "{" ExpressionBox* "}")? .
IfThen = "if" ExpressionBox "{" ExpressionBox* "}" .
LiteralExpression = Literal .
Literal = LiteralFloat | LiteralInt | LiteralString | LiteralBool | LiteralNull .
LiteralFloat = <float> .
LiteralInt = "-"? <int> .
LiteralString = <string> .
LiteralBool = "true" | "false" .
LiteralNull = "null" .
ReferenceOrInvocation = Name ArgumentsList? .
ArgumentsList = ("<" TypeAnnotation ("," TypeAnnotation)* ">")? "(" (NamedArgument ("," NamedArgument)*)? ")" .
NamedArgument = (Name "=")? ExpressionBox .
LambdaOrList = LambdaOrListGenerics? (("[" List) | Lambda) .
LambdaOrListGenerics = "<" TypeAnnotation ("," TypeAnnotation)* ">" .
List = (ExpressionBox ("," ExpressionBox)*)? "]" .
Lambda = LambdaSignature "=" ">" (("{" ExpressionBox* "}") | ExpressionBox) .
LambdaSignature = "(" (Parameter ("," Parameter)*)? ")" (":" TypeAnnotation)? .
Parameter = Name (":" TypeAnnotation)? .
AccessOrInvocation = (DotOrArrowName ArgumentsList?) | ArgumentsList .
DotOrArrowName = ("." | ("-" ">")) Name .`
	grammar, err := parser.Grammar()
	assert.NoError(t, err)
	assert.Equal(t, expected, grammar)
}
