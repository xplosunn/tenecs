package standard_library

func tenecs_boolean_not() Function {
	return function(
		params("b"),
		body("return !b.(bool)"),
	)
}
func tenecs_boolean_and() Function {
	return function(
		params("a", "b"),
		body("return a.(bool) && b.(func()any)().(bool)"),
	)
}
func tenecs_boolean_or() Function {
	return function(
		params("a", "b"),
		body("return a.(bool) || b.(func()any)().(bool)"),
	)
}
