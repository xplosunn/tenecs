package testcode

const When TestCodeCategory = "When"

var WhenExplicitExhaustive = Create(When, "WhenExplicitExhaustive", `
package main

asString := (arg: Boolean | String): String => {
  when arg {
    is Boolean => {
      if arg {
        "true"
      } else {
        "false"
      }
    }
    is String => {
      arg
    }
  }
}
`)

var WhenAnnotatedVariable = Create(When, "WhenAnnotatedVariable", `
package main

asString := (arg: Boolean | String): String => {
  result: String = when arg {
    is Boolean => {
      if arg {
        "true"
      } else {
        "false"
      }
    }
    is String => {
      arg
    }
  }
  result
}
`)

var WhenOtherSingleType = Create(When, "WhenOtherSingleType", `
package main

asString := (arg: Boolean | String): String => {
  when arg {
    is Boolean => {
      if arg {
        "true"
      } else {
        "false"
      }
    }
    other => {
      arg
    }
  }
}
`)

var WhenOtherMultipleTypes = Create(When, "WhenOtherMultipleTypes", `
package main

yeetString := (arg: Boolean | String | Void): Boolean | Void => {
  when arg {
    is String => {
      false
    }
    other => {
      arg
    }
  }
}
`)

var WhenStruct = Create(When, "WhenStruct", `
package main

import tenecs.string.join

struct Post(title: String)

struct BlogPost(title: String)

toString := (input: String | Post | BlogPost): String => {
  when input {
    is String => {
      input
    }
    is Post => {
      join("post:", input.title)
    }
    is BlogPost => {
      join("blogpost:", input.title)
    }
  }
}
`)
