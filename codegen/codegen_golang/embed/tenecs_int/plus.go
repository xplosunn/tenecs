package embed

// embed start

func tenecs_int__plus() any {
	return func(generics any, _a any, _b any) any {
		return map[string]any{
			"$type": "Int",
			"value": _a.(map[string]any)["value"].(int) + _b.(map[string]any)["value"].(int),
		}
	}
}
