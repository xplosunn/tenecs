package standard_library

func tenecs_string_join() Function {
	return function(
		params("left", "right"),
		body("return left.(string) + right.(string)"),
	)
}
