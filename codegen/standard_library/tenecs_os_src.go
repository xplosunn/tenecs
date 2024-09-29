package standard_library

func tenecs_os_Console() Function {
	return function(
		params("log"),
		body(`return map[string]any{
	"$type": "Console",
	"log": log,
}`),
	)
}
func tenecs_os_Main() Function {
	return function(
		params("main"),
		body(`return map[string]any{
	"$type": "Main",
	"main": main,
}`),
	)
}
func tenecs_os_Runtime() Function {
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
