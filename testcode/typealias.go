package testcode

const Typealias TestCodeCategory = "typealias"

var TypealiasSimple = Create(Typealias, "TypealiasSimple", `
package main

typealias UserId = String
`)

var TypealiasSimpleOr = Create(Typealias, "TypealiasSimpleOr", `
package main

typealias MaybeString = String | Void
`)

var TypealiasGeneric = Create(Typealias, "TypealiasGeneric", `
package main

typealias Just<A> = A
`)

var TypealiasGenericOr = Create(Typealias, "TypealiasGenericOr", `
package main

typealias Maybe<A> = A | Void
`)

var TypealiasSimpleUsed = Create(Typealias, "TypealiasSimpleUsed", `
package main

typealias UserId = String

newUserId := (str: String): UserId => {
  str
}
`)

var TypealiasSimpleOrUsed = Create(Typealias, "TypealiasSimpleOrUsed", `
package main

typealias MaybeString = String | Void

caseVoid: MaybeString = null
caseString: MaybeString = ""
`)

var TypealiasGenericUsed = Create(Typealias, "TypealiasGenericUsed", `
package main

typealias Just<A> = A

value: Just<String> = ""
`)

var TypealiasGenericOrUsed = Create(Typealias, "TypealiasGenericOrUsed", `
package main

typealias Maybe<A> = A | Void

strValue: Maybe<String> = ""
strNull: Maybe<String> = null

boolValue: Maybe<Boolean> = true
boolNull: Maybe<Boolean> = null
`)

var TypealiasGenericUsedGeneric = Create(Typealias, "TypealiasGenericUsedGeneric", `
package main

typealias Just<A> = A

identity := <T>(t: T): Just<T> => {
  t
}
`)

var TypealiasNested = Create(Typealias, "TypealiasNested", `
package main

struct Box<A>(value: A)

typealias DoubleBox<T> = Box<Box<T>>

value: DoubleBox<Boolean> = Box(Box(false))
`)
