package standard_library

import "github.com/xplosunn/tenecs/typer/standard_library"

func tenecs_go_Console() Function {
	return structFunction(standard_library.Tenecs_go_Console)
}
func tenecs_go_Main() Function {
	return structFunction(standard_library.Tenecs_go_Main)
}
func tenecs_go_Runtime() Function {
	return structFunction(standard_library.Tenecs_go_Runtime)
}
