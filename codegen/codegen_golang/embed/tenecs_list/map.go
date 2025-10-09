package embed

// embed start

func tenecs_list__map() any {
	return func(generics any, list any, f any) any {
		result := []any{}
		for _, elem := range list.(map[string]any)["value"].([]any) {
			result = append(result, f.(func(any, any) any)([]any{}, elem))
		}

		listTypeArgument := generics.([]any)[1]
		var listOf string = typeNameCollapseRecursive(listTypeArgument)
		return map[string]any{
			"$type": "List<" + listOf + ">",
			"value": result,
		}
	}
}
