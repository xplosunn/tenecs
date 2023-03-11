package backtrack

import (
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/typer/types"
)

type Backtracker struct {
	CursorByReference *immutable.Map[string, Cursor]
}

func NewFromFunctionArguments(functionArguments []types.FunctionArgument) Backtracker {
	backtracker := Backtracker{
		CursorByReference: immutable.NewMap[string, Cursor](nil),
	}
	for _, arg := range functionArguments {
		backtracker = CopyAdding(backtracker, arg.Name, CursorSelf{Name: arg.Name})
	}
	return backtracker
}

func CopyAdding(backtracker Backtracker, name string, cursor Cursor) Backtracker {
	return Backtracker{
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
