function tenecs_json__jsonObject1() {
  return ((generics, build, jsonConverterFieldI0) => {
    return {
      "$type": "tenecs.json.JsonConverter",
      "_fromJson": (generics, input) => {
        try {
          const output = JSON.parse(input.value)
          if (typeof output !== 'object') {
            return {
              "$type": "tenecs.error.Error",
              "_message": {
                "$type": "String",
                "value": "Could not parse object from " + input.value
              }
            }
          }
          
          const i0Name = jsonConverterFieldI0._name.value
          if (!(i0Name in output)) {
            return {
              "$type": "tenecs.error.Error",
              "_message": {
                "$type": "String",
                "value": "Could not find object field \"" + i0Name + "\" in " + input.value
              }
            }
          }
          
          const i0JsonString = {
            "$type": "String",
            "value": JSON.stringify(output[i0Name])
          }
          const i0 = jsonConverterFieldI0._Converter._fromJson([], i0JsonString)
          if (i0.$type === "tenecs.error.Error") {
            return {
              "$type": "tenecs.error.Error",
              "_message": {
                "$type": "String",
                "value": "Could not parse object field \"" + i0Name + "\": " + i0._message.value
              }
            }
          }
          
          return build([], i0)
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
        const output = {}
        
        const fieldI0 = jsonConverterFieldI0._access([], input)
        const i0 = jsonConverterFieldI0._Converter._toJson([], fieldI0)
        output[jsonConverterFieldI0._name.value] = i0.value
        
        const keys = Object.keys(output).sort()
        const result = "{" + keys.map(k => JSON.stringify(k) + ":" + output[k]).join(",") + "}"
        return {
          "$type": "String",
          "value": result
        }
      }
    }
  })
}
