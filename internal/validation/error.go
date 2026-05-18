// Package validation provides domain-neutral validation error types.
// Each product defines its own ErrorCode constants typed as validation.ErrorCode.
package validation

import "fmt"

type ErrorCode string

type Error struct {
	Code    ErrorCode
	Field   string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

type RowError struct {
	RowIndex int
	Errors   []Error
}

func (e RowError) Error() string {
	codes := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		codes = append(codes, string(err.Code))
	}
	return fmt.Sprintf("[%d] %s", e.RowIndex, codes)
}
