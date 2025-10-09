function tenecs_json__jsonBoolean() {
  return ((generics) => {
    return {
      "$type": "tenecs.json.JsonConverter",
      "_fromJson": (generics, input) => {
        try {
          const output = JSON.parse(input.value)
          if (typeof output !== 'boolean') {
            throw new Error("Not a boolean")
          }
          return {
            "$type": "Boolean",
            "value": output
          }
        } catch (err) {
          return {
            "$type": "tenecs.error.Error",
            "_message": {
              "$type": "String",
              "value": "Could not parse Boolean from " + input.value
            }
          }
        }
      },
      "_toJson": (generics, input) => {
        return {
          "$type": "String",
          "value": JSON.stringify(input.value)
        }
      }
    }
  })
}
