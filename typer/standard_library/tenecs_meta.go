package standard_library

var tenecs_meta = packageWith(
	withFunction("typeName", tenecs_meta_typeName),
)

var tenecs_meta_typeName = functionFromType("<T>(value: T) ~> String")
