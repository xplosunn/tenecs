function runTests(implementingUnitTestSuite, implementingUnitTest) {
	const testSummary = {
		runTotal: 0,
		runOk: 0,
		runFail: 0,
		cachedUnitTestOk: 0,
		cachedUnitTestSuiteOk: 0,
	}

	const registry = createTestRegistry(testSummary)

	if (implementingUnitTest.length > 0) {
		process.stdout.write("unit tests:\n")
	}
	for (const implementation of implementingUnitTest) {
		registry["_test"]([], implementation["_name"], implementation["_theTest"])
	}

	for (const implementation of implementingUnitTestSuite) {
		process.stdout.write(implementation["_name"]["value"] + ":\n")
		implementation["_tests"]([], registry)
	}

	process.stdout.write(`\nRan a total of ${testSummary.runTotal} tests\n`)
	process.stdout.write(`  * ${testSummary.runOk} succeeded\n`)
	process.stdout.write(`  * ${testSummary.runFail} failed\n`)
	if (testSummary.cachedUnitTestOk > 0) {
		process.stdout.write(`Skipped ${testSummary.cachedUnitTestOk} successful unit tests cached\n`)
	}
	if (testSummary.cachedUnitTestSuiteOk > 0) {
		process.stdout.write(`Skipped ${testSummary.cachedUnitTestSuiteOk} successful unit test suites cached\n`)
	}
}

function createTestRegistry(testSummary) {
	const assert = {
		"_equal": function (codepoint, generics, expected, value) {
			if (!deepEqual(value, expected)) {
				throw testEqualityErrorMessage(codepoint, generics, value, expected)
			}
			return {"$type": "Void", "value": null}
		},
		"_fail": function (generics, message) {
			throw "@codepointtodo" + ": " + message["value"]
		},
	}

	const refCreator = {
		"_new": function (generics, value) {
			let ref = value
			return {
				"_get": function (generics) {
					return ref
				},
				"_set": function (generics, value) {
					ref = value
					return {"$type": "Void", "value": null}
				},
				"_modify": function (generics, f) {
					ref = f([], ref)
					return {"$type": "Void", "value": null}
				},
			}
		},
	}

	const testkit = {
		"_assert": assert,
		"_ref": refCreator,
	}

	return {
		"_test": function (generics, name, theTest) {
			const testName = name["value"]
			const testFunc = theTest
			let testError = null
			try {
				testFunc([], testkit)
			} catch (err) {
				testError = err
			}
			const ok = "[\u001b[32mOK\u001b[0m]"
			const fail = "[\u001b[31mFAILURE\u001b[0m]"
			const testResultString = testError ? fail : ok
			if (testError) {
				testSummary.runFail += 1
			} else {
				testSummary.runOk += 1
			}
			process.stdout.write(`  ${testResultString} ${testName}\n`)
			if (testError) {
				const errMsg = (testError && typeof testError === 'string') ? testError : ('could not print the failure: ' + typeof testError + ": " + testError)
				process.stdout.write(`    ${errMsg}\n`)
			}
			testSummary.runTotal += 1
			return {"$type": "Void", "value": null}
		},
	}
}

function testEqualityErrorMessage(codepoint, generics, value, expected) {
	const valueIsObject = isPlainObject(value)
	const expectedIsObject = isPlainObject(expected)
	if (valueIsObject && expectedIsObject) {
		try {
			const valueJson = JSON.stringify(value)
			const expectedJson = JSON.stringify(expected)
			return `@${codepoint}: ${expectedJson} is not equal to ${valueJson}`
		} catch (_e) {}
	}
	return `@${codepoint}: ${expected} is not equal to ${value}`
}

function runtime() {
	const log = function (generics, msg) {
		process.stdout.write(msg["value"] + "\n")
		return {"$type": "Void", "value": null}
	}
	const console = {
		"_log": log,
	}
	const refCreator = {
		"_new": function (generics, value) {
			let ref = value
			return {
				"_get": function (generics) {
					return ref
				},
				"_set": function (generics, value) {
					ref = value
					return {"$type": "Void", "value": null}
				},
				"_modify": function (generics, f) {
					ref = f([], ref)
					return {"$type": "Void", "value": null}
				},
			}
		},
	}
	const r = {
		"_console": console,
		"_ref": refCreator,
	}
	return r
}

function isPlainObject(x) {
	return x !== null && typeof x === 'object' && !Array.isArray(x)
}

function deepEqual(a, b) {
	if (a === b) return true
	if (typeof a !== typeof b) return false
	if (isPlainObject(a) && isPlainObject(b)) {
		const aKeys = Object.keys(a)
		const bKeys = Object.keys(b)
		if (aKeys.length !== bKeys.length) return false
		for (const key of aKeys) {
			if (!deepEqual(a[key], b[key])) return false
		}
		return true
	}
	if (Array.isArray(a) && Array.isArray(b)) {
		if (a.length !== b.length) return false
		for (let i = 0; i < a.length; i++) {
			if (!deepEqual(a[i], b[i])) return false
		}
		return true
	}
	return a === b
}
