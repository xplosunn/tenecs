function tenecs_list__atIndexSet() {
  return ((generics, list, index, setTo) => {
    const l = list.value
    const idx = index.value
    if (idx >= 0 && l.length > idx) {
      const result = [...l]
      result[idx] = setTo

      const listTypeArgument = generics[0]
      const listOf = typeNameCollapseRecursive(listTypeArgument)
      return {
        "$type": "List<" + listOf + ">",
        "value": result
      }
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
