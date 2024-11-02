package standard_library

import "github.com/xplosunn/tenecs/typer/standard_library"

func tenecs_test_UnitTestKit() Function {
	return structFunction(standard_library.Tenecs_test_UnitTestKit)
}
func tenecs_test_UnitTestRegistry() Function {
	return structFunction(standard_library.Tenecs_test_UnitTestRegistry)
}
func tenecs_test_UnitTestSuite() Function {
	return structFunction(standard_library.Tenecs_test_UnitTestSuite)
}

func tenecs_test_Assert() Function {
	return structFunction(standard_library.Tenecs_test_Assert)
}
func tenecs_test_UnitTest() Function {
	return structFunction(standard_library.Tenecs_test_UnitTest)
}
