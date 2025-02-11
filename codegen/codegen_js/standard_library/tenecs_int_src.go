// ##################################################################
// # The signatures of this file are generated via code-generation. #
// # Check gen.go                                                   #
// ##################################################################
package standard_library

func tenecs_int_minus() Function {
	return function(
		params("a", "b"),
		body(`return a - b`),
	)
}
func tenecs_int_plus() Function {
	return function(
		params("a", "b"),
		body(`return a + b`),
	)
}
func tenecs_int_times() Function {
	return function(
		params("a", "b"),
		body(`return a * b`),
	)
}
