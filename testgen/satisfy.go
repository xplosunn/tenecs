package testgen

import (
	"fmt"
	"github.com/xplosunn/tenecs/interpreter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
)

type Unsatisfiable struct {
}

func (u Unsatisfiable) Error() string {
	return "could not satisfy"
}

func satisfy(variableType types.VariableType, constraints []valueConstraint) (ast.Expression, error) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := variableType.VariableTypeCases()
	if caseTypeArgument != nil {
		panic("TODO satisfy caseTypeArgument")
	} else if caseStruct != nil {
		panic("TODO satisfy caseStruct")
	} else if caseInterface != nil {
		panic("TODO satisfy caseInterface")
	} else if caseFunction != nil {
		panic("TODO satisfy caseFunction")
	} else if caseBasicType != nil {
		return satisfyBasicType(*caseBasicType, constraints)
	} else if caseVoid != nil {
		panic("TODO satisfy caseVoid")
	} else {
		panic(fmt.Errorf("cases on %v", variableType))
	}
}

func satisfyBasicType(variableType types.BasicType, constraints []valueConstraint) (ast.Expression, error) {
	switch variableType.Type {
	case "Boolean":
		if len(constraints) == 0 {
			return ast.Literal{Literal: parser.LiteralBool{Value: true}}, nil
		}
		value := constraints[0].(valueConstraintEquals).To
		valueBoolean, ok := interpreter.ValueExpect[interpreter.ValueBoolean](value)
		result := valueBoolean.Bool
		if !ok {
			return nil, Unsatisfiable{}
		}
		for _, constraint := range constraints[1:] {
			value := constraint.(valueConstraintEquals).To
			valueBoolean, ok := interpreter.ValueExpect[interpreter.ValueBoolean](value)
			if !ok {
				return nil, Unsatisfiable{}
			}
			if valueBoolean.Bool != result {
				return nil, Unsatisfiable{}
			}
		}
		return ast.Literal{Literal: parser.LiteralBool{Value: result}}, nil
	case "String":
		panic("TODO satisfyBasicType String")
	case "Float":
		panic("TODO satisfyBasicType Float")
	case "Int":
		panic("TODO satisfyBasicType Int")
	default:
		return nil, fmt.Errorf("satisfyBasicType not implemented for %T", variableType)
	}
}
