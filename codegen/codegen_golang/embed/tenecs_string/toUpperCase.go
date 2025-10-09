package embed

import "strings"

// embed start

func tenecs_string__toUpperCase() any {
	return func(generics any, str any) any {
		return map[string]any{"$type": "String", "value": strings.ToUpper(str.(map[string]any)["value"].(string))}
	}
}
