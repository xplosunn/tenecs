package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"text/scanner"
)

func ParseString(s string) (*FileTopLevel, error) {
	p, err := parser()
	if err != nil {
		return nil, err
	}

	res, err := p.ParseString("", s)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Grammar() (string, error) {
	p, err := parser()
	if err != nil {
		return "", err
	}
	return p.String(), nil
}

func parser() (*participle.Parser[FileTopLevel], error) {
	l := lexer.NewTextScannerLexer(func(s *scanner.Scanner) {
		s.Mode = s.Mode - scanner.SkipComments
	})

	return participle.Build[FileTopLevel](participle.Lexer(l), participle.Elide("Comment"), topLevelDeclarationUnion, typeAnnotationElementUnion, literalUnion, expressionUnion)
}

type Node struct {
	Pos    lexer.Position
	EndPos lexer.Position
}

type FileTopLevel struct {
	Tokens               []lexer.Token
	Package              Package               `@@`
	Imports              []Import              `@@*`
	TopLevelDeclarations []TopLevelDeclaration `@@*`
}

func FileTopLevelFields(node FileTopLevel) (Package, []Import, []TopLevelDeclaration) {
	return node.Package, node.Imports, node.TopLevelDeclarations
}

type Name struct {
	Node
	String string `@Ident`
}

type Package struct {
	Node
	DotSeparatedNames []Name `"package" (@@ ("." @@)*)?`
}

type Import struct {
	Node
	DotSeparatedVars []Name `"import" (@@ ("." @@)*)?`
	As               *Name  `("as" @@)?`
}

func ImportFields(node Import) ([]Name, *Name) {
	return node.DotSeparatedVars, node.As
}

type TopLevelDeclaration interface {
	sealedTopLevelDeclaration()
}

func TopLevelDeclarationExhaustiveSwitch(
	topLevelDeclaration TopLevelDeclaration,
	caseDeclaration func(topLevelDeclaration Declaration),
	caseInterface func(topLevelDeclaration Interface),
	caseStruct func(topLevelDeclaration Struct),
	caseTypeAlias func(topLevelDeclaration TypeAlias),
) {
	declaration, ok := topLevelDeclaration.(Declaration)
	if ok {
		caseDeclaration(declaration)
		return
	}
	interf, ok := topLevelDeclaration.(Interface)
	if ok {
		caseInterface(interf)
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

var topLevelDeclarationUnion = participle.Union[TopLevelDeclaration](Struct{}, Interface{}, TypeAlias{}, Declaration{})

type Struct struct {
	Name      Name             `"struct" @@`
	Generics  []Name           `("<" (@@ ("," @@)*)? ">")?`
	Variables []StructVariable `"(" (@@ ("," @@)*)? ")"`
}

func (s Struct) sealedTopLevelDeclaration() {}

func StructFields(struc Struct) (Name, []Name, []StructVariable) {
	return struc.Name, struc.Generics, struc.Variables
}

type StructVariable struct {
	Name Name           `@@ ":"`
	Type TypeAnnotation `@@`
}

func StructVariableFields(structVariable StructVariable) (Name, TypeAnnotation) {
	return structVariable.Name, structVariable.Type
}

type TypeAlias struct {
	Name     Name           `"typealias" @@`
	Generics []Name         `("<" (@@ ("," @@)*)? ">")?`
	Type     TypeAnnotation `"=" @@`
}

func (i TypeAlias) sealedTopLevelDeclaration() {}

func TypeAliasFields(typeAlias TypeAlias) (Name, []Name, TypeAnnotation) {
	return typeAlias.Name, typeAlias.Generics, typeAlias.Type
}

type Interface struct {
	Name      Name                `"interface" @@`
	Generics  []Name              `("<" (@@ ("," @@)*)? ">")?`
	Variables []InterfaceVariable `"{" @@* "}"`
}

func (i Interface) sealedTopLevelDeclaration() {}

func InterfaceFields(interf Interface) (Name, []Name, []InterfaceVariable) {
	return interf.Name, interf.Generics, interf.Variables
}

type InterfaceVariable struct {
	Name Name           `@@`
	Type TypeAnnotation `":" @@`
}

func InterfaceVariableFields(interfaceVariable InterfaceVariable) (Name, TypeAnnotation) {
	return interfaceVariable.Name, interfaceVariable.Type
}

type TypeAnnotation struct {
	Node
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

var typeAnnotationElementUnion = participle.Union[TypeAnnotationElement](SingleNameType{}, FunctionType{})

type SingleNameType struct {
	Node
	TypeName Name             `@@`
	Generics []TypeAnnotation `("<" @@ ("," @@)* ">")?`
}

func (s SingleNameType) sealedTypeAnnotationElement() {}

type FunctionType struct {
	Generics   []Name           `("<" @@ ("," @@)* ">")?`
	Arguments  []TypeAnnotation `"(" (@@ ("," @@)*)? ")"`
	ReturnType TypeAnnotation   `"-" ">" @@`
}

func (f FunctionType) sealedTypeAnnotationElement() {}

type Implementation struct {
	Node
	Implementing Name                        `"implement" @@`
	Generics     []TypeAnnotation            `("<" @@ ("," @@)* ">")?`
	Declarations []ImplementationDeclaration `"{" @@* "}"`
}

func (m Implementation) sealedExpression() {}

func ImplementationFields(node Implementation) (Name, []TypeAnnotation, []ImplementationDeclaration) {
	return node.Implementing, node.Generics, node.Declarations
}

type ImplementationDeclaration struct {
	Name           Name            `@@`
	TypeAnnotation *TypeAnnotation `":" @@?`
	Expression     Expression      `"=" @@`
}

func ImplementationDeclarationFields(node ImplementationDeclaration) (Name, *TypeAnnotation, Expression) {
	return node.Name, node.TypeAnnotation, node.Expression
}

type ArgumentsList struct {
	Node
	Generics  []TypeAnnotation `("<" @@ ("," @@)* ">")?`
	Arguments []NamedArgument  `"(" (@@ ("," @@)*)? ")"`
}

type NamedArgument struct {
	Node
	Name     *Name         `(@@ "=")?`
	Argument ExpressionBox `@@`
}

func NamedArgumentFields(namedArgument NamedArgument) (*Name, ExpressionBox) {
	return namedArgument.Name, namedArgument.Argument
}

type AccessOrInvocation struct {
	VarName   *Name          `("." @@`
	Arguments *ArgumentsList `@@?) | @@`
}

type ExpressionBox struct {
	Node
	Expression              Expression           `@@`
	AccessOrInvocationChain []AccessOrInvocation `@@*`
}

func ExpressionBoxFields(expressionBox ExpressionBox) (Expression, []AccessOrInvocation) {
	return expressionBox.Expression, expressionBox.AccessOrInvocationChain
}

type Expression interface {
	sealedExpression()
}

func ExpressionExhaustiveSwitch(
	expression Expression,
	caseImplementation func(expression Implementation),
	caseLiteralExpression func(expression LiteralExpression),
	caseReferenceOrInvocation func(expression ReferenceOrInvocation),
	caseLambda func(expression Lambda),
	caseDeclaration func(expression Declaration),
	caseIf func(expression If),
	caseList func(expression List),
	caseWhen func(expression When),
) {
	implementation, ok := expression.(Implementation)
	if ok {
		caseImplementation(implementation)
		return
	}
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
		return
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
	list, ok := expression.(List)
	if ok {
		caseList(list)
		return
	}
	when, ok := expression.(When)
	if ok {
		caseWhen(when)
		return
	}
}

var expressionUnion = participle.Union[Expression](When{}, Implementation{}, If{}, Declaration{}, LiteralExpression{}, ReferenceOrInvocation{}, Lambda{}, List{})

type List struct {
	Node
	Generic     *TypeAnnotation `"[" @@? "]"`
	Expressions []ExpressionBox `"(" (@@ ("," @@)*)? ")"`
}

func (a List) sealedExpression() {}

type When struct {
	Node
	Over  ExpressionBox `"when" @@ "{"`
	Is    []WhenIs      `@@*`
	Other *WhenOther    `@@? "}"`
}

func (w When) sealedExpression() {}

type WhenIs struct {
	Node
	Name      *Name           `"is" (@@ ":")?`
	Type      TypeAnnotation  `@@`
	ThenBlock []ExpressionBox `"=" ">" "{" @@* "}"`
}

type WhenOther struct {
	Node
	Name      *Name           `"other" (@@)?`
	ThenBlock []ExpressionBox `"=" ">" "{" @@* "}"`
}

type If struct {
	Node
	Condition ExpressionBox   `"if" @@`
	ThenBlock []ExpressionBox `"{" @@* "}"`
	ElseIfs   []IfThen        `("else" @@)*`
	ElseBlock []ExpressionBox `("else" "{" @@* "}")?`
}

type IfThen struct {
	Node
	Condition ExpressionBox   `"if" @@`
	ThenBlock []ExpressionBox `"{" @@* "}"`
}

func (i If) sealedExpression() {}

func IfFields(parserIf If) (ExpressionBox, []ExpressionBox, []IfThen, []ExpressionBox) {
	return parserIf.Condition, parserIf.ThenBlock, parserIf.ElseIfs, parserIf.ElseBlock
}

type DeclarationShortCircuit struct {
	TypeAnnotation *TypeAnnotation `"?" @@?`
}

type Declaration struct {
	Name           Name                     `@@`
	TypeAnnotation *TypeAnnotation          `":" @@?`
	ShortCircuit   *DeclarationShortCircuit `@@?`
	ExpressionBox  ExpressionBox            `"=" @@`
}

func (d Declaration) sealedTopLevelDeclaration() {}

func (d Declaration) sealedExpression() {}

func DeclarationFields(node Declaration) (Name, *TypeAnnotation, *DeclarationShortCircuit, ExpressionBox) {
	return node.Name, node.TypeAnnotation, node.ShortCircuit, node.ExpressionBox
}

type LiteralExpression struct {
	Node
	Literal Literal `@@`
}

func (l LiteralExpression) sealedExpression() {}

type Lambda struct {
	Node
	Generics   []Name          `("<" @@ ("," @@)* ">")?`
	Parameters []Parameter     `"(" (@@ ("," @@)*)? ")"`
	ReturnType *TypeAnnotation `(":" @@)?`
	Block      []ExpressionBox `"=" ">" (("{" @@* "}") | @@)`
}

func (l Lambda) sealedExpression() {}

func LambdaFields(node Lambda) ([]Name, []Parameter, *TypeAnnotation, []ExpressionBox) {
	return node.Generics, node.Parameters, node.ReturnType, node.Block
}

type Parameter struct {
	Name Name            `@@`
	Type *TypeAnnotation `(":" @@)?`
}

func ParameterFields(node Parameter) (Name, *TypeAnnotation) {
	return node.Name, node.Type
}

type ReferenceOrInvocation struct {
	Var       Name           `@@`
	Arguments *ArgumentsList `@@?`
}

func (r ReferenceOrInvocation) sealedExpression() {}

func ReferenceOrInvocationFields(node ReferenceOrInvocation) (Name, *ArgumentsList) {
	return node.Var, node.Arguments
}
