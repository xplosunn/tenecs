function tenecs_string__length() {
  return ((generics, str) => {
    return {"$type": "Int", "value": str.value.length}
  })
}
