package godsl

type Statement interface {
	sealedGoDSL()
	sealedStatementCases() (*goVariableDeclaration, *goReturn, *goIf, *goNativeFunctionInvocation)
}

func VariableDeclaration(name string, value Expression) Statement {
	return goVariableDeclaration{name, value}
}

type goVariableDeclaration struct {
	name  string
	value Expression
}

func (g goVariableDeclaration) sealedGoDSL() {}

func (g goVariableDeclaration) sealedStatementCases() (*goVariableDeclaration, *goReturn, *goIf, *goNativeFunctionInvocation) {
	return &g, nil, nil, nil
}

func Return(exp Expression) Statement {
	return goReturn{exp}
}

type goReturn struct {
	exp Expression
}

func (g goReturn) sealedGoDSL() {}

func (g goReturn) sealedStatementCases() (*goVariableDeclaration, *goReturn, *goIf, *goNativeFunctionInvocation) {
	return nil, &g, nil, nil
}

type IfBuilder interface {
	Then(statements []Statement) ElseBuilder
}

type ElseBuilder interface {
	ElseIf(condition Expression) IfBuilder
	Else(statements []Statement) Statement
}

func If(condition Expression) IfBuilder {
	return goIf{
		thenBranches: []goIfBranch{goIfBranch{
			condition: condition,
			then:      nil,
		}},
		elseBranch: nil,
	}
}

func (g goIf) Then(statements []Statement) ElseBuilder {
	return goIf{
		thenBranches: append(g.thenBranches[0:len(g.thenBranches)-1], goIfBranch{
			condition: g.thenBranches[len(g.thenBranches)-1].condition,
			then:      statements,
		}),
		elseBranch: nil,
	}
}

func (g goIf) ElseIf(condition Expression) IfBuilder {
	return goIf{
		thenBranches: append(g.thenBranches, goIfBranch{
			condition: condition,
			then:      nil,
		}),
		elseBranch: nil,
	}
}

func (g goIf) Else(statements []Statement) Statement {
	return goIf{
		thenBranches: g.thenBranches,
		elseBranch:   statements,
	}
}

type goIfBranch struct {
	condition Expression
	then      []Statement
}

type goIf struct {
	thenBranches []goIfBranch
	elseBranch   []Statement
}

func (g goIf) sealedGoDSL() {}

func (g goIf) sealedStatementCases() (*goVariableDeclaration, *goReturn, *goIf, *goNativeFunctionInvocation) {
	return nil, nil, &g, nil
}

type NativeFunctionInvocationBuilder interface {
	DeclaringVariables(vars ...string) NativeFunctionInvocationImportBuilder
	Import(pkg string) NativeFunctionInvocationNameBuilder
	Name(name string) NativeFunctionInvocationParametersBuilder
}

type NativeFunctionInvocationImportBuilder interface {
	Import(pkg string) NativeFunctionInvocationNameBuilder
	Name(name string) NativeFunctionInvocationParametersBuilder
}

type NativeFunctionInvocationNameBuilder interface {
	Name(name string) NativeFunctionInvocationParametersBuilder
}

type NativeFunctionInvocationParametersBuilder interface {
	Parameters(params ...Expression) Statement
}

func NativeFunctionInvocation() NativeFunctionInvocationBuilder {
	return goNativeFunctionInvocation{}
}

type goNativeFunctionInvocation struct {
	vars   []string
	pkg    *string
	name   string
	params []Expression
}

func (g goNativeFunctionInvocation) sealedGoDSL() {}

func (g goNativeFunctionInvocation) sealedStatementCases() (*goVariableDeclaration, *goReturn, *goIf, *goNativeFunctionInvocation) {
	return nil, nil, nil, &g
}

func (g goNativeFunctionInvocation) DeclaringVariables(vars ...string) NativeFunctionInvocationImportBuilder {
	return goNativeFunctionInvocation{
		vars:   vars,
		pkg:    g.pkg,
		name:   g.name,
		params: g.params,
	}
}

func (g goNativeFunctionInvocation) Import(pkg string) NativeFunctionInvocationNameBuilder {
	return goNativeFunctionInvocation{
		vars:   g.vars,
		pkg:    &pkg,
		name:   g.name,
		params: g.params,
	}
}

func (g goNativeFunctionInvocation) Name(name string) NativeFunctionInvocationParametersBuilder {
	return goNativeFunctionInvocation{
		vars:   g.vars,
		pkg:    g.pkg,
		name:   name,
		params: g.params,
	}
}

func (g goNativeFunctionInvocation) Parameters(params ...Expression) Statement {
	return goNativeFunctionInvocation{
		vars:   g.vars,
		pkg:    g.pkg,
		name:   g.name,
		params: params,
	}
}
