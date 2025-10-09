package embed

// embed start

func tenecs_list__atIndexGet() any {
	return func(generics any, list any, index any) any {
		l := list.(map[string]any)["value"].([]any)
		idx := index.(map[string]any)["value"].(int)
		if idx >= 0 && len(l) > idx {
			return l[idx]
		}
		return map[string]any{
			"$type": "tenecs.error.Error",
			"_message": map[string]any{
				"$type": "String",
				"value": "Out of bounds",
			},
		}
	}
}
