package embed

// embed start

func createMainRuntime() any {
	log := func(generics any, msg any) any {
		println(msg.(map[string]any)["value"].(string))
		return map[string]any{"$type": "Void", "value": nil}
	}
	console := map[string]any{
		"$type": "tenecs.go.Console",
		"_log":  log,
	}
	refCreator := map[string]any{
		"$type": "tenecs.ref.RefCreator",
		"_new": func(generics any, value any) any {
			var ref any = value
			return map[string]any{
				"$type":    "tenecs.ref.Ref",
				"generics": []any{},
				"_get": func(generics any) any {
					return ref
				},
				"_set": func(generics any, value any) any {
					ref = value
					return map[string]any{"$type": "Void", "value": nil}
				},
				"_modify": func(generics any, f any) any {
					ref = f.(func(any, any) any)([]any{}, ref)
					return map[string]any{"$type": "Void", "value": nil}
				},
			}
		},
	}
	runtime := map[string]any{
		"$type":    "tenecs.go.Runtime",
		"_console": console,
		"_ref":     refCreator,
	}
	return runtime
}
