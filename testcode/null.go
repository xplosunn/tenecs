package testcode

const Null TestCodeCategory = "Null"

var NullValue = Create(Null, "NullValue", `
package main

value := null
`)

var NullFunction = Create(Null, "NullFunction", `
package main

returnsNull := (): Void => {
  null
}
`)
