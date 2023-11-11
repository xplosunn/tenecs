package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_execution = packageWith(
	withInterface("BlockingOperation", tenecs_execution_BlockingOperation, tenecs_execution_BlockingOperation_Fields),
)

var tenecs_execution_BlockingOperation = types.Interface(
	"tenecs.execution",
	"BlockingOperation",
	[]string{"R"},
)

func tenecs_execution_BlockingOperation_Of(varType types.VariableType) *types.KnownType {
	return types.UncheckedApplyGenerics(tenecs_execution_BlockingOperation, []types.VariableType{varType})
}

var tenecs_execution_BlockingOperation_Fields = map[string]types.VariableType{
	"fakeRun": &types.Function{
		Arguments:  []types.FunctionArgument{},
		ReturnType: &types.TypeArgument{Name: "R"},
	},
}
