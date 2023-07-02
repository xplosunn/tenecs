package ast

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/types"
)

type Program struct {
	Declarations           []*Declaration
	StructFunctions        map[string]*types.Function
	NativeFunctions        map[string]*types.Function
	NativeFunctionPackages map[string]string
}

type Expression interface {
	sealedExpression()
	ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When)
}

type Module struct {
	Implements *types.Interface
	Variables  map[string]Expression
}

func (m Module) sealedExpression() {}
func (m Module) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return &m, nil, nil, nil, nil, nil, nil, nil, nil, nil
}

type If struct {
	VariableType types.VariableType
	Condition    Expression
	ThenBlock    []Expression
	ElseBlock    []Expression
}

func (i If) sealedExpression() {}
func (i If) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, nil, nil, &i, nil, nil
}

type Declaration struct {
	Name       string
	Expression Expression
}

func (d Declaration) sealedExpression() {}
func (d Declaration) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, nil, &d, nil, nil, nil
}

type Literal struct {
	VariableType types.VariableType
	Literal      parser.Literal
}

func (l Literal) sealedExpression() {}
func (l Literal) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, &l, nil, nil, nil, nil, nil, nil, nil, nil
}

type Function struct {
	VariableType *types.Function
	Block        []Expression
}

func (f Function) sealedExpression() {}
func (f Function) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, &f, nil, nil, nil, nil
}

type Reference struct {
	VariableType types.VariableType
	Name         string
}

func (r Reference) sealedExpression() {}
func (r Reference) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, &r, nil, nil, nil, nil, nil, nil, nil
}

type Access struct {
	VariableType types.VariableType
	Over         Expression
	Access       string
}

func (a Access) sealedExpression() {}
func (a Access) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, &a, nil, nil, nil, nil, nil, nil
}

type Invocation struct {
	VariableType types.VariableType
	Over         Expression
	Generics     []types.StructFieldVariableType
	Arguments    []Expression
}

func (i Invocation) sealedExpression() {}
func (i Invocation) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, &i, nil, nil, nil, nil, nil
}

type Array struct {
	ContainedVariableType types.StructFieldVariableType
	Arguments             []Expression
}

func (a Array) sealedExpression() {}
func (a Array) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, nil, nil, nil, &a, nil
}

type When struct {
	VariableType types.VariableType
	Over         Expression
	Cases        map[types.VariableType][]Expression
}

func (w When) sealedExpression() {}
func (w When) ExpressionCases() (*Module, *Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *Array, *When) {
	return nil, nil, nil, nil, nil, nil, nil, nil, nil, &w
}

func VariableTypeOfExpression(expression Expression) types.VariableType {
	if expression == nil {
		panic("nil expression in VariableTypeOfExpression")
	}
	caseModule, caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseArray, caseWhen := expression.ExpressionCases()
	if caseModule != nil {
		return caseModule.Implements
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
		return &types.Void{}
	} else if caseIf != nil {
		return caseIf.VariableType
	} else if caseArray != nil {
		return &types.Array{
			OfType: caseArray.ContainedVariableType,
		}
	} else if caseWhen != nil {
		return caseWhen.VariableType
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}
