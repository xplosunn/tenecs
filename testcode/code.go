package testcode

type TestCodeCategory string

type TestCode struct {
	Name    string
	Content string
}

var testCodes = map[TestCodeCategory][]TestCode{}

func GetAll() []TestCode {
	result := []TestCode{}
	for _, cases := range testCodes {
		result = append(result, cases...)
	}
	return result
}

func Create(category TestCodeCategory, name string, content string) string {
	testCategory := testCodes[category]
	if testCategory == nil {
		testCategory = []TestCode{}
	}
	testCode := TestCode{
		Name:    name,
		Content: content,
	}
	testCategory = append(testCategory, testCode)
	testCodes[category] = testCategory
	return testCode.Content
}
