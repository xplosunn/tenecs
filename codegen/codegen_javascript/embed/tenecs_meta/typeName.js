function tenecs_meta__typeName() {
  return ((generics, a) => {
    return {"$type": "String", "value": a["$type"]}
  })
}
