package testgen

import (
	"errors"
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/interpreter"
	"github.com/xplosunn/tenecs/typer/ast"
)

type testCaseConstraints struct {
	argsConstraints *immutable.Map[string, []valueConstraint]
}

type valueConstraint interface {
	sealedValueConstraint()
}

type valueConstraintEquals struct {
	To interpreter.Value
}

func (v valueConstraintEquals) sealedValueConstraint() {}

type valueConstraintFunctionInvocationResult struct {
	Constraint valueConstraint
}

func (v valueConstraintFunctionInvocationResult) sealedValueConstraint() {}

func findConstraints(function *ast.Function) ([]testCaseConstraints, error) {
	backtracker := NewScopeBacktrackerFromFunctionArguments(function.VariableType.Arguments)
	return findConstraintsOverExpressions(backtracker, function.Block)
}

func findConstraintsOverExpressions(backtracker scopeBacktracker, expressions []ast.Expression) ([]testCaseConstraints, error) {
	if len(expressions) == 0 {
		return []testCaseConstraints{}, nil
	}
	expression, remainingExpressions := expressions[0], expressions[1:]
	caseModule, caseLiteral, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseModule != nil {
		return nil, errors.New("todo findConstraintsOverExpressions caseModule")
	} else if caseLiteral != nil {
		return findConstraintsOverExpressions(backtracker, remainingExpressions)
	} else if caseReferenceAndMaybeInvocation != nil {
		if caseReferenceAndMaybeInvocation.ArgumentsList == nil {
			return findConstraintsOverExpressions(backtracker, remainingExpressions)
		} else {
			return nil, errors.New("todo findConstraintsOverExpressions caseReferenceAndMaybeInvocation")
		}
	} else if caseWithAccessAndMaybeInvocation != nil {
		return nil, errors.New("todo findConstraintsOverExpressions caseWithAccessAndMaybeInvocation")
	} else if caseFunction != nil {
		return nil, errors.New("todo findConstraintsOverExpressions caseFunction")
	} else if caseDeclaration != nil {
		constraintsOverExp, err := findConstraintsOverExpressions(backtracker, []ast.Expression{caseDeclaration.Expression})
		if err != nil {
			return nil, err
		}
		cursor, err := findCursorOverExpression(backtracker, expression)
		if err != nil {
			return nil, err
		}
		if cursor != nil {
			backtracker = BacktrackerCopyAdding(backtracker, caseDeclaration.Name, *cursor)
		}
		remainingConstraints, err := findConstraintsOverExpressions(backtracker, remainingExpressions)
		if err != nil {
			return nil, err
		}
		if len(constraintsOverExp) == 0 {
			return remainingConstraints, nil
		} else if len(remainingConstraints) == 0 {
			return constraintsOverExp, nil
		} else {
			resultConstraints := []testCaseConstraints{}
			for _, testCase := range constraintsOverExp {
				for _, remainingTestCase := range remainingConstraints {
					resultConstraints = append(resultConstraints, testCaseConstraintsMerge(testCase, remainingTestCase))
				}

			}
			return resultConstraints, nil
		}
	} else if caseIf != nil {
		trueConstraint, err := applyConstraintToExpression(backtracker, valueConstraintEquals{To: interpreter.ValueBoolean{Bool: true}}, caseIf.Condition)
		if err != nil {
			return nil, err
		}

		thenConstraint, err := findConstraintsOverExpressions(backtracker, caseIf.ThenBlock)
		if err != nil {
			return nil, err
		}

		falseConstraint, err := applyConstraintToExpression(backtracker, valueConstraintEquals{To: interpreter.ValueBoolean{Bool: false}}, caseIf.Condition)
		if err != nil {
			return nil, err
		}

		elseConstraint, err := findConstraintsOverExpressions(backtracker, caseIf.ElseBlock)
		if err != nil {
			return nil, err
		}

		remainingConstraints, err := findConstraintsOverExpressions(backtracker, remainingExpressions)

		return append(
			testCaseConstraintsCombine(
				testCaseConstraintMergeWithEach(thenConstraint, trueConstraint),
				remainingConstraints,
			),
			testCaseConstraintsCombine(
				testCaseConstraintMergeWithEach(elseConstraint, falseConstraint),
				remainingConstraints,
			)...,
		), nil
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func findCursorOverExpression(backtracker scopeBacktracker, expression ast.Expression) (*Cursor, error) {
	caseModule, caseLiteral, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseModule != nil {
		return nil, nil
	} else if caseLiteral != nil {
		return nil, nil
	} else if caseReferenceAndMaybeInvocation != nil {
		return nil, errors.New("todo findCursorOverExpression caseReferenceAndMaybeInvocation")
	} else if caseWithAccessAndMaybeInvocation != nil {
		return nil, errors.New("todo findCursorOverExpression caseWithAccessAndMaybeInvocation")
	} else if caseFunction != nil {
		return nil, nil
	} else if caseDeclaration != nil {
		return nil, nil
	} else if caseIf != nil {
		return nil, errors.New("todo findCursorOverExpression caseIf")
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func applyConstraintToCursor(cursor Cursor, constraint valueConstraint, functionApplication bool) testCaseConstraints {
	cursorSelf, ok := cursor.(CursorSelf)
	if !ok {
		panic("applyConstraintToCursor")
	}

	builder := immutable.NewMapBuilder[string, []valueConstraint](nil)
	if functionApplication {
		builder.Set(cursorSelf.Name, []valueConstraint{valueConstraintFunctionInvocationResult{constraint}})
	} else {
		builder.Set(cursorSelf.Name, []valueConstraint{constraint})
	}

	return testCaseConstraints{
		argsConstraints: builder.Map(),
	}
}

func applyConstraintToExpression(backtracker scopeBacktracker, constraint valueConstraint, expression ast.Expression) (testCaseConstraints, error) {
	var emptyResult testCaseConstraints
	caseModule, caseLiteral, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseFunction, caseDeclaration, caseIf := expression.ExpressionCases()
	if caseModule != nil {
		return emptyResult, errors.New("todo applyConstraintToExpression caseModule")
	} else if caseLiteral != nil {
		return emptyResult, errors.New("todo applyConstraintToExpression caseLiteral")
	} else if caseReferenceAndMaybeInvocation != nil {
		cursor, ok := backtracker.CursorByReference.Get(caseReferenceAndMaybeInvocation.Name)
		if !ok {
			return emptyResult, errors.New("no cursor found on applyConstraintToExpression caseReferenceAndMaybeInvocation")
		}
		return applyConstraintToCursor(cursor, constraint, caseReferenceAndMaybeInvocation.ArgumentsList != nil), nil

	} else if caseWithAccessAndMaybeInvocation != nil {
		return emptyResult, errors.New("todo applyConstraintToExpression caseWithAccessAndMaybeInvocation")
	} else if caseFunction != nil {
		return emptyResult, errors.New("todo applyConstraintToExpression caseFunction")
	} else if caseDeclaration != nil {
		return emptyResult, errors.New("todo applyConstraintToExpression caseDeclaration")
	} else if caseIf != nil {
		return emptyResult, errors.New("todo applyConstraintToExpression caseIf")
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func testCaseConstraintMergeWithEach(constraints []testCaseConstraints, toMerge testCaseConstraints) []testCaseConstraints {
	if len(constraints) == 0 {
		return []testCaseConstraints{toMerge}
	}
	result := []testCaseConstraints{}
	for _, testCase := range constraints {
		result = append(result, testCaseConstraintsMerge(testCase, toMerge))
	}
	return result
}

func testCaseConstraintsCombine(constraints []testCaseConstraints, otherConstraints []testCaseConstraints) []testCaseConstraints {
	if len(constraints) == 0 {
		return otherConstraints
	}
	if len(otherConstraints) == 0 {
		return constraints
	}
	result := []testCaseConstraints{}
	for _, testCase := range constraints {
		for _, otherConstraints := range otherConstraints {
			testCaseConstraintsMerge(testCase, otherConstraints)
		}
	}
	return result
}

func testCaseConstraintsMerge(constraints testCaseConstraints, others ...testCaseConstraints) testCaseConstraints {
	if len(others) == 0 {
		return constraints
	}
	for _, otherConstraints := range others {
		iterator := otherConstraints.argsConstraints.Iterator()
		for !iterator.Done() {
			argToConstraint, valueConstraints, _ := iterator.Next()
			existingConstraints, ok := constraints.argsConstraints.Get(argToConstraint)
			if ok {
				constraints = testCaseConstraints{
					argsConstraints: constraints.argsConstraints.Set(argToConstraint, append(existingConstraints, valueConstraints...)),
				}
			} else {
				constraints = testCaseConstraints{
					argsConstraints: constraints.argsConstraints.Set(argToConstraint, valueConstraints),
				}
			}
		}
	}
	return constraints
}
