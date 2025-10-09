function tenecs_string__reverse() {
  return ((generics, str) => {
    const runes = Array.from(str.value)
    for (let i = 0, j = runes.length - 1; i < j; i++, j--) {
      [runes[i], runes[j]] = [runes[j], runes[i]]
    }
    return {"$type": "String", "value": runes.join("")}
  })
}
