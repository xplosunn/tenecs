function tenecs_string__stripSuffix() {
  return ((generics, str, suffix) => {
    const result = str.value.endsWith(suffix.value) ? str.value.substring(0, str.value.length - suffix.value.length) : str.value
    return {"$type": "String", "value": result}
  })
}
