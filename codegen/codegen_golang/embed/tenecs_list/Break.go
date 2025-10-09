package embed

// embed start

func tenecs_list__Break() any {
	return func(generics any, value any) any {
		return map[string]any{
			"$type":  "tenecs.list.Break<" + typeNameCollapseRecursive(generics.([]any)[0]) + ">",
			"_value": value,
		}
	}
}
