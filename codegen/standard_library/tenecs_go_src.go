package standard_library

func tenecs_go_Console() Function {
	return function(
		params("log"),
		body(`return map[string]any{
	"$type": "Console",
	"log": log,
}`),
	)
}
func tenecs_go_Main() Function {
	return function(
		params("main"),
		body(`return map[string]any{
	"$type": "Main",
	"main": main,
}`),
	)
}
func tenecs_go_Runtime() Function {
	return function(
		params("console", "http", "ref"),
		body(`return map[string]any{
	"$type": "Runtime",
	"console": console,
	"http": http,
	"ref": ref,
}`),
	)
}
