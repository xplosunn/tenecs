function tenecs_string__padLeft() {
  return ((generics, str, length, padChar) => {
    const padCount = length.value - Array.from(str.value).length
    if (padCount <= 0) {
      return str
    }
    return {
      "$type": "String",
      "value": padChar.value.repeat(padCount) + str.value
    }
  })
}
