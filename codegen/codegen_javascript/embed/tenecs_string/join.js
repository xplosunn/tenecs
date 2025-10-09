function tenecs_string__join() {
  return ((generics, a, b) => {
    return {"$type": "String", "value": a.value + b.value}
  })
}
