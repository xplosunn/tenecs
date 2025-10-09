function tenecs_list__flatten() {
  return ((generics, list) => {
    const result = []
    for (const innerList of list.value) {
      result.push(...innerList.value)
    }
    return {
      "$type": "List<" + typeNameCollapseRecursive(generics[0]) + ">",
      "value": result
    }
  })
}
