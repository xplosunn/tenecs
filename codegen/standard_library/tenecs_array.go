package standard_library

import "github.com/xplosunn/tenecs/godsl"

func tenecs_array_append() Function {
	return function(
		params("array", "newElement"),
		bodyDsl(
			godsl.NativeFunctionInvocation().
				DeclaringVariables("result").
				Name("append").
				Parameters(
					godsl.Cast(godsl.VariableReference("array"), godsl.TypeAnyArray()),
					godsl.VariableReference("newElement"),
				),
			godsl.Return(godsl.VariableReference("result")),
		),
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
func tenecs_array_mapNotNull() Function {
	return function(
		params("array", "f"),
		body(`result := []any{}
for _, elem := range array.([]any) {
maybeNull := f.(func(any)any)(elem)
if maybeNull != nil {
result = append(result, maybeNull)
}
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
func tenecs_array_filter() Function {
	return function(
		params("array", "keep"),
		body(`result := []any{}
for _, elem := range array.([]any) {
if keep.(func(any)any)(elem).(bool) {
result = append(result, elem)
}
}
return result`),
	)
}
func tenecs_array_flatMap() Function {
	return function(
		params("array", "f"),
		body(`result := []any{}
for _, elem := range array.([]any) {
result = append(result, f.(func(any)any)(elem).([]any)...)
}
return result
`),
	)
}
func tenecs_array_fold() Function {
	return function(
		params("array", "zero", "f"),
		body(`result := zero
for _, elem := range array.([]any) {
result = f.(func(any,any)any)(result, elem)
}
return result
`),
	)
}
func tenecs_array_forEach() Function {
	return function(
		params("array", "f"),
		body(`for _, elem := range array.([]any) {
f.(func(any)any)(elem)
}
`),
	)
}
