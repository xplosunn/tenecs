package codegen_golang

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xplosunn/tenecs/ir"
)

func stdLibDynamicFunction(ref ir.NativeFunctionRef) *FunctionCodeAndImports {
	if ref.Package == "tenecs_json" && strings.HasPrefix(ref.Name, "jsonObject") {
		return tenecs_json_jsonObject(ref.Name)
	}
	return nil
}

func tenecs_json_jsonObject(functionNameString string) *FunctionCodeAndImports {
	numberString := strings.TrimPrefix(functionNameString, "jsonObject")
	numberOfFields, err := strconv.Atoi(numberString)
	if err != nil {
		panic("tenecs_json_jsonObject: invalid function name: " + functionNameString)
	}

	fromJsonCode := `func(generics any, input any) any {
	jsonString := input.(map[string]any)["value"].(string)
	var output map[string]json.RawMessage
	err := json.Unmarshal([]byte(jsonString), &output)
	if err != nil {
		return map[string]any{
			"$type": "tenecs.error.Error",
			"_message": map[string]any{
				"$type": "String",
				"value": "Could not parse object from " + jsonString,
			},
		}
	}
`

	for i := 0; i < numberOfFields; i++ {
		fromJsonCode += fmt.Sprintf(`
	i%dName := jsonConverterFieldI%d.(map[string]any)["_name"].(map[string]any)["value"].(string)
	i%dJsonRawMessage := output[i%dName]
	if i%dJsonRawMessage == nil {
		return map[string]any{
			"$type": "tenecs.error.Error",
			"_message": map[string]any{
				"$type": "String",
				"value": "Could not find object field \"" + i%dName + "\" in " + jsonString,
			},
		}
	}
	i%dJsonBytes, _ := json.Marshal(&i%dJsonRawMessage)
	i%dJsonString := map[string]any{
		"$type": "String",
		"value": string(i%dJsonBytes),
	}
	i%d := jsonConverterFieldI%d.(map[string]any)["_Converter"].(map[string]any)["_fromJson"].(func(any, any) any)([]any{}, i%dJsonString)
	i%dMap, isMap := i%d.(map[string]any)
	if isMap && i%dMap["$type"] == "tenecs.error.Error" {
		return map[string]any{
			"$type": "tenecs.error.Error",
			"_message": map[string]any{
				"$type": "String",
				"value": "Could not parse object field \"" + i%dName + "\": " + i%dMap["_message"].(map[string]any)["value"].(string),
			},
		}
	}
`, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i, i)
	}

	fromJsonCodeReturnArgs := ""
	fromJsonCodeReturnAnys := ""
	for i := 0; i < numberOfFields; i++ {
		if i > 0 {
			fromJsonCodeReturnArgs += ", "
			fromJsonCodeReturnAnys += ", "
		}
		fromJsonCodeReturnArgs += fmt.Sprintf("i%d", i)
		fromJsonCodeReturnAnys += "any"
	}

	fromJsonCode += fmt.Sprintf(`
	return build.(func(any, %s) any)([]any{}, %s)
}`, fromJsonCodeReturnAnys, fromJsonCodeReturnArgs)

	toJsonCode := `func(generics any, input any) any {
	output := map[string]string{}
`
	for i := 0; i < numberOfFields; i++ {
		toJsonCode += fmt.Sprintf(`
	fieldI%d := jsonConverterFieldI%d.(map[string]any)["_access"].(func(any, any) any)([]any{}, input)
	i%d := jsonConverterFieldI%d.(map[string]any)["_Converter"].(map[string]any)["_toJson"].(func(any, any) any)([]any{}, fieldI%d)
	output[jsonConverterFieldI%d.(map[string]any)["_name"].(map[string]any)["value"].(string)] = i%d.(map[string]any)["value"].(string)
`, i, i, i, i, i, i, i)
	}

	toJsonCode += `
	result := "{"
	i := 0
	outputKeysSorted := []string{}
	for k, _ := range output {
		outputKeysSorted = append(outputKeysSorted, k)
	}
	sort.Strings(outputKeysSorted)
	for _, k := range outputKeysSorted {
		v := output[k]
		nameBytes, _ := json.Marshal(k)
		result += string(nameBytes) + ":" + v
		i += 1
		if i < len(output) {
			result += ","
		}
	}
	result += "}"
	return map[string]any{
		"$type": "String",
		"value": result,
	}
}`

	jsonConverterArgs := ""
	for i := 0; i < numberOfFields; i++ {
		if i > 0 {
			jsonConverterArgs += ", "
		}
		jsonConverterArgs += fmt.Sprintf("jsonConverterFieldI%d any", i)
	}

	code := fmt.Sprintf(`func tenecs_json__jsonObject%d() any {
	return func(generics any, build any, %s) any {
		return map[string]any{
			"$type": "tenecs.json.JsonConverter",
			"_fromJson": %s,
			"_toJson": %s,
		}
	}
}`, numberOfFields, jsonConverterArgs, fromJsonCode, toJsonCode)

	return &FunctionCodeAndImports{
		Imports: []Import{
			"encoding/json",
			"sort",
		},
		Code: code,
	}
}
