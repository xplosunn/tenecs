function tenecs_int__div() {
  return ((generics, a, b) => {
    if (b.value === 0) {
      return {
        "$type": "tenecs.error.Error",
        "_message": {
          "$type": "String",
          "value": "Division by zero"
        }
      }
    }
    const resultUnrounded = a.value / b.value
    if (resultUnrounded > 0) {
      return {"$type": "Int", "value": Math.floor(resultUnrounded)}
    } else {
      return {"$type": "Int", "value": Math.ceil(resultUnrounded)}
    }
  })
}
