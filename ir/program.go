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
	Name           string
	ParameterNames []string
	Body           []Statement
}

type Statement interface {
	sealedStatement()
}

type Expression interface {
	sealedStatement()
	sealedExpression()
}

type Return struct {
	ReturnExpression Expression
}

func (s Return) sealedStatement()  {}

type VariableDeclaration struct {
	ReturnExpression Expression
}

func (s VariableDeclaration) sealedStatement()  {}
func (s VariableDeclaration) sealedExpression() {}

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

type Invocation struct {
	Over      Expression
	Arguments []Expression
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
