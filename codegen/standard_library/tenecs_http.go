package standard_library

import (
	"fmt"
)

func tenecs_http_newServer() Function {

	restHandlerGet := function(
		params("responseJsonSchema", "route", "handler"),
		body(`
if configurationErr != "" {
	return nil
}
if endpoints[route.(string)] == nil {
	endpoints[route.(string)] = map[string]func(http.ResponseWriter, *http.Request){}
}
if endpoints[route.(string)][http.MethodGet] != nil {
	configurationErr = "Configured multiple GET handlers for " + route.(string)
	return nil
}
endpoints[route.(string)][http.MethodGet] = func (w http.ResponseWriter, r *http.Request) {
	responseStatusRef := refCreator.(map[string]any)["new"].(func(any)any)(200)

	responseBody := handler.(func(any)any)(responseStatusRef)

	w.WriteHeader(responseStatusRef.(map[string]any)["get"].(func()any)().(int))
	fmt.Fprint(w, responseJsonSchema.(map[string]any)["toJson"].(func(any)any)(responseBody).(string))
}
`),
	)
	restHandlerPost := function(
		params("requestJsonSchema", "responseJsonSchema", "route", "handler"),
		body(`
if configurationErr != "" {
	return nil
}
if endpoints[route.(string)] == nil {
	endpoints[route.(string)] = map[string]func(http.ResponseWriter, *http.Request){}
}
if endpoints[route.(string)][http.MethodPost] != nil {
	configurationErr = "Configured multiple POST handlers for " + route.(string)
	return nil
}
endpoints[route.(string)][http.MethodPost] = func (w http.ResponseWriter, r *http.Request) {
	responseStatusRef := refCreator.(map[string]any)["new"].(func(any)any)(200)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Error parsing request")
		return		
	}
	requestBody := requestJsonSchema.(map[string]any)["fromJson"].(func(any)any)(string(bodyBytes))
	bodyMap, isMap := requestBody.(map[string]any)
	if isMap && bodyMap["$type"] == "Error" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Error parsing request: " + bodyMap["message"].(string))
		return
	}

	responseBody := handler.(func(any,any)any)(requestBody, responseStatusRef)

	w.WriteHeader(responseStatusRef.(map[string]any)["get"].(func()any)().(int))
	fmt.Fprint(w, responseJsonSchema.(map[string]any)["toJson"].(func(any)any)(responseBody).(string))
}
`),
	)

	runRestPostWithBody := function(
		params("route", "requestBody"),
		body(`
if configurationErr != "" {
	return "configuration error: " + configurationErr
}
if endpoints[route.(string)] == nil {
	return "Not found"
}
if endpoints[route.(string)][http.MethodPost] == nil {
	return "Not found"
}

responseRecorder := httptest.NewRecorder()

httpRequest, err := http.NewRequest("POST", route.(string), bytes.NewBuffer([]byte(requestBody.(string))))
if err != nil {
	return "error: " + err.Error()
}

endpoints[route.(string)][http.MethodPost](responseRecorder, httpRequest)

responseBytes, err := io.ReadAll(responseRecorder.Body)
if err != nil {
	return "error: " + err.Error()
}

return string(responseBytes)
`),
	)

	serveRun := function(
		params("address"),
		body(`
if configurationErr != "" {
	return map[string]any{
		"$type": "ServerError",
		"message": configurationErr,
	}
}
serverMux := http.NewServeMux()
for route, m := range endpoints {
	stableM := m
	serverMux.HandleFunc(route, func (w http.ResponseWriter, r *http.Request) {
		handler := stableM[r.Method]
		if handler != nil {
			handler(w, r)
		} else {
			w.WriteHeader(404)
		}
	})
}
err := http.ListenAndServe(address.(string), serverMux)
if err != nil {
	return map[string]any{
		"$type": "ServerError",
		"message": err.Error(),
	}
}`),
	)
	return function(
		imports("net/http", "net/http/httptest", "io", "bytes"),
		params("refCreator"),
		body(fmt.Sprintf(`
var configurationErr string
endpoints := map[string]map[string]func(http.ResponseWriter, *http.Request){}
return map[string]any{
	"restHandlerGet": %s,
	"restHandlerPost": %s,
	"runRestPostWithBody": %s,
	"__hiddenServe": %s,
}`, restHandlerGet.Code, restHandlerPost.Code, runRestPostWithBody.Code, serveRun.Code)),
	)
}
