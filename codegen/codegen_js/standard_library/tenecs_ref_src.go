package standard_library

func tenecs_ref_Ref() Function {
	return function(
		params("get", "set", "modify"),
		body(`return ({
  "$type": "Ref",
  "get": get,
  "set": set,
  "modify": modify
})`),
	)
}
func tenecs_ref_RefCreator() Function {
	return function(
		params("_new"),
		body(`return ({
  "$type": "RefCreator",
  "new": _new,
})`),
	)
}
