function tenecs_boolean__and() {
  return ((generics, a, b) => {
    return {"$type": "Boolean", "value": a.value && b().value}
  })
}
