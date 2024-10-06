package apperrors

import "fmt"

type ErrorRequiringRestart struct {
	Message string
}

func (err ErrorRequiringRestart) Error() string {
	return fmt.Sprintf("error requiring restart: %v", err.Message)
}
