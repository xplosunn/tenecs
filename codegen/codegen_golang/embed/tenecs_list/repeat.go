package embed

// embed start

func tenecs_list__repeat() any {
	return func(generics any, elem any, times any) any {
		result := []any{}
		for i := 0; i < times.(map[string]any)["value"].(int); i++ {
			result = append(result, elem)
		}
		return map[string]any{
			"$type": "List<" + typeNameCollapseRecursive(generics.([]any)[0]) + ">",
			"value": result,
		}
	}
}
