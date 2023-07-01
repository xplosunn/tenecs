package standard_library

func tenecs_array_append() Function {
	return function(
		params("array", "newElement"),
		body("return append(array.([]any{}), newElement)"),
	)
}
func tenecs_array_map() Function {
	return function(
		params("array", "f"),
		body(`result := []any{}
for _, elem := range array.([]any) {
result = append(result, f.(func(any)any)(elem))
}
return result
`),
	)
}
func tenecs_array_repeat() Function {
	return function(
		params("elem", "times"),
		body(`result := []any{}
for i := 0; i < times.(int); i++ {
result = append(result, elem)
}
return result
`),
	)
}
func tenecs_array_length() Function {
	return function(
		params("array"),
		body(`return len(array.([]any))`),
	)
}
