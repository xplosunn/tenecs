package testgen

import (
	"fmt"
	"github.com/xplosunn/tenecs/interpreter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
)

type Satisfier interface {
	impl() *satisfierImpl
}

var sampleStrings = []string{"fizz", "foo", "bar", "lorem", "ipsum"}

type satisfierImpl struct {
	strIndex int
}

func (s *satisfierImpl) impl() *satisfierImpl {
	return s
}

func NewSatisfier() Satisfier {
	return &satisfierImpl{}
}

func satisfy(satisfier Satisfier, argName string, variableType types.VariableType, constraints []valueConstraint) (ast.Expression, error) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid, caseArray := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO satisfy caseTypeArgument")
	} else if caseStruct != nil {
		return satisfyStruct(satisfier, argName, caseStruct, constraints)
	} else if caseInterface != nil {
		panic("TODO satisfy caseInterface")
	} else if caseFunction != nil {
		return satisfyFunction(satisfier, argName, caseFunction, constraints)
	} else if caseBasicType != nil {
		return satisfyBasicType(satisfier, argName, caseBasicType, constraints)
	} else if caseVoid != nil {
		panic("TODO satisfy caseVoid")
	} else if caseArray != nil {
		panic("TODO satisfy caseArray")
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

func satisfyStruct(satisfier Satisfier, argName string, variableType *types.Struct, constraints []valueConstraint) (ast.Expression, error) {
	if len(constraints) != 0 {
		panic(fmt.Sprintf("TODO satisfyStruct"))
	}
	constructorArgs := []ast.Expression{}
	for _, fieldVarType := range variableType.Fields {
		arg, err := satisfy(satisfier, argName, types.VariableTypeFromStructFieldVariableType(fieldVarType), []valueConstraint{})
		if err != nil {
			return nil, err
		}
		constructorArgs = append(constructorArgs, arg)
	}
	return ast.ReferenceAndMaybeInvocation{
		VariableType: variableType,
		Name:         variableType.Name,
		ArgumentsList: &ast.ArgumentsList{
			Arguments: constructorArgs,
		},
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

func satisfyBasicType(satisfier Satisfier, argName string, variableType *types.BasicType, constraints []valueConstraint) (ast.Expression, error) {
	switch variableType.Type {
	case "Boolean":
		if len(constraints) == 0 {
			return ast.Literal{Literal: parser.LiteralBool{Value: true}}, nil
		}
		value := constraints[0].(valueConstraintEquals).To
		valueBoolean, ok := interpreter.ValueExpect[interpreter.ValueBoolean](value)
		if !ok {
			return nil, unsatisfiableError(argName, variableType, constraints, "can only do eq for bool")
		}
		result := valueBoolean.Bool
		for _, constraint := range constraints[1:] {
			value := constraint.(valueConstraintEquals).To
			valueBoolean, ok := interpreter.ValueExpect[interpreter.ValueBoolean](value)
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
