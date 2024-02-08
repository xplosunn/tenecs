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

var GenericsInferHigherOrderFunction = Create(GenericsInfer, "GenericsInferHigherOrderFunction", `
package main

import tenecs.array.map

usage := (): Array<String> => {
  map([String](), (str) => str)
}
`)

var GenericsInferHigherOrderFunctionOr = Create(GenericsInfer, "GenericsInferHigherOrderFunctionOr", `
package main

import tenecs.array.mapNotNull
import tenecs.compare.eq

usage := (): Array<String> => {
  mapNotNull([]("!", "a", "!", "b"), (str): String | Void => {
    if eq(str, "!") {
      null
    } else {
      str
    }
  })
}
`)

var GenericsInferHigherOrderFunctionOr2 = Create(GenericsInfer, "GenericsInferHigherOrderFunctionOr2", `
package main

import tenecs.array.mapNotNull
import tenecs.compare.eq

usage := (): Array<String> => {
  mapNotNull([]("!", "a", "!", "b"), (str) => {
    if eq(str, "!") {
      null
    } else {
      str
    }
  })
}
`)

var GenericsInferTypeParameter = Create(GenericsInfer, "GenericsInferTypeParameter", `
package main

import tenecs.http.newServer
import tenecs.json.jsonString
import tenecs.test.UnitTestKit

usage := (testkit: UnitTestKit): Void => {
  server := newServer(testkit.runtime.ref)
  server.restHandlerPost(jsonString(), jsonString(), "/echo", (req, statusRef) => {
    req
  })
}
`)
