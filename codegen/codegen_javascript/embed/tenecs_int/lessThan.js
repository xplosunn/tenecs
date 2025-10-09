function tenecs_int__lessThan() {
  return ((generics, a, b) => {
    return {"$type": "Boolean", "value": a.value < b.value}
  })
}
