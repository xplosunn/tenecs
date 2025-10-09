package embed

// embed start

func tenecs_list__append() any {
	return func(generics any, list any, newElement any) any {
		return map[string]any{
			"$type": list.(map[string]any)["$type"],
			"value": append(list.(map[string]any)["value"].([]any), newElement),
		}
	}
}
