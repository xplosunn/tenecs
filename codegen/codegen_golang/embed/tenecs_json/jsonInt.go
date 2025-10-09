package embed

import "encoding/json"

// embed start

func tenecs_json__jsonInt() any {
	return func(generics any) any {
		return map[string]any{
			"$type": "tenecs.json.JsonConverter",
			"_fromJson": func(generics any, input any) any {
				jsonString := input.(map[string]any)["value"].(string)
				var output float64
				err := json.Unmarshal([]byte(jsonString), &output)
				if err != nil || float64(int(output)) != output {
					return map[string]any{
						"$type": "tenecs.error.Error",
						"_message": map[string]any{
							"$type": "String",
							"value": "Could not parse Int from " + jsonString,
						},
					}
				}
				return map[string]any{
					"$type": "Int",
					"value": int(output),
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
