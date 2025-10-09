package embed

// embed start

func tenecs_time__plusYears() any {
	return func(generics any, _date any, _years any) any {
		date := _date.(map[string]any)
		years := _years.(map[string]any)
		return map[string]any{
			"$type":  "tenecs.time.Date",
			"_year":  map[string]any{"$type": "Int", "value": date["_year"].(map[string]any)["value"].(int) + years["value"].(int)},
			"_month": date["_month"],
			"_day":   date["_day"],
		}
	}
}
