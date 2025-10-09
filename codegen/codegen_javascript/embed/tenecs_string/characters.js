function tenecs_string__characters() {
  return ((generics, str) => {
    const chars = Array.from(str.value).map(char => ({
      "$type": "String",
      "value": char
    }))
    return {"$type": "List<String>", "value": chars}
  })
}
