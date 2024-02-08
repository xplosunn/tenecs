package testcode

const Implementation TestCodeCategory = "implementation"

var ImplementationWithConstructorEmpty = Create(Implementation, "ImplementationWithConstructorEmpty", `
package main

interface A {
	public a: () -> String
}

app := (): A => implement A {
	public a := () => ""
}
`)

var ImplementationWithConstructorWithArgUnused = Create(Implementation, "ImplementationWithConstructorWithArgUnused", `
package main

interface A {
	public a: () -> String
}

app := (str: String): A => implement A {
	public a := () => ""
}
`)

var ImplementationWithConstructorWithArgUsed = Create(Implementation, "ImplementationWithConstructorWithArgUsed", `
package main

interface A {
	public a: () -> String
}

app := (str: String): A => implement A {
	public a := () => str
}
`)

var ImplementationCreation1 = Create(Implementation, "ImplementationCreation1", `
package main

interface Goods {
	public name: () -> String
}

food := (): Goods => implement Goods {
	public name := () => "food"
}

interface Factory {
	public produce: () -> Goods
}

foodFactory := (): Factory => implement Factory {
	public produce := (): Goods => {
		food()
	}
}
`)

var ImplementationCreation2 = Create(Implementation, "ImplementationCreation2", `
package main

food := (): Goods => implement Goods {
	public name := () => "food"
}

interface Goods {
	public name: () -> String
}

foodFactory := (): Factory => implement Factory {
	public produce := (): Goods => {
		food()
	}
}

interface Factory {
	public produce: () -> Goods
}
`)

var ImplementationCreation3 = Create(Implementation, "ImplementationCreation3", `
package main

foodFactory := (): Factory => implement Factory {
	public produce := (): Goods => {
		food()
	}
}

interface Factory {
	public produce: () -> Goods
}

food := (): Goods => implement Goods {
	public name := () => "food"
}

interface Goods {
	public name: () -> String
}
`)

var ImplementationSelfCreation = Create(Implementation, "ImplementationSelfCreation", `
package main

interface Clone {
	public copy: () -> Clone
}

clone := (): Clone => implement Clone {
	public copy := () => {
		clone()
	}
}
`)

var ImplementationWithAnnotatedVariable = Create(Implementation, "ImplementationWithAnnotatedVariable", `
package main

interface A {
	public a: () -> String
}

app := (): A => implement A {
	public a: () -> String = () => ""
}
`)
