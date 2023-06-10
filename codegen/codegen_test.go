package codegen_test

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

func TestGenerateAndRunMain(t *testing.T) {
	program := `package main

import tenecs.os.Runtime
import tenecs.os.Main

app := implement Main {
	public main := (runtime: Runtime) => {
		runtime.console.log("Hello world!")
	}
}`

	expectedGo := `package main

import (
	"fmt"
)

var Papp any = map[string]any{
"main": func (Pruntime any) any {
Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any)any)("Hello world!")
return nil
},
}


func main() {
r := runtime()
Papp.(map[string]any)["main"].(func(any)any)(r)
}

func runtime() map[string]any {
return map[string]any{
"console": map[string]any{
"log": func (Pmessage any) any {
fmt.Println(Pmessage)
return nil
},
},
}
}
`

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)

	generated := codegen.Generate(typed)
	assert.Equal(t, expectedGo, generated)

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
