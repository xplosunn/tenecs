package testcode

const ArrowInvocation TestCodeCategory = "ArrowInvocation"

var ArrowInvocationOneArg = Create(ArrowInvocation, "ArrowInvocationOneArg", `package main


f := (str: String): String => {
  str
}

usage := (): String => {
  str := "foo"
  str->f()
}
`)

var ArrowInvocationOneArgChain = Create(ArrowInvocation, "ArrowInvocationOneArgChain", `package main


f := (str: String): String => {
  str
}

g := (str: String): String => {
  str
}

h := (str: String): String => {
  str
}

usage := (): String => {
  str := "foo"
  str->f()->g()->h()
}
`)

var ArrowInvocationTwoArg = Create(ArrowInvocation, "ArrowInvocationTwoArg", `package main


f := (str: String, str2: String): String => {
  str
}

usage := (): String => {
  str := "foo"

  str2 := "foo"
  str->f(str2)
}
`)

var ArrowInvocationThreeArg = Create(ArrowInvocation, "ArrowInvocationThreeArg", `package main


f := (str: String, str2: String, str3: String): String => {
  str
}

usage := (): String => {
  str := "foo"

  str2 := "foo"

  str3 := "foo"
  str->f(str2, str3)
}
`)
