package standard_library

func tenecs_error_error() Function {
	return tenecs_error_Error()
}
func tenecs_error_Error() Function {
	return function(
		params("message"),
		body(`return ({
  "$type": "Error",
  "message": message
})`),
	)
}
