function tenecs_string__startsWith() {
  return ((generics, str, prefix) => {
    return {"$type": "Boolean", "value": str.value.startsWith(prefix.value)}
  })
}
