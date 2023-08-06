package standard_library

import (
	"fmt"
	"strings"
)

func tenecs_json_toJson() Function {
	return function(
		imports("encoding/json"),
		params("input"),
		body(`if inputMap, ok := input.(map[string]any); ok {
copy := map[string]any{}
for k, v := range inputMap {
copy[k] = v
}
delete(copy, "$type")
result, _ := json.Marshal(copy)
return string(result)
}
result, _ := json.Marshal(input)
return string(result)`),
	)
}

func tenecs_json_jsonError() Function {
	return function(
		params("message"),
		body(`return map[string]any{
	"$type": "JsonError",
	"message": message,
}`),
	)
}

func tenecs_json_parseBoolean() Function {
	return function(
		imports("encoding/json"),
		body(`return map[string]any{
	"$type": "FromJson",
	"parse": func(input any) any {
		jsonString := input.(string)
		var output bool
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not parse Boolean from " + jsonString,
			} 
		}
		return output
	},
}`),
	)
}

func tenecs_json_parseInt() Function {
	return function(
		imports("encoding/json"),
		body(`return map[string]any{
	"$type": "FromJson",
	"parse": func(input any) any {
		jsonString := input.(string)
		var output float64
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil || float64(int(output)) != output {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not parse Int from " + jsonString,
			} 
		}
		return int(output)
	},
}`),
	)
}

func tenecs_json_parseOr() Function {
	return function(
		imports("encoding/json"),
		params("fromA", "fromB"),
		body(`return map[string]any{
	"$type": "FromJson",
	"parse": func(input any) any {
		resultA := fromA.(map[string]any)["parse"].(func(any)any)(input)
		resultAMap, isMap := resultA.(map[string]any)
		if isMap && resultAMap["$type"] == "JsonError" {
			resultB := fromB.(map[string]any)["parse"].(func(any)any)(input)
			resultBMap, isMap := resultB.(map[string]any)
			if isMap && resultBMap["$type"] == "JsonError" {
				jsonString := input.(string)
				return map[string]any{
					"$type": "JsonError",
					"message": "Could not parse from " + jsonString,
				}
			}
			return resultB
		}
		return resultA
	},
}`),
	)
}

func tenecs_json_parseString() Function {
	return function(
		imports("encoding/json"),
		body(`return map[string]any{
	"$type": "FromJson",
	"parse": func(input any) any {
		jsonString := input.(string)
		var output string
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not parse String from " + jsonString,
			} 
		}
		return output
	},
}`),
	)
}

func tenecs_json_parseArray() Function {
	return function(
		imports("encoding/json"),
		params("of"),
		body(`return map[string]any{
	"$type": "FromJson",
	"parse": func(input any) any {
		jsonString := input.(string)
		var output []json.RawMessage
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not parse Array from " + jsonString,
			} 
		}
		if len(output) == 0 {
			return []any{}
		}
		ofParse := of.(map[string]any)["parse"].(func(any)any)
		outputArray := []any{}
		for _, elem := range output {
			elemJsonBytes, _ := json.Marshal(&elem)
			result := ofParse(string(elemJsonBytes))
			resultMap, isMap := result.(map[string]any)
			if isMap && resultMap["$type"] == "JsonError" {
				return result
			}
			outputArray = append(outputArray, result)
		}
		return outputArray
	},
}`),
	)
}

func tenecs_json_field() Function {
	return function(
		params("name, fromJson"),
		body(`return map[string]any{
	"$type": "FromJsonField",
	"name": name,
	"fromJson": fromJson,
}`),
	)
}

func tenecs_json_parseObject0() Function {
	return function(
		imports("encoding/json"),
		params("build"),
		body(`return map[string]any{
	"$type": "FromJson",
	"parse": func(input any) any {
		jsonString := input.(string)
		var output map[string]json.RawMessage
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not parse object from " + jsonString,
			} 
		}

		return build.(func()any)()
	},
}`),
	)
}

func tenecs_json_parseObject_X(x int) Function {
	paramNames := []string{"build"}
	for i := 0; i < x; i++ {
		paramNames = append(paramNames, fmt.Sprintf("fromJsonFieldI%d", i))
	}
	bodyStr := `return map[string]any{
	"$type": "FromJson",
	"parse": func(input any) any {
		jsonString := input.(string)
		var output map[string]json.RawMessage
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not parse object from " + jsonString,
			} 
		}
`
	for i := 0; i < x; i++ {
		bodyStr += fmt.Sprintf(`
		i%dName := fromJsonFieldI%d.(map[string]any)["name"].(string)
		i%dJsonRawMessage := output[i%dName]
		if i%dJsonRawMessage == nil {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not find object field \"" + i%dName + "\" in " + jsonString,
			}
		}
		i%dJsonBytes, _ := json.Marshal(&i%dJsonRawMessage)
		i%d := fromJsonFieldI%d.(map[string]any)["fromJson"].(map[string]any)["parse"].(func(any)any)(string(i%dJsonBytes))
		i%dMap, isMap := i%d.(map[string]any)
		if isMap && i%dMap["$type"] == "JsonError" {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not parse object field \"" + i%dName + "\": " + i%dMap["message"].(string),
			} 
		}`, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i)
	}
	anys := []string{}
	buildArgs := []string{}
	for i := 0; i < x; i++ {
		anys = append(anys, "any")
		buildArgs = append(buildArgs, fmt.Sprintf("i%d", i))
	}
	bodyStr += fmt.Sprintf(`

		return build.(func(%s)any)(%s)
	},
}`, strings.Join(anys, ","), strings.Join(buildArgs, ","))
	return function(
		imports("encoding/json"),
		params(paramNames...),
		body(bodyStr),
	)
}

func tenecs_json_parseObject1() Function {
	return tenecs_json_parseObject_X(1)
}
func tenecs_json_parseObject7() Function {
	return tenecs_json_parseObject_X(7)
}
func tenecs_json_parseObject12() Function {
	return tenecs_json_parseObject_X(12)
}
func tenecs_json_parseObject13() Function {
	return tenecs_json_parseObject_X(13)
}
func tenecs_json_parseObject3() Function {
	return tenecs_json_parseObject_X(3)
}
func tenecs_json_parseObject6() Function {
	return tenecs_json_parseObject_X(6)
}
func tenecs_json_parseObject9() Function {
	return tenecs_json_parseObject_X(9)
}
func tenecs_json_parseObject2() Function {
	return tenecs_json_parseObject_X(2)
}
func tenecs_json_parseObject11() Function {
	return tenecs_json_parseObject_X(11)
}
func tenecs_json_parseObject14() Function {
	return tenecs_json_parseObject_X(14)
}
func tenecs_json_parseObject8() Function {
	return tenecs_json_parseObject_X(8)
}
func tenecs_json_parseObject10() Function {
	return tenecs_json_parseObject_X(10)
}
func tenecs_json_parseObject22() Function {
	return tenecs_json_parseObject_X(22)
}
func tenecs_json_parseObject4() Function {
	return tenecs_json_parseObject_X(4)
}
func tenecs_json_parseObject15() Function {
	return tenecs_json_parseObject_X(15)
}
func tenecs_json_parseObject19() Function {
	return tenecs_json_parseObject_X(19)
}
func tenecs_json_parseObject21() Function {
	return tenecs_json_parseObject_X(21)
}
func tenecs_json_parseObject17() Function {
	return tenecs_json_parseObject_X(17)
}
func tenecs_json_parseObject18() Function {
	return tenecs_json_parseObject_X(18)
}
func tenecs_json_parseObject5() Function {
	return tenecs_json_parseObject_X(5)
}
func tenecs_json_parseObject16() Function {
	return tenecs_json_parseObject_X(16)
}
func tenecs_json_parseObject20() Function {
	return tenecs_json_parseObject_X(20)
}
