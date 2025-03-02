package codegen_golang

import "fmt"

func GenerateTestRunner() ([]Import, string) {
	imports, runtime := GenerateRuntime()

	imports = append(imports, "fmt", "reflect")

	ref := runtimeRefCreator()

	result := fmt.Sprintf(`
func runtime() tenecs_go_Runtime{
return %s
}

`, runtime) + `type testSummaryStruct struct {
total int
ok int
fail int
}

var testSummary = testSummaryStruct{}

func runTests(implementingUnitTestSuite []any, implementingUnitTest []any, implementingGoIntegrationTest []any) {
	registry := createTestRegistry()

	if len(implementingUnitTest) > 0 {
		fmt.Printf("unit tests:\n")
	}
	for _, implementation := range implementingUnitTest {
		registry._test.(func(any, any) any)(implementation.(tenecs_test_UnitTest)._name, implementation.(tenecs_test_UnitTest)._theTest)
	}

	for _, implementation := range implementingUnitTestSuite {
		fmt.Println(implementation.(tenecs_test_UnitTestSuite)._name.(string) + ":")
		implementation.(tenecs_test_UnitTestSuite)._tests.(func(any) any)(registry)
	}

	if len(implementingGoIntegrationTest) > 0 {
		fmt.Printf("integration tests:\n")
	}
	for _, implementation := range implementingGoIntegrationTest {
		r := runtime()
		testkit := createGoIntegrationTestKit()
		registry._test.(func(any, any) any)(implementation.(tenecs_test_GoIntegrationTest)._name, func (_ any) any {
			implementation.(tenecs_test_GoIntegrationTest)._theTest.(func(any,any) any)(testkit, r)
			return nil
		})
	}

	fmt.Printf("\nRan a total of %d tests\n", testSummary.total)
	fmt.Printf("  * %d succeeded\n", testSummary.ok)
	fmt.Printf("  * %d failed\n", testSummary.fail)
}

func createGoIntegrationTestKit() tenecs_test_GoIntegrationTestKit {
	assert := tenecs_test_Assert{
		_equal: func(expected any, value any) any {
			if !reflect.DeepEqual(value, expected) {
				panic(testEqualityErrorMessage(value, expected))
			}
			return nil
		},
		_fail: func(message any) any {
			panic(message)
		},
	}

	testkit := tenecs_test_GoIntegrationTestKit{
		_assert: assert,
	}
	return testkit
}

func createTestRegistry() tenecs_test_UnitTestRegistry {
	assert := tenecs_test_Assert{
		_equal: func(expected any, value any) any {
			if !reflect.DeepEqual(value, expected) {
				panic(testEqualityErrorMessage(value, expected))
			}
			return nil
		},
		_fail: func(message any) any {
			panic(message)
		},
	}` + fmt.Sprintf(`

	testkit := tenecs_test_UnitTestKit{
		_assert: assert,
		_ref: %s,
	}`, ref) + `

	return tenecs_test_UnitTestRegistry{
		_test: func(name any, theTest any) any {
			testName := name.(string)
			testFunc := theTest.(func(any) any)
			testSuccess := true
			defer func() {
				errMsg := "could not print the failure"
				if err := recover(); err != nil {
					testSuccess = false
					errMsg = err.(string)
				}
				testResultString := "[\u001b[32mOK\u001b[0m]"
				if !testSuccess {
					testResultString = "[\u001b[31mFAILURE\u001b[0m]"
					testSummary.fail += 1
				} else {
					testSummary.ok += 1
				}
				fmt.Printf("  %s %s\n", testResultString, testName)
				if !testSuccess {
					fmt.Printf("    %s\n", errMsg)
				}
				testSummary.total += 1
			}()

			return testFunc(testkit)
		},
	}
}

func testEqualityErrorMessage(value any, expected any) string {
	return fmt.Sprintf("%+v is not equal to %+v", expected, value)
}
`

	return imports, result
}
