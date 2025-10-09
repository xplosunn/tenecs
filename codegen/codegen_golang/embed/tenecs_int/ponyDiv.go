package embed

// embed start

func tenecs_int__ponyDiv() any {
	return func(generics any, _a any, _b any) any {
		b := _b.(map[string]any)["value"].(int)
		if b == 0 {
			return map[string]any{
				"$type": "Int",
				"value": 0,
			}
		}
		return map[string]any{
			"$type": "Int",
			"value": _a.(map[string]any)["value"].(int) / b,
		}
	}
}
