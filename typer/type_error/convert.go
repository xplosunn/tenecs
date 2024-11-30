package type_error

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/types"
)

func FromResolutionError(node parser.Node, err *binding.ResolutionError) *TypecheckError {
	if err == nil {
		return nil
	}

	message := err.Problem
	if err.VariableType != nil {
		message += ": " + types.PrintableName(*err.VariableType)
	}
	return &TypecheckError{
		Node:    node,
		Message: message,
	}
}

func FromScopeCheckError(err scopecheck.ScopeCheckError) *TypecheckError {
	if err == nil {
		return nil
	}

	return &TypecheckError{
		Node:    err.Node(),
		Message: err.Error(),
	}
}
