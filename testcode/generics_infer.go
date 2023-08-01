package testcode

const GenericsInfer TestCodeCategory = "generics_infer"

var GenericsInferIdentity = Create(GenericsInfer, "GenericsInferIdentity", `
package main

identity := <T>(arg: T): T => {
  arg
}

usage := (): String => {
  identity("foo")
}
`)

var GenericsInferOrSecondArgument = Create(GenericsInfer, "GenericsInferOrSecondArgument", `
package main

pickSecond := <T>(a: T, b: T): T => {
  b
}

stringOrBoolean := (): String | Boolean => {
  ""
}

usage := (): Void => {
  pickSecond("foo", stringOrBoolean())
  null
}
`)

var GenericsInferArray = Create(GenericsInfer, "GenericsInferArray", `
package main

import tenecs.array.length
import tenecs.compare.eq

nonEmpty := <T>(array: Array<T>): Array<T> | Void => {
  if eq(length(array), 0) {
    null
  } else {
    array
  }
}

usage := (): Array<String> | Void => {
  nonEmpty([String]())
}
`)
