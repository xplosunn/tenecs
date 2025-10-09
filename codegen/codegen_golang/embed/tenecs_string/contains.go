package embed

import "strings"

// embed start

func tenecs_string__contains() any {
	return func(generics any, str any, subStr any) any {
		result := strings.Contains(str.(map[string]any)["value"].(string), subStr.(map[string]any)["value"].(string))
		return map[string]any{"$type": "Boolean", "value": result}
	}
}
