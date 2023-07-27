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
	var x float64
	if float64(int(x)) == x {

	}
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
