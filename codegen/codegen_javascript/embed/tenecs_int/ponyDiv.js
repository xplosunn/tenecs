function tenecs_int__ponyDiv() {
  return ((generics, a, b) => {
    if (b.value === 0) {
      return {"$type": "Int", "value": 0}
    }
    const resultUnrounded = a.value / b.value
    if (resultUnrounded > 0) {
      return {"$type": "Int", "value": Math.floor(resultUnrounded)}
    } else {
      return {"$type": "Int", "value": Math.ceil(resultUnrounded)}
    }
  })
}
