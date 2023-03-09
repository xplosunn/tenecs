package parser

import (
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

type Identifier struct {
	Node
	Name string `@Ident`
}

type Package struct {
	Identifier Identifier `"package" @@`
}

type Import struct {
	DotSeparatedVars []string `"import" (@Ident ("." @Ident)*)?`
}

func ImportFields(node Import) []string {
	return node.DotSeparatedVars
}

type TopLevelDeclaration interface {
	sealedTopLevelDeclaration()
	TopLevelDeclarationCases() (*Declaration, *Interface, *Struct)
}

var topLevelDeclarationUnion = participle.Union[TopLevelDeclaration](Struct{}, Interface{}, Declaration{})

type Struct struct {
	Name      string           `"struct" @Ident`
	Generics  []string         `("<" (@Ident ("," @Ident)*)? ">")?`
	Variables []StructVariable `"(" (@@ ("," @@)*)? ")"`
}

func (s Struct) sealedTopLevelDeclaration() {}
func (s Struct) TopLevelDeclarationCases() (*Declaration, *Interface, *Struct) {
	return nil, nil, &s
}

func StructFields(struc Struct) (string, []string, []StructVariable) {
	return struc.Name, struc.Generics, struc.Variables
}

type StructVariable struct {
	Name string         `@Ident`
	Type TypeAnnotation `":" @@`
}

func StructVariableFields(structVariable StructVariable) (string, TypeAnnotation) {
	return structVariable.Name, structVariable.Type
}

type Interface struct {
	Name      string              `"interface" @Ident`
	Variables []InterfaceVariable `"{" @@* "}"`
}

func (i Interface) sealedTopLevelDeclaration() {}
func (i Interface) TopLevelDeclarationCases() (*Declaration, *Interface, *Struct) {
	return nil, &i, nil
}

func InterfaceFields(interf Interface) (string, []InterfaceVariable) {
	return interf.Name, interf.Variables
}

type InterfaceVariable struct {
	Name string         `"public" @Ident`
	Type TypeAnnotation `":" @@`
}

func InterfaceVariableFields(interfaceVariable InterfaceVariable) (string, TypeAnnotation) {
	return interfaceVariable.Name, interfaceVariable.Type
}

type TypeAnnotation interface {
	sealedTypeAnnotation()
	TypeAnnotationCases() (*SingleNameType, *FunctionType)
}

var typeAnnotationUnion = participle.Union[TypeAnnotation](SingleNameType{}, FunctionType{})

type SingleNameType struct {
	TypeName string `@Ident`
}

func (s SingleNameType) sealedTypeAnnotation() {}
func (s SingleNameType) TypeAnnotationCases() (*SingleNameType, *FunctionType) {
	return &s, nil
}

type FunctionType struct {
	Generics   []string         `("<" @Ident ("," @Ident)* ">")?`
	Arguments  []TypeAnnotation `"(" (@@ ("," @@)*)? ")"`
	ReturnType TypeAnnotation   `"-" ">" @@`
}

func (f FunctionType) sealedTypeAnnotation() {}
func (f FunctionType) TypeAnnotationCases() (*SingleNameType, *FunctionType) {
	return nil, &f
}

type Module struct {
	Implementing string              `"implement" @Ident`
	Declarations []ModuleDeclaration `"{" @@* "}"`
}

func (m Module) sealedExpression() {}
func (m Module) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return &m, nil, nil, nil, nil, nil
}

func ModuleFields(node Module) (string, []ModuleDeclaration) {
	return node.Implementing, node.Declarations
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
	ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If)
}

var expressionUnion = participle.Union[Expression](Module{}, If{}, Declaration{}, LiteralExpression{}, ReferenceOrInvocation{}, Lambda{})

type If struct {
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
	Name          string        `@Ident`
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
func DeclarationFields(node Declaration) (string, ExpressionBox) {
	return node.Name, node.ExpressionBox
}

type LiteralExpression struct {
	Literal Literal `@@`
}

func (l LiteralExpression) sealedExpression() {}
func (l LiteralExpression) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, &l, nil, nil, nil, nil
}

type Lambda struct {
	Generics   []string        `("<" @Ident ("," @Ident)* ">")?`
	Parameters []Parameter     `"(" (@@ ("," @@)*)? ")"`
	ReturnType *TypeAnnotation `(":" @@)?`
	Block      []ExpressionBox `"=" ">" (("{" @@* "}") | @@)`
}

func (l Lambda) sealedExpression() {}
func (l Lambda) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, nil, &l, nil, nil
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
func (r ReferenceOrInvocation) ExpressionCases() (*Module, *LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, &r, nil, nil, nil
}

func ReferenceOrInvocationFields(node ReferenceOrInvocation) (string, *ArgumentsList) {
	return node.Var, node.Arguments
}
