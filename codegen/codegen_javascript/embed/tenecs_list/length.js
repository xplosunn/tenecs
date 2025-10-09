function tenecs_list__length() {
  return ((generics, list) => {
    return {"$type": "Int", "value": list.value.length}
  })
}
