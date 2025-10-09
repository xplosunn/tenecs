package embed

// embed start

func tenecs_list__filter() any {
	return func(generics any, list any, keep any) any {
		result := []any{}
		for _, elem := range list.(map[string]any)["value"].([]any) {
			if keep.(func(any, any) any)([]any{}, elem).(map[string]any)["value"].(bool) {
				result = append(result, elem)
			}
		}
		return map[string]any{
			"$type": list.(map[string]any)["$type"],
			"value": result,
		}
	}
}
