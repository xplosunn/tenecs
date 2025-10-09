package embed

// embed start

func tenecs_error__Error() any {
	return func(generics any, _message any) any {
		return map[string]any{
			"$type":    "tenecs.error.Error",
			"_message": _message,
		}
	}
}
