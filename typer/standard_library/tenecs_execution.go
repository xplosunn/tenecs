package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_execution = packageWith(
	withInterface("Blocker", tenecs_execution_Blocker, tenecs_execution_Blocker_Fields),
)

var tenecs_execution_Blocker = types.Interface(
	"tenecs.execution",
	"Blocker",
	nil,
)

var tenecs_execution_Blocker_Fields = map[string]types.VariableType{}
