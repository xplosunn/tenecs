package testcode

const Import TestCodeCategory = "import"

var ImportAliasMain = Create(Import, "ImportAliasMain", `
package main

import tenecs.os.Runtime as Lib
import tenecs.os.Main as App
import tenecs.string.join as concat

app := implement App {
  main := (runtime: Lib) => {
		runtime.console.log(concat("Hello ", "world!"))
	}
}
`)
