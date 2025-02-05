package standard_library

import "github.com/xplosunn/tenecs/typer/standard_library"

func tenecs_time_Date() Function {
	return structFunction(standard_library.Tenecs_time_Date)
}
func tenecs_time_atStartOfMonth() Function {
	return function(
		params("date"),
		body(`d := date.(tenecs_time_Date)
return tenecs_time_Date{
  year: d.year,
  month: d.month,
  day: 1,
}`),
	)
}
func tenecs_time_plusYears() Function {
	return function(
		params("date", "years"),
		body(`d := date.(tenecs_time_Date)
return tenecs_time_Date{
  year: d.year.(int) + years.(int),
  month: d.month,
  day: d.day,
}`),
	)
}
