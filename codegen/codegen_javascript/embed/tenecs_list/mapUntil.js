function tenecs_list__mapUntil() {
  return ((generics, list, f) => {
    const result = []
    for (const elem of list.value) {
      const maybeBreak = f([], elem)
      if (maybeBreak["$type"].startsWith("tenecs.list.Break<")) {
        const left = generics[2]
        const right = maybeBreak._value["$type"]
        if (Array.isArray(left) && left.includes(right) || (!Array.isArray(left) && left === right)) {
          return maybeBreak._value
        } else {
          result.push(maybeBreak)
        }
      } else {
        result.push(maybeBreak)
      }
    }
    return {
      "$type": "List<" + typeNameCollapseRecursive(generics[1]) + ">",
      "value": result
    }
  })
}
