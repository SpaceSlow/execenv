package flags

import "errors"

var (
	ErrIncorrectNetAddress = errors.New("need address in a form host:port")
	ErrIncorrectPort       = errors.New("error occurred when parsing an incorrect port. The port requires a decimal number in the range 0-65535")
)
