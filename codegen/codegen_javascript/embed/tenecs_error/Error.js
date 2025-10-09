function tenecs_error__Error() {
  return ((generics, _message) => {
    return {
      "$type": "tenecs.error.Error",
      "_message": _message
    }
  })
}
