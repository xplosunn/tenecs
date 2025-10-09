function tenecs_int__negate() {
  return ((generics, a) => {
    return {"$type": "Int", "value": -a.value}
  })
}
