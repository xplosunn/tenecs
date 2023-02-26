package testcode

const Test TestCodeCategory = "test"

var TestsUnit = Create(Test, "TestsUnit", `
package test

import tenecs.test.UnitTests
import tenecs.test.UnitTestRegistry
import tenecs.test.Assert

myUnitTests := (): UnitTests => implement UnitTests { 
	public tests := (registry: UnitTestRegistry): Void => {
		registry.test("My test name", (assert: Assert): Void => {
			assert.equal<String>("a", "b")
		})
	}

	myTest := (assert: Assert): Void => {
		assert.equal<String>("a", "b")
	}
}
`)
