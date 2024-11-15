package codegen_golang_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_golang"
	golang2 "github.com/xplosunn/tenecs/external/golang"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testcode"
	"github.com/xplosunn/tenecs/typer"
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
		"http": map[string]any{
			"serve": func(server any, port any) any {

				server.(map[string]any)["__hiddenServe"].(func(any) any)(port)

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

import tenecs.test.UnitTest
import tenecs.test.UnitTestSuite
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry

helloWorld := (): String => {
  "hello world!"
}

_ := UnitTestSuite(
  "My Tests",
  tests = (registry: UnitTestRegistry): Void => {
    registry.test("hello world function", testCaseHelloworld)
  }
)

_ := UnitTest("unitHello", testCaseHelloworld)

testCaseHelloworld := (testkit: UnitTestKit): Void => {
  result := helloWorld()
  expected := "hello world!"
  testkit.assert.equal<String>(result, expected)
}
`

	expectedGo := `package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var test__helloWorld any
var _ = func() any {
	test__helloWorld = func() any {
		return "hello world!"
	}
	return nil
}()

var test__syntheticName_1 any
var _ = func() any {
	test__syntheticName_1 = tenecs_test__UnitTestSuite.(func(any, any) any)("My Tests", func(_registry any) any {
		_registry.(map[string]any)["test"].(func(any, any) any)("hello world function", test__testCaseHelloworld)
		return nil
	})
	return nil
}()

var test__syntheticName_2 any
var _ = func() any {
	test__syntheticName_2 = tenecs_test__UnitTest.(func(any, any) any)("unitHello", test__testCaseHelloworld)
	return nil
}()

var test__testCaseHelloworld any
var _ = func() any {
	test__testCaseHelloworld = func(_testkit any) any {
		var _result any
		var _ = func() any {
			_result = test__helloWorld.(func() any)()
			return nil
		}()
		_ = _result

		var _expected any
		var _ = func() any {
			_expected = "hello world!"
			return nil
		}()
		_ = _expected

		_testkit.(map[string]any)["assert"].(map[string]any)["equal"].(func(any, any) any)(_result, _expected)
		return nil
	}
	return nil
}()

var tenecs_test__UnitTest any = func(name any, theTest any) any {
	return map[string]any{
		"$type":   "UnitTest",
		"name":    name,
		"theTest": theTest,
	}
	return nil
}
var tenecs_test__UnitTestKit any = func(assert any, ref any) any {
	return map[string]any{
		"$type":  "UnitTestKit",
		"assert": assert,
		"ref":    ref,
	}
	return nil
}
var tenecs_test__UnitTestRegistry any = func(test any) any {
	return map[string]any{
		"$type": "UnitTestRegistry",
		"test":  test,
	}
	return nil
}
var tenecs_test__UnitTestSuite any = func(name any, tests any) any {
	return map[string]any{
		"$type": "UnitTestSuite",
		"name":  name,
		"tests": tests,
	}
	return nil
}

func main() {
	runUnitTests([]any{test__syntheticName_1}, []any{test__syntheticName_2})
}

type testSummaryStruct struct {
	total int
	ok    int
	fail  int
}

var testSummary = testSummaryStruct{}

func runUnitTests(implementingUnitTestSuite []any, implementingUnitTest []any) {
	registry := createTestRegistry()

	if len(implementingUnitTest) > 0 {
		fmt.Printf("unit tests:\n")
	}
	for _, implementation := range implementingUnitTest {
		registry["test"].(func(any, any) any)(implementation.(map[string]any)["name"], implementation.(map[string]any)["theTest"])
	}

	for _, implementation := range implementingUnitTestSuite {
		fmt.Println(implementation.(map[string]any)["name"].(string) + ":")
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

	expectedRunResult := fmt.Sprintf(`unit tests:
  [%s] unitHello
My Tests:
  [%s] hello world function

Ran a total of 2 tests
  * 2 succeeded
  * 0 failed
`, codegen_golang.Green("OK"), codegen_golang.Green("OK"))

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramTest(typed, codegen.FindTests(typed))
	assert.Equal(t, expectedGo, golang2.Fmt(t, generated))

	output := golang2.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithStandardLibraryFunction(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main
import tenecs.string.join

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log(join("Hello ", "world!"))
  }
)`

	expectedGo := `package main

import (
	"fmt"
)

var main__app any
var _ = func() any {
	main__app = tenecs_go__Main.(func(any) any)(func(_runtime any) any {
		_runtime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(tenecs_string__join.(func(any, any) any)("Hello ", "world!"))
		return nil
	})
	return nil
}()

var tenecs_go__Main any = func(main any) any {
	return map[string]any{
		"$type": "Main",
		"main":  main,
	}
	return nil
}
var tenecs_go__Runtime any = func(console any, http any, ref any) any {
	return map[string]any{
		"$type":   "Runtime",
		"console": console,
		"http":    http,
		"ref":     ref,
	}
	return nil
}
var tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}

func main() {
	r := runtime()
	main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, expectedGo, golang2.Fmt(t, generated))

	output := golang2.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithImportAlias(t *testing.T) {
	program := testcode.ImportAliasMain

	expectedGo := `package main

import (
	"fmt"
)

var main__app any
var _ = func() any {
	main__app = tenecs_go__Main.(func(any) any)(func(_runtime any) any {
		_runtime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(tenecs_string__join.(func(any, any) any)("Hello ", "world!"))
		return nil
	})
	return nil
}()

var tenecs_go__Main any = func(main any) any {
	return map[string]any{
		"$type": "Main",
		"main":  main,
	}
	return nil
}
var tenecs_go__Runtime any = func(console any, http any, ref any) any {
	return map[string]any{
		"$type":   "Runtime",
		"console": console,
		"http":    http,
		"ref":     ref,
	}
	return nil
}
var tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}

func main() {
	r := runtime()
	main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, expectedGo, golang2.Fmt(t, generated))

	output := golang2.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithStruct(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

struct Post(title: String)

app := Main(
  main = (runtime: Runtime) => {
    post := Post("the title")
    runtime.console.log(post.title)
  }
)`

	expectedGo := `package main

import (
	"fmt"
)

var main__app any
var _ = func() any {
	main__app = tenecs_go__Main.(func(any) any)(func(_runtime any) any {
		var _post any
		var _ = func() any {
			_post = main__Post.(func(any) any)("the title")
			return nil
		}()
		_ = _post

		_runtime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(_post.(map[string]any)["title"])
		return nil
	})
	return nil
}()

var main__Post any = func(title any) any {
	return map[string]any{
		"$type": "Post",
		"title": title,
	}
}
var tenecs_go__Main any = func(main any) any {
	return map[string]any{
		"$type": "Main",
		"main":  main,
	}
	return nil
}
var tenecs_go__Runtime any = func(console any, http any, ref any) any {
	return map[string]any{
		"$type":   "Runtime",
		"console": console,
		"http":    http,
		"ref":     ref,
	}
	return nil
}

func main() {
	r := runtime()
	main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "the title\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, expectedGo, golang2.Fmt(t, generated))

	output := golang2.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMain(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log("Hello world!")
  }
)`

	expectedGo := `package main

import (
	"fmt"
)

var main__app any
var _ = func() any {
	main__app = tenecs_go__Main.(func(any) any)(func(_runtime any) any {
		_runtime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)("Hello world!")
		return nil
	})
	return nil
}()

var tenecs_go__Main any = func(main any) any {
	return map[string]any{
		"$type": "Main",
		"main":  main,
	}
	return nil
}
var tenecs_go__Runtime any = func(console any, http any, ref any) any {
	return map[string]any{
		"$type":   "Runtime",
		"console": console,
		"http":    http,
		"ref":     ref,
	}
	return nil
}

func main() {
	r := runtime()
	main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, expectedGo, golang2.Fmt(t, generated))

	output := golang2.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithRecursion(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main
import tenecs.int.times
import tenecs.int.minus
import tenecs.compare.eq
import tenecs.json.jsonInt

factorial := (i: Int): Int => {
  if eq<Int>(i, 0) {
    1
  } else {
    times(i, factorial(minus(i, 1)))
  }
}

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log(jsonInt().toJson(factorial(5)))
  }
)`

	expectedGo := `package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var main__app any
var _ = func() any {
	main__app = tenecs_go__Main.(func(any) any)(func(_runtime any) any {
		_runtime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(tenecs_json__jsonInt.(func() any)().(map[string]any)["toJson"].(func(any) any)(main__factorial.(func(any) any)(5)))
		return nil
	})
	return nil
}()

var main__factorial any
var _ = func() any {
	main__factorial = func(_i any) any {
		return func() any {
			if func() any { return tenecs_compare__eq.(func(any, any) any)(_i, 0) }().(bool) {
				return 1
			} else {
				return tenecs_int__times.(func(any, any) any)(_i, main__factorial.(func(any) any)(tenecs_int__minus.(func(any, any) any)(_i, 1)))
			}
		}()
	}
	return nil
}()

var tenecs_go__Main any = func(main any) any {
	return map[string]any{
		"$type": "Main",
		"main":  main,
	}
	return nil
}
var tenecs_go__Runtime any = func(console any, http any, ref any) any {
	return map[string]any{
		"$type":   "Runtime",
		"console": console,
		"http":    http,
		"ref":     ref,
	}
	return nil
}
var tenecs_compare__eq any = func(first any, second any) any {
	return reflect.DeepEqual(first, second)
	return nil
}
var tenecs_json__jsonInt any = func() any {
	return map[string]any{
		"$type": "JsonSchema",
		"fromJson": func(input any) any {
			jsonString := input.(string)
			var output float64
			err := json.Unmarshal([]byte(jsonString), &output)
			if err != nil || float64(int(output)) != output {
				return map[string]any{
					"$type":   "Error",
					"message": "Could not parse Int from " + jsonString,
				}
			}
			return int(output)
		},
		"toJson": func(input any) any {
			result, _ := json.Marshal(input)
			return string(result)
		},
	}
	return nil
}
var tenecs_int__minus any = func(a any, b any) any {
	return a.(int) - b.(int)
	return nil
}
var tenecs_int__times any = func(a any, b any) any {
	return a.(int) * b.(int)
	return nil
}

func main() {
	r := runtime()
	main__app.(map[string]any)["main"].(func(any) any)(r)
}

` + runtime

	expectedRunResult := "120\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, expectedGo, golang2.Fmt(t, generated))

	output := golang2.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateAndRunMainWithWhen(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main
import tenecs.json.jsonInt
import tenecs.string.join

struct Post(title: String)

struct BlogPost(title: String)

toString := (input: Int | String | Post | BlogPost): String => {
  when input {
    is i: Int => {
      jsonInt().toJson(i)
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

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log(toString("is it 10?"))
    runtime.console.log(toString(10))
    runtime.console.log(toString(Post("wee")))
    runtime.console.log(toString(BlogPost("wee2")))
  }
)`

	expectedGo := `package main

import (
	"encoding/json"
	"fmt"
)

var main__app any
var _ = func() any {
	main__app = tenecs_go__Main.(func(any) any)(func(_runtime any) any {
		_runtime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(main__toString.(func(any) any)("is it 10?"))
		_runtime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(main__toString.(func(any) any)(10))
		_runtime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(main__toString.(func(any) any)(main__Post.(func(any) any)("wee")))
		_runtime.(map[string]any)["console"].(map[string]any)["log"].(func(any) any)(main__toString.(func(any) any)(main__BlogPost.(func(any) any)("wee2")))
		return nil
	})
	return nil
}()

var main__toString any
var _ = func() any {
	main__toString = func(_input any) any {
		return func() any {
			var over any = _input
			if _, ok := over.(int); ok {
				_i := over
				return tenecs_json__jsonInt.(func() any)().(map[string]any)["toJson"].(func(any) any)(_i)
			}
			if _, ok := over.(string); ok {
				_s := over
				return _s
			}
			if value, okObj := over.(map[string]any); okObj && value["$type"] == "BlogPost" {
				_b := over
				return tenecs_string__join.(func(any, any) any)("blogpost:", _b.(map[string]any)["title"])
			}
			if value, okObj := over.(map[string]any); okObj && value["$type"] == "Post" {
				_p := over
				return tenecs_string__join.(func(any, any) any)("post:", _p.(map[string]any)["title"])
			}
			return nil
		}()
	}
	return nil
}()

var main__BlogPost any = func(title any) any {
	return map[string]any{
		"$type": "BlogPost",
		"title": title,
	}
}
var main__Post any = func(title any) any {
	return map[string]any{
		"$type": "Post",
		"title": title,
	}
}
var tenecs_go__Main any = func(main any) any {
	return map[string]any{
		"$type": "Main",
		"main":  main,
	}
	return nil
}
var tenecs_go__Runtime any = func(console any, http any, ref any) any {
	return map[string]any{
		"$type":   "Runtime",
		"console": console,
		"http":    http,
		"ref":     ref,
	}
	return nil
}
var tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}
var tenecs_json__jsonInt any = func() any {
	return map[string]any{
		"$type": "JsonSchema",
		"fromJson": func(input any) any {
			jsonString := input.(string)
			var output float64
			err := json.Unmarshal([]byte(jsonString), &output)
			if err != nil || float64(int(output)) != output {
				return map[string]any{
					"$type":   "Error",
					"message": "Could not parse Int from " + jsonString,
				}
			}
			return int(output)
		},
		"toJson": func(input any) any {
			result, _ := json.Marshal(input)
			return string(result)
		},
	}
	return nil
}

func main() {
	r := runtime()
	main__app.(map[string]any)["main"].(func(any) any)(r)
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

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, expectedGo, golang2.Fmt(t, generated))

	output := golang2.RunCodeUnlessCached(t, generated)
	assert.Equal(t, expectedRunResult, output)
}

func TestGenerateShortCircuitTwice(t *testing.T) {
	program := testcode.ShortCircuitTwice

	expectedGo := `package main

import ()

var main__stringOrInt any
var _ = func() any {
	main__stringOrInt = func() any {
		return 3
	}
	return nil
}()

var main__usage any
var _ = func() any {
	main__usage = func() any {
		return func() any {
			var over any = main__stringOrInt.(func() any)()
			if _, ok := over.(string); ok {
				_str := over
				return func() any {
					var over any = main__stringOrInt.(func() any)()
					if _, ok := over.(int); ok {
						_strAgain := over
						return _strAgain
					}
					_strAgain := over
					return tenecs_string__join.(func(any, any) any)(_str, _strAgain)
				}()
			}
			_str := over
			return _str
		}()
	}
	return nil
}()

var tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramNonRunnable(typed)
	assert.Equal(t, expectedGo, golang2.Fmt(t, generated))
}
