package standard_library

import (
	"github.com/xplosunn/tenecs/typer/standard_library"
)

func tenecs_list_append() Function {
	return function(
		params("list", "newElement"),
		body(`return append(list.([]any), newElement)`),
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
func tenecs_list_first() Function {
	return function(
		params("list"),
		body(`l := list.([]any)
if len(l) > 0 {
return l[0]
}
return nil
`),
	)
}
func tenecs_list_atIndexGet() Function {
	return function(
		params("list", "index"),
		body(`l := list.([]any)
idx := index.(int)
if idx >= 0 && len(l) > idx {
return l[idx]
}
return tenecs_error_Error{
_message: "Out of bounds",
}
`),
	)
}
func tenecs_list_atIndexSet() Function {
	return function(
		params("list", "index", "setTo"),
		body(`l := list.([]any)
idx := index.(int)
if idx >= 0 && len(l) > idx {
result := make([]any, len(l))
copy(result, l)
result[idx] = setTo
return result
}
return tenecs_error_Error{
_message: "Out of bounds",
}
`),
	)
}
func tenecs_list_appendAll() Function {
	return function(
		params("list", "newElements"),
		body(`return append(list.([]any), newElements.([]any)...)`),
	)
}
func tenecs_list_flatten() Function {
	return function(
		params("list"),
		body(`result := []any{}
for _, innerList := range list.([]any) {
result = append(result, innerList.([]any)...)
}
return result
`),
	)
}
