package typer

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func TypecheckErrorFromResolutionError(node parser.Node, err *binding.ResolutionError) *type_error.TypecheckError {
	if err == nil {
		return nil
	}

	message := err.Problem
	if err.VariableType != nil {
		message += ": " + types.PrintableName(*err.VariableType)
	}
	return &type_error.TypecheckError{
		Node:    node,
		Message: message,
	}
}

func TypecheckErrorFromScopeCheckError(err scopecheck.ScopeCheckError) *type_error.TypecheckError {
	if err == nil {
		return nil
	}

	return &type_error.TypecheckError{
		Node:    err.Node(),
		Message: err.Error(),
	}
}
