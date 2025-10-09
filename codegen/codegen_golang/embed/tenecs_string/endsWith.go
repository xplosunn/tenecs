package embed

import "strings"

// embed start

func tenecs_string__endsWith() any {
	return func(generics any, str any, suffix any) any {
		return map[string]any{"$type": "Boolean", "value": strings.HasSuffix(str.(map[string]any)["value"].(string), suffix.(map[string]any)["value"].(string))}
	}
}
