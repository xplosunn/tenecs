package standard_library

var tenecs_string = packageWith(
	withFunction("characters", tenecs_string_characters),
	withFunction("contains", tenecs_string_contains),
	withFunction("endsWith", tenecs_string_endsWith),
	withFunction("firstChar", tenecs_string_firstChar),
	withFunction("firstCharCode", tenecs_string_firstCharCode),
	withFunction("isBlank", tenecs_string_isBlank),
	withFunction("isEmpty", tenecs_string_isEmpty),
	withFunction("join", tenecs_string_join),
	withFunction("length", tenecs_string_length),
	withFunction("padLeft", tenecs_string_padLeft),
	withFunction("padRight", tenecs_string_padRight),
	withFunction("repeat", tenecs_string_repeat),
	withFunction("reverse", tenecs_string_reverse),
	withFunction("startsWith", tenecs_string_startsWith),
	withFunction("stripPrefix", tenecs_string_stripPrefix),
	withFunction("stripSuffix", tenecs_string_stripSuffix),
	withFunction("toLowerCase", tenecs_string_toLowerCase),
	withFunction("toUpperCase", tenecs_string_toUpperCase),
	withFunction("trim", tenecs_string_trim),
	withFunction("trimLeft", tenecs_string_trimLeft),
	withFunction("trimRight", tenecs_string_trimRight),
)

var tenecs_string_characters = functionFromType("(str: String) ~> List<String>")

var tenecs_string_contains = functionFromType("(str: String, subStr: String) ~> Boolean")

var tenecs_string_join = functionFromType("(left: String, right: String) ~> String")

var tenecs_string_endsWith = functionFromType("(str: String, suffix: String) ~> Boolean")

var tenecs_string_startsWith = functionFromType("(str: String, prefix: String) ~> Boolean")

var tenecs_string_stripPrefix = functionFromType("(str: String, prefix: String) ~> String")

var tenecs_string_stripSuffix = functionFromType("(str: String, suffix: String) ~> String")

var tenecs_string_length = functionFromType("(str: String) ~> Int")

var tenecs_string_toLowerCase = functionFromType("(str: String) ~> String")

var tenecs_string_toUpperCase = functionFromType("(str: String) ~> String")

var tenecs_string_trim = functionFromType("(str: String) ~> String")

var tenecs_string_trimLeft = functionFromType("(str: String) ~> String")

var tenecs_string_trimRight = functionFromType("(str: String) ~> String")

var tenecs_string_isEmpty = functionFromType("(str: String) ~> Boolean")

var tenecs_string_isBlank = functionFromType("(str: String) ~> Boolean")

var tenecs_string_repeat = functionFromType("(str: String, count: Int) ~> String")

var tenecs_string_reverse = functionFromType("(str: String) ~> String")

var tenecs_string_padLeft = functionFromType("(str: String, length: Int, padChar: String) ~> String")

var tenecs_string_padRight = functionFromType("(str: String, length: Int, padChar: String) ~> String")

var tenecs_string_firstCharCode = functionFromType("(str: String) ~> Int")

var tenecs_string_firstChar = functionFromType("(str: String) ~> String")
