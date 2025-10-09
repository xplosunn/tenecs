function tenecs_list__filter() {
  return ((generics, list, keep) => {
    const result = []
    for (const elem of list.value) {
      if (keep([], elem).value) {
        result.push(elem)
      }
    }
    return {
      "$type": list["$type"],
      "value": result
    }
  })
}
