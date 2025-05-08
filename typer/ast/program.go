package ast

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/types"
	"sort"
)

type Program struct {
	Declarations    map[Ref]Expression
	TypeAliases     map[Ref]TypeAlias
	StructFunctions map[Ref]*types.Function
	NativeFunctions map[Ref]*types.Function
	FieldsByType    map[Ref]map[string]types.VariableType
}

type TypeAlias struct {
	Generics     []string
	VariableType types.VariableType
}

type Ref struct {
	Package string
	Name    string
}

func SortRefs(refs []Ref) {
	sort.Slice(refs, func(i, j int) bool {
		return refs[i].Package+"_"+refs[i].Name < refs[j].Package+"_"+refs[j].Name
	})
}

type CodePoint struct {
	FileName string
	Line     int
}

type Expression interface {
	sealedExpression()
	SourceCodePoint() CodePoint
	ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When)
}
type If struct {
	CodePoint    CodePoint
	VariableType types.VariableType
	Condition    Expression
	ThenBlock    []Expression
	ElseBlock    []Expression
}

func (i If) sealedExpression()          {}
func (i If) SourceCodePoint() CodePoint { return i.CodePoint }
func (i If) ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When) {
	return nil, nil, nil, nil, nil, nil, &i, nil, nil
}

type Declaration struct {
	CodePoint  CodePoint
	Name       string
	Expression Expression
}

func (d Declaration) sealedExpression()          {}
func (d Declaration) SourceCodePoint() CodePoint { return d.CodePoint }
func (d Declaration) ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When) {
	return nil, nil, nil, nil, nil, &d, nil, nil, nil
}

type Literal struct {
	CodePoint    CodePoint
	VariableType types.VariableType
	Literal      parser.Literal
}

func (l Literal) sealedExpression()          {}
func (l Literal) SourceCodePoint() CodePoint { return l.CodePoint }
func (l Literal) ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When) {
	return &l, nil, nil, nil, nil, nil, nil, nil, nil
}

type Function struct {
	CodePoint    CodePoint
	VariableType *types.Function
	Block        []Expression
}

func (f Function) sealedExpression()          {}
func (f Function) SourceCodePoint() CodePoint { return f.CodePoint }
func (f Function) ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When) {
	return nil, nil, nil, nil, &f, nil, nil, nil, nil
}

type Reference struct {
	CodePoint    CodePoint
	VariableType types.VariableType
	PackageName  *string
	Name         string
}

func (r Reference) sealedExpression()          {}
func (r Reference) SourceCodePoint() CodePoint { return r.CodePoint }
func (r Reference) ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When) {
	return nil, &r, nil, nil, nil, nil, nil, nil, nil
}

type Access struct {
	CodePoint    CodePoint
	VariableType types.VariableType
	Over         Expression
	Access       string
}

func (a Access) sealedExpression()          {}
func (a Access) SourceCodePoint() CodePoint { return a.CodePoint }
func (a Access) ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When) {
	return nil, nil, &a, nil, nil, nil, nil, nil, nil
}

type Invocation struct {
	CodePoint    CodePoint
	VariableType types.VariableType
	Over         Expression
	Generics     []types.VariableType
	Arguments    []Expression
}

func (i Invocation) sealedExpression()          {}
func (i Invocation) SourceCodePoint() CodePoint { return i.CodePoint }
func (i Invocation) ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When) {
	return nil, nil, nil, &i, nil, nil, nil, nil, nil
}

type List struct {
	CodePoint             CodePoint
	ContainedVariableType types.VariableType
	Arguments             []Expression
}

func (a List) sealedExpression()          {}
func (a List) SourceCodePoint() CodePoint { return a.CodePoint }
func (a List) ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When) {
	return nil, nil, nil, nil, nil, nil, nil, &a, nil
}

type WhenCase struct {
	Name         *string
	VariableType types.VariableType
	Block        []Expression
}

type When struct {
	CodePoint     CodePoint
	VariableType  types.VariableType
	Over          Expression
	Cases         []WhenCase
	OtherCase     []Expression
	OtherCaseName *string
}

func (w When) sealedExpression()          {}
func (w When) SourceCodePoint() CodePoint { return w.CodePoint }
func (w When) ExpressionCases() (*Literal, *Reference, *Access, *Invocation, *Function, *Declaration, *If, *List, *When) {
	return nil, nil, nil, nil, nil, nil, nil, nil, &w
}

func VariableTypeOfExpression(expression Expression) types.VariableType {
	if expression == nil {
		panic("nil expression in VariableTypeOfExpression")
	}
	caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseList, caseWhen := expression.ExpressionCases()
	if caseLiteral != nil {
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
	} else if caseList != nil {
		return &types.List{
			Generic: caseList.ContainedVariableType,
		}
	} else if caseWhen != nil {
		return caseWhen.VariableType
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}
