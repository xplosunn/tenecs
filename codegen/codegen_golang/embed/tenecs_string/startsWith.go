package embed

import "strings"

// embed start

func tenecs_string__startsWith() any {
	return func(generics any, str any, prefix any) any {
		return map[string]any{"$type": "Boolean", "value": strings.HasPrefix(str.(map[string]any)["value"].(string), prefix.(map[string]any)["value"].(string))}
	}
}
