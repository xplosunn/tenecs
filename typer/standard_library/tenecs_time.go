package standard_library

import "github.com/xplosunn/tenecs/typer/types"

var tenecs_time = packageWith(
	withStruct(Tenecs_time_Date),
	withFunction("atStartOfMonth", Tenecs_time_atStartOfMonth),
	withFunction("plusYears", Tenecs_time_plusYears),
)

var Tenecs_time_Date = structWithFields("Date", tenecs_time_Date, tenecs_time_Date_Fields...)

var tenecs_time_Date = types.Struct(
	"tenecs.time",
	"Date",
	nil,
)

var tenecs_time_Date_Fields = []func(fields *StructWithFields){
	structField("year", types.Int()),
	structField("month", types.Int()),
	structField("day", types.Int()),
}

var Tenecs_time_atStartOfMonth = functionFromType("(date: Date) ~> Date", Tenecs_time_Date)

var Tenecs_time_plusYears = functionFromType("(date: Date, years: Int) ~> Date", Tenecs_time_Date)
