package embed

// embed start

func tenecs_int__abs() any {
	return func(generics any, _a any) any {
		a := _a.(map[string]any)["value"].(int)
		if a < 0 {
			return map[string]any{
				"$type": "Int",
				"value": -a,
			}
		}
		return map[string]any{
			"$type": "Int",
			"value": a,
		}
	}
}
