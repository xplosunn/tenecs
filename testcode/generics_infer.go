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

var GenericsInferList = Create(GenericsInfer, "GenericsInferList", `
package main

import tenecs.list.length
import tenecs.compare.eq

nonEmpty := <T>(list: List<T>): List<T> | Void => {
  if eq(length(list), 0) {
    null
  } else {
    list
  }
}

usage := (): List<String> | Void => {
  nonEmpty([String]())
}
`)

var GenericsInferHigherOrderFunction = Create(GenericsInfer, "GenericsInferHigherOrderFunction", `
package main

import tenecs.list.map

usage := (): List<String> => {
  map([String](), (str) => str)
}
`)

var GenericsInferHigherOrderFunctionOr = Create(GenericsInfer, "GenericsInferHigherOrderFunctionOr", `
package main

import tenecs.list.mapNotNull
import tenecs.compare.eq

usage := (): List<String> => {
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

import tenecs.list.mapNotNull
import tenecs.compare.eq

usage := (): List<String> => {
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
  server := newServer(testkit.ref)
  server.restHandlerPost(jsonString(), jsonString(), "/echo", (req, statusRef) => {
    req
  })
}
`)

var GenericsInferTypeParameterPartialLeft = Create(GenericsInfer, "GenericsInferTypeParameterPartialLeft", `
package main

pickRight := <L, R>(left: L, right: R): R => {
  right
}

usage := (): Void => {
  str := pickRight<_, String>("", "")
}
`)
