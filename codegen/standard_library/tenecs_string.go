package standard_library

func tenecs_string_join() Function {
	return function(
		params("Pleft", "Pright"),
		body("return Pleft.(string) + Pright.(string)"),
	)
}
