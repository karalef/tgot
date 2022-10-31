package api

import (
	"fmt"
	"strings"

	"github.com/karalef/tgot/api/tg"
)

func makeError[T error](method string, d Data, err T) baseError[T] {
	return baseError[T]{
		Method: method,
		Data:   d,
		Err:    err,
	}
}

type baseError[T error] struct {
	Method string
	Data   Data
	Err    T
}

func (e baseError[T]) Error() string {
	return fmt.Sprintf("%s\n%s %s", e.Err.Error(), e.Method, e.formatData())
}

func (e baseError[T]) Unwrap() error {
	return e.Err
}

func (e baseError[T]) formatData() string {
	if len(e.Data.Files) == 0 {
		return e.Data.Params.Encode()
	}
	var sb strings.Builder
	sb.WriteString(e.Data.Params.Encode())
	for _, f := range e.Data.Files {
		sb.WriteByte('\n')
		name, r := f.FileData()
		sb.WriteString("(file) " + f.Field + ": " + name)
		if r != nil {
			sb.WriteString(" (upload data)")
		}
	}
	return sb.String()
}

// Error represents a telegram api error and also contains method and data.
type Error struct {
	baseError[tg.APIError]
}

// Is implements errors.Is interface.
func (e Error) Is(err error) bool {
	if tge, ok := err.(tg.APIError); ok {
		return e.Err.Code != tge.Code
	}
	return false
}

// HTTPError represents http error.
type HTTPError struct {
	baseError[error]
}

// JSONError represents JSON error.
type JSONError struct {
	baseError[error]
	Response []byte
}

func (e *JSONError) Error() string {
	return e.baseError.Error() + "\nresponse:\n" + string(e.Response)
}
