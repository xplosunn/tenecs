package codegen_js

import "fmt"

func generateNodeTestRunner() string {
	ref := runtimeRefCreator()

	result := fmt.Sprintf(`

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
  "ref": %s
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
`, ref)

	return result
}
