// ##################################################################
// # The signatures of this file are generated via code-generation. #
// # Check gen.go                                                   #
// ##################################################################
package standard_library

func tenecs_boolean_not() Function {
	return function(
		params("b"),
		body("return !b"),
	)
}
func tenecs_boolean_and() Function {
	return function(
		params("a", "b"),
		body("return a && b()"),
	)
}
func tenecs_boolean_or() Function {
	return function(
		params("a", "b"),
		body("return a || b()"),
	)
}
