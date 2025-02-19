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
func tenecs_int_ponyDiv() Function {
	return function(
		params("a", "b"),
		body(`if (b == 0) {
  return 0
} else {
  return Math.trunc(a / b)
}
`),
	)
}
func tenecs_int_div() Function {
	return function(
		params("a", "b"),
		body(`if (b == 0) {
  return ({
    "$type": "Error",
    "message": "Division by zero"
  })
} else {
  return Math.trunc(a / b)
}
`),
	)
}
func tenecs_int_mod() Function {
	return function(
		params("a", "b"),
		body(`if (b == 0) {
  return ({
    "$type": "Error",
    "message": "Division by zero"
  })
} else {
  return a % b
}
`),
	)
}
func tenecs_int_ponyMod() Function {
	return function(
		params("a", "b"),
		body(`if (b == 0) {
  return 0
} else {
  return a % b
}
`),
	)
}
func tenecs_int_greaterThan() Function {
	return function(
		params("a", "b"),
		body(`return a > b`),
	)
}
func tenecs_int_lessThan() Function {
	return function(
		params("a", "b"),
		body(`return a < b`),
	)
}
