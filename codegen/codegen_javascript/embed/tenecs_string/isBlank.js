function tenecs_string__isBlank() {
  return ((generics, str) => {
    return {"$type": "Boolean", "value": str.value.trim().length === 0}
  })
}
