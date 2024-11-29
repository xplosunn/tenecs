package standard_library

var tenecs_boolean = packageWith(
	withFunction("and", tenecs_boolean_and),
	withFunction("not", tenecs_boolean_not),
	withFunction("or", tenecs_boolean_or),
)

var tenecs_boolean_and = functionFromSignature("(a: Boolean, b: () ~> Boolean): Boolean")

var tenecs_boolean_not = functionFromSignature("(b: Boolean): Boolean")

var tenecs_boolean_or = functionFromSignature("(a: Boolean, b: () ~> Boolean): Boolean")
