package api

import (
	"fmt"
	"strings"

	"github.com/karalef/tgot/api/tg"
)

func makeError[T error](method string, d *Data, err T) baseError[T] {
	e := baseError[T]{
		Method: method,
		Err:    err,
	}
	if d != nil {
		e.Params = make(map[string]string, len(d.Params))
		for k, v := range d.Params {
			e.Params[k] = v
		}
		e.Files = make(map[string]string, len(d.Files))
		for field, f := range d.Files {
			s, _ := f.FileData()
			e.Files[field] = s
		}
	}
	return e
}

type baseError[T error] struct {
	Method string
	Params map[string]string
	Files  map[string]string
	Err    T
}

func (e baseError[T]) Error() string {
	return fmt.Sprintf("%s\n%s %s", e.Err.Error(), e.Method, e.formatData())
}

func (e baseError[T]) Unwrap() error {
	return e.Err
}

func (e baseError[T]) formatData() string {
	var sb strings.Builder
	for k, v := range e.Params {
		sb.WriteString(k + "=" + v + " ")
	}
	for field, f := range e.Files {
		sb.WriteByte('\n')
		sb.WriteString("[file] " + field + ": " + f)
	}
	return sb.String()
}

// Error represents a telegram api error and also contains method and data.
type Error struct {
	baseError[*tg.APIError]
}

// Is implements errors.Is interface.
func (e Error) Is(err error) bool {
	if tge, ok := err.(*tg.APIError); ok {
		return e.Err.Code == tge.Code
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

func (e JSONError) Error() string {
	return e.baseError.Error() + "\nresponse:\n" + string(e.Response)
}
