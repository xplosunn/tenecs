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
