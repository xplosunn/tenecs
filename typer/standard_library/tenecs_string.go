package standard_library

var tenecs_string = packageWith(
	withFunction("join", tenecs_string_join),
	withFunction("endsWith", tenecs_string_endsWith),
	withFunction("startsWith", tenecs_string_startsWith),
)

var tenecs_string_join = functionFromSignature("(left: String, right: String): String")

var tenecs_string_endsWith = functionFromSignature("(str: String, suffix: String): Boolean")

var tenecs_string_startsWith = functionFromSignature("(str: String, prefix: String): Boolean")
