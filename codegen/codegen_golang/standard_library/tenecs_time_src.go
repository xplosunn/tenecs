package standard_library

import (
	"github.com/xplosunn/tenecs/typer/standard_library"
)

func tenecs_time_Date() Function {
	return structFunction(standard_library.Tenecs_time_Date)
}
func tenecs_time_atStartOfMonth() Function {
	return function(
		params("date"),
		body(`d := date.(tenecs_time_Date)
return tenecs_time_Date{
  _year: d._year,
  _month: d._month,
  _day: 1,
}`),
	)
}
func tenecs_time_plusYears() Function {
	return function(
		params("date", "years"),
		body(`d := date.(tenecs_time_Date)
return tenecs_time_Date{
  _year: d._year.(int) + years.(int),
  _month: d._month,
  _day: d._day,
}`),
	)
}
func tenecs_time_plusDays() Function {
	return function(
		imports("time"),
		params("date", "days"),
		body(`tenecsDate := date.(tenecs_time_Date)
d := time.Date(tenecsDate._year.(int), time.Month(tenecsDate._month.(int)), tenecsDate._day.(int), 1, 10, 30, 0, time.UTC)
d = d.AddDate(0, 0, days.(int))

return tenecs_time_Date{
  _year: d.Year(),
  _month: int(d.Month()),
  _day: d.Day(),
}`),
	)
}
