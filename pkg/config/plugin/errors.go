package plugin

// ErrorConfiguring is the error returned when plugin configuring fails
type ErrorConfiguring struct {
	field string
}

// NewErrorConfiguring creates new ErrorConfiguring instance
func NewErrorConfiguring(field string) *ErrorConfiguring {
	return &ErrorConfiguring{field: field}
}

// Field returns configuring error field name
func (e *ErrorConfiguring) Field() string {
	return e.field
}

func (e *ErrorConfiguring) Error() string {
	return "Failed to configure plugin, required field missing"
}
