package testcode

const Or TestCodeCategory = "or"

var OrVariableWithEmptyList = Create(Or, "OrVariableWithEmptyList", `package main


empty := <String | Boolean>[]
`)

var OrVariableWithTwoElementList = Create(Or, "OrVariableWithTwoElementList", `package main


hasStuff := <Boolean | String>["first", false]
`)

var OrFunction = Create(Or, "OrFunction", `package main


strOrBool := (): String | Boolean => {
  ""
}
`)

var OrListFunction = Create(Or, "OrListFunction", `package main


strOrBool := (): List<String | Boolean> => {
  <String | Boolean>[]
}
`)
