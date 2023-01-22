package parser

import (
	"github.com/alecthomas/participle/v2"
)

func ParseString(s string) (*FileTopLevel, error) {
	p, err := participle.Build[FileTopLevel](topLevelDeclarationUnion, typeAnnotationUnion, literalUnion, expressionUnion)
	if err != nil {
		return nil, err
	}

	res, err := p.ParseString("", s)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type FileTopLevel struct {
	Package              Package               `@@`
	Imports              []Import              `@@*`
	TopLevelDeclarations []TopLevelDeclaration `@@*`
}

func FileTopLevelFields(node FileTopLevel) (Package, []Import, []TopLevelDeclaration) {
	return node.Package, node.Imports, node.TopLevelDeclarations
}

type Package struct {
	Identifier string `"package" @Ident`
}

func PackageFields(node Package) string {
	return node.Identifier
}

type Import struct {
	DotSeparatedVars []string `"import" (@Ident ("." @Ident)*)?`
}

func ImportFields(node Import) []string {
	return node.DotSeparatedVars
}

type TopLevelDeclaration interface {
	sealedTopLevelDeclaration()
	Cases() (*Module, *Interface, *Struct)
}

var topLevelDeclarationUnion = participle.Union[TopLevelDeclaration](Struct{}, Module{}, Interface{})

type Struct struct {
	Name      string           `"struct" @Ident`
	Generics  []string         `("<" (@Ident ("," @Ident)*)? ">")?`
	Variables []StructVariable `"(" (@@ ("," @@)*)? ")"`
}

func (s Struct) sealedTopLevelDeclaration() {}
func (s Struct) Cases() (*Module, *Interface, *Struct) {
	return nil, nil, &s
}

func StructFields(struc Struct) (string, []string, []StructVariable) {
	return struc.Name, struc.Generics, struc.Variables
}

type StructVariable struct {
	Name string         `@Ident`
	Type TypeAnnotation `":" @@`
}

type Interface struct {
	Name      string              `"interface" @Ident`
	Variables []InterfaceVariable `"{" @@* "}"`
}

func (i Interface) sealedTopLevelDeclaration() {}
func (i Interface) Cases() (*Module, *Interface, *Struct) {
	return nil, &i, nil
}

func InterfaceFields(interf Interface) (string, []InterfaceVariable) {
	return interf.Name, interf.Variables
}

type InterfaceVariable struct {
	Name string         `"public" @Ident`
	Type TypeAnnotation `":" @@`
}

type TypeAnnotation interface {
	sealedTypeAnnotation()
	Cases() (*SingleNameType, *FunctionType)
}

var typeAnnotationUnion = participle.Union[TypeAnnotation](SingleNameType{}, FunctionType{})

type SingleNameType struct {
	TypeName string `@Ident`
}

func (s SingleNameType) sealedTypeAnnotation() {}
func (s SingleNameType) Cases() (*SingleNameType, *FunctionType) {
	return &s, nil
}

type FunctionType struct {
	Generics   []string         `("<" @Ident ("," @Ident)* ">")?`
	Arguments  []TypeAnnotation `"(" (@@ ("," @@)*)? ")"`
	ReturnType TypeAnnotation   `"-" ">" @@`
}

func (f FunctionType) sealedTypeAnnotation() {}
func (f FunctionType) Cases() (*SingleNameType, *FunctionType) {
	return nil, &f
}

type Module struct {
	Implementing    string              `"implementing" @Ident`
	Name            string              `"module" @Ident`
	ConstructorArgs []ModuleParameter   `("(" (@@ ("," @@)*)? ")")?`
	Declarations    []ModuleDeclaration `"{" @@* "}"`
}

func (m Module) sealedTopLevelDeclaration() {}
func (m Module) Cases() (*Module, *Interface, *Struct) {
	return &m, nil, nil
}

func ModuleFields(node Module) (string, string, []ModuleParameter, []ModuleDeclaration) {
	return node.Implementing, node.Name, node.ConstructorArgs, node.Declarations
}

type ModuleParameter struct {
	Public bool           `@"public"?`
	Name   string         `@Ident`
	Type   TypeAnnotation `":" @@`
}

type ModuleDeclaration struct {
	Public     bool       `@"public"?`
	Name       string     `@Ident`
	Expression Expression `":" "=" @@`
}

func ModuleDeclarationFields(node ModuleDeclaration) (bool, string, Expression) {
	return node.Public, node.Name, node.Expression
}

type ArgumentsList struct {
	Generics  []string        `("<" @Ident ("," @Ident)* ">")?`
	Arguments []ExpressionBox `"(" (@@ ("," @@)*)? ")"`
}

type AccessOrInvocation struct {
	VarName   string         `"." @Ident`
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
	Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If)
}

var expressionUnion = participle.Union[Expression](If{}, Declaration{}, LiteralExpression{}, ReferenceOrInvocation{}, Lambda{})

type If struct {
	Condition ExpressionBox   `"if" @@`
	ThenBlock []ExpressionBox `"{" @@* "}"`
	ElseBlock []ExpressionBox `("else" "{" @@* "}")?`
}

func (i If) sealedExpression() {}
func (i If) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, nil, nil, &i
}

type Declaration struct {
	Name          string        `@Ident`
	ExpressionBox ExpressionBox `":" "=" @@`
}

func DeclarationFields(node Declaration) (string, ExpressionBox) {
	return node.Name, node.ExpressionBox
}
func (d Declaration) sealedExpression() {}
func (d Declaration) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, nil, &d, nil
}

type LiteralExpression struct {
	Literal Literal `@@`
}

func (l LiteralExpression) sealedExpression() {}
func (l LiteralExpression) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return &l, nil, nil, nil, nil
}

type Lambda struct {
	Generics   []string        `("<" @Ident ("," @Ident)* ">")?`
	Parameters []Parameter     `"(" (@@ ("," @@)*)? ")"`
	ReturnType *TypeAnnotation `(":" @@)?`
	Block      []ExpressionBox `"=" ">" "{" @@* "}"`
}

func (l Lambda) sealedExpression() {}
func (l Lambda) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, &l, nil, nil
}

func LambdaFields(node Lambda) ([]string, []Parameter, *TypeAnnotation, []ExpressionBox) {
	return node.Generics, node.Parameters, node.ReturnType, node.Block
}

type Parameter struct {
	Name string          `@Ident`
	Type *TypeAnnotation `(":" @@)?`
}

func ParameterFields(node Parameter) (string, *TypeAnnotation) {
	return node.Name, node.Type
}

type ReferenceOrInvocation struct {
	Var       string         `@Ident`
	Arguments *ArgumentsList `@@?`
}

func (r ReferenceOrInvocation) sealedExpression() {}
func (r ReferenceOrInvocation) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, &r, nil, nil, nil
}

func ReferenceOrInvocationFields(node ReferenceOrInvocation) (string, *ArgumentsList) {
	return node.Var, node.Arguments
}
