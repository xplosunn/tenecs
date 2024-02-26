package godsl

type Expression interface {
	sealedGoDSL()
	sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess, *goVariableReference, *goCast, *goLiteral)
}

func FunctionCreation(parameterNames ...string) Expression {
	return goFunctionCreation{parameterNames}
}

type goFunctionCreation struct {
	parameterNames []string
}

func (g goFunctionCreation) sealedGoDSL() {}

func (g goFunctionCreation) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess, *goVariableReference, *goCast, *goLiteral) {
	return &g, nil, nil, nil, nil, nil, nil
}

func FunctionInvocation(over Expression, arguments ...Expression) Expression {
	return goFunctionInvocation{over, arguments}
}

type goFunctionInvocation struct {
	over      Expression
	arguments []Expression
}

func (g goFunctionInvocation) sealedGoDSL() {}

func (g goFunctionInvocation) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess, *goVariableReference, *goCast, *goLiteral) {
	return nil, &g, nil, nil, nil, nil, nil
}

func ObjectField(name string, value Expression) func(*goObjectCreation) {
	return func(objectCreation *goObjectCreation) {
		objectCreation.fields[name] = value
	}
}

func ObjectCreation(typeName string, fields ...func(*goObjectCreation)) Expression {
	result := goObjectCreation{typeName, map[string]Expression{}}
	for _, fieldFunc := range fields {
		fieldFunc(&result)
	}
	return result
}

type goObjectCreation struct {
	typeName string
	fields   map[string]Expression
}

func (g goObjectCreation) sealedGoDSL() {}

func (g goObjectCreation) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess, *goVariableReference, *goCast, *goLiteral) {
	return nil, nil, &g, nil, nil, nil, nil
}

func ObjectAccess(over Expression, fieldName string) Expression {
	return goObjectAccess{over, fieldName}
}

type goObjectAccess struct {
	over      Expression
	fieldName string
}

func (g goObjectAccess) sealedGoDSL() {}

func (g goObjectAccess) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess, *goVariableReference, *goCast, *goLiteral) {
	return nil, nil, nil, &g, nil, nil, nil
}

func VariableReference(name string) Expression {
	return goVariableReference{name}
}

type goVariableReference struct {
	name string
}

func (g goVariableReference) sealedGoDSL() {}

func (g goVariableReference) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess, *goVariableReference, *goCast, *goLiteral) {
	return nil, nil, nil, nil, &g, nil, nil
}

func Cast(expression Expression, toType Type) Expression {
	return goCast{expression, toType}
}

type goCast struct {
	expression Expression
	toType     Type
}

func (g goCast) sealedGoDSL() {}

func (g goCast) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess, *goVariableReference, *goCast, *goLiteral) {
	return nil, nil, nil, nil, nil, &g, nil
}

func Literal(value string) Expression {
	return goLiteral{value}
}

type goLiteral struct {
	value string
}

func (g goLiteral) sealedGoDSL() {}

func (g goLiteral) sealedExpressionCases() (*goFunctionCreation, *goFunctionInvocation, *goObjectCreation, *goObjectAccess, *goVariableReference, *goCast, *goLiteral) {
	return nil, nil, nil, nil, nil, nil, &g
}
