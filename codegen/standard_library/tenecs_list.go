package standard_library

import "github.com/xplosunn/tenecs/godsl"

func tenecs_list_append() Function {
	return function(
		params("list", "newElement"),
		bodyDsl(
			godsl.NativeFunctionInvocation().
				DeclaringVariables("result").
				Name("append").
				Parameters(
					godsl.Cast(godsl.VariableReference("list"), godsl.TypeAnyList()),
					godsl.VariableReference("newElement"),
				),
			godsl.Return(godsl.VariableReference("result")),
		),
	)
}
func tenecs_list_map() Function {
	return function(
		params("list", "f"),
		body(`result := []any{}
for _, elem := range list.([]any) {
result = append(result, f.(func(any)any)(elem))
}
return result
`),
	)
}
func tenecs_list_mapNotNull() Function {
	return function(
		params("list", "f"),
		body(`result := []any{}
for _, elem := range list.([]any) {
maybeNull := f.(func(any)any)(elem)
if maybeNull != nil {
result = append(result, maybeNull)
}
}
return result
`),
	)
}
func tenecs_list_repeat() Function {
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
func tenecs_list_length() Function {
	return function(
		params("list"),
		body(`return len(list.([]any))`),
	)
}
func tenecs_list_filter() Function {
	return function(
		params("list", "keep"),
		body(`result := []any{}
for _, elem := range list.([]any) {
if keep.(func(any)any)(elem).(bool) {
result = append(result, elem)
}
}
return result`),
	)
}
func tenecs_list_flatMap() Function {
	return function(
		params("list", "f"),
		body(`result := []any{}
for _, elem := range list.([]any) {
result = append(result, f.(func(any)any)(elem).([]any)...)
}
return result
`),
	)
}
func tenecs_list_fold() Function {
	return function(
		params("list", "zero", "f"),
		body(`result := zero
for _, elem := range list.([]any) {
result = f.(func(any,any)any)(result, elem)
}
return result
`),
	)
}
func tenecs_list_forEach() Function {
	return function(
		params("list", "f"),
		body(`for _, elem := range list.([]any) {
f.(func(any)any)(elem)
}
`),
	)
}
func tenecs_list_mapUntil() Function {
	return function(
		params("list", "f"),
		body(`result := []any{}
for _, elem := range list.([]any) {
maybeBreak := f.(func(any)any)(elem)
if maybeBreak != nil {
obj, okObj := maybeBreak.(map[string]any)
if okObj && obj["$type"] == "Break" {
return obj["value"]
}
}
result = append(result, maybeBreak)
}
return result
`),
	)
}
func tenecs_list_Break() Function {
	return function(
		params("value"),
		body(`return map[string]any{
	"$type": "Break",
	"value": value,
}`),
	)
}
