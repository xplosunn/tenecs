function tenecs_json__jsonObject0() {
  return ((generics, build) => {
    return {
      "$type": "tenecs.json.JsonConverter",
      "_fromJson": (generics, input) => {
        try {
          const output = JSON.parse(input.value)
          if (typeof output !== 'object' || Array.isArray(output)) {
            throw new Error("Not an object")
          }
          return build([])
        } catch (err) {
          return {
            "$type": "tenecs.error.Error",
            "_message": {
              "$type": "String",
              "value": "Could not parse object from " + input.value
            }
          }
        }
      },
      "_toJson": (generics, input) => {
        return {
          "$type": "String",
          "value": "{}"
        }
      }
    }
  })
}
