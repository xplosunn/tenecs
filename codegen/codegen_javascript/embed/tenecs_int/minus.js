function tenecs_int__minus() {
  return ((generics, a, b) => {
    return {"$type": "Int", "value": a.value - b.value}
  })
}
