package codegen_javascript

import (
	_ "embed"
	"fmt"

	"github.com/xplosunn/tenecs/ir"
)

func generateNativeFunction(nativeFunctionRef ir.NativeFunctionRef) string {
	stdLibFunc, ok := stdLibStaticFunctions[nativeFunctionRef]
	if ok {
		return stdLibFunc
	}

	stdLibDynamicFunc := stdLibDynamicFunction(nativeFunctionRef)
	if stdLibDynamicFunc != nil {
		return *stdLibDynamicFunc
	}
	panic("TODO generateNativeFunction" + fmt.Sprintf("%+v", nativeFunctionRef))
}

//go:embed embed/tenecs_boolean/and.js
var tenecs_boolean__and string

//go:embed embed/tenecs_boolean/not.js
var tenecs_boolean__not string

//go:embed embed/tenecs_boolean/or.js
var tenecs_boolean__or string

//go:embed embed/tenecs_compare/eq.js
var tenecs_compare__eq string

//go:embed embed/tenecs_error/Error.js
var tenecs_error__Error string

//go:embed embed/tenecs_int/abs.js
var tenecs_int__abs string

//go:embed embed/tenecs_int/div.js
var tenecs_int__div string

//go:embed embed/tenecs_int/greaterThan.js
var tenecs_int__greaterThan string

//go:embed embed/tenecs_int/lessThan.js
var tenecs_int__lessThan string

//go:embed embed/tenecs_int/minus.js
var tenecs_int__minus string

//go:embed embed/tenecs_int/mod.js
var tenecs_int__mod string

//go:embed embed/tenecs_int/negate.js
var tenecs_int__negate string

//go:embed embed/tenecs_int/plus.js
var tenecs_int__plus string

//go:embed embed/tenecs_int/ponyDiv.js
var tenecs_int__ponyDiv string

//go:embed embed/tenecs_int/ponyMod.js
var tenecs_int__ponyMod string

//go:embed embed/tenecs_int/times.js
var tenecs_int__times string

//go:embed embed/tenecs_time/atStartOfMonth.js
var tenecs_time__atStartOfMonth string

//go:embed embed/tenecs_time/plusDays.js
var tenecs_time__plusDays string

//go:embed embed/tenecs_time/plusYears.js
var tenecs_time__plusYears string

//go:embed embed/tenecs_json/jsonBoolean.js
var tenecs_json__jsonBoolean string

//go:embed embed/tenecs_json/jsonInt.js
var tenecs_json__jsonInt string

//go:embed embed/tenecs_json/jsonString.js
var tenecs_json__jsonString string

//go:embed embed/tenecs_json/jsonList.js
var tenecs_json__jsonList string

//go:embed embed/tenecs_json/jsonOr.js
var tenecs_json__jsonOr string

//go:embed embed/tenecs_json/jsonObject0.js
var tenecs_json__jsonObject0 string

//go:embed embed/tenecs_json/jsonObject1.js
var tenecs_json__jsonObject1 string

//go:embed embed/tenecs_list/append.js
var tenecs_list__append string

//go:embed embed/tenecs_list/first.js
var tenecs_list__first string

//go:embed embed/tenecs_list/length.js
var tenecs_list__length string

//go:embed embed/tenecs_list/map.js
var tenecs_list__map string

//go:embed embed/tenecs_list/appendAll.js
var tenecs_list__appendAll string

//go:embed embed/tenecs_list/atIndexGet.js
var tenecs_list__atIndexGet string

//go:embed embed/tenecs_list/atIndexSet.js
var tenecs_list__atIndexSet string

//go:embed embed/tenecs_list/Break.js
var tenecs_list__Break string

//go:embed embed/tenecs_list/filter.js
var tenecs_list__filter string

//go:embed embed/tenecs_list/find.js
var tenecs_list__find string

//go:embed embed/tenecs_list/flatMap.js
var tenecs_list__flatMap string

//go:embed embed/tenecs_list/flatten.js
var tenecs_list__flatten string

//go:embed embed/tenecs_list/fold.js
var tenecs_list__fold string

//go:embed embed/tenecs_list/forEach.js
var tenecs_list__forEach string

//go:embed embed/tenecs_list/mapNotNull.js
var tenecs_list__mapNotNull string

//go:embed embed/tenecs_list/mapUntil.js
var tenecs_list__mapUntil string

//go:embed embed/tenecs_list/repeat.js
var tenecs_list__repeat string

//go:embed embed/tenecs_meta/typeName.js
var tenecs_meta__typeName string

//go:embed embed/tenecs_string/characters.js
var tenecs_string__characters string

//go:embed embed/tenecs_string/contains.js
var tenecs_string__contains string

//go:embed embed/tenecs_string/firstChar.js
var tenecs_string__firstChar string

//go:embed embed/tenecs_string/firstCharCode.js
var tenecs_string__firstCharCode string

//go:embed embed/tenecs_string/endsWith.js
var tenecs_string__endsWith string

//go:embed embed/tenecs_string/isBlank.js
var tenecs_string__isBlank string

//go:embed embed/tenecs_string/isEmpty.js
var tenecs_string__isEmpty string

//go:embed embed/tenecs_string/join.js
var tenecs_string__join string

//go:embed embed/tenecs_string/length.js
var tenecs_string__length string

//go:embed embed/tenecs_string/padLeft.js
var tenecs_string__padLeft string

//go:embed embed/tenecs_string/padRight.js
var tenecs_string__padRight string

//go:embed embed/tenecs_string/repeat.js
var tenecs_string__repeat string

//go:embed embed/tenecs_string/reverse.js
var tenecs_string__reverse string

//go:embed embed/tenecs_string/startsWith.js
var tenecs_string__startsWith string

//go:embed embed/tenecs_string/stripPrefix.js
var tenecs_string__stripPrefix string

//go:embed embed/tenecs_string/stripSuffix.js
var tenecs_string__stripSuffix string

//go:embed embed/tenecs_string/toLowerCase.js
var tenecs_string__toLowerCase string

//go:embed embed/tenecs_string/toUpperCase.js
var tenecs_string__toUpperCase string

//go:embed embed/tenecs_string/trim.js
var tenecs_string__trim string

//go:embed embed/tenecs_string/trimLeft.js
var tenecs_string__trimLeft string

//go:embed embed/tenecs_string/trimRight.js
var tenecs_string__trimRight string

var stdLibStaticFunctions map[ir.NativeFunctionRef]string = func() map[ir.NativeFunctionRef]string {
	stdLibFunctions := map[ir.NativeFunctionRef]string{}
	addToStdLibFunctions := func(pkgName string, funcName string, funcCode string) {
		stdLibFunctions[ir.NativeFunctionRef{
			Package: pkgName,
			Name:    funcName,
		}] = funcCode
	}
	addToStdLibFunctions("tenecs_boolean", "and", tenecs_boolean__and)
	addToStdLibFunctions("tenecs_boolean", "not", tenecs_boolean__not)
	addToStdLibFunctions("tenecs_boolean", "or", tenecs_boolean__or)
	addToStdLibFunctions("tenecs_compare", "eq", tenecs_compare__eq)
	addToStdLibFunctions("tenecs_error", "Error", tenecs_error__Error)
	addToStdLibFunctions("tenecs_int", "abs", tenecs_int__abs)
	addToStdLibFunctions("tenecs_int", "div", tenecs_int__div)
	addToStdLibFunctions("tenecs_int", "greaterThan", tenecs_int__greaterThan)
	addToStdLibFunctions("tenecs_int", "lessThan", tenecs_int__lessThan)
	addToStdLibFunctions("tenecs_int", "minus", tenecs_int__minus)
	addToStdLibFunctions("tenecs_int", "mod", tenecs_int__mod)
	addToStdLibFunctions("tenecs_int", "negate", tenecs_int__negate)
	addToStdLibFunctions("tenecs_int", "plus", tenecs_int__plus)
	addToStdLibFunctions("tenecs_int", "ponyDiv", tenecs_int__ponyDiv)
	addToStdLibFunctions("tenecs_int", "ponyMod", tenecs_int__ponyMod)
	addToStdLibFunctions("tenecs_int", "times", tenecs_int__times)
	addToStdLibFunctions("tenecs_time", "atStartOfMonth", tenecs_time__atStartOfMonth)
	addToStdLibFunctions("tenecs_time", "plusDays", tenecs_time__plusDays)
	addToStdLibFunctions("tenecs_time", "plusYears", tenecs_time__plusYears)
	addToStdLibFunctions("tenecs_json", "jsonBoolean", tenecs_json__jsonBoolean)
	addToStdLibFunctions("tenecs_json", "jsonInt", tenecs_json__jsonInt)
	addToStdLibFunctions("tenecs_json", "jsonString", tenecs_json__jsonString)
	addToStdLibFunctions("tenecs_json", "jsonList", tenecs_json__jsonList)
	addToStdLibFunctions("tenecs_json", "jsonOr", tenecs_json__jsonOr)
	addToStdLibFunctions("tenecs_json", "jsonObject0", tenecs_json__jsonObject0)
	addToStdLibFunctions("tenecs_json", "jsonObject1", tenecs_json__jsonObject1)
	addToStdLibFunctions("tenecs_list", "append", tenecs_list__append)
	addToStdLibFunctions("tenecs_list", "first", tenecs_list__first)
	addToStdLibFunctions("tenecs_list", "length", tenecs_list__length)
	addToStdLibFunctions("tenecs_list", "map", tenecs_list__map)
	addToStdLibFunctions("tenecs_list", "appendAll", tenecs_list__appendAll)
	addToStdLibFunctions("tenecs_list", "atIndexGet", tenecs_list__atIndexGet)
	addToStdLibFunctions("tenecs_list", "atIndexSet", tenecs_list__atIndexSet)
	addToStdLibFunctions("tenecs_list", "Break", tenecs_list__Break)
	addToStdLibFunctions("tenecs_list", "filter", tenecs_list__filter)
	addToStdLibFunctions("tenecs_list", "find", tenecs_list__find)
	addToStdLibFunctions("tenecs_list", "flatMap", tenecs_list__flatMap)
	addToStdLibFunctions("tenecs_list", "flatten", tenecs_list__flatten)
	addToStdLibFunctions("tenecs_list", "fold", tenecs_list__fold)
	addToStdLibFunctions("tenecs_list", "forEach", tenecs_list__forEach)
	addToStdLibFunctions("tenecs_list", "mapNotNull", tenecs_list__mapNotNull)
	addToStdLibFunctions("tenecs_list", "mapUntil", tenecs_list__mapUntil)
	addToStdLibFunctions("tenecs_list", "repeat", tenecs_list__repeat)
	addToStdLibFunctions("tenecs_meta", "typeName", tenecs_meta__typeName)
	addToStdLibFunctions("tenecs_string", "characters", tenecs_string__characters)
	addToStdLibFunctions("tenecs_string", "contains", tenecs_string__contains)
	addToStdLibFunctions("tenecs_string", "firstChar", tenecs_string__firstChar)
	addToStdLibFunctions("tenecs_string", "firstCharCode", tenecs_string__firstCharCode)
	addToStdLibFunctions("tenecs_string", "endsWith", tenecs_string__endsWith)
	addToStdLibFunctions("tenecs_string", "isBlank", tenecs_string__isBlank)
	addToStdLibFunctions("tenecs_string", "isEmpty", tenecs_string__isEmpty)
	addToStdLibFunctions("tenecs_string", "join", tenecs_string__join)
	addToStdLibFunctions("tenecs_string", "length", tenecs_string__length)
	addToStdLibFunctions("tenecs_string", "padLeft", tenecs_string__padLeft)
	addToStdLibFunctions("tenecs_string", "padRight", tenecs_string__padRight)
	addToStdLibFunctions("tenecs_string", "repeat", tenecs_string__repeat)
	addToStdLibFunctions("tenecs_string", "reverse", tenecs_string__reverse)
	addToStdLibFunctions("tenecs_string", "startsWith", tenecs_string__startsWith)
	addToStdLibFunctions("tenecs_string", "stripPrefix", tenecs_string__stripPrefix)
	addToStdLibFunctions("tenecs_string", "stripSuffix", tenecs_string__stripSuffix)
	addToStdLibFunctions("tenecs_string", "toLowerCase", tenecs_string__toLowerCase)
	addToStdLibFunctions("tenecs_string", "toUpperCase", tenecs_string__toUpperCase)
	addToStdLibFunctions("tenecs_string", "trim", tenecs_string__trim)
	addToStdLibFunctions("tenecs_string", "trimLeft", tenecs_string__trimLeft)
	addToStdLibFunctions("tenecs_string", "trimRight", tenecs_string__trimRight)
	return stdLibFunctions
}()
