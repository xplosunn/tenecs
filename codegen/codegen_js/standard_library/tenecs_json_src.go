// ##################################################################
// # The signatures of this file are generated via code-generation. #
// # Check gen.go                                                   #
// ##################################################################
package standard_library

import "github.com/xplosunn/tenecs/typer/standard_library"

func tenecs_json_jsonBoolean() Function {
	return function(
		body(`return ({
  "$type": "JsonSchema",
  "fromJson": (input) => {
    try {
      let parsed = JSON.parse(input)
      if (typeof parsed == "boolean") {
        return parsed
      }
    } catch (e) {}
    return ({
      "$type": "Error",
      "message": "Could not parse Boolean from " + input
    })
  },
  "toJson": (input) => {
    return JSON.stringify(input)
  },
})`),
	)
}

func tenecs_json_jsonInt() Function {
	return function(
		body(`return ({
  "$type": "JsonSchema",
  "fromJson": (input) => {
    try {
      let parsed = JSON.parse(input)
      if (typeof parsed == "number" && parsed % 1 == 0) {
        return parsed
      }
    } catch (e) {}
    return ({
      "$type": "Error",
      "message": "Could not parse Int from " + input
    })
  },
  "toJson": (input) => {
    return JSON.stringify(input)
  },
})`),
	)
}

func tenecs_json_jsonOr() Function {
	return function(
		params("schemaA", "schemaB", "toJsonSchemaPicker"),
		body(`return ({
  "$type": "JsonSchema",
  "fromJson": (input) => {
    let resultA = schemaA["fromJson"](input)
    if (resultA && resultA["$type"] && resultA["$type"] == "Error") {
      let resultB = schemaB["fromJson"](input)
      if (resultB && resultB["$type"] && resultB["$type"] == "Error") {
        return ({
          "$type": "Error",
          "message": "Could not parse from " + input
        })
      }
      return resultB
    }
    return resultA
  },
  "toJson": (input) => {
    let schema = toJsonSchemaPicker(input)
    return schema["toJson"](input)
  }
})`),
	)
}

func tenecs_json_jsonString() Function {
	return function(
		body(`return ({
  "$type": "JsonSchema",
  "fromJson": (input) => {
    try {
      let parsed = JSON.parse(input)
      if (typeof parsed == "string") {
        return parsed
      }
    } catch (e) {}
    return ({
      "$type": "Error",
      "message": "Could not parse String from " + input
    })
  },
  "toJson": (input) => {
    return JSON.stringify(input)
  },
})`),
	)
}

func tenecs_json_jsonList() Function {
	return function(
		params("of"),
		body(`return ({
  "$type": "JsonSchema",
  "fromJson": (input) => {
    const result = []
    let fullParsed = JSON.parse(input)
    if (!Array.isArray(fullParsed)) {
      return ({
        "$type": "Error",
        "message": "Could not parse list from " + input
      })
    }
    for (const elem of fullParsed) {
      let field = of.fromJson(JSON.stringify(elem))
      if (field && typeof field == "object" && field["$type"] == "Error") {
        return ({
          "$type": "Error",
          "message": field.message
        })
      }
      result.push(field)
    }
    return result
  },
  "toJson": (input) => {
    let result = "["
    for (let i = 0; i < input.length; i++) {
      if (i > 0) {
        result += ","
      }
      result += of["toJson"](input[i])
    }
    result += "]"
    return result
  },
})`),
	)
}

func tenecs_json_jsonObject0() Function {
	return function(
		params("f"),
		body(`return ({
  "$type": "JsonSchema",
  "fromJson": (input) => {
    try {
      let parsed = JSON.parse(input)
      if (typeof parsed == "object") {
        return f()
      }
    } catch (e) {}
    return ({
      "$type": "Error",
      "message": "Could not parse object from " + input
    })
  },
  "toJson": (input) => {
    return "{}"
  },
})`),
	)
}

func tenecs_json_jsonObject_X(x int) Function {
	return function(
		params("f"), // the others are not listed
		body(`
let fieldParsers = []
for (let i = 1; i < arguments.length; i++) {
  fieldParsers.push(arguments[i])
}
return ({
  "$type": "JsonSchema",
  "fromJson": (input) => {
    let fullParsed = JSON.parse(input)
    if (typeof fullParsed != "object") {
      return ({
        "$type": "Error",
        "message": "Could not parse object from " + input
      })
    }
    let resultArguments = []
    for (const fieldParser of fieldParsers) {
      if (fullParsed[fieldParser.name] == undefined) {
        return ({
          "$type": "Error",
          "message": "Could not find object field \"" + fieldParser.name + "\" in " + input
        })
      } 
      let field = fieldParser.schema.fromJson(JSON.stringify(fullParsed[fieldParser.name]))
      if (field && typeof field == "object" && field["$type"] == "Error") {
        return ({
          "$type": "Error",
          "message": "Could not parse object field \"" + fieldParser.name + "\": " + field.message
        })
      }
      resultArguments.push(field)
    }
    return f(...resultArguments)
  },
  "toJson": (input) => {
    const sortedFieldParsers = [...fieldParsers]
    sortedFieldParsers.sort((a,b) => a.name.localeCompare(b.name))
    let result = "{"
    for (let i = 0; i < sortedFieldParsers.length; i++) {
      if (i > 0) {
        result += ","
      }
      const fieldParser = sortedFieldParsers[i] 
      result += "\"" + fieldParser.name + "\":" + fieldParser.schema.toJson(fieldParser.access(input))
    }
    result += "}"
    return result
  },
})`),
	)
}

func tenecs_json_jsonObject1() Function {
	return tenecs_json_jsonObject_X(1)
}
func tenecs_json_jsonObject7() Function {
	return tenecs_json_jsonObject_X(7)
}
func tenecs_json_jsonObject12() Function {
	return tenecs_json_jsonObject_X(12)
}
func tenecs_json_jsonObject13() Function {
	return tenecs_json_jsonObject_X(13)
}
func tenecs_json_jsonObject3() Function {
	return tenecs_json_jsonObject_X(3)
}
func tenecs_json_jsonObject6() Function {
	return tenecs_json_jsonObject_X(6)
}
func tenecs_json_jsonObject9() Function {
	return tenecs_json_jsonObject_X(9)
}
func tenecs_json_jsonObject2() Function {
	return tenecs_json_jsonObject_X(2)
}
func tenecs_json_jsonObject11() Function {
	return tenecs_json_jsonObject_X(11)
}
func tenecs_json_jsonObject14() Function {
	return tenecs_json_jsonObject_X(14)
}
func tenecs_json_jsonObject8() Function {
	return tenecs_json_jsonObject_X(8)
}
func tenecs_json_jsonObject10() Function {
	return tenecs_json_jsonObject_X(10)
}
func tenecs_json_jsonObject22() Function {
	return tenecs_json_jsonObject_X(22)
}
func tenecs_json_jsonObject4() Function {
	return tenecs_json_jsonObject_X(4)
}
func tenecs_json_jsonObject15() Function {
	return tenecs_json_jsonObject_X(15)
}
func tenecs_json_jsonObject19() Function {
	return tenecs_json_jsonObject_X(19)
}
func tenecs_json_jsonObject21() Function {
	return tenecs_json_jsonObject_X(21)
}
func tenecs_json_jsonObject17() Function {
	return tenecs_json_jsonObject_X(17)
}
func tenecs_json_jsonObject18() Function {
	return tenecs_json_jsonObject_X(18)
}
func tenecs_json_jsonObject5() Function {
	return tenecs_json_jsonObject_X(5)
}
func tenecs_json_jsonObject16() Function {
	return tenecs_json_jsonObject_X(16)
}
func tenecs_json_jsonObject20() Function {
	return tenecs_json_jsonObject_X(20)
}
func tenecs_json_JsonField() Function {
	return structFunction(standard_library.Tenecs_json_JsonField)
}
func tenecs_json_JsonSchema() Function {
	return structFunction(standard_library.Tenecs_json_JsonSchema)
}
