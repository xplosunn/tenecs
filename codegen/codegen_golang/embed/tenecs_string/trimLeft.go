package embed

import "strings"

// embed start

func tenecs_string__trimLeft() any {
	return func(generics any, str any) any {
		return map[string]any{"$type": "String", "value": strings.TrimLeft(str.(map[string]any)["value"].(string), " \t\n\r")}
	}
}
