package embed

// embed start

func tenecs_string__join() any {
	return func(generics any, left any, right any) any {
		leftStr := left.(map[string]any)["value"].(string)
		rightStr := right.(map[string]any)["value"].(string)
		return map[string]any{"$type": "String", "value": leftStr + rightStr}
	}
}
