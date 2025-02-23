package testcode

const ShortCircuit TestCodeCategory = "shortcircuit"

var ShortCircuitExplicit = Create(ShortCircuit, "ShortCircuitExplicit", `package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  str: String ? Int = stringOrInt()
  join(str, "!")
}
`)

var ShortCircuitInferLeft = Create(ShortCircuit, "ShortCircuitInferLeft", `package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  str :? Int = stringOrInt()
  join(str, "!")
}
`)

var ShortCircuitInferRight = Create(ShortCircuit, "ShortCircuitInferRight", `package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  str: String ?= stringOrInt()
  join(str, "!")
}
`)

var ShortCircuitTwice = Create(ShortCircuit, "ShortCircuitTwice", `package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  str: String ?= stringOrInt()

  strAgain :? Int = stringOrInt()
  join(str, strAgain)
}
`)

var ShortCircuitInsideFunction = Create(ShortCircuit, "ShortCircuitInsideFunction", `package main

import tenecs.error.Error
import tenecs.list.atIndexGet
import tenecs.list.find
import tenecs.string.join

head := <T>(list: List<T>): T | Void => {
  when atIndexGet(list, 0) {
    is Error => {
      null
    }
    other result => {
      result
    }
  }
}

usage := (): String | Void => {
  list := <List<String> | Void>[null, ["a"], null, ["b"]]
  list->find((maybeString) => {
    strings :? Void = maybeString

    void: Void ?= head<String>(strings)
    void
  })
}
`)
