package embed

import (
	"slices"
	"strings"
)

// embed start

func tenecs_list__mapUntil() any {
	return func(generics any, list any, f any) any {
		result := []any{}
		for _, elem := range list.(map[string]any)["value"].([]any) {
			maybeBreak := f.(func(any, any) any)([]any{}, elem)
			if strings.HasPrefix(maybeBreak.(map[string]any)["$type"].(string), "tenecs.list.Break<") {
				left := generics.([]any)[2]
				right := maybeBreak.(map[string]any)["_value"].(map[string]any)["$type"]
				if _, ok := any(left).([]any); ok && slices.Contains(any(left).([]any), right) || (!ok && (left == right)) {
					return maybeBreak.(map[string]any)["_value"]
				} else {
					result = append(result, maybeBreak)
				}
			} else {
				result = append(result, maybeBreak)
			}
		}
		return map[string]any{
			"$type": "List<" + typeNameCollapseRecursive(generics.([]any)[1]) + ">",
			"value": result,
		}
	}
}
