package embed

// embed start

func tenecs_string__firstChar() any {
	return func(generics any, str any) any {
		runes := []rune(str.(map[string]any)["value"].(string))
		if len(runes) > 0 {
			return map[string]any{"$type": "String", "value": string(runes[0])}
		} else {
			return map[string]any{"$type": "String", "value": ""}
		}
	}
}
