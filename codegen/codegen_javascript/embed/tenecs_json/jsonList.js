function tenecs_json__jsonList() {
  return ((topGenerics, of) => {
    return {
      "$type": "tenecs.json.JsonConverter",
      "_fromJson": (generics, input) => {
        try {
          const output = JSON.parse(input.value)
          if (!Array.isArray(output)) {
            throw new Error("Not an array")
          }
          
          const listTypeArgument = topGenerics[0]
          const listOf = typeNameCollapseRecursive(listTypeArgument)
          
          const ofParse = of._fromJson
          const outputList = []
          for (const elem of output) {
            const result = ofParse([], {
              "$type": "String",
              "value": JSON.stringify(elem)
            })
            if (result.$type === "tenecs.error.Error") {
              return result
            }
            outputList.push(result)
          }
          return {
            "$type": "List<" + listOf + ">",
            "value": outputList
          }
        } catch (err) {
          return {
            "$type": "tenecs.error.Error",
            "_message": {
              "$type": "String",
              "value": "Could not parse List from " + input.value
            }
          }
        }
      },
      "_toJson": (generics, input) => {
        const results = []
        const ofToJson = of._toJson
        for (const elem of input.value) {
          const result = ofToJson([], elem)
          results.push(result.value)
        }
        return {
          "$type": "String",
          "value": "[" + results.join(",") + "]"
        }
      }
    }
  })
}
