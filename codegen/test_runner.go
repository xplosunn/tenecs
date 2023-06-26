package codegen

func GenerateTestRunner() ([]Import, string) {
	imports := []Import{"fmt", "reflect"}

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
				panic("equal was not equal")
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
				if err := recover(); err != nil {
					testSuccess = false
				}
				testResultString := "[OK]"
				if !testSuccess {
					testResultString = "[FAILURE]"
				}
				fmt.Printf("  %s %s\n", testResultString, testName)
			}()

			return testFunc(assert)
		},
	}
}
`

	return imports, result
}
