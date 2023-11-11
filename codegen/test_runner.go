package codegen

import "fmt"

func GenerateTestRunner() ([]Import, string) {
	imports := []Import{"fmt", "reflect", "encoding/json"}

	runtimeImports, runtime := generateRuntime()
	imports = append(imports, runtimeImports...)

	result := `type testSummaryStruct struct {
total int
ok int
fail int
}

var testSummary = testSummaryStruct{}

func runTests(varNames []string, implementingUnitTests []any) {
	registry := createTestRegistry()

	for i, implementation := range implementingUnitTests {
		fmt.Println(varNames[i] + ":")
		implementation.(map[string]any)["tests"].(func(any) any)(registry)
	}

	fmt.Printf("\nRan a total of %d tests\n", testSummary.total)
	fmt.Printf("  * %d succeeded\n", testSummary.ok)
	fmt.Printf("  * %d failed\n", testSummary.fail)
}

func createTestRegistry() map[string]any {
	assert := map[string]any{
		"equal": func(expected any, value any) any {
			if !reflect.DeepEqual(value, expected) {
				panic(testEqualityErrorMessage(value, expected))
			}
			return nil
		},
		"fail": func(message any) any {
			panic(message)
		},
	}` + fmt.Sprintf(`

	testkit := map[string]any{
		"assert": assert,
		"runtime": %s,
	}`, runtime) + `

	return map[string]any{
		"test": func(name any, theTest any) any {
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
	toJson := func(input any) string {
		if inputMap, ok := input.(map[string]any); ok {
			copy := map[string]any{}
			for k, v := range inputMap {
				copy[k] = v
			}
			delete(copy, "$type")
			result, _ := json.Marshal(copy)
			return string(result)
		}
		result, _ := json.Marshal(input)
		return string(result)
	}
	return "expected " + toJson(expected) + " but got " + toJson(value)
}
`

	return imports, result
}

func generateRuntime() ([]Import, string) {
	imports := []Import{}

	imports = append(imports, "fmt")
	console := ofMap(map[string]string{
		"log": function(params("Pmessage"), body(``)),
	})

	execution := ofMap(map[string]string{
		"runBlocking": function(params("blockingOp"), body(`return blockingOp.(map[string]any)["fakeRun"].(func()any)()`)),
	})

	runtime := ofMap(map[string]string{
		"console":   console,
		"execution": execution,
		"ref":       runtimeRefCreator(),
	})

	return imports, runtime
}
