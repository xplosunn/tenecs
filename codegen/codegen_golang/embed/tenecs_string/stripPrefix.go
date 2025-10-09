package embed

import "strings"

// embed start

func tenecs_string__stripPrefix() any {
	return func(generics any, str any, prefix any) any {
		result := strings.TrimPrefix(str.(map[string]any)["value"].(string), prefix.(map[string]any)["value"].(string))
		return map[string]any{"$type": "String", "value": result}
	}
}
