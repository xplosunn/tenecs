function tenecs_string__trimRight() {
  return ((generics, str) => {
    return {"$type": "String", "value": str.value.replace(/[\s\t\n\r]+$/, "")}
  })
}
