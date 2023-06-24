package testcode

const Or TestCodeCategory = "or"

var OrVariableWithEmptyArray = Create(Or, "OrVariableWithEmptyArray", `
package main

empty := [ String | Boolean ] ( )
`)

var OrVariableWithTwoElementArray = Create(Or, "OrVariableWithTwoElementArray", `
package main

hasStuff := [ Boolean | String ] ( "first", false )
`)

var OrFunction = Create(Or, "OrFunction", `
package main

strOrBool := (): String | Boolean => {
  ""
}
`)

var OrArrayFunction = Create(Or, "OrArrayFunction", `
package main

strOrBool := (): Array<String | Boolean> => {
  [Boolean]()
}
`)
