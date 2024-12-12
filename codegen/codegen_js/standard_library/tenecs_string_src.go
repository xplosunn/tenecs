package standard_library

func tenecs_string_join() Function {
	return function(
		params("left", "right"),
		body("return left + right"),
	)
}

func tenecs_string_startsWith() Function {
	return function(
		params("str", "prefix"),
		body(`return str.startsWith(prefix)`),
	)
}
func tenecs_string_endsWith() Function {
	return function(
		params("str", "suffix"),
		body(`return str.endsWith(suffix)`),
	)
}
func tenecs_string_contains() Function {
	return function(
		params("str", "subStr"),
		body(`return str.includes(subStr)`),
	)
}
