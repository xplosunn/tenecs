package embed

// embed start

func tenecs_string__reverse() any {
	return func(generics any, str any) any {
		runes := []rune(str.(map[string]any)["value"].(string))
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return map[string]any{"$type": "String", "value": string(runes)}
	}
}
