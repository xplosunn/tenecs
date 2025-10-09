package embed

import "time"

// embed start

func tenecs_time__plusDays() any {
	return func(generics any, _date any, _days any) any {
		date := _date.(map[string]any)
		days := _days.(map[string]any)

		// Convert to Go time.Date
		goDate := time.Date(
			date["_year"].(map[string]any)["value"].(int),
			time.Month(date["_month"].(map[string]any)["value"].(int)),
			date["_day"].(map[string]any)["value"].(int),
			0, 0, 0, 0, time.UTC,
		)

		// Add days
		resultDate := goDate.AddDate(0, 0, days["value"].(int))

		return map[string]any{
			"$type":  "tenecs.time.Date",
			"_year":  map[string]any{"$type": "Int", "value": resultDate.Year()},
			"_month": map[string]any{"$type": "Int", "value": int(resultDate.Month())},
			"_day":   map[string]any{"$type": "Int", "value": resultDate.Day()},
		}
	}
}
