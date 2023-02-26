package testgen

import (
	"fmt"
	"github.com/xplosunn/tenecs/interpreter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
)

func satisfy(argName string, variableType types.VariableType, constraints []valueConstraint) (ast.Expression, error) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO satisfy caseTypeArgument")
	} else if caseStruct != nil {
		return satisfyStruct(argName, caseStruct, constraints)
	} else if caseInterface != nil {
		panic("TODO satisfy caseInterface")
	} else if caseFunction != nil {
		return satisfyFunction(argName, caseFunction, constraints)
	} else if caseBasicType != nil {
		return satisfyBasicType(argName, caseBasicType, constraints)
	} else if caseVoid != nil {
		panic("TODO satisfy caseVoid")
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

func satisfyStruct(argName string, variableType *types.Struct, constraints []valueConstraint) (ast.Expression, error) {
	if len(constraints) != 0 {
		panic(fmt.Sprintf("TODO satisfyStruct"))
	}
	constructorArgs := []ast.Expression{}
	for _, fieldVarType := range variableType.Fields {
		arg, err := satisfy(argName, types.VariableTypeFromStructFieldVariableType(fieldVarType), []valueConstraint{})
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

func satisfyFunction(argName string, variableType *types.Function, constraints []valueConstraint) (ast.Expression, error) {
	resultConstraints := []valueConstraint{}
	for _, constraint := range constraints {
		c, ok := constraint.(valueConstraintFunctionInvocationResult)
		if !ok {
			panic(fmt.Sprintf("TODO satisfyFunction %T", constraint))
		}
		resultConstraints = append(resultConstraints, c.Constraint)
	}
	resultExp, err := satisfy(argName, variableType.ReturnType, resultConstraints)
	if err != nil {
		return nil, err
	}
	return ast.Function{
		VariableType: variableType,
		Block:        []ast.Expression{resultExp},
	}, nil
}

func satisfyBasicType(argName string, variableType *types.BasicType, constraints []valueConstraint) (ast.Expression, error) {
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
			return ast.Literal{Literal: parser.LiteralString{Value: "\"foo\""}}, nil
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
