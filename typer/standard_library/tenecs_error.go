package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_error = packageWith(
	withStruct(Tenecs_error_Error),
)

var Tenecs_error_Error = structWithFields("Error", tenecs_error_Error, tenecs_error_Error_Fields...)

var tenecs_error_Error = types.Struct(
	"tenecs.error",
	"Error",
	nil,
)

var tenecs_error_Error_Fields = []func(fields *StructWithFields){
	structField("message", types.String()),
}
