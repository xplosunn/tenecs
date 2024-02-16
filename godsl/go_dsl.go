package godsl

//TODO
// - for
// - array

type GoDSL interface {
	sealedGoDSL()
}

func exhaustiveSwitch(goDSL GoDSL) (*Expression, *Statement, *TopLevelStatement) {
	caseExpression, ok := goDSL.(Expression)
	if ok {
		return &caseExpression, nil, nil
	}
	caseStatement, ok := goDSL.(Statement)
	if ok {
		return nil, &caseStatement, nil
	}
	caseTopLevelStatement, ok := goDSL.(TopLevelStatement)
	if ok {
		return nil, nil, &caseTopLevelStatement
	}
	return nil, nil, nil
}

type Expression interface {
	sealedGoDSL()
	sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess)
}
type Statement interface {
	sealedGoDSL()
	sealedStatementCases() (*goVariableDeclaration, *goReturn, *goIf, *goNativeFunctionInvocation)
}

type TopLevelStatement interface {
	sealedGoDSL()
	sealedTopLevelStatementCases() *goNativeFunctionDeclaration
}

func FunctionCreation(parameterNames ...string) Expression {
	return goFunctionCreation{parameterNames}
}

type goFunctionCreation struct {
	parameterNames []string
}

func (g goFunctionCreation) sealedGoDSL() {}

func (g goFunctionCreation) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess) {
	return &g, nil, nil, nil
}

func FunctionInvocation(over Expression, arguments ...Expression) Expression {
	return goFunctionInvocation{over, arguments}
}

type goFunctionInvocation struct {
	over      Expression
	arguments []Expression
}

func (g goFunctionInvocation) sealedGoDSL() {}

func (g goFunctionInvocation) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess) {
	return nil, &g, nil, nil
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

func ObjectField(name string, value Expression) func(*goObjectCreation) {
	return func(objectCreation *goObjectCreation) {
		objectCreation.fields[name] = value
	}
}

func ObjectCreation(typeName string, fields ...func(*goObjectCreation)) Expression {
	result := goObjectCreation{typeName, map[string]Expression{}}
	for _, fieldFunc := range fields {
		fieldFunc(&result)
	}
	return result
}

type goObjectCreation struct {
	typeName string
	fields   map[string]Expression
}

func (g goObjectCreation) sealedGoDSL() {}

func (g goObjectCreation) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess) {
	return nil, nil, &g, nil
}

func ObjectAccess(over Expression, fieldName string) Expression {
	return goObjectAccess{over, fieldName}
}

type goObjectAccess struct {
	over      Expression
	fieldName string
}

func (g goObjectAccess) sealedGoDSL() {}

func (g goObjectAccess) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess) {
	return nil, nil, nil, &g
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

type NativeFunctionInvocationBuilder interface {
	DeclaringVariables(vars ...string) NativeFunctionInvocationImportBuilder
	Import(pkg string) NativeFunctionInvocationNameBuilder
	Name(name string) NativeFunctionInvocationParametersBuilder
}

type NativeFunctionInvocationImportBuilder interface {
	Import(pkg string) NativeFunctionInvocationNameBuilder
}

type NativeFunctionInvocationNameBuilder interface {
	Name(name string) NativeFunctionInvocationParametersBuilder
}

type NativeFunctionInvocationParametersBuilder interface {
	Parameters(params ...string) Statement
}

func NativeFunctionInvocation() NativeFunctionInvocationBuilder {
	return goNativeFunctionInvocation{}
}

type goNativeFunctionInvocation struct {
	vars   []string
	pkg    *string
	name   string
	params []string
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

func (g goNativeFunctionInvocation) Parameters(params ...string) Statement {
	return goNativeFunctionInvocation{
		vars:   g.vars,
		pkg:    g.pkg,
		name:   g.name,
		params: params,
	}
}
