package standard_library

func tenecs_compare_eq() Function {
	return function(
		imports("reflect"),
		params("first", "second"),
		body(`return reflect.DeepEqual(first, second)`),
	)
}
