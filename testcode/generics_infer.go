package testcode

const GenericsInfer TestCodeCategory = "generics_infer"

var GenericsInferEmptyListArg = Create(GenericsInfer, "GenericsInferEmptyListArg", `package main


f := (arg: List<String>): Void => {
  null
}

usage := (): Void => {
  f([]())
}
`)

var GenericsInferEmptyListAssignment = Create(GenericsInfer, "GenericsInferEmptyListAssignment", `package main


usage := (): Void => {
  strings: List<String> = []()
}
`)

var GenericsInferIdentity = Create(GenericsInfer, "GenericsInferIdentity", `package main


identity := <T>(arg: T): T => {
  arg
}

usage := (): String => {
  identity("foo")
}
`)

var GenericsInferOrSecondArgument = Create(GenericsInfer, "GenericsInferOrSecondArgument", `package main


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

var GenericsInferList = Create(GenericsInfer, "GenericsInferList", `package main

import tenecs.compare.eq
import tenecs.list.length

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

var GenericsInferHigherOrderFunction = Create(GenericsInfer, "GenericsInferHigherOrderFunction", `package main

import tenecs.list.map

usage := (): List<String> => {
  map([String](), (str) => {
    str
  })
}
`)

var GenericsInferHigherOrderFunctionOr = Create(GenericsInfer, "GenericsInferHigherOrderFunctionOr", `package main

import tenecs.compare.eq
import tenecs.list.mapNotNull

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

var GenericsInferHigherOrderFunctionOr2 = Create(GenericsInfer, "GenericsInferHigherOrderFunctionOr2", `package main

import tenecs.compare.eq
import tenecs.list.mapNotNull

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

var GenericsInferTypeParameter = Create(GenericsInfer, "GenericsInferTypeParameter", `package main

import tenecs.json.JsonSchema
import tenecs.json.jsonString
import tenecs.ref.Ref

restHandlerPost := <RequestBody, ResponseBody>(fromJson: JsonSchema<RequestBody>, toJson: JsonSchema<ResponseBody>, route: String, handler: (RequestBody, Ref<Int>) ~> ResponseBody): Void => {
  null
}

usage := (): Void => {
  restHandlerPost(jsonString(), jsonString(), "/echo", (req, statusRef) => {
    req
  })
}
`)

var GenericsInferTypeParameterPartialLeft = Create(GenericsInfer, "GenericsInferTypeParameterPartialLeft", `package main


pickRight := <L, R>(left: L, right: R): R => {
  right
}

usage := (): Void => {
  str := pickRight<_, String>("", "")
}
`)

var GenericsInferFunctionResult = Create(GenericsInfer, "GenericsInferFunctionResult", `package main

import tenecs.boolean.and
import tenecs.test.UnitTest
import tenecs.test.UnitTestKit

_ := UnitTest("and", (testkit: UnitTestKit): Void => {
  testkit.assert.equal(false, and(false, () => {
    testkit.assert.fail("invoked")
  }))
})
`)
