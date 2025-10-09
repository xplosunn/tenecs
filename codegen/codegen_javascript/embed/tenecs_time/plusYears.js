function tenecs_time__plusYears() {
  return ((generics, _date, _years) => {
    const date = _date
    const years = _years
    return {
      "$type": "tenecs.time.Date",
      "_year": {"$type": "Int", "value": date._year.value + years.value},
      "_month": date._month,
      "_day": date._day
    }
  })
}
