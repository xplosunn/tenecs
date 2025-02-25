// ##################################################################
// # The signatures of this file are generated via code-generation. #
// # Check gen.go                                                   #
// ##################################################################
package standard_library

import "github.com/xplosunn/tenecs/typer/standard_library"

func tenecs_time_Date() Function {
	return structFunction(standard_library.Tenecs_time_Date)
}
func tenecs_time_atStartOfMonth() Function {
	return function(
		params("date"),
		body(`return ({
  "$type": "Date",
  "year": date.year,
  "month": date.month,
  "day": 1
})`),
	)
}
func tenecs_time_plusYears() Function {
	return function(
		params("date", "years"),
		body(`return ({
  "$type": "Date",
  "year": date.year + years,
  "month": date.month,
  "day": date.day
})`),
	)
}
func tenecs_time_plusDays() Function {
	return function(
		params("date", "days"),
		body(`let d = new Date(date.year, date.month - 1, date.day, 0, 0, 0, 0);
d.setDate(d.getDate() + days);
return ({
  "$type": "Date",
  "year": d.getFullYear(),
  "month": d.getMonth() + 1,
  "day": d.getDate()
})`),
	)
}
