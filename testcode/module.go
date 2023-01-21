package testcode

const Module TestCodeCategory = "module"

var ModuleWithConstructorEmpty = Create(Module, "ModuleWithConstructorEmpty", `
package main

interface A {
	public a: String
}

implementing A module app() {
	public a := ""
}
`)

var ModuleWithConstructorWithArgUnused = Create(Module, "ModuleWithConstructorWithArgUnused", `
package main

interface A {
	public a: String
}

implementing A module app(str: String) {
	public a := ""
}
`)

var ModuleWithConstructorWithArgUsed = Create(Module, "ModuleWithConstructorWithArgUsed", `
package main

interface A {
	public a: String
}

implementing A module app(str: String) {
	public a := str
}
`)

var ModuleWithConstructorWithArgPublic = Create(Module, "ModuleWithConstructorWithArgPublic", `
package main

interface A {
	public a: String
}

implementing A module app(public a: String) {
	
}
`)

var ModuleCreation1 = Create(Module, "ModuleCreation1", `
package main

interface Goods {
	public name: String
}

implementing Goods module food {
	public name := "food"
}

interface Factory {
	public produce: () -> Goods
}

implementing Factory module foodFactory() {
	public produce := (): Goods => {
		food()
	}
}
`)

var ModuleCreation2 = Create(Module, "ModuleCreation2", `
package main

implementing Goods module food {
	public name := "food"
}

interface Goods {
	public name: String
}

implementing Factory module foodFactory() {
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

implementing Factory module foodFactory() {
	public produce := (): Goods => {
		food()
	}
}

interface Factory {
	public produce: () -> Goods
}

implementing Goods module food {
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

implementing Clone module clone {
	public copy := () => {
		clone()
	}
}
`)
