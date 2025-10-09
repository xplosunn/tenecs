package embed

// embed start

func tenecs_json__jsonOr() any {
	return func(generics any, ConverterA any, ConverterB any, toJsonConverterPicker any) any {
		return map[string]any{
			"$type": "tenecs.json.JsonConverter",
			"_fromJson": func(generics any, input any) any {
				resultA := ConverterA.(map[string]any)["_fromJson"].(func(any, any) any)([]any{}, input)
				resultAMap, isMap := resultA.(map[string]any)
				if isMap && resultAMap["$type"] == "tenecs.error.Error" {
					resultB := ConverterB.(map[string]any)["_fromJson"].(func(any, any) any)([]any{}, input)
					resultBMap, isMap := resultB.(map[string]any)
					if isMap && resultBMap["$type"] == "tenecs.error.Error" {
						jsonString := input.(map[string]any)["value"].(string)
						return map[string]any{
							"$type": "tenecs.error.Error",
							"_message": map[string]any{
								"$type": "String",
								"value": "Could not parse from " + jsonString,
							},
						}
					}
					return resultB
				}
				return resultA
			},
			"_toJson": func(generics any, input any) any {
				Converter := toJsonConverterPicker.(func(any, any) any)([]any{}, input)
				return Converter.(map[string]any)["_toJson"].(func(any, any) any)([]any{}, input)
			},
		}
	}
}
