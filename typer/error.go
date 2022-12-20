package typer

import "fmt"

type TypecheckError struct {
	Message string
}

func (t TypecheckError) Error() string {
	return t.Message
}

func TypeCheckErrorf(format string, a ...any) TypecheckError {
	return TypecheckError{fmt.Sprintf(format, a...)}
}

func PtrTypeCheckErrorf(format string, a ...any) *TypecheckError {
	return &TypecheckError{fmt.Sprintf(format, a...)}
}
