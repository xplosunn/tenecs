package ast

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/types"
)

type Program struct {
	Declarations    []*Declaration
	StructFunctions map[string]*types.Function
	NativeFunctions map[string]*types.Function
}

type Expression interface {
	sealedExpression()
	ExpressionCases() (*Module, *Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If, *Array)
}

type Module struct {
	Implements *types.Interface
	Variables  map[string]Expression
}

func (m Module) sealedExpression() {}
func (m Module) ExpressionCases() (*Module, *Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If, *Array) {
	return &m, nil, nil, nil, nil, nil, nil, nil
}

type If struct {
	VariableType types.VariableType
	Condition    Expression
	ThenBlock    []Expression
	ElseBlock    []Expression
}

func (i If) sealedExpression() {}
func (i If) ExpressionCases() (*Module, *Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If, *Array) {
	return nil, nil, nil, nil, nil, nil, &i, nil
}

type Declaration struct {
	VariableType types.VariableType
	Name         string
	Expression   Expression
}

func (d Declaration) sealedExpression() {}
func (d Declaration) ExpressionCases() (*Module, *Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If, *Array) {
	return nil, nil, nil, nil, nil, &d, nil, nil
}

type Literal struct {
	VariableType *types.BasicType
	Literal      parser.Literal
}

func (l Literal) sealedExpression() {}
func (l Literal) ExpressionCases() (*Module, *Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If, *Array) {
	return nil, &l, nil, nil, nil, nil, nil, nil
}

type Function struct {
	VariableType *types.Function
	Block        []Expression
}

func (f Function) sealedExpression() {}
func (f Function) ExpressionCases() (*Module, *Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If, *Array) {
	return nil, nil, nil, nil, &f, nil, nil, nil
}

type ArgumentsList struct {
	Generics  []types.StructFieldVariableType
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
func (r ReferenceAndMaybeInvocation) ExpressionCases() (*Module, *Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If, *Array) {
	return nil, nil, &r, nil, nil, nil, nil, nil
}

type WithAccessAndMaybeInvocation struct {
	VariableType types.VariableType
	Over         Expression
	AccessChain  []AccessAndMaybeInvocation
}

func (w WithAccessAndMaybeInvocation) sealedExpression() {}
func (w WithAccessAndMaybeInvocation) ExpressionCases() (*Module, *Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If, *Array) {
	return nil, nil, nil, &w, nil, nil, nil, nil
}

type Array struct {
	VariableType types.StructFieldVariableType
	Arguments    []Expression
}

func (a Array) sealedExpression() {}
func (a Array) ExpressionCases() (*Module, *Literal, *ReferenceAndMaybeInvocation, *WithAccessAndMaybeInvocation, *Function, *Declaration, *If, *Array) {
	return nil, nil, nil, nil, nil, nil, nil, &a
}

func VariableTypeOfExpression(expression Expression) types.VariableType {
	caseModule, caseLiteral, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf, caseArray := expression.ExpressionCases()
	if caseModule != nil {
		return caseModule.Implements
	} else if caseLiteral != nil {
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
	} else if caseArray != nil {
		return types.VariableTypeFromStructFieldVariableType(caseArray.VariableType)
	} else {
		panic("code")
	}
}
