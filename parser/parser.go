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
	Cases() (*Module, *Interface)
}

var topLevelDeclarationUnion = participle.Union[TopLevelDeclaration](Module{}, Interface{})

type Interface struct {
	Name      string              `"interface" @Ident`
	Variables []InterfaceVariable `"{" @@* "}"`
}

func (i Interface) sealedTopLevelDeclaration() {}
func (i Interface) Cases() (*Module, *Interface) {
	return nil, &i
}

func InterfaceFields(interf Interface) (string, []InterfaceVariable) {
	return interf.Name, interf.Variables
}

type InterfaceVariable struct {
	Name string         `@Ident`
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
	Arguments  []TypeAnnotation `"(" (@@ ("," @@)*)? ")"`
	ReturnType TypeAnnotation   `"-" ">" @@`
}

func (f FunctionType) sealedTypeAnnotation() {}
func (f FunctionType) Cases() (*SingleNameType, *FunctionType) {
	return nil, &f
}

type Module struct {
	Name         string              `"module" @Ident`
	Implements   []string            `(":" @Ident ("," @Ident)*)?`
	Declarations []ModuleDeclaration `"{" @@* "}"`
}

func (m Module) sealedTopLevelDeclaration() {}
func (m Module) Cases() (*Module, *Interface) {
	return &m, nil
}

func ModuleFields(node Module) (string, []string, []ModuleDeclaration) {
	return node.Name, node.Implements, node.Declarations
}

type ModuleDeclaration struct {
	Public     bool       `@"public"?`
	Name       string     `@Ident`
	Expression Expression `":" "=" @@`
}

func ModuleDeclarationFields(node ModuleDeclaration) (bool, string, Expression) {
	return node.Public, node.Name, node.Expression
}

type Expression interface {
	sealedExpression()
	Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If)
}

var expressionUnion = participle.Union[Expression](If{}, Declaration{}, LiteralExpression{}, ReferenceOrInvocation{}, Lambda{})

type If struct {
	Condition Expression   `"if" @@`
	ThenBlock []Expression `"{" @@* "}"`
	ElseBlock []Expression `("else" "{" @@* "}")?`
}

func (i If) sealedExpression() {}
func (i If) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, nil, nil, &i
}

type Declaration struct {
	Name       string     `@Ident`
	Expression Expression `":" "=" @@`
}

func DeclarationFields(node Declaration) (string, Expression) {
	return node.Name, node.Expression
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
	Parameters []Parameter     `"(" (@@ ("," @@)*)? ")"`
	ReturnType *TypeAnnotation `(":" @@)?`
	Block      []Expression    `"=" ">" "{" @@* "}"`
}

func (l Lambda) sealedExpression() {}
func (l Lambda) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, nil, &l, nil, nil
}

func LambdaFields(node Lambda) ([]Parameter, *TypeAnnotation, []Expression) {
	return node.Parameters, node.ReturnType, node.Block
}

type Parameter struct {
	Name string          `@Ident`
	Type *TypeAnnotation `(":" @@)?`
}

func ParameterFields(node Parameter) (string, *TypeAnnotation) {
	return node.Name, node.Type
}

type ArgumentsList struct {
	Arguments []Expression `"(" (@@ ("," @@)*)? ")"`
}

type ReferenceOrInvocation struct {
	DotSeparatedVars []string       `@Ident ("." @Ident)*`
	Arguments        *ArgumentsList `@@?`
}

func (r ReferenceOrInvocation) sealedExpression() {}
func (r ReferenceOrInvocation) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda, *Declaration, *If) {
	return nil, &r, nil, nil, nil
}

func ReferenceOrInvocationFields(node ReferenceOrInvocation) ([]string, *[]Expression) {
	if node.Arguments == nil {
		return node.DotSeparatedVars, nil
	}
	return node.DotSeparatedVars, &node.Arguments.Arguments
}
