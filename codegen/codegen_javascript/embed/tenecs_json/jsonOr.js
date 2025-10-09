function tenecs_json__jsonOr() {
  return ((generics, ConverterA, ConverterB, toJsonConverterPicker) => {
    return {
      "$type": "tenecs.json.JsonConverter",
      "_fromJson": (generics, input) => {
        const resultA = ConverterA._fromJson([], input)
        if (resultA.$type === "tenecs.error.Error") {
          const resultB = ConverterB._fromJson([], input)
          if (resultB.$type === "tenecs.error.Error") {
            return {
              "$type": "tenecs.error.Error",
              "_message": {
                "$type": "String",
                "value": "Could not parse from " + input.value
              }
            }
          }
          return resultB
        }
        return resultA
      },
      "_toJson": (generics, input) => {
        const Converter = toJsonConverterPicker([], input)
        return Converter._toJson([], input)
      }
    }
  })
}
