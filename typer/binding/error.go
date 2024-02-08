package binding

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/types"
)

type ResolutionError struct {
	VariableType *types.VariableType
	Problem      string
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

func (err *ResolutionError) Error() string {
	return err.Problem
}
