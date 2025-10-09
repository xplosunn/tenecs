package embed

// embed start

func tenecs_string__isEmpty() any {
	return func(generics any, str any) any {
		return map[string]any{"$type": "Boolean", "value": len(str.(map[string]any)["value"].(string)) == 0}
	}
}
