package embed

// embed start

func tenecs_meta__typeName() any {
	return func(generics any, value any) any {
		result := value.(map[string]any)["$type"]
		return map[string]any{
			"$type": "String",
			"value": result,
		}
	}
}
