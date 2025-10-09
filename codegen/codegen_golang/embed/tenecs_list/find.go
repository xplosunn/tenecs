package embed

// embed start

func tenecs_list__find() any {
	return func(generics any, list any, f any) any {
		for _, elem := range list.(map[string]any)["value"].([]any) {
			maybeReturn := f.(func(any, any) any)([]any{}, elem)
			if maybeReturn.(map[string]any)["$type"] != "Void" {
				return maybeReturn
			}
		}
		return map[string]any{"$type": "Void", "value": nil}
	}
}
