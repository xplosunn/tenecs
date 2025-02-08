package standard_library

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"strings"
)

func tenecs_json_jsonBoolean() Function {
	return function(
		imports("encoding/json"),
		body(`return tenecs_json_JsonSchema{
	_fromJson: func(input any) any {
		jsonString := input.(string)
		var output bool
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return tenecs_error_Error{
				_message: "Could not parse Boolean from " + jsonString,
			} 
		}
		return output
	},
	_toJson: func(input any) any {
		result, _ := json.Marshal(input)
		return string(result)
	},
}`),
	)
}

func tenecs_json_jsonInt() Function {
	return function(
		imports("encoding/json"),
		body(`return tenecs_json_JsonSchema{
	_fromJson: func(input any) any {
		jsonString := input.(string)
		var output float64
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil || float64(int(output)) != output {
			return tenecs_error_Error{
				_message: "Could not parse Int from " + jsonString,
			} 
		}
		return int(output)
	},
	_toJson: func(input any) any {
		result, _ := json.Marshal(input)
		return string(result)
	},
}`),
	)
}

func tenecs_json_jsonOr() Function {
	return function(
		imports("encoding/json"),
		params("schemaA", "schemaB", "toJsonSchemaPicker"),
		body(`return tenecs_json_JsonSchema{
	_fromJson: func(input any) any {
		resultA := schemaA.(tenecs_json_JsonSchema)._fromJson.(func(any)any)(input)
		_, isMap := resultA.(tenecs_error_Error)
		if isMap {
			resultB := schemaB.(tenecs_json_JsonSchema)._fromJson.(func(any)any)(input)
			_, isMap := resultB.(tenecs_error_Error)
			if isMap {
				jsonString := input.(string)
				return tenecs_error_Error{
					_message: "Could not parse from " + jsonString,
				}
			}
			return resultB
		}
		return resultA
	},
	_toJson: func(input any) any {
		schema := toJsonSchemaPicker.(func(any)any)(input)
		return schema.(tenecs_json_JsonSchema)._toJson.(func(any)any)(input)
	},
}`),
	)
}

func tenecs_json_jsonString() Function {
	return function(
		imports("encoding/json"),
		body(`return tenecs_json_JsonSchema{
	_fromJson: func(input any) any {
		jsonString := input.(string)
		var output string
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return tenecs_error_Error{
				_message: "Could not parse String from " + jsonString,
			} 
		}
		return output
	},
	_toJson: func(input any) any {
		result, _ := json.Marshal(input)
		return string(result)
	},
}`),
	)
}

func tenecs_json_jsonList() Function {
	return function(
		imports("encoding/json", "strings"),
		params("of"),
		body(`return tenecs_json_JsonSchema{
	_fromJson: func(input any) any {
		jsonString := input.(string)
		var output []json.RawMessage
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return tenecs_error_Error{
				_message: "Could not parse List from " + jsonString,
			} 
		}
		if len(output) == 0 {
			return []any{}
		}
		ofParse := of.(tenecs_json_JsonSchema)._fromJson.(func(any)any)
		outputList := []any{}
		for _, elem := range output {
			elemJsonBytes, _ := json.Marshal(&elem)
			result := ofParse(string(elemJsonBytes))
			resultMap, isMap := result.(tenecs_error_Error)
			if isMap {
				return resultMap
			}
			outputList = append(outputList, result)
		}
		return outputList
	},
	_toJson: func(input any) any {
		results := []string{}
		ofToJson := of.(tenecs_json_JsonSchema)._toJson.(func(any)any)
		for _, elem := range input.([]any) {
			result := ofToJson(elem)
			results = append(results, result.(string))
		}
		return "[" + strings.Join(results, ",") + "]"
	},
}`),
	)
}

func tenecs_json_jsonObject0() Function {
	return function(
		imports("encoding/json"),
		params("build"),
		body(`return tenecs_json_JsonSchema{
	_fromJson: func(input any) any {
		jsonString := input.(string)
		var output map[string]json.RawMessage
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return tenecs_error_Error{
				_message: "Could not parse object from " + jsonString,
			} 
		}

		return build.(func()any)()
	},
	_toJson: func(input any) any {
		return "{}"
	},
}`),
	)
}

func tenecs_json_jsonObject_X(x int) Function {
	paramNames := []string{"build"}
	for i := 0; i < x; i++ {
		paramNames = append(paramNames, fmt.Sprintf("jsonSchemaFieldI%d", i))
	}
	bodyStr := `return tenecs_json_JsonSchema{
	_fromJson: func(input any) any {
		jsonString := input.(string)
		var output map[string]json.RawMessage
		err := json.Unmarshal([]byte(jsonString), &output)
		if err != nil {
			return tenecs_error_Error{
				_message: "Could not parse object from " + jsonString,
			} 
		}
`
	for i := 0; i < x; i++ {
		bodyStr += fmt.Sprintf(`
		i%dName := jsonSchemaFieldI%d.(tenecs_json_JsonField)._name.(string)
		i%dJsonRawMessage := output[i%dName]
		if i%dJsonRawMessage == nil {
			return tenecs_error_Error{
				_message: "Could not find object field \"" + i%dName + "\" in " + jsonString,
			}
		}
		i%dJsonBytes, _ := json.Marshal(&i%dJsonRawMessage)
		i%d := jsonSchemaFieldI%d.(tenecs_json_JsonField)._schema.(tenecs_json_JsonSchema)._fromJson.(func(any)any)(string(i%dJsonBytes))
		i%dMap, isMap := i%d.(tenecs_error_Error)
		if isMap {
			return tenecs_error_Error{
				_message: "Could not parse object field \"" + i%dName + "\": " + i%dMap._message.(string),
			} 
		}`, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i)
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
`, strings.Join(anys, ","), strings.Join(buildArgs, ","))

	bodyStr += `
	_toJson: func(input any) any {
		output := map[string]string{}
`
	for i := 0; i < x; i++ {
		bodyStr += fmt.Sprintf(`
		fieldI%d := jsonSchemaFieldI%d.(tenecs_json_JsonField)._access.(func(any)any)(input)
		i%d := jsonSchemaFieldI%d.(tenecs_json_JsonField)._schema.(tenecs_json_JsonSchema)._toJson.(func(any)any)(fieldI%d)
		output[jsonSchemaFieldI%d.(tenecs_json_JsonField)._name.(string)] = i%d.(string)
`, i, i, i, i, i, i, i)
	}
	bodyStr += `
		result := "{"
		i := 0
		outputKeysSorted := []string{}
		for k, _ := range output {
			outputKeysSorted = append(outputKeysSorted, k)
		}
		sort.Strings(outputKeysSorted)
		for _, k := range outputKeysSorted {
			v := output[k]
			nameBytes, _ := json.Marshal(k)
			result += string(nameBytes) + ":" + v
			i += 1
			if i < len(output) {
				result += ","
			}
		}
		result += "}"
		return result
	},
`

	bodyStr += `
}`

	return function(
		imports("encoding/json", "sort"),
		params(paramNames...),
		body(bodyStr),
	)
}

func tenecs_json_jsonObject1() Function {
	return tenecs_json_jsonObject_X(1)
}
func tenecs_json_jsonObject7() Function {
	return tenecs_json_jsonObject_X(7)
}
func tenecs_json_jsonObject12() Function {
	return tenecs_json_jsonObject_X(12)
}
func tenecs_json_jsonObject13() Function {
	return tenecs_json_jsonObject_X(13)
}
func tenecs_json_jsonObject3() Function {
	return tenecs_json_jsonObject_X(3)
}
func tenecs_json_jsonObject6() Function {
	return tenecs_json_jsonObject_X(6)
}
func tenecs_json_jsonObject9() Function {
	return tenecs_json_jsonObject_X(9)
}
func tenecs_json_jsonObject2() Function {
	return tenecs_json_jsonObject_X(2)
}
func tenecs_json_jsonObject11() Function {
	return tenecs_json_jsonObject_X(11)
}
func tenecs_json_jsonObject14() Function {
	return tenecs_json_jsonObject_X(14)
}
func tenecs_json_jsonObject8() Function {
	return tenecs_json_jsonObject_X(8)
}
func tenecs_json_jsonObject10() Function {
	return tenecs_json_jsonObject_X(10)
}
func tenecs_json_jsonObject22() Function {
	return tenecs_json_jsonObject_X(22)
}
func tenecs_json_jsonObject4() Function {
	return tenecs_json_jsonObject_X(4)
}
func tenecs_json_jsonObject15() Function {
	return tenecs_json_jsonObject_X(15)
}
func tenecs_json_jsonObject19() Function {
	return tenecs_json_jsonObject_X(19)
}
func tenecs_json_jsonObject21() Function {
	return tenecs_json_jsonObject_X(21)
}
func tenecs_json_jsonObject17() Function {
	return tenecs_json_jsonObject_X(17)
}
func tenecs_json_jsonObject18() Function {
	return tenecs_json_jsonObject_X(18)
}
func tenecs_json_jsonObject5() Function {
	return tenecs_json_jsonObject_X(5)
}
func tenecs_json_jsonObject16() Function {
	return tenecs_json_jsonObject_X(16)
}
func tenecs_json_jsonObject20() Function {
	return tenecs_json_jsonObject_X(20)
}
func tenecs_json_JsonField() Function {
	return structFunction(standard_library.Tenecs_json_JsonField)
}
func tenecs_json_JsonSchema() Function {
	return structFunction(standard_library.Tenecs_json_JsonSchema)
}
