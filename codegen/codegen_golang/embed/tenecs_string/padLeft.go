package embed

import "strings"

// embed start

func tenecs_string__padLeft() any {
	return func(generics any, str any, length any, padChar any) any {
		padCount := length.(map[string]any)["value"].(int) - len([]rune(str.(map[string]any)["value"].(string)))
		if padCount <= 0 {
			return str
		}
		return map[string]any{
			"$type": "String",
			"value": strings.Repeat(padChar.(map[string]any)["value"].(string), padCount) + str.(map[string]any)["value"].(string),
		}
	}
}
