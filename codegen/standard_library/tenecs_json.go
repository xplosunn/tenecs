package standard_library

func tenecs_json_toJson() Function {
	return function(
		imports("encoding/json"),
		params("input"),
		body(`result, _ := json.Marshal(input)
return string(result)`),
	)
}
