package embed

// embed start

func tenecs_list__flatMap() any {
	return func(generics any, list any, f any) any {
		result := []any{}
		for _, elem := range list.(map[string]any)["value"].([]any) {
			result = append(result, f.(func(any, any) any)([]any{}, elem).(map[string]any)["value"].([]any)...)
		}
		return map[string]any{
			"$type": "List<" + typeNameCollapseRecursive(generics.([]any)[1]) + ">",
			"value": result,
		}
	}
}
