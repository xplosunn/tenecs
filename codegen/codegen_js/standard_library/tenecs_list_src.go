// ##################################################################
// # The signatures of this file are generated via code-generation. #
// # Check gen.go                                                   #
// ##################################################################
package standard_library

import "github.com/xplosunn/tenecs/typer/standard_library"

func tenecs_list_append() Function {
	return function(
		params("list", "newElement"),
		body(`return list.concat([newElement])`),
	)
}
func tenecs_list_map() Function {
	return function(
		params("list", "f"),
		body(`return list.map(f)`),
	)
}
func tenecs_list_mapNotNull() Function {
	return function(
		params("list", "f"),
		body(`return list.flatMap((e) => {
  let mapped = f(e)
  if (mapped == null) {
    return []
  } else {
    return [mapped]
  }
})`),
	)
}
func tenecs_list_repeat() Function {
	return function(
		params("elem", "times"),
		body(`return Array(times).fill(elem)`),
	)
}
func tenecs_list_length() Function {
	return function(
		params("list"),
		body(`return list.size`),
	)
}
func tenecs_list_filter() Function {
	return function(
		params("list", "keep"),
		body(`return list.filter(keep)`),
	)
}
func tenecs_list_flatMap() Function {
	return function(
		params("list", "f"),
		body(`return list.flatMap(f)`),
	)
}
func tenecs_list_fold() Function {
	return function(
		params("list", "zero", "f"),
		body(`
let acc = zero;
for (const elem of list) {
  acc = f(acc, elem);
}
return acc;
`),
	)
}
func tenecs_list_forEach() Function {
	return function(
		params("list", "f"),
		body(`list.forEach(f)`),
	)
}
func tenecs_list_mapUntil() Function {
	return function(
		params("list", "f"),
		body(`
let result = [];
for (const elem of list) {
  let e = f(elem);
  if (e && e["$type"] && e["$type"] === "Break") {
    return e.value;
  }
  result.push(e);
}
return result;
`),
	)
}
func tenecs_list_Break() Function {
	return structFunction(standard_library.Tenecs_list_Break)
}
func tenecs_list_find() Function {
	return function(
		params("list", "f"),
		body(`
for (const elem of list) {
  let e = f(elem);
  if (e != null) {
    return e;
  }
}
return null;
`),
	)
}
func tenecs_list_first() Function {
	return function(
		params("list"),
		body(`
if (list.length > 0) {
  return list[0];
}
return null;
`),
	)
}
func tenecs_list_atIndexGet() Function {
	return function(
		params("list", "index"),
		body(`
if (index >= 0 && list.length > index) {
  return list[index];
}
return ({
  "$type": "Error",
  "message": "Out of bounds"
});
`),
	)
}
func tenecs_list_atIndexSet() Function {
	return function(
		params("list", "index", "setTo"),
		body(`
if (index >= 0 && list.length > index) {
  const result = list.slice();
  result[index] = setTo;
  return result;
}
return ({
  "$type": "Error",
  "message": "Out of bounds"
});
`),
	)
}
func tenecs_list_appendAll() Function {
	return function(
		params("list", "newElements"),
		body(`return list.concat(newElements);`),
	)
}
func tenecs_list_flatten() Function {
	return function(
		params("list"),
		body(`return list.flat();`),
	)
}
