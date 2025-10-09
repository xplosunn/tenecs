function tenecs_compare__eq() {
  return ((generics, a, b) => {
    if (typeof a === "function" || typeof b === "function") {
      return {"$type": "Boolean", "value": false}
    }
    return {"$type": "Boolean", "value": JSON.stringify(a) === JSON.stringify(b)}
  })
}
