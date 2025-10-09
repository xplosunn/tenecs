package embed

import "encoding/json"

// embed start

func tenecs_json__jsonObject0() any {
	return func(generics any, build any) any {
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

				return build.(func(any) any)([]any{})
			},
			"_toJson": func(generics any, input any) any {
				return map[string]any{
					"$type": "String",
					"value": "{}",
				}
			},
		}
	}
}
