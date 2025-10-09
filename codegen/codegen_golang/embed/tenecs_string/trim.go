package embed

import "strings"

// embed start

func tenecs_string__trim() any {
	return func(generics any, str any) any {
		return map[string]any{"$type": "String", "value": strings.TrimSpace(str.(map[string]any)["value"].(string))}
	}
}
