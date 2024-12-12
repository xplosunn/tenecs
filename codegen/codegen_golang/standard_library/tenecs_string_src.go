package standard_library

func tenecs_string_join() Function {
	return function(
		params("Pleft", "Pright"),
		body("return Pleft.(string) + Pright.(string)"),
	)
}

func tenecs_string_startsWith() Function {
	return function(
		imports("strings"),
		params("Pstr", "Pprefix"),
		body("return strings.HasPrefix(Pstr.(string), Pprefix.(string))"),
	)
}
func tenecs_string_endsWith() Function {
	return function(
		imports("strings"),
		params("Pstr", "Psuffix"),
		body("return strings.HasSuffix(Pstr.(string), Psuffix.(string))"),
	)
}
func tenecs_string_contains() Function {
	return function(
		imports("strings"),
		params("str", "subStr"),
		body("return strings.Contains(str.(string), subStr.(string))"),
	)
}
