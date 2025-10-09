function tenecs_string__isEmpty() {
  return ((generics, str) => {
    return {"$type": "Boolean", "value": str.value.length === 0}
  })
}
