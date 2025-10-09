function tenecs_int__times() {
  return ((generics, a, b) => {
    return {"$type": "Int", "value": a.value * b.value}
  })
}
