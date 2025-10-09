function tenecs_list__append() {
  return ((generics, list, newElement) => {
    return {
      "$type": list["$type"],
      "value": [...list.value, newElement]
    }
  })
}
