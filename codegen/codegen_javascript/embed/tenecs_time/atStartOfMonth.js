function tenecs_time__atStartOfMonth() {
  return ((generics, _date) => {
    const date = _date
    return {
      "$type": "tenecs.time.Date",
      "_year": date._year,
      "_month": date._month,
      "_day": {"$type": "Int", "value": 1}
    }
  })
}
