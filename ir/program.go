package ir

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/types"
)

type Program struct {
	Declarations    map[Reference]TopLevelFunction
	StructFunctions map[Reference]*types.Function
	NativeFunctions map[NativeFunctionRef]*types.Function
}

type NativeFunctionRef struct {
	Package string
	Name    string
}

type TopLevelFunction struct {
	ParameterNames []string
	Body           []Statement
}

type Statement interface {
	sealedStatement()
}

type Expression interface {
	//TODO probably shouldn't treat this as an extension of Statement
	// as I'm not sure if all expressions are statements
	sealedStatement()
	sealedExpression()
}

type Return struct {
	ReturnExpression Expression
}

func (s Return) sealedStatement() {}

type VariableDeclaration struct {
	Name       string
	Expression Expression
}

func (s VariableDeclaration) sealedStatement() {}

type ObjectInstantiation struct {
	Fields map[string]Expression
}

func (s ObjectInstantiation) sealedStatement()  {}
func (s ObjectInstantiation) sealedExpression() {}

type FieldAccess struct {
	Over      Expression
	FieldName string
}

func (s FieldAccess) sealedStatement()  {}
func (s FieldAccess) sealedExpression() {}

type InvocationOverTopLevelFunction struct {
	Over Expression
}

func (s InvocationOverTopLevelFunction) sealedStatement()  {}
func (s InvocationOverTopLevelFunction) sealedExpression() {}

type Invocation struct {
	Over           Expression
	Arguments      []Expression
	GenericsPassed []string
}

func (s Invocation) sealedStatement()  {}
func (s Invocation) sealedExpression() {}

type LocalFunction struct {
	ParameterNames []string
	Block          []Statement
}

func (s LocalFunction) sealedStatement()  {}
func (s LocalFunction) sealedExpression() {}

type Reference struct {
	Name string
}

func (s Reference) sealedStatement()  {}
func (s Reference) sealedExpression() {}

type Literal struct {
	Value parser.Literal
}

func (s Literal) sealedStatement()  {}
func (s Literal) sealedExpression() {}

type If struct {
	Condition Expression
	ThenBlock []Statement
	ElseBlock []Statement
}

func (s If) sealedStatement()  {}
func (s If) sealedExpression() {}

type EqualityComparison struct {
	Left  Expression
	Right Expression
}

func (s EqualityComparison) sealedStatement()  {}
func (s EqualityComparison) sealedExpression() {}
