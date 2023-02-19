package interpreter

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
)

func EvalBlock(scope Scope, expressions []ast.Expression) (Scope, Value, error) {
	unchangedScope := scope
	var value Value = ValueVoid{}
	var err error
	for _, expression := range expressions {
		scope, value, err = EvalExpression(scope, expression)
		if err != nil {
			return nil, nil, err
		}
	}
	return unchangedScope, value, nil
}

func EvalExpression(scope Scope, expression ast.Expression) (Scope, Value, error) {
	caseModule, caseLiteral, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseModule != nil {
		panic("TODO EvalExpression caseModule")
	} else if caseLiteral != nil {
		return EvalLiteral(scope, *caseLiteral)
	} else if caseReferenceAndMaybeInvocation != nil {
		return EvalReferenceAndMaybeInvocation(scope, *caseReferenceAndMaybeInvocation)
	} else if caseWithAccessAndMaybeInvocation != nil {
		panic("TODO EvalExpression caseWithAccessAndMaybeInvocation")
	} else if caseFunction != nil {
		return EvalFunction(scope, *caseFunction)
	} else if caseDeclaration != nil {
		return EvalDeclaration(scope, *caseDeclaration)
	} else if caseIf != nil {
		return EvalIf(scope, *caseIf)
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func EvalReferenceAndMaybeInvocation(scope Scope, expression ast.ReferenceAndMaybeInvocation) (Scope, Value, error) {
	referencedValue, err := Resolve(scope, expression.Name)
	if err != nil {
		return nil, nil, err
	}
	if expression.ArgumentsList == nil {
		return scope, referencedValue, nil
	}
	referencedFunction, ok := referencedValue.(ValueFunction)
	if !ok {
		return nil, nil, fmt.Errorf("expected %s to be a function so an invocation can be made but it's %T", expression.Name, referencedValue)
	}
	return EvalFunctionInvocation(scope, referencedFunction, expression.ArgumentsList.Arguments)
}

func EvalFunctionInvocation(scope Scope, function ValueFunction, arguments []ast.Expression) (Scope, Value, error) {
	invocationScope := scope
	for i, argument := range arguments {
		_, value, err := EvalExpression(scope, argument)
		if err != nil {
			return nil, nil, err
		}
		invocationScope = CopyAdding(invocationScope, function.AstFunction.VariableType.Arguments[i].Name, value)
	}
	return EvalBlock(invocationScope, function.AstFunction.Block)
}

func EvalFunction(scope Scope, expression ast.Function) (Scope, Value, error) {
	value := ValueFunction{
		Scope:       scope,
		AstFunction: expression,
	}
	return scope, value, nil
}

func EvalDeclaration(scope Scope, expression ast.Declaration) (Scope, Value, error) {
	_, value, err := EvalExpression(scope, expression.Expression)
	if err != nil {
		return nil, nil, err
	}
	scope = CopyAdding(scope, expression.Name, value)
	return scope, ValueVoid{}, nil
}

func EvalLiteral(scope Scope, expression ast.Literal) (Scope, Value, error) {
	value := parser.LiteralFold[Value](
		expression.Literal,
		func(arg float64) Value { return ValueFloat{Float: arg} },
		func(arg int) Value { return ValueInt{Int: arg} },
		func(arg string) Value { return ValueString{String: arg} },
		func(arg bool) Value { return ValueBoolean{Bool: arg} },
	)
	return scope, value, nil
}

func EvalIf(scope Scope, expression ast.If) (Scope, Value, error) {
	_, conditionValue, err := EvalExpression(scope, expression.Condition)
	if err != nil {
		return nil, nil, err
	}
	conditionBoolean, ok := conditionValue.(ValueBoolean)
	if !ok {
		return nil, nil, fmt.Errorf("expected to eval Boolean on if condition but got %T", conditionValue)
	}
	if conditionBoolean.Bool {
		return EvalBlock(scope, expression.ThenBlock)
	} else {
		return EvalBlock(scope, expression.ElseBlock)
	}
}
