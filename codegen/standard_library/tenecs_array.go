package standard_library

func tenecs_array_emptyArray() Function {
	return function(
		params(),
		body("return []any{}"),
	)
}
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
