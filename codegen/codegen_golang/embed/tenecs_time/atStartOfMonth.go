package embed

// embed start

func tenecs_time__atStartOfMonth() any {
	return func(generics any, _date any) any {
		date := _date.(map[string]any)
		return map[string]any{
			"$type":  "tenecs.time.Date",
			"_year":  date["_year"],
			"_month": date["_month"],
			"_day":   map[string]any{"$type": "Int", "value": 1},
		}
	}
}
