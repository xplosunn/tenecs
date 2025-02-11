package testcode

const Functions TestCodeCategory = "functions"

var MainProgramWithSingleExpression = Create(Functions, "MainProgramWithSingleExpression", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime) => {
  runtime.console.log("Hello world!")
})
`)

var MainProgramAnnotatedType = Create(Functions, "MainProgramAnnotatedType", `package main.program

import tenecs.go.Main
import tenecs.go.Runtime

app: Main = Main((runtime: Runtime) => {
  runtime.console.log("Hello world!")
})
`)

var MainProgramWithInnerFunction = Create(Functions, "MainProgramWithInnerFunction", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime) => {
  go := (): Void => {
    runtime.console.log("Hello world!")
  }
  go()
})
`)

var MainProgramWithVariableWithFunction = Create(Functions, "MainProgramWithVariableWithFunction", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime) => {
  output := (): String => {
    "Hello world!"
  }
  runtime.console.log(output())
})
`)

var MainProgramWithVariableWithFunctionTakingFunction = Create(Functions, "MainProgramWithVariableWithFunctionTakingFunction", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime) => {
  output := (): String => {
    "Hello world!"
  }

  run := (f: () ~> String): String => {
    f()
  }
  runtime.console.log(run(output))
})
`)

var MainProgramWithVariableWithFunctionTakingFunctionFromStdLib1 = Create(Functions, "MainProgramWithVariableWithFunctionTakingFunctionFromStdLib1", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime) => {
  applyToString := (f: (String) ~> Void, str: String): Void => {
    f(str)
  }

  output := (): String => {
    "Hello world!"
  }
  applyToString(runtime.console.log, output())
})
`)

var MainProgramWithVariableWithFunctionTakingFunctionFromStdLib2 = Create(Functions, "MainProgramWithVariableWithFunctionTakingFunctionFromStdLib2", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime) => {
  applyToString := (f: (String) ~> Void, strF: () ~> String): Void => {
    f(strF())
  }

  output := (): String => {
    "Hello world!"
  }
  applyToString(runtime.console.log, output)
})
`)

var MainProgramWithVariableWithFunctionWithTypeInferred = Create(Functions, "MainProgramWithVariableWithFunctionWithTypeInferred", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime) => {
  applyToString := (f: (String) ~> Void, strF: () ~> String): Void => {
    f(strF())
  }
  applyToString(runtime.console.log, () => {
    "Hello World!"
  })
})
`)

var MainProgramWithAnotherFunctionTakingConsole = Create(Functions, "MainProgramWithAnotherFunctionTakingConsole", `package main

import tenecs.go.Console
import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime) => {
  mainRun(runtime.console)
})

mainRun := (console: Console): Void => {
  console.log("Hello world!")
}
`)

var MainProgramWithAnotherFunctionTakingConsoleAndMessage = Create(Functions, "MainProgramWithAnotherFunctionTakingConsoleAndMessage", `package main

import tenecs.go.Console
import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime) => {
  mainRun(runtime.console, "Hello world!")
})

mainRun := (console: Console, message: String): Void => {
  console.log(message)
}
`)

var MainProgramWithAnotherFunctionTakingConsoleAndMessageFromAnotherFunction = Create(Functions, "MainProgramWithAnotherFunctionTakingConsoleAndMessageFromAnotherFunction", `package main

import tenecs.go.Console
import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime) => {
  mainRun(runtime.console, helloWorld())
})

mainRun := (console: Console, message: String): Void => {
  console.log(message)
}

helloWorld := (): String => {
  "Hello world!"
}
`)

var MainProgramWithArgAnnotatedArg = Create(Functions, "MainProgramWithArgAnnotatedArg", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime) => {
  runtime.console.log("Hello world!")
})
`)

var MainProgramWithArgAnnotatedReturn = Create(Functions, "MainProgramWithArgAnnotatedReturn", `package main

import tenecs.go.Main

app := Main((runtime): Void => {
  runtime.console.log("Hello world!")
})
`)

var MainProgramWithArgAnnotatedArgAndReturn = Create(Functions, "MainProgramWithArgAnnotatedArgAndReturn", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime: Runtime): Void => {
  runtime.console.log("Hello world!")
})
`)

var MainProgramWithAnotherFunctionTakingRuntime = Create(Functions, "MainProgramWithAnotherFunctionTakingRuntime", `package main

import tenecs.go.Main
import tenecs.go.Runtime

app := Main((runtime) => {
  mainRun(runtime)
})

mainRun := (runtime: Runtime): Void => {
  runtime.console.log("Hello world!")
}
`)

var FunctionsCallAndThenCall = Create(Functions, "FunctionsCallAndThenCall", `package main


f := (): () ~> String => {
  () => {
    ""
  }
}

usage := (): String => {
  f()()
}
`)

var FunctionsNamedArg = Create(Functions, "FunctionsNamedArg", `package main


f := (a: String, b: String): String => {
  a
}

usage := (): String => {
  f("", "")
  f(
    a = "",
    ""
  )
  f(
    "",
    b = ""
  )
  f(
    a = "",
    b = ""
  )
}
`)

var FunctionsNamedArgWithTypeAnnotation = Create(Functions, "FunctionsNamedArgWithTypeAnnotation", `package main


f: (a: String, b: String) ~> String = (a: String, b: String): String => {
  a
}

usage := (): String => {
  f("", "")
  f(
    a = "",
    ""
  )
  f(
    "",
    b = ""
  )
  f(
    a = "",
    b = ""
  )
}
`)
