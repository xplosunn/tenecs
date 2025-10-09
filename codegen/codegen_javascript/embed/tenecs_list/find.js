function tenecs_list__find() {
  return ((generics, list, f) => {
    for (const elem of list.value) {
      const maybeReturn = f([], elem)
      if (maybeReturn["$type"] !== "Void") {
        return maybeReturn
      }
    }
    return {"$type": "Void", "value": null}
  })
}
