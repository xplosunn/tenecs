package embed

// embed start

func tenecs_list__atIndexSet() any {
	return func(generics any, list any, index any, setTo any) any {
		l := list.(map[string]any)["value"].([]any)
		idx := index.(map[string]any)["value"].(int)
		if idx >= 0 && len(l) > idx {
			result := make([]any, len(l))
			copy(result, l)
			result[idx] = setTo

			listTypeArgument := generics.([]any)[0]
			var listOf string = typeNameCollapseRecursive(listTypeArgument)
			return map[string]any{
				"$type": "List<" + listOf + ">",
				"value": result,
			}
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
