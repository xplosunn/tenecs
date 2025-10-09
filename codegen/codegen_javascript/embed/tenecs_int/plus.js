function tenecs_int__plus() {
  return ((generics, a, b) => {
    return {"$type": "Int", "value": a.value + b.value}
  })
}
