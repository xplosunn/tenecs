package embed

// embed start

func tenecs_string__length() any {
	return func(generics any, str any) any {
		return map[string]any{"$type": "Int", "value": len([]rune(str.(map[string]any)["value"].(string)))}
	}
}
