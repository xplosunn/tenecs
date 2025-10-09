package embed

import "strings"

// embed start

func tenecs_string__trimRight() any {
	return func(generics any, str any) any {
		return map[string]any{"$type": "String", "value": strings.TrimRight(str.(map[string]any)["value"].(string), " \t\n\r")}
	}
}
