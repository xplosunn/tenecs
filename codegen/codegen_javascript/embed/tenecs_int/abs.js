function tenecs_int__abs() {
  return ((generics, a) => {
    return {"$type": "Int", "value": Math.abs(a.value)}
  })
}
