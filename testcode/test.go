package testcode

const Test TestCodeCategory = "test"

var TestsUnit = Create(Test, "TestsUnit", `
package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry

myUnitTests := UnitTests( 
  tests = (registry: UnitTestRegistry): Void => {
    registry.test("My test name", myTest)
  }
)

myTest := (testkit: UnitTestKit): Void => {
  testkit.assert.equal<String>("a", "b")
}

`)
