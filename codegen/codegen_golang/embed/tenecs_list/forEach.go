package embed

// embed start

func tenecs_list__forEach() any {
	return func(generics any, list any, f any) any {
		for _, elem := range list.(map[string]any)["value"].([]any) {
			f.(func(any, any) any)([]any{}, elem)
		}
		return map[string]any{
			"$type": "Void",
			"value": nil,
		}
	}
}
