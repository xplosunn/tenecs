package desugar

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/xplosunn/tenecs/parser"
)

type FileTopLevel struct {
	Tokens               []lexer.Token
	Package              Package               `@@`
	Imports              []Import              `@@*`
	TopLevelDeclarations []TopLevelDeclaration `@@*`
}

type Name struct {
	parser.Node
	String string
}

type Package struct {
	parser.Node
	DotSeparatedNames []Name
}

type Import struct {
	parser.Node
	DotSeparatedVars []Name
	As               *Name
}

type TopLevelDeclaration interface {
	sealedTopLevelDeclaration()
}

func TopLevelDeclarationExhaustiveSwitch(
	topLevelDeclaration TopLevelDeclaration,
	caseDeclaration func(topLevelDeclaration Declaration),
	caseStruct func(topLevelDeclaration Struct),
	caseTypeAlias func(topLevelDeclaration TypeAlias),
) {
	declaration, ok := topLevelDeclaration.(Declaration)
	if ok {
		caseDeclaration(declaration)
		return
	}
	struc, ok := topLevelDeclaration.(Struct)
	if ok {
		caseStruct(struc)
		return
	}
	typeAlias, ok := topLevelDeclaration.(TypeAlias)
	if ok {
		caseTypeAlias(typeAlias)
		return
	}
}

type Struct struct {
	Name      Name
	Generics  []Name
	Variables []StructVariable
}

func (s Struct) sealedTopLevelDeclaration() {}

type StructVariable struct {
	Name Name
	Type TypeAnnotation
}

type TypeAlias struct {
	Name     Name
	Generics []Name
	Type     TypeAnnotation
}

func (i TypeAlias) sealedTopLevelDeclaration() {}

type TypeAnnotation struct {
	parser.Node
	OrTypes []TypeAnnotationElement `@@ ("|" @@)*`
}

type TypeAnnotationElement interface {
	sealedTypeAnnotationElement()
}

func TypeAnnotationElementExhaustiveSwitch(
	typeAnnotationElement TypeAnnotationElement,
	caseUnderscore func(underscoreTypeAnnotation SingleNameType),
	caseSingleNameType func(typeAnnotation SingleNameType),
	caseFunctionType func(typeAnnotation FunctionType),
) {
	singleNameType, ok := typeAnnotationElement.(SingleNameType)
	if ok {
		if singleNameType.TypeName.String == "_" {
			caseUnderscore(singleNameType)
			return
		}
		caseSingleNameType(singleNameType)
		return
	}
	functionType, ok := typeAnnotationElement.(FunctionType)
	if ok {
		caseFunctionType(functionType)
		return
	}
}

type SingleNameType struct {
	parser.Node
	TypeName Name
	Generics []TypeAnnotation
}

func (s SingleNameType) sealedTypeAnnotationElement() {}

type FunctionTypeArgument struct {
	Name *Name
	Type TypeAnnotation
}

type FunctionType struct {
	Generics   []Name
	Arguments  []FunctionTypeArgument
	ReturnType TypeAnnotation
}

func (f FunctionType) sealedTypeAnnotationElement() {}

type ArgumentsList struct {
	parser.Node
	Generics  []TypeAnnotation
	Arguments []NamedArgument
}

type NamedArgument struct {
	parser.Node
	Name     *Name
	Argument ExpressionBox
}

type DotName struct {
	parser.Node
	VarName Name
}

type AccessOrInvocation struct {
	parser.Node
	DotName   *DotName
	Arguments *ArgumentsList
}

type ExpressionBox struct {
	parser.Node
	Expression              Expression
	AccessOrInvocationChain []AccessOrInvocation
}

type Expression interface {
	sealedExpression()
}

func ExpressionExhaustiveSwitch(
	expression Expression,
	caseLiteralExpression func(expression LiteralExpression),
	caseReferenceOrInvocation func(expression ReferenceOrInvocation),
	caseLambda func(expression Lambda),
	caseDeclaration func(expression Declaration),
	caseIf func(expression If),
	caseList func(expression List),
	caseWhen func(expression When),
) {
	literalExpression, ok := expression.(LiteralExpression)
	if ok {
		caseLiteralExpression(literalExpression)
		return
	}
	referenceOrInvocation, ok := expression.(ReferenceOrInvocation)
	if ok {
		caseReferenceOrInvocation(referenceOrInvocation)
		return
	}
	lambda, ok := expression.(Lambda)
	if ok {
		caseLambda(lambda)
	}
	list, ok := expression.(List)
	if ok {
		caseList(list)
	}
	declaration, ok := expression.(Declaration)
	if ok {
		caseDeclaration(declaration)
		return
	}
	ifExp, ok := expression.(If)
	if ok {
		caseIf(ifExp)
		return
	}
	when, ok := expression.(When)
	if ok {
		caseWhen(when)
		return
	}
}

type When struct {
	parser.Node
	Over  ExpressionBox
	Is    []WhenIs
	Other *WhenOther
}

func (w When) sealedExpression() {}

type WhenIs struct {
	parser.Node
	Name      *Name
	Type      TypeAnnotation
	ThenBlock []ExpressionBox
}

type WhenOther struct {
	parser.Node
	Name      *Name
	ThenBlock []ExpressionBox
}

type If struct {
	parser.Node
	Condition ExpressionBox
	ThenBlock []ExpressionBox
	ElseIfs   []IfThen
	ElseBlock []ExpressionBox
}

type IfThen struct {
	parser.Node
	Condition ExpressionBox
	ThenBlock []ExpressionBox
}

func (i If) sealedExpression() {}

type Declaration struct {
	Name           Name
	TypeAnnotation *TypeAnnotation
	ExpressionBox  ExpressionBox
}

func (d Declaration) sealedTopLevelDeclaration() {}

func (d Declaration) sealedExpression() {}

type LiteralExpression struct {
	parser.Node
	Literal parser.Literal
}

func (l LiteralExpression) sealedExpression() {}

type List struct {
	parser.Node
	Generics    []TypeAnnotation
	Expressions []ExpressionBox
}

func (l List) sealedExpression() {}

type LambdaSignature struct {
	parser.Node
	Parameters []Parameter
	ReturnType *TypeAnnotation
}

type Lambda struct {
	parser.Node
	Generics  []TypeAnnotation
	Signature LambdaSignature
	Block     []ExpressionBox
}

func (l Lambda) sealedExpression() {}

type Parameter struct {
	Name Name
	Type *TypeAnnotation
}

type ReferenceOrInvocation struct {
	Var       Name
	Arguments *ArgumentsList
}

func (r ReferenceOrInvocation) sealedExpression() {}
