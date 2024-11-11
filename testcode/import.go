package testcode

const Import TestCodeCategory = "import"

var ImportAliasMain = Create(Import, "ImportAliasMain", `package main

import tenecs.go.Main as App
import tenecs.go.Runtime as Lib
import tenecs.string.join as concat

app := App(main = (runtime: Lib) => {
  runtime.console.log(concat("Hello ", "world!"))
})
`)
