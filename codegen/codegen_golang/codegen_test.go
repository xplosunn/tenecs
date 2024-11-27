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

var runtime = `func runtime() tenecs_go_Runtime {
	return tenecs_go_Runtime{
		console: tenecs_go_Console{
			log: func(Pmessage any) any {
				fmt.Println(Pmessage)
				return nil
			},
		},
		ref: tenecs_ref_RefCreator{
			new: func(Pvalue any) any {
				var ref any = Pvalue
				return tenecs_ref_Ref{
					get: func() any {
						return ref
					},
					set: func(value any) any {
						ref = value
						return nil
					},
					modify: func(f any) any {
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
		_registry.(tenecs_test_UnitTestRegistry).test.(func(any, any) any)("hello world function", test__testCaseHelloworld)
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

		_testkit.(tenecs_test_UnitTestKit).assert.(tenecs_test_Assert).equal.(func(any, any) any)(_result, _expected)
		return nil
	}
	return nil
}()

var tenecs_test__UnitTest any = func(name any, theTest any) any {
	return tenecs_test_UnitTest{
		name,
		theTest,
	}
}
var tenecs_test__UnitTestKit any = func(assert any, ref any) any {
	return tenecs_test_UnitTestKit{
		assert,
		ref,
	}
}
var tenecs_test__UnitTestRegistry any = func(test any) any {
	return tenecs_test_UnitTestRegistry{
		test,
	}
}
var tenecs_test__UnitTestSuite any = func(name any, tests any) any {
	return tenecs_test_UnitTestSuite{
		name,
		tests,
	}
}
` + codegen_golang.GenerateStdLibStructs() + `

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
		registry.test.(func(any, any) any)(implementation.(tenecs_test_UnitTest).name, implementation.(tenecs_test_UnitTest).theTest)
	}

	for _, implementation := range implementingUnitTestSuite {
		fmt.Println(implementation.(tenecs_test_UnitTestSuite).name.(string) + ":")
		implementation.(tenecs_test_UnitTestSuite).tests.(func(any) any)(registry)
	}

	fmt.Printf("\nRan a total of %d tests\n", testSummary.total)
	fmt.Printf("  * %d succeeded\n", testSummary.ok)
	fmt.Printf("  * %d failed\n", testSummary.fail)
}

func createTestRegistry() tenecs_test_UnitTestRegistry {
	assert := tenecs_test_Assert{
		equal: func(expected any, value any) any {
			if !reflect.DeepEqual(value, expected) {
				panic(testEqualityErrorMessage(value, expected))
			}
			return nil
		},
		fail: func(message any) any {
			panic(message)
		},
	}

	testkit := tenecs_test_UnitTestKit{
		assert: assert,
		ref: tenecs_ref_RefCreator{
			new: func(Pvalue any) any {
				var ref any = Pvalue
				return tenecs_ref_Ref{
					get: func() any {
						return ref
					},
					set: func(value any) any {
						ref = value
						return nil
					},
					modify: func(f any) any {
						ref = f.(func(any) any)(ref)
						return nil
					},
				}

				return nil
			},
		},
	}

	return tenecs_test_UnitTestRegistry{
		test: func(name any, theTest any) any {
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
	assert.Equal(t, golang2.Fmt(t, expectedGo), golang2.Fmt(t, generated))

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
		_runtime.(tenecs_go_Runtime).console.(tenecs_go_Console).log.(func(any) any)(tenecs_string__join.(func(any, any) any)("Hello ", "world!"))
		return nil
	})
	return nil
}()

var tenecs_go__Main any = func(main any) any {
	return tenecs_go_Main{
		main,
	}
}
var tenecs_go__Runtime any = func(console any, ref any) any {
	return tenecs_go_Runtime{
		console,
		ref,
	}
}
var tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}

` + codegen_golang.GenerateStdLibStructs() + `
func main() {
	r := runtime()
	main__app.(tenecs_go_Main).main.(func(any) any)(r)
}

` + runtime

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, golang2.Fmt(t, expectedGo), golang2.Fmt(t, generated))

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
		_runtime.(tenecs_go_Runtime).console.(tenecs_go_Console).log.(func(any) any)(tenecs_string__join.(func(any, any) any)("Hello ", "world!"))
		return nil
	})
	return nil
}()

var tenecs_go__Main any = func(main any) any {
	return tenecs_go_Main{
		main,
	}
}
var tenecs_go__Runtime any = func(console any, ref any) any {
	return tenecs_go_Runtime{
		console,
		ref,
	}
}
var tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}

` + codegen_golang.GenerateStdLibStructs() + `
func main() {
	r := runtime()
	main__app.(tenecs_go_Main).main.(func(any) any)(r)
}

` + runtime

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, golang2.Fmt(t, expectedGo), golang2.Fmt(t, generated))

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

		_runtime.(tenecs_go_Runtime).console.(tenecs_go_Console).log.(func(any) any)(_post.(main_Post).title)
		return nil
	})
	return nil
}()

type main_Post struct {
	title any
}

var main__Post any = func(title any) any {
	return main_Post{
		title,
	}
}
var tenecs_go__Main any = func(main any) any {
	return tenecs_go_Main{
		main,
	}
}
var tenecs_go__Runtime any = func(console any, ref any) any {
	return tenecs_go_Runtime{
		console,
		ref,
	}
}

` + codegen_golang.GenerateStdLibStructs() + `
func main() {
	r := runtime()
	main__app.(tenecs_go_Main).main.(func(any) any)(r)
}

` + runtime

	expectedRunResult := "the title\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, golang2.Fmt(t, expectedGo), golang2.Fmt(t, generated))

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
		_runtime.(tenecs_go_Runtime).console.(tenecs_go_Console).log.(func(any) any)("Hello world!")
		return nil
	})
	return nil
}()

var tenecs_go__Main any = func(main any) any {
	return tenecs_go_Main{
		main,
	}
}
var tenecs_go__Runtime any = func(console any, ref any) any {
	return tenecs_go_Runtime{
		console,
		ref,
	}
}
` + codegen_golang.GenerateStdLibStructs() + `
func main() {
	r := runtime()
	main__app.(tenecs_go_Main).main.(func(any) any)(r)
}

` + runtime

	expectedRunResult := "Hello world!\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, golang2.Fmt(t, expectedGo), golang2.Fmt(t, generated))

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
		_runtime.(tenecs_go_Runtime).console.(tenecs_go_Console).log.(func(any) any)(tenecs_json__jsonInt.(func() any)().(tenecs_json_JsonSchema).toJson.(func(any) any)(main__factorial.(func(any) any)(5)))
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
	return tenecs_go_Main{
		main,
	}
}
var tenecs_go__Runtime any = func(console any, ref any) any {
	return tenecs_go_Runtime{
		console,
		ref,
	}
}
var tenecs_compare__eq any = func(first any, second any) any {
	return reflect.DeepEqual(first, second)
	return nil
}
var tenecs_json__jsonInt any = func() any {
	return tenecs_json_JsonSchema{
		fromJson: func(input any) any {
			jsonString := input.(string)
			var output float64
			err := json.Unmarshal([]byte(jsonString), &output)
			if err != nil || float64(int(output)) != output {
				return tenecs_error_Error{
					message: "Could not parse Int from " + jsonString,
				}
			}
			return int(output)
		},
		toJson: func(input any) any {
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

` + codegen_golang.GenerateStdLibStructs() + `
func main() {
	r := runtime()
	main__app.(tenecs_go_Main).main.(func(any) any)(r)
}

` + runtime

	expectedRunResult := "120\n"

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramMain(typed, "app")
	assert.Equal(t, golang2.Fmt(t, expectedGo), golang2.Fmt(t, generated))

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
		_runtime.(tenecs_go_Runtime).console.(tenecs_go_Console).log.(func(any) any)(main__toString.(func(any) any)("is it 10?"))
		_runtime.(tenecs_go_Runtime).console.(tenecs_go_Console).log.(func(any) any)(main__toString.(func(any) any)(10))
		_runtime.(tenecs_go_Runtime).console.(tenecs_go_Console).log.(func(any) any)(main__toString.(func(any) any)(main__Post.(func(any) any)("wee")))
		_runtime.(tenecs_go_Runtime).console.(tenecs_go_Console).log.(func(any) any)(main__toString.(func(any) any)(main__BlogPost.(func(any) any)("wee2")))
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
				return tenecs_json__jsonInt.(func() any)().(tenecs_json_JsonSchema).toJson.(func(any) any)(_i)
			}
			if _, ok := over.(string); ok {
				_s := over
				return _s
			}
			if _, okObj := over.(main_BlogPost); okObj {
				_b := over
				return tenecs_string__join.(func(any, any) any)("blogpost:", _b.(main_BlogPost).title)
			}
			if _, okObj := over.(main_Post); okObj {
				_p := over
				return tenecs_string__join.(func(any, any) any)("post:", _p.(main_Post).title)
			}
			return nil
		}()
	}
	return nil
}()

type main_BlogPost struct {
	title any
}

var main__BlogPost any = func(title any) any {
	return main_BlogPost{
		title,
	}
}

type main_Post struct {
	title any
}

var main__Post any = func(title any) any {
	return main_Post{
		title,
	}
}
var tenecs_go__Main any = func(main any) any {
	return tenecs_go_Main{
		main,
	}
}
var tenecs_go__Runtime any = func(console any, ref any) any {
	return tenecs_go_Runtime{
		console,
		ref,
	}
}
var tenecs_string__join any = func(Pleft any, Pright any) any {
	return Pleft.(string) + Pright.(string)
	return nil
}
var tenecs_json__jsonInt any = func() any {
	return tenecs_json_JsonSchema{
		fromJson: func(input any) any {
			jsonString := input.(string)
			var output float64
			err := json.Unmarshal([]byte(jsonString), &output)
			if err != nil || float64(int(output)) != output {
				return tenecs_error_Error{
					message: "Could not parse Int from " + jsonString,
				}
			}
			return int(output)
		},
		toJson: func(input any) any {
			result, _ := json.Marshal(input)
			return string(result)
		},
	}
	return nil
}
` + codegen_golang.GenerateStdLibStructs() + `
func main() {
	r := runtime()
	main__app.(tenecs_go_Main).main.(func(any) any)(r)
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
	assert.Equal(t, golang2.Fmt(t, expectedGo), golang2.Fmt(t, generated))

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
` + codegen_golang.GenerateStdLibStructs()

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_golang.GenerateProgramNonRunnable(typed)
	assert.Equal(t, golang2.Fmt(t, expectedGo), golang2.Fmt(t, generated))
}
