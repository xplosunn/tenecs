package testcode

const BasicType TestCodeCategory = "BasicType"

var BasicTypeTrue = Create(BasicType, "BasicTypeTrue", `package main


value := true
`)

var BasicTypeFalse = Create(BasicType, "BasicTypeFalse", `package main


value := false
`)
