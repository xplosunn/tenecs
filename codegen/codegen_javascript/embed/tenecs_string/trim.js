function tenecs_string__trim() {
  return ((generics, str) => {
    return {"$type": "String", "value": str.value.trim()}
  })
}
