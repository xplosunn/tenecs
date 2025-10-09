function tenecs_string__stripPrefix() {
  return ((generics, str, prefix) => {
    const result = str.value.startsWith(prefix.value) ? str.value.slice(prefix.value.length) : str.value
    return {"$type": "String", "value": result}
  })
}
