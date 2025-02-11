package standard_library

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/typer/types"
	"testing"
)

func TestFunctionFromSignature(t *testing.T) {
	assert.Equal(t, functionFromType("() ~> Void"), &types.Function{
		Generics:   []string{},
		Arguments:  []types.FunctionArgument{},
		ReturnType: types.Void(),
	})

}
