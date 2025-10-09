function tenecs_list__Break() {
  return ((generics, value) => {
    return {
      "$type": "tenecs.list.Break<" + typeNameCollapseRecursive(generics[0]) + ">",
      "_value": value
    }
  })
}
