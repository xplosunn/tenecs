function tenecs_string__toLowerCase() {
  return ((generics, str) => {
    return {"$type": "String", "value": str.value.toLowerCase()}
  })
}
