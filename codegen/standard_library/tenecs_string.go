package standard_library

func tenecs_string_join() Function {
	return function(
		params("Pleft", "Pright"),
		body("return Pleft.(string) + Pright.(string)"),
	)
}

func tenecs_string_hasPrefix() Function {
	return function(
		imports("strings"),
		params("Pstr", "Pprefix"),
		body("return strings.HasPrefix(Pstr.(string), Pprefix.(string))"),
	)
}
