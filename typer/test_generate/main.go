package main

import (
	"fmt"
	"github.com/xplosunn/tenecs/testcode"
	"log"
	"os"
	"sort"
)

func main() {
	fmt.Println("Starting typer test from testcode generation")
	allTests := testcode.GetAll()

	filePath := fmt.Sprintf("../test/%s.go", "testcode_test")

	fileContent := `package parser_typer_test

// ###############################################
// # This file is generated via code-generation. #
// # Check gen_test.go                           #
// ###############################################

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

`

	sort.Slice(allTests, func(i, j int) bool {
		return allTests[i].Name < allTests[j].Name
	})
	for _, test := range allTests {
		fileContent += fmt.Sprintf(`func Test%s(t *testing.T) {
	validProgram(t, testcode.%s)
}

`, test.Name, test.Name)
	}

	if fileExists(filePath) {
		err := os.Remove(filePath)
		if err != nil {
			fail(err)
		}
	}
	err := os.WriteFile(filePath, []byte(fileContent), os.ModePerm)
	if err != nil {
		fail(err)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func fail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
