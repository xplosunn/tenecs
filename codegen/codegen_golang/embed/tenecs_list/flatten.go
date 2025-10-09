package embed

// embed start

func tenecs_list__flatten() any {
	return func(generics any, list any) any {
		result := []any{}
		for _, innerList := range list.(map[string]any)["value"].([]any) {
			result = append(result, innerList.(map[string]any)["value"].([]any)...)
		}
		return map[string]any{
			"$type": "List<" + typeNameCollapseRecursive(generics.([]any)[0]) + ">",
			"value": result,
		}
	}
}
