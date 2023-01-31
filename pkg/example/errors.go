package example

type ErrorTypeExample struct {
	message string
}

func NewErrorTypeExample() ErrorTypeExample {
	return ErrorTypeExample{
		message: "type error",
	}
}

func (e ErrorTypeExample) Error() string {
	return e.message
}
