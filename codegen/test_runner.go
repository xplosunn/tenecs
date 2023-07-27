package codegen

func GenerateTestRunner() ([]Import, string) {
	imports := []Import{"fmt", "reflect", "encoding/json"}

	result := `func runTests(varNames []string, implementingUnitTests []any) {
	registry := createTestRegistry()

	for i, module := range implementingUnitTests {
		fmt.Println(varNames[i] + ":")
		module.(map[string]any)["tests"].(func(any) any)(registry)
	}
}

func createTestRegistry() map[string]any {
	assert := map[string]any{
		"equal": func(value any, expected any) any {
			if !reflect.DeepEqual(value, expected) {
				panic(testEqualityErrorMessage(value, expected))
			}
			return nil
		},
	}

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
				testResultString := "[OK]"
				if !testSuccess {
					testResultString = "[FAILURE]"
				}
				fmt.Printf("  %s %s\n", testResultString, testName)
				if !testSuccess {
					fmt.Printf("    %s\n", errMsg)
				}
			}()

			return testFunc(assert)
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
