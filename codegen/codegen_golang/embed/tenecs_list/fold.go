package embed

// embed start

func tenecs_list__fold() any {
	return func(generics any, list any, zero any, f any) any {
		result := zero
		for _, elem := range list.(map[string]any)["value"].([]any) {
			result = f.(func(any, any, any) any)([]any{}, result, elem)
		}
		return result
	}
}
