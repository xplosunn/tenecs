package codegen

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/ast"
)

func Generate(program *ast.Program) string {
	result := ""
	result += builtins()
	for _, declaration := range program.Declarations {
		result += GenerateDeclaration(declaration) + "\n"
	}
	return result
}

func builtins() string {
	return ""
}

func GenerateDeclaration(declaration *ast.Declaration) string {
	return "var P" + declaration.Name + " any = " + GenerateExpression(declaration.Expression) + "\n"
}

func GenerateExpression(expression ast.Expression) string {
	caseModule, caseLiteral, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf, caseArray, caseWhen := expression.ExpressionCases()
	if caseModule != nil {
		panic("TODO GenerateExpression caseModule")
	} else if caseLiteral != nil {
		panic("TODO GenerateExpression caseLiteral")
	} else if caseReferenceAndMaybeInvocation != nil {
		panic("TODO GenerateExpression caseReferenceAndMaybeInvocation")
	} else if caseWithAccessAndMaybeInvocation != nil {
		panic("TODO GenerateExpression caseWithAccessAndMaybeInvocation")
	} else if caseFunction != nil {
		panic("TODO GenerateExpression caseFunction")
	} else if caseDeclaration != nil {
		panic("TODO GenerateExpression caseDeclaration")
	} else if caseIf != nil {
		panic("TODO GenerateExpression caseIf")
	} else if caseArray != nil {
		panic("TODO GenerateExpression caseArray")
	} else if caseWhen != nil {
		panic("TODO GenerateExpression caseWhen")
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}
