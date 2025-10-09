package embed

import "strings"

// embed start

func tenecs_string__repeat() any {
	return func(generics any, str any, count any) any {
		return map[string]any{
			"$type": "String",
			"value": strings.Repeat(
				str.(map[string]any)["value"].(string),
				count.(map[string]any)["value"].(int),
			),
		}
	}
}
