package standard_library

import (
	"fmt"
	"github.com/xplosunn/tenecs/typer/standard_library"
)

func tenecs_http_newServer() Function {
	restHandlerGet := `(responseJsonSchema, route, handler) => {
  if (configurationErr) {
    return null
  }
  if (!endpoints[route]) {
    endpoints[route] = {}
  }
  if (endpoints[route]["GET"]) {
    configurationErr = "Configured multiple GET handlers for " + route
    return null
  }
  endpoints[route]["GET"] = () => {
    let responseStatusRef = refCreator["new"](200)
    let responseBody = handler(responseStatusRef)
    let responseJsonBody = responseJsonSchema["toJson"](responseBody)
    return responseJsonBody
  }
}
`
	restHandlerPost := `(requestJsonSchema, responseJsonSchema, route, handler) => {
  if (configurationErr) {
    return null
  }
  if (!endpoints[route]) {
    endpoints[route] = {}
  }
  if (endpoints[route]["POST"]) {
    configurationErr = "Configured multiple GET handlers for " + route
    return null
  }
  endpoints[route]["POST"] = (bodyString) => {
    let responseStatusRef = refCreator["new"](200)
	let requestBody = requestJsonSchema["fromJson"](bodyString)
    if (requestBody && requestBody["$type"] && requestBody["$type"] == "Error") {
      responseStatusRef["set"](400)
      return "Error parsing request: " + requestBody["message"]
    }

    let responseBody = handler(requestBody, responseStatusRef)
    let responseJsonBody = responseJsonSchema["toJson"](responseBody)
    return responseJsonBody
  }
}
`
	runRestPostWithBody := `(route, body) => {
  if (configurationErr) {
    return "configuration error: " + configurationErr
  }
  if (!endpoints[route]) {
    return "Not found"
  }
  if (!endpoints[route]["POST"]) {
    return "Not found"
  }
  return endpoints[route]["POST"](body)
}
`
	return function(
		params("refCreator"),
		body(fmt.Sprintf(`
let configurationErr = null
const endpoints = {}

return ({
	"restHandlerGet": %s,
	"restHandlerPost": %s,
	"runRestPostWithBody": %s
})
`, restHandlerGet, restHandlerPost, runRestPostWithBody)),
	)
}

func tenecs_http_ServerError() Function {
	return structFunction(standard_library.Tenecs_http_ServerError)
}

func tenecs_http_RuntimeServer() Function {
	return structFunction(standard_library.Tenecs_http_RuntimeServer)
}

func tenecs_http_Server() Function {
	return structFunction(standard_library.Tenecs_http_Server)
}
