package parser

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func ParseString(s string) (*FileTopLevel, error) {
	p, err := parser()
	if err != nil {
		return nil, err
	}

	res, err := p.ParseString("", s)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Grammar() (string, error) {
	p, err := parser()
	if err != nil {
		return "", err
	}
	return p.String(), nil
}

func parser() (*participle.Parser[FileTopLevel], error) {
	return participle.Build[FileTopLevel](topLevelDeclarationUnion, typeAnnotationUnion, literalUnion, expressionUnion)
}

type Node struct {
	Pos    lexer.Position
	EndPos lexer.Position
}

type FileTopLevel struct {
	Package              Package               `@@`
	Imports              []Import              `@@*`
	TopLevelDeclarations []TopLevelDeclaration `@@*`
}

func FileTopLevelFields(node FileTopLevel) (Package, []Import, []TopLevelDeclaration) {
	return node.Package, node.Imports, node.TopLevelDeclarations
}

type Name struct {
	Node
	String string `@Ident`
}

type Package struct {
	Identifier Name `"package" @@`
}

type Import struct {
	Node
	DotSeparatedVars []Name `"import" (@@ ("." @@)*)?`
}

func ImportFields(node Import) []Name {
	return node.DotSeparatedVars
}

type TopLevelDeclaration interface {
	sealedTopLevelDeclaration()
	TopLevelDeclarationCases() (*Declaration, *Interface, *Struct)
}

var topLevelDeclarationUnion = participle.Union[TopLevelDeclaration](Struct{}, Interface{}, Declaration{})

type Struct struct {
	Name      Name             `"struct" @@`
	Generics  []Name           `("<" (@@ ("," @@)*)? ">")?`
	Variables []StructVariable `"(" (@@ ("," @@)*)? ")"`
}

func (s Struct) sealedTopLevelDeclaration() {}
func (s Struct) TopLevelDeclarationCases() (*Declaration, *Interface, *Struct) {
	return nil, nil, &s
}

func StructFields(struc Struct) (Name, []Name, []StructVariable) {
	return struc.Name, struc.Generics, struc.Variables
}

type StructVariable struct {
	Name Name           `@@`
	Type TypeAnnotation `":" @@`
}

func StructVariableFields(structVariable StructVariable) (Name, TypeAnnotation) {
	return structVariable.Name, structVariable.Type
}

type Interface struct {
	Name      Name                `"interface" @@`
	Variables []InterfaceVariable `"{" @@* "}"`
}

func (i Interface) sealedTopLevelDeclaration() {}
func (i Interface) TopLevelDeclarationCases() (*Declaration, *Interface, *Struct) {
	return nil, &i, nil
}

func InterfaceFields(interf Interface) (Name, []InterfaceVariable) {
	return interf.Name, interf.Variables
}

type InterfaceVariable struct {
	Name Name           `"public" @@`
	Type TypeAnnotation `":" @@`
}

func InterfaceVariableFields(interfaceVariable InterfaceVariable) (Name, TypeAnnotation) {
	return interfaceVariable.Name, interfaceVariable.Type
}

type TypeAnnotation interface {
	sealedTypeAnnotation()
	TypeAnnotationCases() (*SingleNameType, *FunctionType)
}

var typeAnnotationUnion = participle.Union[TypeAnnotation](SingleNameType{}, FunctionType{})

type SingleNameType struct {
	TypeName Name `@@`
}

func (s SingleNameType) sealedTypeAnnotation() {}
func (s SingleNameType) TypeAnnotationCases() (*SingleNameType, *FunctionType) {
	return &s, nil
}

type FunctionType struct {
	Generics   []Name           `("<" @@ ("," @@)* ">")?`
	Arguments  []TypeAnnotation `"(" (@@ ("," @@)*)? ")"`
	ReturnType TypeAnnotation   `"-" ">" @@`
}

func (f FunctionType) sealedTypeAnnotation() {}
func (f FunctionType) TypeAnnotationCases() (*SingleNameType, *FunctionType) {
	return nil, &f
}

type Module struct {
	Node
	Implementing Name                `"implement" @@`
	Declarations []ModuleDeclaration `"{" @@* "}"`
}

func (m Module) sealedExpression() {}
func (m Module) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return &m, nil, nil, nil, nil, nil
}

func ModuleFields(node Module) (Name, []ModuleDeclaration) {
	return node.Implementing, node.Declarations
}

type ModuleDeclaration struct {
	Public     bool       `@"public"?`
	Name       Name       `@@`
	Expression Expression `":" "=" @@`
}

func ModuleDeclarationFields(node ModuleDeclaration) (bool, Name, Expression) {
	return node.Public, node.Name, node.Expression
}

type ArgumentsList struct {
	Node
	Generics  []Name          `("<" @@ ("," @@)* ">")?`
	Arguments []ExpressionBox `"(" (@@ ("," @@)*)? ")"`
}

type AccessOrInvocation struct {
	VarName   Name           `"." @@`
	Arguments *ArgumentsList `@@?`
}

type ExpressionBox struct {
	Expression              Expression           `@@`
	AccessOrInvocationChain []AccessOrInvocation `@@*`
}

func ExpressionBoxFields(expressionBox ExpressionBox) (Expression, []AccessOrInvocation) {
	return expressionBox.Expression, expressionBox.AccessOrInvocationChain
}

type Expression interface {
	sealedExpression()
	ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If)
}

func GetExpressionNode(expression Expression) Node {
	caseModule, caseLiteral, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseModule != nil {
		return caseModule.Node
	} else if caseLiteral != nil {
		return caseLiteral.Node
	} else if caseReferenceOrInvocation != nil {
		return caseReferenceOrInvocation.Var.Node
	} else if caseLambda != nil {
		return caseLambda.Node
	} else if caseDeclaration != nil {
		return caseDeclaration.Name.Node
	} else if caseIf != nil {
		return caseIf.Node
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

var expressionUnion = participle.Union[Expression](Module{}, If{}, Declaration{}, LiteralExpression{}, ReferenceOrInvocation{}, Lambda{})

type If struct {
	Node
	Condition ExpressionBox   `"if" @@`
	ThenBlock []ExpressionBox `"{" @@* "}"`
	ElseBlock []ExpressionBox `("else" "{" @@* "}")?`
}

func (i If) sealedExpression() {}
func (i If) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, nil, nil, nil, &i
}

func IfFields(parserIf If) (ExpressionBox, []ExpressionBox, []ExpressionBox) {
	return parserIf.Condition, parserIf.ThenBlock, parserIf.ElseBlock
}

type Declaration struct {
	Name          Name          `@@`
	ExpressionBox ExpressionBox `":" "=" @@`
}

func (d Declaration) sealedTopLevelDeclaration() {}
func (d Declaration) TopLevelDeclarationCases() (*Declaration, *Interface, *Struct) {
	return &d, nil, nil
}
func (d Declaration) sealedExpression() {}
func (d Declaration) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, nil, nil, &d, nil
}
func DeclarationFields(node Declaration) (Name, ExpressionBox) {
	return node.Name, node.ExpressionBox
}

type LiteralExpression struct {
	Node
	Literal Literal `@@`
}

func (l LiteralExpression) sealedExpression() {}
func (l LiteralExpression) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, &l, nil, nil, nil, nil
}

type Lambda struct {
	Node
	Generics   []Name          `("<" @@ ("," @@)* ">")?`
	Parameters []Parameter     `"(" (@@ ("," @@)*)? ")"`
	ReturnType *TypeAnnotation `(":" @@)?`
	Block      []ExpressionBox `"=" ">" (("{" @@* "}") | @@)`
}

func (l Lambda) sealedExpression() {}
func (l Lambda) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, nil, &l, nil, nil
}

func LambdaFields(node Lambda) ([]Name, []Parameter, *TypeAnnotation, []ExpressionBox) {
	return node.Generics, node.Parameters, node.ReturnType, node.Block
}

type Parameter struct {
	Name Name            `@@`
	Type *TypeAnnotation `(":" @@)?`
}

func ParameterFields(node Parameter) (Name, *TypeAnnotation) {
	return node.Name, node.Type
}

type ReferenceOrInvocation struct {
	Var       Name           `@@`
	Arguments *ArgumentsList `@@?`
}

func (r ReferenceOrInvocation) sealedExpression() {}
func (r ReferenceOrInvocation) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, &r, nil, nil, nil
}

func ReferenceOrInvocationFields(node ReferenceOrInvocation) (Name, *ArgumentsList) {
	return node.Var, node.Arguments
}
