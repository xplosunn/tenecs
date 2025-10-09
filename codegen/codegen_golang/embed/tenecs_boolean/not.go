package embed

// embed start

func tenecs_boolean__not() any {
	return func(generics any, _b any) any {
		return map[string]any{
			"$type": "Boolean",
			"value": !(_b.(map[string]any)["value"].(bool)),
		}
	}
}
