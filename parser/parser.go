package parser

import (
	"github.com/alecthomas/participle/v2"
)

func ParseString(s string) (*FileTopLevel, error) {
	p, err := participle.Build[FileTopLevel](literalUnion)
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
	Public bool   `@"public"?`
	Name   string `@Ident`
	Lambda Lambda `":" "=" @@*`
}

func DeclarationFields(node Declaration) (bool, string, Lambda) {
	return node.Public, node.Name, node.Lambda
}

type Lambda struct {
	Parameters []Parameter  `"(" (@@ ("," @@)*)? ")"`
	ReturnType string       `(":" @Ident)?`
	Block      []Invocation `"=" ">" "{" @@* "}"`
}

func LambdaFields(node Lambda) ([]Parameter, string, []Invocation) {
	return node.Parameters, node.ReturnType, node.Block
}

type Parameter struct {
	Name string `@Ident`
	Type string `(":" @Ident)?`
}

func ParameterFields(node Parameter) (string, string) {
	return node.Name, node.Type
}

type Invocation struct {
	DotSeparatedVars []string  `(@Ident ("." @Ident)*)?`
	Argument         []Literal `"(" (@@ ("," @@)*)? ")"`
}

func InvocationFields(node Invocation) ([]string, []Literal) {
	return node.DotSeparatedVars, node.Argument
}
