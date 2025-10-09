package ir

import (
	"github.com/xplosunn/tenecs/parser"
)

type Program struct {
	Declarations          map[Reference]TopLevelFunction
	StructFunctions       map[Reference]*StructFunction
	NativeFunctions       []NativeFunctionRef
	StdLibStructFunctions map[Reference]*StructFunction
}

type StructFunction struct {
	GenericsArgumentName string
	OtherArgumentNames   []string
	ReturnType           NonFunctionTypeName
}

type NativeFunctionRef struct {
	Package string
	Name    string
}

type TopLevelFunction struct {
	Body Statement
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

type CodePoint struct {
	FileName string
	Line     int
}

type Invocation struct {
	CodePointAsFirstArgument *CodePoint
	Over                     Expression
	Arguments                []Expression
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

func (s If) sealedStatement() {}

type IfArrayIncludesElseEqual struct {
	Left  Expression
	Right Expression
}

type IfForTypeMatching struct {
	Condition IfArrayIncludesElseEqual
	ThenBlock []Statement
	ElseBlock []Statement
}

func (s IfForTypeMatching) sealedStatement() {}

type TypeNamesOrConcat struct {
	Expressions []Expression
}

func (s TypeNamesOrConcat) sealedStatement()  {}
func (s TypeNamesOrConcat) sealedExpression() {}

type NonFunctionTypeName struct {
	Prefix              string
	CommaSpaceSeparated []TypeNamesOrConcat
	Suffix              string
}

func (s NonFunctionTypeName) sealedStatement()  {}
func (s NonFunctionTypeName) sealedExpression() {}

type ArrayInstantiation struct {
	Arguments []Expression
}

func (s ArrayInstantiation) sealedStatement()  {}
func (s ArrayInstantiation) sealedExpression() {}

type ArrayIndexAccess struct {
	Over  Expression
	Index int
}

func (s ArrayIndexAccess) sealedStatement()  {}
func (s ArrayIndexAccess) sealedExpression() {}
