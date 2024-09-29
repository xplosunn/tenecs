package standard_library

func tenecs_test_UnitTestKit() Function {
	return function(
		params("assert", "runtime"),
		body(`return map[string]any{
	"$type": "UnitTestKit",
	"assert": assert,
	"runtime": runtime,
}`),
	)
}
func tenecs_test_UnitTestRegistry() Function {
	return function(
		params("tests"),
		body(`return map[string]any{
	"$type": "UnitTestRegistry",
	"tests": tests,
}`),
	)
}
func tenecs_test_UnitTests() Function {
	return function(
		params("tests"),
		body(`return map[string]any{
	"$type": "UnitTests",
	"tests": tests,
}`),
	)
}

func tenecs_test_Assert() Function {
	return function(
		params("equal", "fail"),
		body(`return map[string]any{
	"$type": "Assert",
	"equal": equal,
	"fail": fail,
}`),
	)
}
