package standard_library

var tenecs_boolean = packageWith(
	withFunction("and", tenecs_boolean_and),
	withFunction("not", tenecs_boolean_not),
	withFunction("or", tenecs_boolean_or),
)

var tenecs_boolean_and = functionFromType("(a: Boolean, b: () ~> Boolean) ~> Boolean")

var tenecs_boolean_not = functionFromType("(b: Boolean) ~> Boolean")

var tenecs_boolean_or = functionFromType("(a: Boolean, b: () ~> Boolean) ~> Boolean")
