function tenecs_list__map() {
  return ((generics, list, f) => {
    const result = list.value.map(elem => f([], elem))
    
    const listTypeArgument = generics[1]
    const listOf = typeNameCollapseRecursive(listTypeArgument)
    return {
      "$type": "List<" + listOf + ">",
      "value": result
    }
  })
}
