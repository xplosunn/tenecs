package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_error = packageWith(
	withStruct("Error", tenecs_error_Error, tenecs_error_Error_Fields...),
	withFunction("error", tenecs_error_error),
)

var tenecs_error_Error = types.Struct(
	"tenecs.error",
	"Error",
	nil,
)

var tenecs_error_Error_Fields = []func(fields *StructWithFields){
	structField("message", types.String()),
}

var tenecs_error_error = &types.Function{
	Arguments: []types.FunctionArgument{
		types.FunctionArgument{
			Name:         "message",
			VariableType: types.String(),
		},
	},
	ReturnType: tenecs_error_Error,
}
