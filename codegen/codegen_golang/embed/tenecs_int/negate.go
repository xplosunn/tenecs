package embed

// embed start

func tenecs_int__negate() any {
	return func(generics any, _a any) any {
		return map[string]any{
			"$type": "Int",
			"value": -_a.(map[string]any)["value"].(int),
		}
	}
}
