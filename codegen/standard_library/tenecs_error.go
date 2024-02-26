package standard_library

func tenecs_error_error() Function {
	return function(
		params("message"),
		body(`return map[string]any{
	"$type": "Error",
	"message": message,
}`),
	)
}
