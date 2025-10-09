package embed

// embed start

func tenecs_list__mapNotNull() any {
	return func(generics any, list any, f any) any {
		result := []any{}
		for _, elem := range list.(map[string]any)["value"].([]any) {
			maybeNull := f.(func(any, any) any)([]any{}, elem)
			if maybeNull.(map[string]any)["$type"] != "Void" {
				result = append(result, maybeNull)
			}
		}
		return map[string]any{
			"$type": "List<" + typeNameCollapseRecursive(generics.([]any)[1]) + ">",
			"value": result,
		}
	}
}
