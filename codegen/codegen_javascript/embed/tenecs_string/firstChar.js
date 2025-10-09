function tenecs_string__firstChar() {
  return ((generics, str) => {
    const runes = Array.from(str.value)
    if (runes.length > 0) {
      return {"$type": "String", "value": runes[0]}
    } else {
      return {"$type": "String", "value": ""}
    }
  })
}
