package standard_library

import "fmt"

func tenecs_http_newServer() Function {
	//imports from this function are ignored, they should be added below
	restHandler := function(
		params("fromJson", "route", "handler"),
		body(`responseStatusRef := refCreator.(map[string]any)["new"].(func(any)any)(200)
serverMux.HandleFunc(route.(string), func (w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Error parsing request")
		return		
	}
	requestBody := fromJson.(map[string]any)["parse"].(func(any)any)(string(bodyBytes))
	bodyMap, isMap := requestBody.(map[string]any)
	if isMap && bodyMap["$type"] == "JsonError" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Error parsing request: " + bodyMap["message"].(string))
		return
	}

	responseBody := handler.(func(any,any)any)(requestBody, responseStatusRef)

	w.WriteHeader(responseStatusRef.(map[string]any)["get"].(func()any)().(int))
	fmt.Fprint(w, toJson(responseBody).(string))
})
`),
	)
	//imports from this function are ignored, they should be added below
	serve := function(
		params("address", "blocker"),
		body(`
err := http.ListenAndServe(address.(string), serverMux)
if err != nil {
	return map[string]any{
		"$type": "ServerError",
		"message": err.Error(),
	}
}`),
	)

	toJsonFunction := tenecs_json_toJson()

	return function(
		imports(append(toJsonFunction.Imports, "net/http", "io")...),
		params("refCreator"),
		body(fmt.Sprintf(`serverMux := http.NewServeMux()
toJson := %s
return map[string]any{
	"restHandler": %s,
	"serve": %s,
}`, toJsonFunction.Code, restHandler.Code, serve.Code)),
	)
}
