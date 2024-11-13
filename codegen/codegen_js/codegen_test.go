package codegen_js_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/codegen/codegen_js"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestGenerateProgramNonRunnable(t *testing.T) {
	program := `package main

import tenecs.go.Runtime
import tenecs.go.Main

app := Main(
  main = (runtime: Runtime) => {
    runtime.console.log("Hello world")
  }
)`

	expectedJs := `let main__app = tenecs_go__Main((main__runtime) => {
return main__runtime.console.log("Hello world")
})
function tenecs_go__Main(main) {
return ({
  "$type": "Main",
  "main": main,
})
return null
}
function tenecs_go__Runtime(console, http, ref) {
return ({
  "$type": "Runtime",
  "console": console,
  "http": http,
  "ref": ref,
})
return null
}

`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_js.GenerateProgramNonRunnable(typed)
	assert.Equal(t, expectedJs, generated)
}

func TestGenerateProgramTest(t *testing.T) {
	program := `package test

import tenecs.test.UnitTestKit
import tenecs.test.UnitTest

_ := UnitTest("and", (testkit: UnitTestKit): Void => {
  testkit.assert.equal("", "")
})`

	expectedJs := `let test__syntheticName_0 = tenecs_test__UnitTest("and", (test__testkit) => {
return test__testkit.assert.equal("", "")
})
function tenecs_test__UnitTest(name, theTest) {
return ({
  "$type": "UnitTest",
  "name": name,
  "theTest": theTest,
})
return null
}
function tenecs_test__UnitTestKit(assert, ref) {
return ({
  "$type": "UnitTestKit",
  "assert": assert,
  "ref": ref,
})
return null
}




let testSummary = {
  "total": 0,
  "ok": 0,
  "fail": 0,
}

function runUnitTests(implementingUnitTestSuite, implementingUnitTest) {
  let registry = createTestRegistry()
  if (implementingUnitTest.length > 0) {
    console.log("unit tests:")
  }
  for (const implementation of implementingUnitTest) {
    registry.test(implementation.name, implementation.theTest)
  }

  for (const implementation of implementingUnitTestSuite) {
    console.log(implementation.name + ":")
    implementation.tests(registry)
  }

  console.log("Ran a total of", testSummary.total, "tests")
  console.log("  *", testSummary.ok, "succeeded")
  console.log("  *", testSummary.fail, "failed")
}

function areDeeplyEqual(obj1, obj2) {
  if (obj1 === obj2) return true;

  if (Array.isArray(obj1) && Array.isArray(obj2)) {

    if(obj1.length !== obj2.length) return false;
    
    return obj1.every((elem, index) => {
      return areDeeplyEqual(elem, obj2[index]);
    })


  }

  if(typeof obj1 === "object" && typeof obj2 === "object" && obj1 !== null && obj2 !== null) {
    if(Array.isArray(obj1) || Array.isArray(obj2)) return false;
    
    const keys1 = Object.keys(obj1)
    const keys2 = Object.keys(obj2)

    if(keys1.length !== keys2.length || !keys1.every(key => keys2.includes(key))) return false;
    
    for(let key in obj1) {
       let isEqual = areDeeplyEqual(obj1[key], obj2[key])
       if (!isEqual) { return false; }
    }

    return true;
    
  }

  return false;
}

let testkit = {
  "assert": {
    "equal": (expected, value) => {
      if (!areDeeplyEqual(expected, value)) {
        throw new Error(testEqualityErrorMessage(expected, value))
      }
      return null
    },
    "fail": (message) => {
      throw new Error(message)
    }
  },
  "ref": ({
  "new": (value) => {
    let ref = value
    return ({
      "$type": "Ref",
      "get": () => {
        return ref
      },
      "set": (value) => {
        ref = value
        return null
      },
      "modify": (f) => {
        ref = f(ref)
        return null
      }
    })
  }
})
}

function createTestRegistry() {
  return ({
    "test": (name, theTest) => {
      try {
        theTest(testkit)
        console.log("  [\u001b[32mOK\u001b[0m]", name)
        testSummary.ok += 1
      } catch (e) {
        let errMsg = "could not print the failure"
        if (e.message) {
          errMsg = e.message
        }
        console.log("  [\u001b[31mFAILURE\u001b[0m]", name)
        console.log("    " + errMsg)
        testSummary.fail += 1
      }
      testSummary.total += 1
    }
  })
}

function testEqualityErrorMessage(value, expected) {
  return JSON.stringify(value) + " is not equal to " + JSON.stringify(expected)
}

runUnitTests([], [test__syntheticName_0])
`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_js.GenerateProgramTest(typed, codegen.FindTests(typed))
	assert.Equal(t, expectedJs, generated)
}

func TestGenerateHtmlPageForWebApp(t *testing.T) {
	program := `package mypage

import tenecs.web.WebApp
import tenecs.web.HtmlElement
import tenecs.web.HtmlElementProperty

struct State()
struct Event()

webApp := WebApp<State, Event>(
  init = () => State(),
  update = update,
  view = view
)

update := (model: State, event: Event): State => {
  model
}

view := (model: State): HtmlElement<Event> => {
  HtmlElement("p", [HtmlElementProperty<Event>](), "Hello world!")
}
`

	expectedHtml := `<html><body><script>function mypage__update(mypage__model, mypage__event) {
return mypage__model
}
function mypage__view(mypage__model) {
return tenecs_web__HtmlElement("p", [], "Hello world!")
}
let mypage__webApp = tenecs_web__WebApp(() => {
return mypage__State()
}, mypage__update, mypage__view)
function mypage__Event() {
return ({  "$type": "Event"})
}
function mypage__State() {
return ({  "$type": "State"})
}
function tenecs_web__HtmlElement(name, properties, children) {
return ({
  "$type": "HtmlElement",
  "name": name,
  "properties": properties,
  "children": children,
})
return null
}
function tenecs_web__HtmlElementProperty(name, value) {
return ({
  "$type": "HtmlElementProperty",
  "name": name,
  "value": value,
})
return null
}
function tenecs_web__WebApp(init, update, view) {
return ({
  "$type": "WebApp",
  "init": init,
  "update": update,
  "view": view,
})
return null
}


const webApp = mypage__webApp

function render(htmlElement) {
  let result = "<" + htmlElement.name
  for (const property of htmlElement.properties) {
    result += " " + property.name + "="
    if (typeof property.value == "string") {
      result += "\"" + property.value + "\"" 
    } else {
      alert("todo render function")
    }  
  }
  result += ">"
  if (typeof htmlElement.children == "string") {
    result += htmlElement.children
  } else {
    for(const child of htmlElement.children) {
      result += render(child)
    }
  }
  result += "</" + htmlElement.name + ">"
  return result
}

document.body.innerHTML = render(webApp.view(webApp.init()))
</script></body></html>`

	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)

	typed, err := typer.TypecheckSingleFile(*parsed)
	assert.NoError(t, err)

	generated := codegen_js.GenerateHtmlPageForWebApp(typed, "webApp")
	assert.Equal(t, expectedHtml, generated)
}
