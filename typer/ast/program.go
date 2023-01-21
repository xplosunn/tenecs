package ast

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/types"
)

type Program struct {
	Modules []*Module
}

type Module struct {
	Name                 string
	Implements           types.Interface
	ConstructorArguments []ConstructorArgument
	Variables            map[string]Expression
}

type ConstructorArgument struct {
	Name         string
	VariableType types.VariableType
}

type Expression interface {
	sealedExpression()
	Cases() (*Literal, *ReferenceOrInvocation, *Function, *Declaration, *If)
}

type If struct {
	VariableType types.VariableType
	Condition    Expression
	ThenBlock    []Expression
	ElseBlock    []Expression
}

func (i If) sealedExpression() {}
func (i If) Cases() (*Literal, *ReferenceOrInvocation, *Function, *Declaration, *If) {
	return nil, nil, nil, nil, &i
}

type Declaration struct {
	VariableType types.VariableType
	Name         string
	Expression   Expression
}

func (d Declaration) sealedExpression() {}
func (d Declaration) Cases() (*Literal, *ReferenceOrInvocation, *Function, *Declaration, *If) {
	return nil, nil, nil, &d, nil
}

type Literal struct {
	VariableType types.BasicType
	Literal      parser.Literal
}

func (l Literal) sealedExpression() {}
func (l Literal) Cases() (*Literal, *ReferenceOrInvocation, *Function, *Declaration, *If) {
	return &l, nil, nil, nil, nil
}

type Function struct {
	VariableType types.Function
	Block        []Expression
}

func (f Function) sealedExpression() {}
func (f Function) Cases() (*Literal, *ReferenceOrInvocation, *Function, *Declaration, *If) {
	return nil, nil, &f, nil, nil
}

type ArgumentsList struct {
	Arguments []Expression
}

type ReferenceOrInvocation struct {
	VariableType     types.VariableType
	DotSeparatedVars []string       `@Ident ("." @Ident)*`
	Arguments        *ArgumentsList `@@?`
}

func (r ReferenceOrInvocation) sealedExpression() {}
func (r ReferenceOrInvocation) Cases() (*Literal, *ReferenceOrInvocation, *Function, *Declaration, *If) {
	return nil, &r, nil, nil, nil
}

func VariableTypeOfExpression(expression Expression) types.VariableType {
	caseLiteral, caseReferenceOrInvocation, caseFunction, caseDeclaration, caseIf := expression.Cases()
	if caseLiteral != nil {
		return caseLiteral.VariableType
	} else if caseReferenceOrInvocation != nil {
		return caseReferenceOrInvocation.VariableType
	} else if caseFunction != nil {
		return caseFunction.VariableType
	} else if caseDeclaration != nil {
		return caseDeclaration.VariableType
	} else if caseIf != nil {
		return caseIf.VariableType
	} else {
		panic("code")
	}
}
