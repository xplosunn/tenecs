package embed

// embed start

func tenecs_boolean__or() any {
	return func(generics any, _a any, _b any) any {
		return map[string]any{
			"$type": "Boolean",
			"value": _a.(map[string]any)["value"].(bool) || _b.(func(any) any)([]any{}).(map[string]any)["value"].(bool),
		}
	}
}
