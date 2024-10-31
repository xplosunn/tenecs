package standard_library

func tenecs_go_Console() Function {
	return function(
		params("log"),
		body(`return ({
  "$type": "Console",
  "log": log
})`),
	)
}
func tenecs_go_Main() Function {
	return function(
		params("main"),
		body(`return ({
  "$type": "Main",
  "main": main
})`),
	)
}
func tenecs_go_Runtime() Function {
	return function(
		params("console", "http", "ref"),
		body(`return ({
  "$type": "Runtime",
  "console": console,
  "http": http,
  "ref": ref
})`),
	)
}
