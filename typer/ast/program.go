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
	ExpressionCases() (*Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If)
}

type If struct {
	VariableType types.VariableType
	Condition    Expression
	ThenBlock    []Expression
	ElseBlock    []Expression
}

func (i If) sealedExpression() {}
func (i If) ExpressionCases() (*Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If) {
	return nil, nil, nil, nil, nil, &i
}

type Declaration struct {
	VariableType types.VariableType
	Name         string
	Expression   Expression
}

func (d Declaration) sealedExpression() {}
func (d Declaration) ExpressionCases() (*Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If) {
	return nil, nil, nil, nil, &d, nil
}

type Literal struct {
	VariableType types.BasicType
	Literal      parser.Literal
}

func (l Literal) sealedExpression() {}
func (l Literal) ExpressionCases() (*Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If) {
	return &l, nil, nil, nil, nil, nil
}

type Function struct {
	VariableType types.Function
	Block        []Expression
}

func (f Function) sealedExpression() {}
func (f Function) ExpressionCases() (*Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If) {
	return nil, nil, nil, &f, nil, nil
}

type ArgumentsList struct {
	Arguments []Expression
}

type AccessAndMaybeInvocation struct {
	VariableType  types.VariableType
	Access        string
	ArgumentsList *ArgumentsList
}

type ReferenceAndMaybeInvocation struct {
	VariableType  types.VariableType
	Name          string
	ArgumentsList *ArgumentsList
}

func (r ReferenceAndMaybeInvocation) sealedExpression() {}
func (r ReferenceAndMaybeInvocation) ExpressionCases() (*Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If) {
	return nil, &r, nil, nil, nil, nil
}

type WithAccessAndMaybeInvocation struct {
	VariableType types.VariableType
	Over         Expression
	AccessChain  []AccessAndMaybeInvocation
}

func (w WithAccessAndMaybeInvocation) sealedExpression() {}
func (w WithAccessAndMaybeInvocation) ExpressionCases() (*Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If) {
	return nil, nil, &w, nil, nil, nil
}

func VariableTypeOfExpression(expression Expression) types.VariableType {
	caseLiteral, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseLiteral != nil {
		return caseLiteral.VariableType
	} else if caseReferenceAndMaybeInvocation != nil {
		return caseReferenceAndMaybeInvocation.VariableType
	} else if caseWithAccessAndMaybeInvocation != nil {
		return caseWithAccessAndMaybeInvocation.VariableType
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
