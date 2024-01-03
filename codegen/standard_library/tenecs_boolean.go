package standard_library

func tenecs_boolean_not() Function {
	return function(
		params("b"),
		body("return !b.(bool)"),
	)
}
