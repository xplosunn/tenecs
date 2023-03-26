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
Package = "package" Name .
Name = <ident> .
Import = "import" (Name ("." Name)*)? .
TopLevelDeclaration = Struct | Interface | Declaration .
Struct = "struct" Name ("<" (Name ("," Name)*)? ">")? "(" (StructVariable ("," StructVariable)*)? ")" .
StructVariable = Name ":" TypeAnnotation .
TypeAnnotation = TypeAnnotationElement ("|" TypeAnnotationElement)* .
TypeAnnotationElement = SingleNameType | FunctionType .
SingleNameType = Name ("<" Name ("," Name)* ">")? .
FunctionType = ("<" Name ("," Name)* ">")? "(" (TypeAnnotation ("," TypeAnnotation)*)? ")" "-" ">" TypeAnnotation .
Interface = "interface" Name "{" InterfaceVariable* "}" .
InterfaceVariable = "public" Name ":" TypeAnnotation .
Declaration = Name ":" "=" ExpressionBox .
ExpressionBox = Expression AccessOrInvocation* .
Expression = Module | If | Declaration | LiteralExpression | ReferenceOrInvocation | Lambda | Array .
Module = "implement" Name "{" ModuleDeclaration* "}" .
ModuleDeclaration = "public"? Name ":" "=" Expression .
If = "if" ExpressionBox "{" ExpressionBox* "}" ("else" "{" ExpressionBox* "}")? .
LiteralExpression = Literal .
Literal = LiteralFloat | LiteralInt | LiteralString | LiteralBool .
LiteralFloat = <float> .
LiteralInt = <int> .
LiteralString = <string> .
LiteralBool = "true" | "false" .
ReferenceOrInvocation = Name ArgumentsList? .
ArgumentsList = ("<" TypeAnnotation ("," TypeAnnotation)* ">")? "(" (ExpressionBox ("," ExpressionBox)*)? ")" .
Lambda = ("<" Name ("," Name)* ">")? "(" (Parameter ("," Parameter)*)? ")" (":" TypeAnnotation)? "=" ">" (("{" ExpressionBox* "}") | ExpressionBox) .
Parameter = Name (":" TypeAnnotation)? .
Array = "[" TypeAnnotation? "]" "(" (ExpressionBox ("," ExpressionBox)*)? ")" .
AccessOrInvocation = "." Name ArgumentsList? .`
	grammar, err := parser.Grammar()
	assert.NoError(t, err)
	assert.Equal(t, expected, grammar)
}
