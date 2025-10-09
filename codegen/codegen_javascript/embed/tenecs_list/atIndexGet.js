function tenecs_list__atIndexGet() {
  return ((generics, list, index) => {
    const l = list.value
    const idx = index.value
    if (idx >= 0 && l.length > idx) {
      return l[idx]
    }
    return {
      "$type": "tenecs.error.Error",
      "_message": {
        "$type": "String",
        "value": "Out of bounds"
      }
    }
  })
}
