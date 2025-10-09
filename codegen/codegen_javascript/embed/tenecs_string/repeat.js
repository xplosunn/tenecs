function tenecs_string__repeat() {
  return ((generics, str, count) => {
    return {
      "$type": "String",
      "value": str.value.repeat(count.value)
    }
  })
}
