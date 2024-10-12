package testcode

const Generics TestCodeCategory = "generics"

var GenericFunctionDeclared = Create(Generics, "GenericFunctionDeclared", `
package main

import tenecs.go.Main

app := Main(
  main = (runtime): Void => {
  }
)

identity := <T>(arg: T): T => {
  arg
}
`)

var GenericFunctionInvoked1 = Create(Generics, "GenericFunctionInvoked1", `
package main

import tenecs.go.Main

app := Main(
  main = (runtime): Void => {
		output := "Hello world!"
		hw := identity<String>(output)
		runtime.console.log(hw)
	}
)

identity := <T>(arg: T): T => {
  arg
}
`)

var GenericFunctionInvoked2 = Create(Generics, "GenericFunctionInvoked2", `
package main

import tenecs.go.Main

app := Main(
  main = (runtime): Void => {
    hw := identity<String>("Hello world!")
    runtime.console.log(hw)
  }
)

identity := <T>(arg: T): T => {
  arg
}
`)

var GenericFunctionInvoked3 = Create(Generics, "GenericFunctionInvoked3", `
package main

import tenecs.go.Main

app := Main(
  main = (runtime): Void => {
    runtime.console.log(identity<String>("Hello world!"))
  }
)

identity := <T>(arg: T): T => {
  arg
}
`)

var GenericFunctionInvoked4 = Create(Generics, "GenericFunctionInvoked4", `
package main

import tenecs.go.Main

app := Main(
  main = (runtime): Void => {
    output := "Hello world!"
    hw := identity<String>(output)
    runtime.console.log(hw)
  }
)

identity := <T>(arg: T): T => {
  result := arg
  result
}
`)

var GenericFunctionDoubleInvoked = Create(Generics, "GenericFunctionDoubleInvoked", `
package main

import tenecs.go.Main

app := Main(
  main = (runtime): Void => {
    runtime.console.log(identity<String>("ciao"))
  }
)

identity := <T>(arg: T): T => {
  output := identityFn<T>(arg)
  output
}
identityFn := <A>(arg: A): A => {
  result := arg
  result
}
`)

var GenericStruct = Create(Generics, "GenericStruct", `
package main

struct Box<T>(inside: T)
`)

var GenericStructInstance = Create(Generics, "GenericStructInstance", `
package main

import tenecs.go.Main

struct Box<T>(inside: T)

app := Main(
  main = (runtime) => {
    box := Box<String>("Hello world!")
    runtime.console.log(box.inside)
  }
)
`)

var GenericImplementedStructFunctionAllAnnotated = Create(Generics, "GenericImplementedStructFunctionAllAnnotated", `
package main

struct IdentityFunction(
  identity: <T>(T) ~> T
)

id := (): IdentityFunction => IdentityFunction(
  identity = <T>(t: T): T => {
		t
	}
)
`)

var GenericImplementedStructFunctionAnnotatedReturnType = Create(Generics, "GenericImplementedStructFunctionAnnotatedReturnType", `
package main

struct IdentityFunction(
  identity: <T>(T) ~> T
)

id := (): IdentityFunction => IdentityFunction(
  identity = <T>(t): T => {
		t
	}
)
`)

var GenericImplementedStructFunctionAnnotatedArg = Create(Generics, "GenericImplementedStructFunctionAnnotatedArg", `
package main

struct IdentityFunction(
  identity: <T>(T) ~> T
)

id := (): IdentityFunction => IdentityFunction(
  identity = <T>(t: T) => {
		t
	}
)
`)

var GenericImplementedStructFunctionNotAnnotated = Create(Generics, "GenericImplementedStructFunctionNotAnnotated", `
package main

struct IdentityFunction(
  identity: <T>(T) ~> T
)

id := (): IdentityFunction => IdentityFunction(
  identity = <T>(t) => {
		t
	}
)
`)

var GenericFunctionFixingList = Create(Generics, "GenericFunctionFixingList", `
package mypackage

emptyStringList := (): List<String> => {
  [String]()
}
`)

var GenericFunctionSingleElementList = Create(Generics, "GenericFunctionSingleElementList", `
package mypackage

import tenecs.list.append

listOf := (elem: String): List<String> => {
  append<String>([String](), elem)
}
`)

var GenericFunctionTakingList = Create(Generics, "GenericFunctionTakingList", `
package mypackage

toJson := <T>(t: T): String => {
  "not actually implemented"
}

doStuff := (): String => {
  arr := [String]("a", "b")
  toJson<List<String>>(arr)
}
`)

var GenericStructFunction = Create(Generics, "GenericStructFunction", `
package mypackage

struct Box<T>(elem: T)

f := <T>(): Box<String> => {
  b := Box<String>("wee")
  b
}
`)

// TODO FIXME change _map to map
var GenericIO = Create(Generics, "GenericIO", `
package mypackage

struct IO<A>(
  run: () ~> A,
  _map: <B>((A) ~> B) ~> IO<B>
)

make := <A>(a: () ~> A): IO<A> => {
  IO<A>(
    run = () => {
      a()
    },
    _map = <B>(f: (A) ~> B): IO<B> => {
      make<B>(() => { f(a()) })
    }
  )
}
`)

var GenericFromJson = Create(Generics, "GenericFromJson", `
package mypackage

import tenecs.compare.eq
import tenecs.string.join

struct Error(message: String)

struct FromJson<A>(
  parse: (String) ~> A | Error
)

parseBoolean := FromJson<Boolean>(
  parse = (input: String): Boolean | Error => {
    if eq<String>(input, "true") {
      true
    } else {
      if eq<String>(input, "false") {
        false
      } else {
        Error(join("Couldn't parse boolean from '", join(input, "'")))
      }
    }
  }
)
`)
