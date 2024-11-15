package testcode

const When TestCodeCategory = "When"

var WhenNonOr = Create(When, "WhenNonOr", `package main


asString := (arg: Boolean): String => {
  when arg {
    is a: Boolean => {
      if a {
        "true"
      } else {
        "false"
      }
    }
  }
}
`)

var WhenExplicitExhaustive = Create(When, "WhenExplicitExhaustive", `package main


asString := (arg: Boolean | String): String => {
  when arg {
    is a: Boolean => {
      if a {
        "true"
      } else {
        "false"
      }
    }
    is b: String => {
      b
    }
  }
}
`)

var WhenAnnotatedVariable = Create(When, "WhenAnnotatedVariable", `package main


asString := (arg: Boolean | String): String => {
  result: String = when arg {
    is a: Boolean => {
      if a {
        "true"
      } else {
        "false"
      }
    }
    is s: String => {
      s
    }
  }
  result
}
`)

var WhenOtherSingleType = Create(When, "WhenOtherSingleType", `package main


asString := (arg: Boolean | String): String => {
  when arg {
    is a: Boolean => {
      if a {
        "true"
      } else {
        "false"
      }
    }
    other a => {
      a
    }
  }
}
`)

var WhenOtherMultipleTypes = Create(When, "WhenOtherMultipleTypes", `package main


yeetString := (arg: Boolean | String | Void): Boolean | Void => {
  when arg {
    is String => {
      false
    }
    other a => {
      a
    }
  }
}
`)

var WhenStruct = Create(When, "WhenStruct", `package main

import tenecs.string.join

struct Post(
  title: String
)

struct BlogPost(
  title: String
)

toString := (input: String | Post | BlogPost): String => {
  when input {
    is s: String => {
      s
    }
    is p: Post => {
      join("post:", p.title)
    }
    is b: BlogPost => {
      join("blogpost:", b.title)
    }
  }
}
`)

var WhenGeneric = Create(When, "WhenGeneric", `package main

import tenecs.string.join

f := <T>(input: String | T): String | T => {
  when input {
    is s: String => {
      join(s, "!")
    }
    other t => {
      t
    }
  }
}

usage := (): Void => {
  str: String = f("")

  strOrInt: String | Int = f(1)
}
`)
