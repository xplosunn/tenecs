package testgen

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
)

type Satisfier interface {
	impl() *satisfierImpl
}

var sampleStrings = []string{"fizz", "foo", "bar", "lorem", "ipsum"}

type satisfierImpl struct {
	program ast.Program

	strIndex int
}

func (s *satisfierImpl) impl() *satisfierImpl {
	return s
}

func NewSatisfier(program ast.Program) Satisfier {
	return &satisfierImpl{
		program: program,
	}
}

func satisfy(satisfier Satisfier, argName string, variableType types.VariableType, constraints []valueConstraint) (ast.Expression, error) {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO satisfy caseTypeArgument")
	} else if caseKnownType != nil {
		if caseKnownType.Package == "" {
			return satisfyBasicType(satisfier, argName, caseKnownType, constraints)
		}
		return satisfyKnownType(satisfier, argName, caseKnownType, constraints)
	} else if caseFunction != nil {
		return satisfyFunction(satisfier, argName, caseFunction, constraints)
	} else if caseOr != nil {
		panic("TODO satisfy caseOr")
	} else {
		panic(fmt.Errorf("cases on %v", variableType))
	}
}

type Unsatisfiable struct {
	argName      string
	variableType types.VariableType
	constraints  []valueConstraint
	reason       string
}

func (u Unsatisfiable) Error() string {
	return fmt.Sprintf("could not satisfy constraints for %s (reason: %s)", u.argName, u.reason)
}

func unsatisfiableError(argName string, variableType types.VariableType, constraints []valueConstraint, reason string) error {
	return Unsatisfiable{
		argName:      argName,
		variableType: variableType,
		constraints:  constraints,
		reason:       reason,
	}
}

func satisfyKnownType(satisfier Satisfier, argName string, variableType *types.KnownType, constraints []valueConstraint) (ast.Expression, error) {
	if len(constraints) != 0 {
		panic(fmt.Sprintf("TODO satisfyKnownType"))
	}
	constructorArgs := []ast.Expression{}
	for _, fieldVarType := range satisfier.impl().program.FieldsByType[variableType.Package+"->"+variableType.Name] {
		arg, err := satisfy(satisfier, argName, fieldVarType, []valueConstraint{})
		if err != nil {
			return nil, err
		}
		constructorArgs = append(constructorArgs, arg)
	}
	return ast.Invocation{
		VariableType: variableType,
		Over: ast.Reference{
			VariableType: nil,
			Name:         variableType.Name,
		},
		Generics:  nil,
		Arguments: constructorArgs,
	}, nil
}

func satisfyFunction(satisfier Satisfier, argName string, variableType *types.Function, constraints []valueConstraint) (ast.Expression, error) {
	resultConstraints := []valueConstraint{}
	for _, constraint := range constraints {
		c, ok := constraint.(valueConstraintFunctionInvocationResult)
		if !ok {
			panic(fmt.Sprintf("TODO satisfyFunction %T", constraint))
		}
		resultConstraints = append(resultConstraints, c.Constraint)
	}
	resultExp, err := satisfy(satisfier, argName, variableType.ReturnType, resultConstraints)
	if err != nil {
		return nil, err
	}
	return ast.Function{
		VariableType: variableType,
		Block:        []ast.Expression{resultExp},
	}, nil
}

func satisfyBasicType(satisfier Satisfier, argName string, variableType *types.KnownType, constraints []valueConstraint) (ast.Expression, error) {
	switch variableType.Name {
	case "Boolean":
		if len(constraints) == 0 {
			return ast.Literal{Literal: parser.LiteralBool{Value: true}}, nil
		}
		value := constraints[0].(valueConstraintEquals).To
		valueBoolean, ok := ValueExpect[ValueBoolean](value)
		if !ok {
			return nil, unsatisfiableError(argName, variableType, constraints, "can only do eq for bool")
		}
		result := valueBoolean.Bool
		for _, constraint := range constraints[1:] {
			value := constraint.(valueConstraintEquals).To
			valueBoolean, ok := ValueExpect[ValueBoolean](value)
			if !ok {
				return nil, unsatisfiableError(argName, variableType, constraints, "can only do eq for bool")
			}
			if valueBoolean.Bool != result {
				return nil, unsatisfiableError(argName, variableType, constraints, "can't satisfy both eq->true and eq->false")
			}
		}
		return ast.Literal{Literal: parser.LiteralBool{Value: result}}, nil
	case "String":
		if len(constraints) == 0 {
			satisfier.impl().strIndex = (satisfier.impl().strIndex + 1) % len(sampleStrings)
			return ast.Literal{Literal: parser.LiteralString{Value: "\"" + sampleStrings[satisfier.impl().strIndex] + "\""}}, nil
		}
		panic("TODO satisfyBasicType String")
	case "Float":
		panic("TODO satisfyBasicType Float")
	case "Int":
		panic("TODO satisfyBasicType Int")
	default:
		return nil, fmt.Errorf("satisfyBasicType not implemented for %T", variableType)
	}
}
