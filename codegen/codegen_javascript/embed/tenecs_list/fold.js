function tenecs_list__fold() {
  return ((generics, list, zero, f) => {
    let result = zero
    for (const elem of list.value) {
      result = f([], result, elem)
    }
    return result
  })
}
