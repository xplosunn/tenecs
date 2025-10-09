package embed

// embed start

func tenecs_string__firstCharCode() any {
	return func(generics any, str any) any {
		runes := []rune(str.(map[string]any)["value"].(string))
		if len(runes) > 0 {
			return map[string]any{"$type": "Int", "value": int(runes[0])}
		} else {
			return map[string]any{"$type": "Int", "value": -1}
		}
	}
}
