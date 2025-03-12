package testcode

const Recursion TestCodeCategory = "recursion"

var RecursionFactorial = Create(Recursion, "RecursionFactorial", `package main

import tenecs.compare.eq
import tenecs.int.minus
import tenecs.int.times

factorial := (i: Int): Int => {
  if eq<Int>(i, 0) {
    1
  } else {
    times(i, factorial(minus(i, 1)))
  }
}
`)

var RecursionLocalFactorial = Create(Recursion, "RecursionLocalFactorial", `package main

import tenecs.compare.eq
import tenecs.int.minus
import tenecs.int.times

factorial := (of: Int): Int => {
  go := (i: Int): Int => {
    if eq<Int>(i, 0) {
      1
    } else {
      times(i, go(minus(i, 1)))
    }
  }
  go(of)
}
`)
