package codegen_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"github.com/xplosunn/tenecs/typer"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var runtime = `func runtime() map[string]any {
	return map[string]any{
		"console": map[string]any{
			"log": func(Pmessage any) any {
				fmt.Println(Pmessage)
				return nil
			},
		},
		"execution": map[string]any{
			"runBlocking": func(blockingOp any) any {
				return blockingOp.(map[string]any)["run"].(func() any)()
				return nil
			},
		},
		"ref": map[string]any{
			"new": func(Pvalue any) any {
				var ref any = Pvalue
				return map[string]any{
					"$type": "Ref",
					"get": func() any {
						return ref
					},
					"set": func(value any) any {
						ref = value
						return nil
					},
					"modify": func(f any) any {
						ref = f.(func(any) any)(ref)
						return nil
					},
				}

				return nil
			},
		},
	}
}
`

func TestGenerateAndRunTest(t *testing.T) {
	program := `package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry

helloWorld := (): String => {
  "hello world!"
}

myTests := implement UnitTests {
  public tests := (registry: UnitTestRegistry): Void => {
    registry.test("hello world function", testCaseHelloworld)
  }

  testCaseHelloworld := (testkit: UnitTestKit): Void => {
    result := helloWorld()
    expected := "hello world!"
    testkit.assert.equal<String>(result, expected)
  }
}`

	expectedGo := `package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var P__test__helloWorld any
var _ = func() any {
	P__test__helloWorld = func() any {
		return "hello world!"
	}
	return nil
}()

var P__test__myTests any
var _ = func() any {
	P__test__myTests = func() any {
		var PmyTests any = map[string]any{}
		var PtestCaseHelloworld any
		var Ptests any
		PtestCaseHelloworld = func(Ptestkit any) any {
			var Presult any
			var _ = func() any {
				Presult = P__test__helloWorld.(func() any)()
				return nil
			}()

			var Pexpected any
			var _ = func() any {
				Pexpected = "hello world!"
				return nil
			}()

			return Ptestkit.(map[string]any)["assert"].(map[string]any)["equal"].(func(any, any) any)(Presult, Pexpected)
		}
		PmyTests.(map[string]any)["testCaseHelloworld"] = PtestCaseHelloworld
		Ptests = func(Pregistry any) any {
			return Pregistry.(map[string]any)["test"].(func(any, any) any)("hello world function", PtestCaseHelloworld)
		}
		PmyTests.(map[string]any)["tests"] = Ptests
		return PmyTests
	}()
	return nil
}()

func main() {
	runTests([]string{"myTests"}, []any{P__test__myTests})
}

type testSummaryStruct struct {
	total int
	ok    int
	fail  int
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
	}

	testkit := map[string]any{
		"assert": assert,
		"runtime": map[string]any{
			"console": map[string]any{
				"log": func(Pmessage any) any {

					return nil
				},
			},
			"execution": map[string]any{
				"runBlocking": func(blockingOp any) any {
					return blockingOp.(map[string]any)["fakeRun"].(func() any)()
					return nil
				},
			},
			"ref": map[string]any{
				"new": func(Pvalue any) any {
					var ref any = Pvalue
					return map[string]any{
						"$type": "Ref",
						"get": func() any {
							return ref
						},
						"set": func(value any) any {
							ref = value
							return nil
						},
						"modify": func(f any) any {
							ref = f.(func(any) any)(ref)
							return nil
						},
					}

					return nil
				},
			},
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
	return toJson(expected) + " is not equal to " + toJson(value)
}
`

	expectedRunResult := fmt.Sprintf(`myTests:
  [%s] hello world function

Ran a total of 1 tests
  * 1 succeeded
  * 0 failed
`, codegen.Green("OK"))

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramTest(typed)
	assert.Equal(t, expectedGo, gofmt(t, generated))

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

var P__main__app any
var _ = func() any {
	P__main__app = func() any {
		var Papp any = map[string]any{}
		var Pmain any
		Pmain = func(Pruntime any) any {
			return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(P__tenecs_string__join.(func(any, any) any)("Hello ", "world!"))
		}
		Papp.(map[string]any)["main"] = Pmain
		return Papp
	}()
	return nil
}()

var P__tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}

func main() {
	r := runtime()
	P__main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramMain(typed, nil)
	assert.Equal(t, expectedGo, gofmt(t, generated))

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithImportAlias(t *testing.T) {
	program := testcode.ImportAliasMain

	expectedGo := `package main

import (
	"fmt"
)

var P__main__app any
var _ = func() any {
	P__main__app = func() any {
		var Papp any = map[string]any{}
		var Pmain any
		Pmain = func(Pruntime any) any {
			return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(P__tenecs_string__join.(func(any, any) any)("Hello ", "world!"))
		}
		Papp.(map[string]any)["main"] = Pmain
		return Papp
	}()
	return nil
}()

var P__tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}

func main() {
	r := runtime()
	P__main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramMain(typed, nil)
	assert.Equal(t, expectedGo, gofmt(t, generated))

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

var P__main__app any
var _ = func() any {
	P__main__app = func() any {
		var Papp any = map[string]any{}
		var Pmain any
		Pmain = func(Pruntime any) any {
			var Ppost any
			var _ = func() any {
				Ppost = P__main__Post.(func(any) any)("the title")
				return nil
			}()

			return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(Ppost.(map[string]any)["title"])
		}
		Papp.(map[string]any)["main"] = Pmain
		return Papp
	}()
	return nil
}()

var P__main__Post any = func(title any) any {
	return map[string]any{
		"$type": "Post",
		"title": title,
	}
}

func main() {
	r := runtime()
	P__main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "the title\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramMain(typed, nil)
	assert.Equal(t, expectedGo, gofmt(t, generated))

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

var P__main__app any
var _ = func() any {
	P__main__app = func() any {
		var Papp any = map[string]any{}
		var Pmain any
		Pmain = func(Pruntime any) any {
			return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)("Hello world!")
		}
		Papp.(map[string]any)["main"] = Pmain
		return Papp
	}()
	return nil
}()

func main() {
	r := runtime()
	P__main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramMain(typed, nil)
	assert.Equal(t, expectedGo, gofmt(t, generated))

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
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var P__main__app any
var _ = func() any {
	P__main__app = func() any {
		var Papp any = map[string]any{}
		var Pmain any
		Pmain = func(Pruntime any) any {
			return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(P__tenecs_json__toJson.(func(any) any)(P__main__factorial.(func(any) any)(5)))
		}
		Papp.(map[string]any)["main"] = Pmain
		return Papp
	}()
	return nil
}()

var P__main__factorial any
var _ = func() any {
	P__main__factorial = func(Pi any) any {
		return func() any {
			if P__tenecs_compare__eq.(func(any, any) any)(Pi, 0).(bool) {
				return 1
			} else {
				return P__tenecs_int__times.(func(any, any) any)(Pi, P__main__factorial.(func(any) any)(P__tenecs_int__minus.(func(any, any) any)(Pi, 1)))
			}
		}()
	}
	return nil
}()

var P__tenecs_compare__eq any = func(first any, second any) any {
	return reflect.DeepEqual(first, second)
	return nil
}
var P__tenecs_int__minus any = func(a any, b any) any {
	return a.(int) - b.(int)
	return nil
}
var P__tenecs_int__times any = func(a any, b any) any {
	return a.(int) * b.(int)
	return nil
}
var P__tenecs_json__toJson any = func(input any) any {
	var toJson func(any) any
	toJson = func(input any) any {
		if inputArray, ok := input.([]any); ok {
			result := []string{}
			for _, elem := range inputArray {
				result = append(result, toJson(elem).(string))
			}
			return "[" + strings.Join(result, ",") + "]"
		}
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
	return toJson(input)
	return nil
}

func main() {
	r := runtime()
	P__main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "120\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramMain(typed, nil)
	assert.Equal(t, expectedGo, gofmt(t, generated))

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}
func TestGenerateAndRunMainWithImportedStruct(t *testing.T) {
	program := `package main

import tenecs.os.Main
import tenecs.json.JsonError
import tenecs.json.toJson

app := implement Main {
	public main := (runtime) => {
		runtime.console.log(toJson(JsonError("fake")))
	}
}`

	expectedGo := `package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

var P__main__app any
var _ = func() any {
	P__main__app = func() any {
		var Papp any = map[string]any{}
		var Pmain any
		Pmain = func(Pruntime any) any {
			return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(P__tenecs_json__toJson.(func(any) any)(P__tenecs_json__JsonError.(func(any) any)("fake")))
		}
		Papp.(map[string]any)["main"] = Pmain
		return Papp
	}()
	return nil
}()

var P__tenecs_json__JsonError any = func(message any) any {
	return map[string]any{
		"$type":   "JsonError",
		"message": message,
	}
	return nil
}
var P__tenecs_json__toJson any = func(input any) any {
	var toJson func(any) any
	toJson = func(input any) any {
		if inputArray, ok := input.([]any); ok {
			result := []string{}
			for _, elem := range inputArray {
				result = append(result, toJson(elem).(string))
			}
			return "[" + strings.Join(result, ",") + "]"
		}
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
	return toJson(input)
	return nil
}

func main() {
	r := runtime()
	P__main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "{\"message\":\"fake\"}\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramMain(typed, nil)
	assert.Equal(t, expectedGo, gofmt(t, generated))

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithWhen(t *testing.T) {
	program := `package main

import tenecs.os.Runtime
import tenecs.os.Main
import tenecs.json.toJson
import tenecs.string.join

struct Post(title: String)

struct BlogPost(title: String)

toString := (input: Int | String | Post | BlogPost): String => {
  when input {
    is i: Int => {
      toJson<Int>(i)
    }
    is s: String => {
      s
    }
    is p: Post => {
      join("post:", p.title)
    }
    is b: BlogPost => {
      join("blogpost:", b.title)
    }
  }
}

app := implement Main {
  public main := (runtime: Runtime) => {
    runtime.console.log(toString("is it 10?"))
    runtime.console.log(toString(10))
    runtime.console.log(toString(Post("wee")))
    runtime.console.log(toString(BlogPost("wee2")))
  }
}`

	expectedGo := `package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

var P__main__app any
var _ = func() any {
	P__main__app = func() any {
		var Papp any = map[string]any{}
		var Pmain any
		Pmain = func(Pruntime any) any {
			Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(P__main__toString.(func(any) any)("is it 10?"))
			Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(P__main__toString.(func(any) any)(10))
			Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(P__main__toString.(func(any) any)(P__main__Post.(func(any) any)("wee")))
			return Pruntime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(P__main__toString.(func(any) any)(P__main__BlogPost.(func(any) any)("wee2")))
		}
		Papp.(map[string]any)["main"] = Pmain
		return Papp
	}()
	return nil
}()

var P__main__toString any
var _ = func() any {
	P__main__toString = func(Pinput any) any {
		return func() any {
			var over any = Pinput
			if _, ok := over.(int); ok {
				Pi := over
				return P__tenecs_json__toJson.(func(any) any)(Pi)
			}
			if _, ok := over.(string); ok {
				Ps := over
				return Ps
			}
			if value, okObj := over.(map[string]any); okObj && value["$type"] == "BlogPost" {
				Pb := over
				return P__tenecs_string__join.(func(any, any) any)("blogpost:", Pb.(map[string]any)["title"])
			}
			if value, okObj := over.(map[string]any); okObj && value["$type"] == "Post" {
				Pp := over
				return P__tenecs_string__join.(func(any, any) any)("post:", Pp.(map[string]any)["title"])
			}
			return nil
		}()
	}
	return nil
}()

var P__main__BlogPost any = func(title any) any {
	return map[string]any{
		"$type": "BlogPost",
		"title": title,
	}
}
var P__main__Post any = func(title any) any {
	return map[string]any{
		"$type": "Post",
		"title": title,
	}
}
var P__tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}
var P__tenecs_json__toJson any = func(input any) any {
	var toJson func(any) any
	toJson = func(input any) any {
		if inputArray, ok := input.([]any); ok {
			result := []string{}
			for _, elem := range inputArray {
				result = append(result, toJson(elem).(string))
			}
			return "[" + strings.Join(result, ",") + "]"
		}
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
	return toJson(input)
	return nil
}

func main() {
	r := runtime()
	P__main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := `is it 10?
10
post:wee
blogpost:wee2
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen.GenerateProgramMain(typed, nil)
	assert.Equal(t, expectedGo, gofmt(t, generated))

	output := createFileAndRun(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func gofmt(t *testing.T, fileContent string) string {
	dir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	filePath := filepath.Join(dir, t.Name()+".go")

	_, err = os.Create(filePath)

	contentBytes := []byte(fileContent)
	err = os.WriteFile(filePath, contentBytes, 0644)
	assert.NoError(t, err)

	cmd := exec.Command("gofmt", "-w", filePath)
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		t.Log(filePath)
	}
	assert.NoError(t, err)

	formatted, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	return string(formatted)
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
