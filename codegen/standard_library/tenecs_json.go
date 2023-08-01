package standard_library

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

func tenecs_json_parseObject1() Function {
	return function(
		imports("encoding/json"),
		params("build", "fromJsonFieldI1"),
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

		i1Name := fromJsonFieldI1.(map[string]any)["name"].(string)
		i1JsonRawMessage := output[i1Name]
		if i1JsonRawMessage == nil {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not find object field \"" + i1Name + "\" in " + jsonString,
			} 
		}
		i1JsonBytes, _ := json.Marshal(&i1JsonRawMessage)
		i1 := fromJsonFieldI1.(map[string]any)["fromJson"].(map[string]any)["parse"].(func(any)any)(string(i1JsonBytes))
		i1Map, isMap := i1.(map[string]any)
		if isMap && i1Map["$type"] == "JsonError" {
			return map[string]any{
				"$type": "JsonError",
				"message": "Could not parse object field \"" + i1Name + "\": " + i1Map["message"].(string),
			} 
		}

		return build.(func(any)any)(i1)
	},
}`),
	)
}
