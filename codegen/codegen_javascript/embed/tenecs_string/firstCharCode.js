function tenecs_string__firstCharCode() {
  return ((generics, str) => {
    const runes = Array.from(str.value)
    if (runes.length > 0) {
      return {"$type": "Int", "value": runes[0].charCodeAt(0)}
    } else {
      return {"$type": "Int", "value": -1}
    }
  })
}
