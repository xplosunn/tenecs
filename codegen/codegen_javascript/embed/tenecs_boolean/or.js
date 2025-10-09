function tenecs_boolean__or() {
  return ((generics, a, b) => {
    return {"$type": "Boolean", "value": a.value || b().value}
  })
}
