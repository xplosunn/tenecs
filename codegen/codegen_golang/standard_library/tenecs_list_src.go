package standard_library

import (
	godsl2 "github.com/xplosunn/tenecs/codegen/codegen_golang/godsl"
	"github.com/xplosunn/tenecs/typer/standard_library"
)

func tenecs_list_append() Function {
	return function(
		params("list", "newElement"),
		bodyDsl(
			godsl2.NativeFunctionInvocation().
				DeclaringVariables("result").
				Name("append").
				Parameters(
					godsl2.Cast(godsl2.VariableReference("list"), godsl2.TypeAnyList()),
					godsl2.VariableReference("newElement"),
				),
			godsl2.Return(godsl2.VariableReference("result")),
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
obj, okObj := maybeBreak.(tenecs_list_Break)
if okObj {
return obj._value
}
}
result = append(result, maybeBreak)
}
return result
`),
	)
}
func tenecs_list_Break() Function {
	return structFunction(standard_library.Tenecs_list_Break)
}
func tenecs_list_find() Function {
	return function(
		params("list", "f"),
		body(`for _, elem := range list.([]any) {
maybeReturn := f.(func(any)any)(elem)
if maybeReturn != nil {
return maybeReturn
}
}
return nil
`),
	)
}
