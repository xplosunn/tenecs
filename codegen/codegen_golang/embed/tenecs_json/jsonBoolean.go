package embed

import "encoding/json"

// embed start

func tenecs_json__jsonBoolean() any {
	return func(generics any) any {
		return map[string]any{
			"$type": "tenecs.json.JsonConverter",
			"_fromJson": func(generics any, input any) any {
				jsonString := input.(map[string]any)["value"].(string)
				var output bool
				err := json.Unmarshal([]byte(jsonString), &output)
				if err != nil {
					return map[string]any{
						"$type": "tenecs.error.Error",
						"_message": map[string]any{
							"$type": "String",
							"value": "Could not parse Boolean from " + jsonString,
						},
					}
				}
				return map[string]any{
					"$type": "Boolean",
					"value": output,
				}
			},
			"_toJson": func(generics any, input any) any {
				result, _ := json.Marshal(input.(map[string]any)["value"])
				return map[string]any{
					"$type": "String",
					"value": string(result),
				}
			},
		}
	}
}
