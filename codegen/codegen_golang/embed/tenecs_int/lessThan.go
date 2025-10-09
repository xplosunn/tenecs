package embed

// embed start

func tenecs_int__lessThan() any {
	return func(generics any, _a any, _b any) any {
		return map[string]any{
			"$type": "Boolean",
			"value": _a.(map[string]any)["value"].(int) < _b.(map[string]any)["value"].(int),
		}
	}
}
