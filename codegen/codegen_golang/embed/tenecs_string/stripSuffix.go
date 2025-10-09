package embed

import "strings"

// embed start

func tenecs_string__stripSuffix() any {
	return func(generics any, str any, suffix any) any {
		result := strings.TrimSuffix(str.(map[string]any)["value"].(string), suffix.(map[string]any)["value"].(string))
		return map[string]any{"$type": "String", "value": result}
	}
}
