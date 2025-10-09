package embed

import "strings"

// embed start

func tenecs_string__isBlank() any {
	return func(generics any, str any) any {
		return map[string]any{"$type": "Boolean", "value": len(strings.TrimSpace(str.(map[string]any)["value"].(string))) == 0}
	}
}
