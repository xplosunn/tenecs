package embed

import "strings"

// embed start

func tenecs_string__toLowerCase() any {
	return func(generics any, str any) any {
		return map[string]any{"$type": "String", "value": strings.ToLower(str.(map[string]any)["value"].(string))}
	}
}
