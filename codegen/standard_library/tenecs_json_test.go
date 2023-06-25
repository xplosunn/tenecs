package standard_library_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestToJson(t *testing.T) {
	program := `package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestRegistry
import tenecs.test.Assert
import tenecs.json.toJson

struct Post(title: String)

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("toJson", (assert: Assert): Void => {
      assert.equal<String>("42", toJson<Int>(42))
      assert.equal<String>("true", toJson<Boolean>(true))
      assert.equal<String>("\"rawr\"", toJson<String>("rawr"))
      assert.equal<String>("{\"title\":\"the title\"}", toJson<Post>(Post("the title")))
    })
  }
}`
	expectedRunResult := `myTests:
  [OK] toJson
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)

	generated := codegen.Generate(true, typed)

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func createFileAndRun(t *testing.T, fileContent string) string {
	dir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	filePath := filepath.Join(dir, t.Name()+".go")

	_, err = os.Create(filePath)

	contentBytes := []byte(fileContent)
	err = os.WriteFile(filePath, contentBytes, 0644)
	assert.NoError(t, err)

	cmd := exec.Command("go", "run", filePath)
	cmd.Dir = dir
	outputBytes, err := cmd.Output()
	t.Log(dir)
	assert.NoError(t, err)
	return string(outputBytes)
}
