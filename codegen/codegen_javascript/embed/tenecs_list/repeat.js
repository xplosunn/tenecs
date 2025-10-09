function tenecs_list__repeat() {
  return ((generics, elem, times) => {
    const result = []
    for (let i = 0; i < times.value; i++) {
      result.push(elem)
    }
    return {
      "$type": "List<" + typeNameCollapseRecursive(generics[0]) + ">",
      "value": result
    }
  })
}
