package testgen

import (
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/typer/types"
)

type scopeBacktracker struct {
	CursorByReference *immutable.Map[string, Cursor]
}

func NewScopeBacktrackerFromFunctionArguments(functionArguments []types.FunctionArgument) scopeBacktracker {
	backtracker := scopeBacktracker{
		CursorByReference: immutable.NewMap[string, Cursor](nil),
	}
	for _, arg := range functionArguments {
		backtracker = BacktrackerCopyAdding(backtracker, arg.Name, CursorSelf{Name: arg.Name})
	}
	return backtracker
}

func BacktrackerCopyAdding(backtracker scopeBacktracker, name string, cursor Cursor) scopeBacktracker {
	return scopeBacktracker{
		CursorByReference: backtracker.CursorByReference.Set(name, cursor),
	}
}

type Cursor interface {
	sealedCursor()
}

type CursorSelf struct {
	Name string
}

func (c CursorSelf) sealedCursor() {}
