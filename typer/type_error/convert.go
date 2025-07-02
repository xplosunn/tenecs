package type_error

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/types"
)

func FromResolutionError(file string, node parser.Node, err *binding.ResolutionError) *TypecheckError {
	if err == nil {
		return nil
	}

	message := err.Problem
	if err.VariableType != nil {
		message += ": " + types.PrintableName(*err.VariableType)
	}
	return &TypecheckError{
		File:    file,
		Node:    node,
		Message: message,
	}
}

func FromScopeCheckError(file string, err scopecheck.ScopeCheckError) *TypecheckError {
	if err == nil {
		return nil
	}

	return &TypecheckError{
		File:    file,
		Node:    err.Node(),
		Message: err.Error(),
	}
}
