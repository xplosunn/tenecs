package godsl

type TopLevelStatement interface {
	sealedGoDSL()
	sealedTopLevelStatementCases() *goNativeFunctionDeclaration
}

type NativeFunctionDeclarationBuilder interface {
	Parameters(params ...string) NativeFunctionDeclarationBodyBuilder
}

type NativeFunctionDeclarationBodyBuilder interface {
	Body(statements ...Statement) TopLevelStatement
}

func NativeFunctionDeclaration(name string) NativeFunctionDeclarationBuilder {
	return goNativeFunctionDeclaration{
		name:   name,
		params: nil,
		body:   nil,
	}
}

type goNativeFunctionDeclaration struct {
	name   string
	params []string
	body   []Statement
}

func (g goNativeFunctionDeclaration) sealedTopLevelStatementCases() *goNativeFunctionDeclaration {
	return &g
}

func (g goNativeFunctionDeclaration) sealedGoDSL() {}

func (g goNativeFunctionDeclaration) Parameters(params ...string) NativeFunctionDeclarationBodyBuilder {
	return goNativeFunctionDeclaration{
		name:   g.name,
		params: params,
		body:   g.body,
	}
}

func (g goNativeFunctionDeclaration) Body(statements ...Statement) TopLevelStatement {
	return goNativeFunctionDeclaration{
		name:   g.name,
		params: g.params,
		body:   statements,
	}
}
