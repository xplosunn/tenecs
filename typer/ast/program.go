package ast

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/types"
)

type Program struct {
	Package                string
	Declarations           []*Declaration
	StructFunctions        map[string]*types.Function
	NativeFunctions        map[string]*types.Function
	NativeFunctionPackages map[string]string
	FieldsByType           map[string]map[string]types.VariableType
}

type Expression interface {
	sealedExpression()
	ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When)
}

type Implementation struct {
	Implements *types.KnownType
	Variables  map[string]Expression
}

func (m Implementation) sealedExpression() {}
func (m Implementation) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return &m, nil, nil, nil, nil, nil, nil, nil, nil, nil
}

type If struct {
	VariableType types.VariableType
	Condition    Expression
	ThenBlock    []Expression
	ElseBlock    []Expression
}

func (i If) sealedExpression() {}
func (i If) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, nil, nil, &i, nil, nil
}

type Declaration struct {
	Name       string
	Expression Expression
}

func (d Declaration) sealedExpression() {}
func (d Declaration) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, nil, &d, nil, nil, nil
}

type Literal struct {
	VariableType types.VariableType
	Literal      parser.Literal
}

func (l Literal) sealedExpression() {}
func (l Literal) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, &l, nil, nil, nil, nil, nil, nil, nil, nil
}

type Function struct {
	VariableType *types.Function
	Block        []Expression
}

func (f Function) sealedExpression() {}
func (f Function) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, &f, nil, nil, nil, nil
}

type Reference struct {
	VariableType types.VariableType
	PackageName  *string
	Name         string
}

func (r Reference) sealedExpression() {}
func (r Reference) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, &r, nil, nil, nil, nil, nil, nil, nil
}

type Access struct {
	VariableType types.VariableType
	Over         Expression
	Access       string
}

func (a Access) sealedExpression() {}
func (a Access) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, &a, nil, nil, nil, nil, nil, nil
}

type Invocation struct {
	VariableType types.VariableType
	Over         Expression
	Generics     []types.VariableType
	Arguments    []Expression
}

func (i Invocation) sealedExpression() {}
func (i Invocation) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, &i, nil, nil, nil, nil, nil
}

type Array struct {
	ContainedVariableType types.VariableType
	Arguments             []Expression
}

func (a Array) sealedExpression() {}
func (a Array) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, nil, nil, nil, &a, nil
}

type When struct {
	VariableType  types.VariableType
	Over          Expression
	Cases         map[types.VariableType][]Expression
	CaseNames     map[types.VariableType]*string
	OtherCase     []Expression
	OtherCaseName *string
}

func (w When) sealedExpression() {}
func (w When) ExpressionCases() (*Implementation, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, nil, nil, nil, nil, &w
}

func VariableTypeOfExpression(expression Expression) types.VariableType {
	if expression == nil {
		panic("nil expression in VariableTypeOfExpression")
	}
	caseImplementation, caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseArray, caseWhen := expression.ExpressionCases()
	if caseImplementation != nil {
		return caseImplementation.Implements
	} else if caseLiteral != nil {
		return caseLiteral.VariableType
	} else if caseReference != nil {
		return caseReference.VariableType
	} else if caseAccess != nil {
		return caseAccess.VariableType
	} else if caseInvocation != nil {
		return caseInvocation.VariableType
	} else if caseFunction != nil {
		return caseFunction.VariableType
	} else if caseDeclaration != nil {
		return types.Void()
	} else if caseIf != nil {
		return caseIf.VariableType
	} else if caseArray != nil {
		return types.Array(caseArray.ContainedVariableType)
	} else if caseWhen != nil {
		return caseWhen.VariableType
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}
