package embed

import "reflect"

// embed start

func tenecs_compare__eq() any {
	return func(generics any, _a any, _b any) any {
		return map[string]any{
			"$type": "Boolean",
			"value": reflect.DeepEqual(_a, _b),
		}
	}
}
