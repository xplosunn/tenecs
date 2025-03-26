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
func tenecs_string_stripPrefix() Function {
	return function(
		imports("strings"),
		params("str", "subStr"),
		body("return strings.TrimPrefix(str.(string), subStr.(string))"),
	)
}
func tenecs_string_stripSuffix() Function {
	return function(
		imports("strings"),
		params("str", "subStr"),
		body("return strings.TrimSuffix(str.(string), subStr.(string))"),
	)
}
func tenecs_string_characters() Function {
	return function(
		params("str"),
		body(`
input := str.(string)
runes := []rune(input)
result := []any{}
for _, r := range runes {
result = append(result, string(r))
}
return result`),
	)
}
func tenecs_string_length() Function {
	return function(
		params("str"),
		body("return len([]rune(str.(string)))"),
	)
}
func tenecs_string_toUpperCase() Function {
	return function(
		imports("strings"),
		params("str"),
		body("return strings.ToUpper(str.(string))"),
	)
}
func tenecs_string_trimRight() Function {
	return function(
		imports("strings"),
		params("str"),
		body("return strings.TrimRight(str.(string), \" \\t\\n\\r\")"),
	)
}
func tenecs_string_isEmpty() Function {
	return function(
		params("str"),
		body("return len(str.(string)) == 0"),
	)
}
func tenecs_string_repeat() Function {
	return function(
		imports("strings"),
		params("str", "count"),
		body("return strings.Repeat(str.(string), count.(int))"),
	)
}
func tenecs_string_padLeft() Function {
	return function(
		imports("strings"),
		params("str", "length", "padChar"),
		body(`
padCount := length.(int) - len([]rune(str.(string)))
if padCount <= 0 {
return str.(string)
}
return strings.Repeat(padChar.(string), padCount) + str.(string)`),
	)
}
func tenecs_string_toLowerCase() Function {
	return function(
		imports("strings"),
		params("str"),
		body("return strings.ToLower(str.(string))"),
	)
}
func tenecs_string_trimLeft() Function {
	return function(
		imports("strings"),
		params("str"),
		body("return strings.TrimLeft(str.(string), \" \\t\\n\\r\")"),
	)
}
func tenecs_string_trim() Function {
	return function(
		imports("strings"),
		params("str"),
		body("return strings.TrimSpace(str.(string))"),
	)
}
func tenecs_string_isBlank() Function {
	return function(
		imports("strings"),
		params("str"),
		body("return len(strings.TrimSpace(str.(string))) == 0"),
	)
}
func tenecs_string_padRight() Function {
	return function(
		imports("strings"),
		params("str", "length", "padChar"),
		body(`
padCount := int(length.(int)) - len([]rune(str.(string)))
if padCount <= 0 {
return str.(string)
}
return str.(string) + strings.Repeat(padChar.(string), padCount)`),
	)
}
func tenecs_string_reverse() Function {
	return function(
		params("str"),
		body(`
runes := []rune(str.(string))
for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
runes[i], runes[j] = runes[j], runes[i]
}
return string(runes)`),
	)
}
