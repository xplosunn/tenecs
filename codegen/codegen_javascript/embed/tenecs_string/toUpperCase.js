function tenecs_string__toUpperCase() {
  return ((generics, str) => {
    return {"$type": "String", "value": str.value.toUpperCase()}
  })
}
