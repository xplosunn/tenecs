package standard_library

var tenecs_string = packageWith(
	withFunction("contains", tenecs_string_contains),
	withFunction("join", tenecs_string_join),
	withFunction("endsWith", tenecs_string_endsWith),
	withFunction("startsWith", tenecs_string_startsWith),
	withFunction("stripPrefix", tenecs_string_stripPrefix),
	withFunction("stripSuffix", tenecs_string_stripSuffix),
)

var tenecs_string_contains = functionFromType("(str: String, subStr: String) ~> Boolean")

var tenecs_string_join = functionFromType("(left: String, right: String) ~> String")

var tenecs_string_endsWith = functionFromType("(str: String, suffix: String) ~> Boolean")

var tenecs_string_startsWith = functionFromType("(str: String, prefix: String) ~> Boolean")

var tenecs_string_stripPrefix = functionFromType("(str: String, prefix: String) ~> String")

var tenecs_string_stripSuffix = functionFromType("(str: String, suffix: String) ~> String")
