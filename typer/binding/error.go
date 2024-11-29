package binding

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/types"
)

type ResolutionError struct {
	VariableType *types.VariableType
	Problem      string
}

func (err ResolutionError) Error() string {
	return err.Problem
}

func ResolutionErrorCouldNotResolve(typeName string) *ResolutionError {
	return &ResolutionError{
		VariableType: nil,
		Problem:      "not found type: " + typeName,
	}
}

func ResolutionErrorWrongNumberOfGenerics(variableType types.VariableType, expected int, got int) *ResolutionError {
	return &ResolutionError{
		VariableType: &variableType,
		Problem:      fmt.Sprintf("wrong number of generics, expected %d but got %d", expected, got),
	}
}

func ResolutionErrorTypeAlreadyExists(variableType types.VariableType) *ResolutionError {
	return &ResolutionError{
		VariableType: &variableType,
		Problem:      "type already exists",
	}
}

func ResolutionErrorTypeFieldsAlreadyExists(typeNameString string) *ResolutionError {
	return &ResolutionError{
		VariableType: nil,
		Problem:      "type fields already exist: " + typeNameString,
	}
}

func ResolutionErrorVariableAlreadyExists(varType types.VariableType, varName string) *ResolutionError {
	return &ResolutionError{
		VariableType: &varType,
		Problem:      "duplicate variable '" + varName + "'",
	}
}
