function tenecs_list__mapNotNull() {
  return ((generics, list, f) => {
    const result = []
    for (const elem of list.value) {
      const maybeNull = f([], elem)
      if (maybeNull["$type"] !== "Void") {
        result.push(maybeNull)
      }
    }
    return {
      "$type": "List<" + typeNameCollapseRecursive(generics[1]) + ">",
      "value": result
    }
  })
}
