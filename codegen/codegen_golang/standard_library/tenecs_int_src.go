package standard_library

func tenecs_int_minus() Function {
	return function(
		params("a", "b"),
		body(`return a.(int) - b.(int)`),
	)
}
func tenecs_int_plus() Function {
	return function(
		params("a", "b"),
		body(`return a.(int) + b.(int)`),
	)
}
func tenecs_int_times() Function {
	return function(
		params("a", "b"),
		body(`return a.(int) * b.(int)`),
	)
}
func tenecs_int_div() Function {
	return function(
		params("a", "b"),
		body(`if (b == 0) {
return tenecs_error_Error{
_message: "Division by zero",
}
} else {
return a.(int) / b.(int)
}`),
	)
}
func tenecs_int_ponyDiv() Function {
	return function(
		params("a", "b"),
		body(`if (b == 0) {
return 0
} else {
return a.(int) / b.(int)
}`),
	)
}
func tenecs_int_mod() Function {
	return function(
		params("a", "b"),
		body(`if (b == 0) {
return tenecs_error_Error{
_message: "Division by zero",
}
} else {
return a.(int) % b.(int)
}`),
	)
}
func tenecs_int_ponyMod() Function {
	return function(
		params("a", "b"),
		body(`if (b == 0) {
return 0
} else {
return a.(int) % b.(int)
}`),
	)
}
func tenecs_int_lessThan() Function {
	return function(
		params("a", "b"),
		body(`return a.(int) < b.(int)`),
	)
}
func tenecs_int_greaterThan() Function {
	return function(
		params("a", "b"),
		body(`return a.(int) > b.(int)`),
	)
}
func tenecs_int_abs() Function {
	return function(
		params("a"),
		body(`if a.(int) < 0 {
return -a.(int)
}
return a.(int)`),
	)
}

func tenecs_int_negate() Function {
	return function(
		params("a"),
		body(`return -a.(int)`),
	)
}
