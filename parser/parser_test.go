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

func TestParserGrammar(t *testing.T) {
	expected := `FileTopLevel = Package Import* TopLevelDeclaration* .
Package = "package" <ident> .
Import = "import" (<ident> ("." <ident>)*)? .
TopLevelDeclaration = Struct | Interface | Declaration .
Struct = "struct" <ident> ("<" (<ident> ("," <ident>)*)? ">")? "(" (StructVariable ("," StructVariable)*)? ")" .
StructVariable = <ident> ":" TypeAnnotation .
TypeAnnotation = SingleNameType | FunctionType .
SingleNameType = <ident> .
FunctionType = ("<" <ident> ("," <ident>)* ">")? "(" (TypeAnnotation ("," TypeAnnotation)*)? ")" "-" ">" TypeAnnotation .
Interface = "interface" <ident> "{" InterfaceVariable* "}" .
InterfaceVariable = "public" <ident> ":" TypeAnnotation .
Declaration = <ident> ":" "=" ExpressionBox .
ExpressionBox = Expression AccessOrInvocation* .
Expression = Module | If | Declaration | LiteralExpression | ReferenceOrInvocation | Lambda .
Module = "implement" <ident> "{" ModuleDeclaration* "}" .
ModuleDeclaration = "public"? <ident> ":" "=" Expression .
If = "if" ExpressionBox "{" ExpressionBox* "}" ("else" "{" ExpressionBox* "}")? .
LiteralExpression = Literal .
Literal = LiteralFloat | LiteralInt | LiteralString | LiteralBool .
LiteralFloat = <float> .
LiteralInt = <int> .
LiteralString = <string> .
LiteralBool = ("true" | "false") .
ReferenceOrInvocation = <ident> ArgumentsList? .
ArgumentsList = ("<" <ident> ("," <ident>)* ">")? "(" (ExpressionBox ("," ExpressionBox)*)? ")" .
Lambda = ("<" <ident> ("," <ident>)* ">")? "(" (Parameter ("," Parameter)*)? ")" (":" TypeAnnotation)? "=" ">" (("{" ExpressionBox* "}") | ExpressionBox) .
Parameter = <ident> (":" TypeAnnotation)? .
AccessOrInvocation = "." <ident> ArgumentsList? .`
	grammar, err := parser.ParserGrammar()
	assert.NoError(t, err)
	assert.Equal(t, expected, grammar)
}
