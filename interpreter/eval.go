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
	caseModule, caseLiteral, caseReference, caseAccess, caseInvocation, caseFunction, caseDeclaration, caseIf, caseArray, caseWhen := expression.ExpressionCases()
	if caseModule != nil {
		panic("TODO EvalExpression caseModule")
	} else if caseLiteral != nil {
		return EvalLiteral(scope, *caseLiteral)
	} else if caseReference != nil {
		return EvalReference(scope, *caseReference)
	} else if caseAccess != nil {
		return EvalAccess(scope, *caseAccess)
	} else if caseInvocation != nil {
		return EvalInvocation(scope, *caseInvocation)
	} else if caseFunction != nil {
		return EvalFunction(scope, *caseFunction)
	} else if caseDeclaration != nil {
		return EvalDeclaration(scope, *caseDeclaration)
	} else if caseIf != nil {
		return EvalIf(scope, *caseIf)
	} else if caseArray != nil {
		panic("TODO EvalExpression array")
	} else if caseWhen != nil {
		panic("TODO EvalExpression when")
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func EvalAccess(scope Scope, expression ast.Access) (Scope, Value, error) {
	_, value, err := EvalExpression(scope, expression.Over)
	if err != nil {
		return nil, nil, err
	}

	valueStruct, ok := value.(ValueStruct)
	if !ok {
		return nil, nil, fmt.Errorf("Eval expected struct for access but got %T", value)
	}
	value = valueStruct.KeyValues[expression.Access]

	return scope, value, nil
}

func EvalReference(scope Scope, expression ast.Reference) (Scope, Value, error) {
	referencedValue, err := Resolve(scope, expression.Name)
	if err != nil {
		return nil, nil, err
	}
	return scope, referencedValue, nil
}

func EvalInvocation(scope Scope, expression ast.Invocation) (Scope, Value, error) {
	scope, referencedValue, err := EvalExpression(scope, expression.Over)
	if err != nil {
		return nil, nil, err
	}

	referencedFunction, ok := referencedValue.(ValueFunction)
	if ok {
		return EvalFunctionInvocation(scope, referencedFunction, expression.Arguments)
	}
	referencedStructFunction, ok := referencedValue.(ValueStructFunction)
	if ok {
		argValues := []Value{}
		for _, argument := range expression.Arguments {
			_, value, err := EvalExpression(scope, argument)
			if err != nil {
				return nil, nil, err
			}
			argValues = append(argValues, value)
		}
		return scope, referencedStructFunction.Create(argValues), nil
	}
	referencedNativeFunction, ok := referencedValue.(ValueNativeFunction)
	if ok {
		argValues := []Value{}
		for _, argument := range expression.Arguments {
			_, value, err := EvalExpression(scope, argument)
			if err != nil {
				return nil, nil, err
			}
			argValues = append(argValues, value)
		}
		return scope, referencedNativeFunction.Invoke(expression.Generics, argValues), nil
	}
	return nil, nil, fmt.Errorf("expected to be a function so an invocation can be made but it's %T", referencedValue)
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
	var value Value
	parser.LiteralExhaustiveSwitch(
		expression.Literal,
		func(literal float64) { value = ValueFloat{Float: literal} },
		func(literal int) { value = ValueInt{Int: literal} },
		func(literal string) { value = ValueString{String: literal} },
		func(literal bool) { value = ValueBoolean{Bool: literal} },
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
