function tenecs_list__appendAll() {
  return ((generics, list, newElements) => {
    return {
      "$type": list["$type"],
      "value": [...list.value, ...newElements.value]
    }
  })
}
