function tenecs_list__first() {
  return ((generics, list) => {
    if (list.value.length > 0) {
      return list.value[0]
    }
    return {"$type": "Void", "value": null}
  })
}
