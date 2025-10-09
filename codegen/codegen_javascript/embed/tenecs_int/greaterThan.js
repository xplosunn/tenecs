function tenecs_int__greaterThan() {
  return ((generics, a, b) => {
    return {"$type": "Boolean", "value": a.value > b.value}
  })
}
