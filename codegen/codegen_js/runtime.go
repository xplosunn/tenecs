package codegen_js

func runtimeRefCreator() string {
	return `({
  "new": (value) => {
    let ref = value
    return ({
      "$type": "Ref",
      "get": () => {
        return ref
      },
      "set": (value) => {
        ref = value
        return null
      },
      "modify": (f) => {
        ref = f(ref)
        return null
      }
    })
  }
})`
}
