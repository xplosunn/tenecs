package embed

// embed start

func tenecs_list__appendAll() any {
	return func(generics any, list any, newElements any) any {
		return map[string]any{
			"$type": list.(map[string]any)["$type"],
			"value": append(list.(map[string]any)["value"].([]any), newElements.(map[string]any)["value"].([]any)...),
		}
	}
}
