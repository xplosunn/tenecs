package testcode

const Test TestCodeCategory = "test"

var TestsUnit = Create(Test, "TestsUnit", `package test

import tenecs.test.UnitTestKit
import tenecs.test.UnitTestRegistry
import tenecs.test.UnitTestSuite

myUnitTestSuite := UnitTestSuite("My Unit Test Suite", (registry: UnitTestRegistry): Void => {
  registry.test("My test name", myTest)
})

myTest := (testkit: UnitTestKit): Void => {
  testkit.assert.equal<String>("a", "b")
}
`)
