package embed

import (
	"encoding/json"
	"sort"
)

// embed start

func tenecs_json__jsonObject1() any {
	return func(generics any, build any, jsonConverterFieldI0 any) any {
		return map[string]any{
			"$type": "tenecs.json.JsonConverter",
			"_fromJson": func(generics any, input any) any {
				jsonString := input.(map[string]any)["value"].(string)
				var output map[string]json.RawMessage
				err := json.Unmarshal([]byte(jsonString), &output)
				if err != nil {
					return map[string]any{
						"$type": "tenecs.error.Error",
						"_message": map[string]any{
							"$type": "String",
							"value": "Could not parse object from " + jsonString,
						},
					}
				}

				i0Name := jsonConverterFieldI0.(map[string]any)["_name"].(map[string]any)["value"].(string)
				i0JsonRawMessage := output[i0Name]
				if i0JsonRawMessage == nil {
					return map[string]any{
						"$type": "tenecs.error.Error",
						"_message": map[string]any{
							"$type": "String",
							"value": "Could not find object field \"" + i0Name + "\" in " + jsonString,
						},
					}
				}
				i0JsonBytes, _ := json.Marshal(&i0JsonRawMessage)
				i0JsonString := map[string]any{
					"$type": "String",
					"value": string(i0JsonBytes),
				}
				i0 := jsonConverterFieldI0.(map[string]any)["_Converter"].(map[string]any)["_fromJson"].(func(any, any) any)([]any{}, i0JsonString)
				i0Map, isMap := i0.(map[string]any)
				if isMap && i0Map["$type"] == "tenecs.error.Error" {
					return map[string]any{
						"$type": "tenecs.error.Error",
						"_message": map[string]any{
							"$type": "String",
							"value": "Could not parse object field \"" + i0Name + "\": " + i0Map["_message"].(map[string]any)["value"].(string),
						},
					}
				}

				return build.(func(any, any) any)([]any{}, i0)
			},
			"_toJson": func(generics any, input any) any {
				output := map[string]string{}

				fieldI0 := jsonConverterFieldI0.(map[string]any)["_access"].(func(any, any) any)([]any{}, input)
				i0 := jsonConverterFieldI0.(map[string]any)["_Converter"].(map[string]any)["_toJson"].(func(any, any) any)([]any{}, fieldI0)
				output[jsonConverterFieldI0.(map[string]any)["_name"].(map[string]any)["value"].(string)] = i0.(map[string]any)["value"].(string)

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
				return map[string]any{
					"$type": "String",
					"value": result,
				}
			},
		}
	}
}
