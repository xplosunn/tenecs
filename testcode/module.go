package testcode

const Module TestCodeCategory = "module"

var ModuleWithConstructorEmpty = Create(Module, "ModuleWithConstructorEmpty", `
package main

interface A {
	public a: String
}

app := (): A => implement A {
	public a := ""
}
`)

var ModuleWithConstructorWithArgUnused = Create(Module, "ModuleWithConstructorWithArgUnused", `
package main

interface A {
	public a: String
}

app := (str: String): A => implement A {
	public a := ""
}
`)

var ModuleWithConstructorWithArgUsed = Create(Module, "ModuleWithConstructorWithArgUsed", `
package main

interface A {
	public a: String
}

app := (str: String): A => implement A {
	public a := str
}
`)

var ModuleCreation1 = Create(Module, "ModuleCreation1", `
package main

interface Goods {
	public name: String
}

food := (): Goods => implement Goods {
	public name := "food"
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

var ModuleCreation2 = Create(Module, "ModuleCreation2", `
package main

food := (): Goods => implement Goods {
	public name := "food"
}

interface Goods {
	public name: String
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

var ModuleCreation3 = Create(Module, "ModuleCreation3", `
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
	public name := "food"
}

interface Goods {
	public name: String
}
`)

var ModuleSelfCreation = Create(Module, "ModuleSelfCreation", `
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
