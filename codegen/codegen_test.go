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

func TestGenerateAndRunTest(t *testing.T) {
	program := `package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestRegistry
import tenecs.test.Assert

helloWorld := (): String => {
  "hello world!"
}

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("hello world function", testCaseHelloworld)
  }

  testCaseHelloworld := (assert: Assert): Void => {
    result := helloWorld()
    expected := "hello world!"
    assert.equal<String>(result, expected)
  }
}`

	expectedGo := `package main

import (
	"fmt"
	"reflect"
)

var PhelloWorld any
var _ = func() any {
PhelloWorld = func () any {
return "hello world!"
}
return nil
}()

var PmyTests any
var _ = func() any {
PmyTests = func() any {
var PmyTests any = map[string]any{}
var PtestCaseHelloworld any
var Ptests any
PtestCaseHelloworld = func (Passert any) any {
var Presult any
var _ = func() any {
Presult = PhelloWorld.(func()any)()
return nil
}()

var Pexpected any
var _ = func() any {
Pexpected = "hello world!"
return nil
}()

return Passert.(map[string]any)["equal"].(func(any,any)any)(Presult, Pexpected)
}
PmyTests.(map[string]any)["testCaseHelloworld"] = PtestCaseHelloworld
Ptests = func (Pregistry any) any {
return Pregistry.(map[string]any)["test"].(func(any,any)any)("hello world function", PtestCaseHelloworld)
}
PmyTests.(map[string]any)["tests"] = Ptests
return PmyTests
}()
return nil
}()


func main() {
runTests([]string{"myTests"}, []any{PmyTests})
}

func runTests(varNames []string, implementingUnitTests []any) {
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

	expectedRunResult := `myTests:
  [OK] hello world function
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)

	generated := codegen.Generate(true, typed)
	assert.Equal(t, expectedGo, generated)

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithStandardLibraryFunction(t *testing.T) {
	program := `package main

import tenecs.os.Runtime
import tenecs.os.Main
import tenecs.string.join

app := implement Main {
	public main := (runtime: Runtime) => {
		runtime.console.log(join("Hello ", "world!"))
	}
}`

	expectedGo := `package main

import (
	"fmt"
)

var Papp any
var _ = func() any {
Papp = func() any {
var Papp any = map[string]any{}
var Pmain any
Pmain = func (Pruntime any) any {
return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any)any)(Pjoin.(func(any,any)any)("Hello ", "world!"))
}
Papp.(map[string]any)["main"] = Pmain
return Papp
}()
return nil
}()

var Pjoin any = func (Pleft any, Pright any) any {
return Pleft.(string) + Pright.(string)
return nil
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

	generated := codegen.Generate(false, typed)
	assert.Equal(t, expectedGo, generated)

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithStruct(t *testing.T) {
	program := `package main

import tenecs.os.Runtime
import tenecs.os.Main

struct Post(title: String)

app := implement Main {
	public main := (runtime: Runtime) => {
        post := Post("the title")
		runtime.console.log(post.title)
	}
}`

	expectedGo := `package main

import (
	"fmt"
)

var Papp any
var _ = func() any {
Papp = func() any {
var Papp any = map[string]any{}
var Pmain any
Pmain = func (Pruntime any) any {
var Ppost any
var _ = func() any {
Ppost = PPost.(func(any)any)("the title")
return nil
}()

return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any)any)(Ppost.(map[string]any)["title"])
}
Papp.(map[string]any)["main"] = Pmain
return Papp
}()
return nil
}()

var PPost any = func (title any) any {
return map[string]any{
"title": title,
}
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

	expectedRunResult := "the title\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)

	generated := codegen.Generate(false, typed)
	assert.Equal(t, expectedGo, generated)

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

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

var Papp any
var _ = func() any {
Papp = func() any {
var Papp any = map[string]any{}
var Pmain any
Pmain = func (Pruntime any) any {
return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any)any)("Hello world!")
}
Papp.(map[string]any)["main"] = Pmain
return Papp
}()
return nil
}()


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

	generated := codegen.Generate(false, typed)
	assert.Equal(t, expectedGo, generated)

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithRecursion(t *testing.T) {
	program := `package main

import tenecs.os.Runtime
import tenecs.os.Main
import tenecs.int.times
import tenecs.int.minus
import tenecs.compare.eq
import tenecs.json.toJson

factorial := (i: Int): Int => {
  if eq<Int>(i, 0) {
    1
  } else {
    times(i, factorial(minus(i, 1)))
  }
}

app := implement Main {
	public main := (runtime: Runtime) => {
		runtime.console.log(toJson<Int>(factorial(5)))
	}
}`

	expectedGo := `package main

import (
	"reflect"
	"encoding/json"
	"fmt"
)

var Papp any
var _ = func() any {
Papp = func() any {
var Papp any = map[string]any{}
var Pmain any
Pmain = func (Pruntime any) any {
return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any)any)(PtoJson.(func(any)any)(Pfactorial.(func(any)any)(5)))
}
Papp.(map[string]any)["main"] = Pmain
return Papp
}()
return nil
}()

var Pfactorial any
var _ = func() any {
Pfactorial = func (Pi any) any {
return func() any {
if Peq.(func(any,any)any)(Pi, 0).(bool) {
return 1
} else {
return Ptimes.(func(any,any)any)(Pi, Pfactorial.(func(any)any)(Pminus.(func(any,any)any)(Pi, 1)))
}
}()
}
return nil
}()

var Ptimes any = func (a any, b any) any {
return a.(int) * b.(int)
return nil
}
var Pminus any = func (a any, b any) any {
return a.(int) - b.(int)
return nil
}
var Peq any = func (first any, second any) any {
return reflect.DeepEqual(first, second)
return nil
}
var PtoJson any = func (input any) any {
result, _ := json.Marshal(input)
return string(result)
return nil
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

	expectedRunResult := "120\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)

	generated := codegen.Generate(false, typed)
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
