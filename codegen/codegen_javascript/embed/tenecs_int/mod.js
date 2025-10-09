function tenecs_int__mod() {
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
    return {"$type": "Int", "value": a.value % b.value}
  })
}
