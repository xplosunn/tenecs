package embed

// embed start

func tenecs_string__characters() any {
	return func(generics any, str any) any {
		input := str.(map[string]any)["value"].(string)
		runes := []rune(input)
		result := []any{}
		for _, r := range runes {
			result = append(result, map[string]any{"$type": "String", "value": string(r)})
		}
		return map[string]any{"$type": "List<String>", "value": result}
	}
}
