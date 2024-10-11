package testcode

const Test TestCodeCategory = "test"

var TestsUnit = Create(Test, "TestsUnit", `
package test

import tenecs.test.UnitTestSuite
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry

myUnitTestSuite := UnitTestSuite(
  "My Unit Test Suite",
  tests = (registry: UnitTestRegistry): Void => {
    registry.test("My test name", myTest)
  }
)

myTest := (testkit: UnitTestKit): Void => {
  testkit.assert.equal<String>("a", "b")
}

`)
