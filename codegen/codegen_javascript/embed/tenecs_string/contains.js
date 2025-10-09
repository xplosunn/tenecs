function tenecs_string__contains() {
  return ((generics, str, subStr) => {
    return {"$type": "Boolean", "value": str.value.includes(subStr.value)}
  })
}
