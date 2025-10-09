function tenecs_string__padRight() {
  return ((generics, str, length, padChar) => {
    const padCount = length.value - Array.from(str.value).length
    if (padCount <= 0) {
      return str
    }
    return {
      "$type": "String",
      "value": str.value + padChar.value.repeat(padCount)
    }
  })
}
