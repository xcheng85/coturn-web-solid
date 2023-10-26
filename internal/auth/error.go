package auth

import (
	"fmt"
)

type UnauthorizedError struct {
	message string
}

func (r *UnauthorizedError) Error() string {
	return fmt.Sprintf(r.message)
}

// assert style in golang
func (s *UnauthorizedError) Is(target error) bool {
	t, ok := target.(*UnauthorizedError)
	if !ok {
		return false
	}
	return t.message == s.message
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		message,
	}
}
