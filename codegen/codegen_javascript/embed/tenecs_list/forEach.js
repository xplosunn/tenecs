function tenecs_list__forEach() {
  return ((generics, list, f) => {
    for (const elem of list.value) {
      f([], elem)
    }
    return {
      "$type": "Void",
      "value": null
    }
  })
}
