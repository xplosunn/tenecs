function tenecs_list__flatMap() {
  return ((generics, list, f) => {
    const result = []
    for (const elem of list.value) {
      result.push(...f([], elem).value)
    }
    return {
      "$type": "List<" + typeNameCollapseRecursive(generics[1]) + ">",
      "value": result
    }
  })
}
