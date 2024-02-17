package testcode

const ShortCircuit TestCodeCategory = "shortcircuit"

var ShortCircuitExplicit = Create(ShortCircuit, "ShortCircuitExplicit", `
package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  str: String ? Int = stringOrInt()
  join(str, "!")
}
`)

var ShortCircuitInferLeft = Create(ShortCircuit, "ShortCircuitInferLeft", `
package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  str :? Int = stringOrInt()
  join(str, "!")
}
`)

var ShortCircuitInferRight = Create(ShortCircuit, "ShortCircuitInferRight", `
package main

import tenecs.string.join

stringOrInt := (): String | Int => {
  3
}

usage := (): String | Int => {
  str: String ?= stringOrInt()
  join(str, "!")
}
`)
