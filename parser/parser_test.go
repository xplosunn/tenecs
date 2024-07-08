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
Package = "package" (Name ("." Name)*)? .
Name = <ident> .
Import = "import" (Name ("." Name)*)? ("as" Name)? .
TopLevelDeclaration = Struct | Interface | TypeAlias | Declaration .
Struct = "struct" Name ("<" (Name ("," Name)*)? ">")? "(" (StructVariable ("," StructVariable)*)? ")" .
StructVariable = Name ":" TypeAnnotation .
TypeAnnotation = TypeAnnotationElement ("|" TypeAnnotationElement)* .
TypeAnnotationElement = SingleNameType | FunctionType .
SingleNameType = Name ("<" TypeAnnotation ("," TypeAnnotation)* ">")? .
FunctionType = ("<" Name ("," Name)* ">")? "(" (TypeAnnotation ("," TypeAnnotation)*)? ")" "-" ">" TypeAnnotation .
Interface = "interface" Name ("<" (Name ("," Name)*)? ">")? "{" InterfaceVariable* "}" .
InterfaceVariable = Name ":" TypeAnnotation .
TypeAlias = "typealias" Name ("<" (Name ("," Name)*)? ">")? "=" TypeAnnotation .
Declaration = Name ":" TypeAnnotation? DeclarationShortCircuit? "=" ExpressionBox .
DeclarationShortCircuit = "?" TypeAnnotation? .
ExpressionBox = Expression AccessOrInvocation* .
Expression = When | Implementation | If | Declaration | LiteralExpression | ReferenceOrInvocation | Lambda | List .
When = "when" ExpressionBox "{" WhenIs* WhenOther? "}" .
WhenIs = "is" (Name ":")? TypeAnnotation "=" ">" "{" ExpressionBox* "}" .
WhenOther = "other" Name? "=" ">" "{" ExpressionBox* "}" .
Implementation = "implement" Name ("<" TypeAnnotation ("," TypeAnnotation)* ">")? "{" ImplementationDeclaration* "}" .
ImplementationDeclaration = Name ":" TypeAnnotation? "=" Expression .
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
Lambda = ("<" Name ("," Name)* ">")? "(" (Parameter ("," Parameter)*)? ")" (":" TypeAnnotation)? "=" ">" (("{" ExpressionBox* "}") | ExpressionBox) .
Parameter = Name (":" TypeAnnotation)? .
List = "[" TypeAnnotation? "]" "(" (ExpressionBox ("," ExpressionBox)*)? ")" .
AccessOrInvocation = ("." Name ArgumentsList?) | ArgumentsList .`
	grammar, err := parser.Grammar()
	assert.NoError(t, err)
	assert.Equal(t, expected, grammar)
}
