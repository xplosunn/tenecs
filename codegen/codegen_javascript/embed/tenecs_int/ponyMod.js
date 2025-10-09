function tenecs_int__ponyMod() {
  return ((generics, a, b) => {
    if (b.value === 0) {
      return {"$type": "Int", "value": 0}
    }
    return {"$type": "Int", "value": a.value % b.value}
  })
}
