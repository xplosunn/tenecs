// ##################################################################
// # The signatures of this file are generated via code-generation. #
// # Check gen.go                                                   #
// ##################################################################
package standard_library

func tenecs_string_join() Function {
	return function(
		params("left", "right"),
		body("return left + right"),
	)
}

func tenecs_string_startsWith() Function {
	return function(
		params("str", "prefix"),
		body(`return str.startsWith(prefix)`),
	)
}
func tenecs_string_endsWith() Function {
	return function(
		params("str", "suffix"),
		body(`return str.endsWith(suffix)`),
	)
}
func tenecs_string_contains() Function {
	return function(
		params("str", "subStr"),
		body(`return str.includes(subStr)`),
	)
}
func tenecs_string_stripSuffix() Function {
	return function(
		params("str", "subStr"),
		body(`return subStr && str.endsWith(subStr) ? str.slice(0, 0 - subStr.length) : str;`),
	)
}
func tenecs_string_stripPrefix() Function {
	return function(
		params("str", "subStr"),
		body(`return str.startsWith(subStr) ? str.slice(subStr.length) : str;`),
	)
}
func tenecs_string_characters() Function {
	return function(
		params("str"),
		body(`return [...str]`),
	)
}
func tenecs_string_padLeft() Function {
	return function(
		params("str", "length", "padChar"),
		body(`return str.padStart(length, padChar || " ")`),
	)
}
func tenecs_string_toLowerCase() Function {
	return function(
		params("str"),
		body(`return str.toLowerCase()`),
	)
}
func tenecs_string_trimLeft() Function {
	return function(
		params("str"),
		body(`return str.trimStart()`),
	)
}
func tenecs_string_reverse() Function {
	return function(
		params("str"),
		body(`return [...str].reverse().join("")`),
	)
}
func tenecs_string_length() Function {
	return function(
		params("str"),
		body(`return str.length`),
	)
}
func tenecs_string_repeat() Function {
	return function(
		params("str", "count"),
		body(`return str.repeat(count)`),
	)
}
func tenecs_string_isBlank() Function {
	return function(
		params("str"),
		body(`return str.trim().length === 0`),
	)
}
func tenecs_string_isEmpty() Function {
	return function(
		params("str"),
		body(`return str.length === 0`),
	)
}
func tenecs_string_padRight() Function {
	return function(
		params("str", "length", "padChar"),
		body(`return str.padEnd(length, padChar || " ")`),
	)
}
func tenecs_string_toUpperCase() Function {
	return function(
		params("str"),
		body(`return str.toUpperCase()`),
	)
}
func tenecs_string_trim() Function {
	return function(
		params("str"),
		body(`return str.trim()`),
	)
}
func tenecs_string_trimRight() Function {
	return function(
		params("str"),
		body(`return str.trimEnd()`),
	)
}
func tenecs_string_firstCharCode() Function {
	return function(
		params("str"),
		body(`return str.length ? str.charCodeAt(0) : -1`),
	)
}
func tenecs_string_firstChar() Function {
	return function(
		params("str"),
		body(`return str.substring(0, 1)`),
	)
}
