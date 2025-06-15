package ir

import (
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/types"
)

func ToIR(program ast.Program) Program {
	declarations := map[Reference]TopLevelFunction{}
	structFunctions := map[Reference]*types.Function{}
	nativeFunctions := map[NativeFunctionRef]*types.Function{}

	return Program{
		Declarations:    declarations,
		StructFunctions: structFunctions,
		NativeFunctions: nativeFunctions,
	}
}
