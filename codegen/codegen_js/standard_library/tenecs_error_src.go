package standard_library

import "github.com/xplosunn/tenecs/typer/standard_library"

func tenecs_error_error() Function {
	return tenecs_error_Error()
}
func tenecs_error_Error() Function {
	return structFunction(standard_library.Tenecs_error_Error)
}
