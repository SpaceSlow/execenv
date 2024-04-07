package flags

type IncorrectNetAddressError struct{}

func (e *IncorrectNetAddressError) Error() string {
	return "need address in a form host:port"
}
