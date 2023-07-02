package testcode

const Recursion TestCodeCategory = "recursion"

var RecursionFactorial = Create(Recursion, "RecursionFactorial", `
package main

import tenecs.int.times
import tenecs.int.minus
import tenecs.compare.eq

factorial := (i: Int): Int => {
  if eq<Int>(i, 0) {
    1
  } else {
    times(i, factorial(minus(i, 1)))
  }
}
`)
