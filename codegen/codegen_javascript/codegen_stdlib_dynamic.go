package codegen_javascript

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xplosunn/tenecs/ir"
)

func stdLibDynamicFunction(ref ir.NativeFunctionRef) *string {
	if ref.Package == "tenecs_json" && strings.HasPrefix(ref.Name, "jsonObject") {
		return tenecs_json_jsonObject(ref.Name)
	}
	return nil
}

func tenecs_json_jsonObject(functionNameString string) *string {
	numberString := strings.TrimPrefix(functionNameString, "jsonObject")
	numberOfFields, err := strconv.Atoi(numberString)
	if err != nil {
		panic("tenecs_json_jsonObject: invalid function name: " + functionNameString)
	}

	fromJsonCode := `(generics, input) => {
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
`

	for i := 0; i < numberOfFields; i++ {
		fromJsonCode += fmt.Sprintf(`
          const i%dName = jsonConverterFieldI%d._name.value
          if (!(i%dName in output)) {
            return {
              "$type": "tenecs.error.Error",
              "_message": {
                "$type": "String",
                "value": "Could not find object field \"" + i%dName + "\" in " + input.value
              }
            }
          }
          
          const i%dJsonString = {
            "$type": "String",
            "value": JSON.stringify(output[i%dName])
          }
          const i%d = jsonConverterFieldI%d._Converter._fromJson([], i%dJsonString)
          if (i%d.$type === "tenecs.error.Error") {
            return {
              "$type": "tenecs.error.Error",
              "_message": {
                "$type": "String",
                "value": "Could not parse object field \"" + i%dName + "\": " + i%d._message.value
              }
            }
          }
`, i, i, i, i, i, i, i, i, i, i, i, i)
	}

	fromJsonCodeReturnArgs := ""
	for i := 0; i < numberOfFields; i++ {
		if i > 0 {
			fromJsonCodeReturnArgs += ", "
		}
		fromJsonCodeReturnArgs += fmt.Sprintf("i%d", i)
	}

	fromJsonCode += fmt.Sprintf(`
          return build([], %s)
        } catch (err) {
          return {
            "$type": "tenecs.error.Error",
            "_message": {
              "$type": "String",
              "value": "Could not parse object from " + input.value
            }
          }
        }
      }`, fromJsonCodeReturnArgs)

	toJsonCode := `(generics, input) => {
        const output = {}
`
	for i := 0; i < numberOfFields; i++ {
		toJsonCode += fmt.Sprintf(`
        const fieldI%d = jsonConverterFieldI%d._access([], input)
        const i%d = jsonConverterFieldI%d._Converter._toJson([], fieldI%d)
        output[jsonConverterFieldI%d._name.value] = i%d.value
`, i, i, i, i, i, i, i)
	}

	toJsonCode += `
        const keys = Object.keys(output).sort()
        const result = "{" + keys.map(k => JSON.stringify(k) + ":" + output[k]).join(",") + "}"
        return {
          "$type": "String",
          "value": result
        }
      }`

	jsonConverterArgs := ""
	for i := 0; i < numberOfFields; i++ {
		if i > 0 {
			jsonConverterArgs += ", "
		}
		jsonConverterArgs += fmt.Sprintf("jsonConverterFieldI%d", i)
	}

	code := fmt.Sprintf(`function tenecs_json__jsonObject%d() {
  return ((generics, build, %s) => {
    return {
      "$type": "tenecs.json.JsonConverter",
      "_fromJson": %s,
      "_toJson": %s,
    }
  })
}`, numberOfFields, jsonConverterArgs, fromJsonCode, toJsonCode)

	return &code
}
