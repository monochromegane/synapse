package synapse

import "fmt"

type UnexpectedTypeError struct{}

func (e UnexpectedTypeError) Error() string {
	return "unexpected type from module symbol"
}

type NoRouteMatchError struct {
	Path string
}

func (e NoRouteMatchError) Error() string {
	return fmt.Sprintf("No route matches %s", e.Path)
}
