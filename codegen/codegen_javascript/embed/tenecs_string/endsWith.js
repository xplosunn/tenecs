function tenecs_string__endsWith() {
  return ((generics, str, suffix) => {
    return {"$type": "Boolean", "value": str.value.endsWith(suffix.value)}
  })
}
