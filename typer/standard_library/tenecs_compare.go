package standard_library

var tenecs_compare = packageWith(
	withFunction("eq", tenecs_compare_eq),
)

var tenecs_compare_eq = functionFromType("<T>(first: T, second: T) ~> Boolean")
