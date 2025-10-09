function tenecs_json__jsonString() {
  return ((generics) => {
    return {
      "$type": "tenecs.json.JsonConverter",
      "_fromJson": (generics, input) => {
        try {
          const output = JSON.parse(input.value)
          if (typeof output !== 'string') {
            throw new Error("Not a string")
          }
          return {
            "$type": "String",
            "value": output
          }
        } catch (err) {
          return {
            "$type": "tenecs.error.Error",
            "_message": {
              "$type": "String",
              "value": "Could not parse String from " + input.value
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
