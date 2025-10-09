function tenecs_json__jsonInt() {
  return ((generics) => {
    return {
      "$type": "tenecs.json.JsonConverter",
      "_fromJson": (generics, input) => {
        try {
          const output = JSON.parse(input.value)
          if (typeof output !== 'number' || !Number.isInteger(output)) {
            throw new Error("Not an integer")
          }
          return {
            "$type": "Int",
            "value": output
          }
        } catch (err) {
          return {
            "$type": "tenecs.error.Error",
            "_message": {
              "$type": "String",
              "value": "Could not parse Int from " + input.value
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
