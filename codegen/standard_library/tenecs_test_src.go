package standard_library

func tenecs_test_UnitTestKit() Function {
	return function(
		params("assert", "ref"),
		body(`return map[string]any{
	"$type": "UnitTestKit",
	"assert": assert,
	"ref": ref,
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
func tenecs_test_UnitTestSuite() Function {
	return function(
		params("name", "tests"),
		body(`return map[string]any{
	"$type": "UnitTestSuite",
	"name": name,
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
func tenecs_test_UnitTest() Function {
	return function(
		params("name", "theTest"),
		body(`return map[string]any{
	"$type": "UnitTest",
	"name": name,
	"theTest": theTest,
}`),
	)
}
