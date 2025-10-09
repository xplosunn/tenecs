function tenecs_boolean__not() {
  return ((generics, a) => {
    return {"$type": "Boolean", "value": !a.value}
  })
}
