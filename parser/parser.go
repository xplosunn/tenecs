package parser

import (
	"github.com/alecthomas/participle/v2"
)

func ParseString(s string) (*FileTopLevel, error) {
	p, err := participle.Build[FileTopLevel](literalUnion, expressionUnion)
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
	Package Package  `@@`
	Imports []Import `@@*`
	Modules []Module `@@*`
}

func FileTopLevelFields(node FileTopLevel) (Package, []Import, []Module) {
	return node.Package, node.Imports, node.Modules
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

type Module struct {
	Name         string        `"module" @Ident`
	Implements   []string      `(":" @Ident ("," @Ident)*)?`
	Declarations []Declaration `"{" @@* "}"`
}

func ModuleFields(node Module) (string, []string, []Declaration) {
	return node.Name, node.Implements, node.Declarations
}

type Declaration struct {
	Public     bool       `@"public"?`
	Name       string     `@Ident`
	Expression Expression `":" "=" @@`
}

func DeclarationFields(node Declaration) (bool, string, Expression) {
	return node.Public, node.Name, node.Expression
}

type Expression interface {
	sealedExpression()
	Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda)
}

var expressionUnion = participle.Union[Expression](LiteralExpression{}, ReferenceOrInvocation{}, Lambda{})

type LiteralExpression struct {
	Literal Literal `@@`
}

func (l LiteralExpression) sealedExpression() {}
func (l LiteralExpression) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda) {
	return &l, nil, nil
}

type Lambda struct {
	Parameters []Parameter             `"(" (@@ ("," @@)*)? ")"`
	ReturnType string                  `(":" @Ident)?`
	Block      []ReferenceOrInvocation `"=" ">" "{" @@* "}"`
}

func (l Lambda) sealedExpression() {}
func (l Lambda) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda) {
	return nil, nil, &l
}

func LambdaFields(node Lambda) ([]Parameter, string, []ReferenceOrInvocation) {
	return node.Parameters, node.ReturnType, node.Block
}

type Parameter struct {
	Name string `@Ident`
	Type string `(":" @Ident)?`
}

func ParameterFields(node Parameter) (string, string) {
	return node.Name, node.Type
}

type ReferenceOrInvocation struct {
	DotSeparatedVars []string       `@Ident ("." @Ident)*`
	Arguments        *ArgumentsList `@@?`
}

type ArgumentsList struct {
	Arguments []Expression `"(" (@@ ("," @@)*)? ")"`
}

func (r ReferenceOrInvocation) sealedExpression() {}
func (r ReferenceOrInvocation) Cases() (*LiteralExpression, *ReferenceOrInvocation, *Lambda) {
	return nil, &r, nil
}

func ReferenceOrInvocationFields(node ReferenceOrInvocation) ([]string, *[]Expression) {
	if node.Arguments == nil {
		return node.DotSeparatedVars, nil
	}
	return node.DotSeparatedVars, &node.Arguments.Arguments
}
