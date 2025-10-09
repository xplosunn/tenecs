package embed

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// embed start

type testSummaryStruct struct {
	runTotal              int
	runOk                 int
	runFail               int
	cachedUnitTestOk      int
	cachedUnitTestSuiteOk int
}

var testSummary = testSummaryStruct{}

func runTests(implementingUnitTestSuite []any, implementingUnitTest []any, implementingGoIntegrationTest []any) {
	registry := createTestRegistry()

	if len(implementingUnitTest) > 0 {
		fmt.Printf("unit tests:\n")
	}
	for _, implementation := range implementingUnitTest {
		registry.(map[string]any)["_test"].(func(any, any, any) any)([]any{}, implementation.(map[string]any)["_name"], implementation.(map[string]any)["_theTest"])
	}

	for _, implementation := range implementingUnitTestSuite {
		fmt.Println(implementation.(map[string]any)["_name"].(map[string]any)["value"].(string) + ":")
		implementation.(map[string]any)["_tests"].(func(any, any) any)([]any{}, registry)
	}

	if len(implementingGoIntegrationTest) > 0 {
		fmt.Printf("integration tests:\n")
	}
	for _, implementation := range implementingGoIntegrationTest {
		r := runtime()
		testkit := createGoIntegrationTestKit()
		registry.(map[string]any)["_test"].(func(any, any, any) any)([]any{}, implementation.(map[string]any)["_name"], func(_ any, _ any) any {
			implementation.(map[string]any)["_theTest"].(func(any, any, any) any)([]any{}, testkit, r)
			return map[string]any{"$type": "Void", "value": nil}
		})
	}

	fmt.Printf("\nRan a total of %d tests\n", testSummary.runTotal)
	fmt.Printf("  * %d succeeded\n", testSummary.runOk)
	fmt.Printf("  * %d failed\n", testSummary.runFail)
	if testSummary.cachedUnitTestOk > 0 {
		fmt.Printf("Skipped %d successful unit tests cached\n", testSummary.cachedUnitTestOk)
	}
	if testSummary.cachedUnitTestSuiteOk > 0 {
		fmt.Printf("Skipped %d successful unit test suites cached\n", testSummary.cachedUnitTestSuiteOk)
	}
}

func createGoIntegrationTestKit() any {
	assert := map[string]any{
		"_equal": func(codepoint string, generics any, expected any, value any) any {
			if !reflect.DeepEqual(value, expected) {
				panic(testEqualityErrorMessage(codepoint, generics, value, expected))
			}
			return map[string]any{"$type": "Void", "value": nil}
		},
		"_fail": func(generics any, message any) any {
			panic("@codepointtodo" + ": " + message.(map[string]any)["value"].(string))
		},
	}

	testkit := map[string]any{
		"_assert": assert,
	}
	return testkit
}

func createTestRegistry() any {
	assert := map[string]any{
		"_equal": func(codepoint string, generics any, expected any, value any) any {
			if !reflect.DeepEqual(value, expected) {
				panic(testEqualityErrorMessage(codepoint, generics, value, expected))
			}
			return map[string]any{"$type": "Void", "value": nil}
		},
		"_fail": func(generics any, message any) any {
			panic("@codepointtodo" + ": " + message.(map[string]any)["value"].(string))
		},
	}

	refCreator := map[string]any{
		"_new": func(generics any, value any) any {
			var ref any = value
			return map[string]any{
				"_get": func(generics any) any {
					return ref
				},
				"_set": func(generics any, value any) any {
					ref = value
					return map[string]any{"$type": "Void", "value": nil}
				},
				"_modify": func(generics any, f any) any {
					ref = f.(func(any, any) any)([]any{}, ref)
					return map[string]any{"$type": "Void", "value": nil}
				},
			}
		},
	}

	testkit := map[string]any{
		"_assert": assert,
		"_ref":    refCreator,
	}

	return map[string]any{
		"_test": func(generics any, name any, theTest any) any {
			testName := name.(map[string]any)["value"].(string)
			testFunc := theTest.(func(any, any) any)
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
					testSummary.runFail += 1
				} else {
					testSummary.runOk += 1
				}
				fmt.Printf("  %s %s\n", testResultString, testName)
				if !testSuccess {
					fmt.Printf("    %s\n", errMsg)
				}
				testSummary.runTotal += 1
			}()

			return testFunc([]any{}, testkit)
		},
	}
}

func testEqualityErrorMessage(codepoint string, generics any, value any, expected any) string {
	valueIsObject, valueIsObjectOk := value.(map[string]any)
	expectedIsObject, expectedIsObjectOk := expected.(map[string]any)
	if valueIsObjectOk && expectedIsObjectOk {
		valueJson, err1 := json.Marshal(valueIsObject)
		expectedJson, err2 := json.Marshal(expectedIsObject)
		if err1 == nil && err2 == nil {
			return fmt.Sprintf("@%s: %+v is not equal to %+v", codepoint, string(expectedJson), string(valueJson))
		}
	}
	return fmt.Sprintf("@%s: %+v is not equal to %+v", codepoint, expected, value)
}

func runtime() map[string]any {
	log := func(generics any, msg any) any {
		println(msg.(map[string]any)["value"].(string))
		return map[string]any{"$type": "Void", "value": nil}
	}
	console := map[string]any{
		"_log": log,
	}
	refCreator := map[string]any{
		"_new": func(generics any, value any) any {
			var ref any = value
			return map[string]any{
				"_get": func(generics any) any {
					return ref
				},
				"_set": func(generics any, value any) any {
					ref = value
					return map[string]any{"$type": "Void", "value": nil}
				},
				"_modify": func(generics any, f any) any {
					ref = f.(func(any, any) any)([]any{}, ref)
					return map[string]any{"$type": "Void", "value": nil}
				},
			}
		},
	}
	runtime := map[string]any{
		"_console": console,
		"_ref":     refCreator,
	}
	return runtime
}
