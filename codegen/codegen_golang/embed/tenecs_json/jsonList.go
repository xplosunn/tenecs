package embed

import (
	"encoding/json"
	"strings"
)

// embed start

func tenecs_json__jsonList() any {
	return func(topGenerics any, of any) any {
		return map[string]any{
			"$type": "tenecs.json.JsonConverter",
			"_fromJson": func(generics any, input any) any {
				jsonString := input.(map[string]any)["value"].(string)
				var output []json.RawMessage
				err := json.Unmarshal([]byte(jsonString), &output)
				if err != nil {
					return map[string]any{
						"$type": "tenecs.error.Error",
						"_message": map[string]any{
							"$type": "String",
							"value": "Could not parse List from " + jsonString,
						},
					}
				}

				listTypeArgument := topGenerics.([]any)[0]
				var listOf string = typeNameCollapseRecursive(listTypeArgument)

				ofParse := of.(map[string]any)["_fromJson"].(func(any, any) any)
				outputList := []any{}
				for _, elem := range output {
					elemJsonBytes, _ := json.Marshal(&elem)
					result := ofParse([]any{}, map[string]any{
						"$type": "String",
						"value": string(elemJsonBytes),
					})
					resultMap, isMap := result.(map[string]any)
					if isMap && resultMap["$type"] == "tenecs.error.Error" {
						return resultMap
					}
					outputList = append(outputList, result)
				}
				return map[string]any{
					"$type": "List<" + listOf + ">",
					"value": outputList,
				}
			},
			"_toJson": func(generics any, input any) any {
				results := []string{}
				ofToJson := of.(map[string]any)["_toJson"].(func(any, any) any)
				for _, elem := range input.(map[string]any)["value"].([]any) {
					result := ofToJson([]any{}, elem)
					results = append(results, result.(map[string]any)["value"].(string))
				}
				return map[string]any{
					"$type": "String",
					"value": "[" + strings.Join(results, ",") + "]",
				}
			},
		}
	}
}
