function tenecs_time__plusDays() {
  return ((generics, _date, _days) => {
    const date = _date
    const days = _days

    // Create a JavaScript Date object
    const jsDate = new Date(
      date._year.value,
      date._month.value - 1, // JavaScript months are 0-indexed
      date._day.value
    )

    // Add days
    jsDate.setDate(jsDate.getDate() + days.value)

    return {
      "$type": "tenecs.time.Date",
      "_year": {"$type": "Int", "value": jsDate.getFullYear()},
      "_month": {"$type": "Int", "value": jsDate.getMonth() + 1}, // Convert back to 1-indexed
      "_day": {"$type": "Int", "value": jsDate.getDate()}
    }
  })
}
