package embed

// embed start

func tenecs_list__length() any {
	return func(generics any, list any) any {
		return map[string]any{
			"$type": "Int",
			"value": len(list.(map[string]any)["value"].([]any)),
		}
	}
}
