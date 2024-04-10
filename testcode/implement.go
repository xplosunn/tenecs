package testcode

const Implementation TestCodeCategory = "implementation"

var ImplementationWithConstructorEmpty = Create(Implementation, "ImplementationWithConstructorEmpty", `
package main

interface A {
  a: () -> String
}

app := (): A => implement A {
  a := () => ""
}
`)

var ImplementationWithConstructorWithArgUnused = Create(Implementation, "ImplementationWithConstructorWithArgUnused", `
package main

interface A {
  a: () -> String
}

app := (str: String): A => implement A {
  a := () => ""
}
`)

var ImplementationWithConstructorWithArgUsed = Create(Implementation, "ImplementationWithConstructorWithArgUsed", `
package main

interface A {
  a: () -> String
}

app := (str: String): A => implement A {
  a := () => str
}
`)

var ImplementationCreation1 = Create(Implementation, "ImplementationCreation1", `
package main

interface Goods {
  name: () -> String
}

food := (): Goods => implement Goods {
  name := () => "food"
}

interface Factory {
  produce: () -> Goods
}

foodFactory := (): Factory => implement Factory {
  produce := (): Goods => {
		food()
	}
}
`)

var ImplementationCreation2 = Create(Implementation, "ImplementationCreation2", `
package main

food := (): Goods => implement Goods {
  name := () => "food"
}

interface Goods {
  name: () -> String
}

foodFactory := (): Factory => implement Factory {
  produce := (): Goods => {
		food()
	}
}

interface Factory {
  produce: () -> Goods
}
`)

var ImplementationCreation3 = Create(Implementation, "ImplementationCreation3", `
package main

foodFactory := (): Factory => implement Factory {
  produce := (): Goods => {
		food()
	}
}

interface Factory {
  produce: () -> Goods
}

food := (): Goods => implement Goods {
  name := () => "food"
}

interface Goods {
  name: () -> String
}
`)

var ImplementationSelfCreation = Create(Implementation, "ImplementationSelfCreation", `
package main

interface Clone {
  copy: () -> Clone
}

clone := (): Clone => implement Clone {
  copy := () => {
		clone()
	}
}
`)

var ImplementationWithAnnotatedVariable = Create(Implementation, "ImplementationWithAnnotatedVariable", `
package main

interface A {
  a: () -> String
}

app := (): A => implement A {
  a: () -> String = () => ""
}
`)
