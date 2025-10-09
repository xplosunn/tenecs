package embed

// embed start

func tenecs_list__first() any {
	return func(generics any, list any) any {
		l := list.(map[string]any)["value"].([]any)
		if len(l) > 0 {
			return l[0]
		}
		return map[string]any{
			"$type": "Void",
			"value": nil,
		}
	}
}
