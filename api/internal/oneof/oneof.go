package oneof

import (
	"encoding/json"
	"errors"

	"github.com/karalef/tgot/api/internal"
)

// Oneof represent the types registry which can allocate object.
type Oneof[Type ~string] interface {
	New(Type) (Value[Type], bool)
}

// Map represent the map of values.
type Map[Type ~string] map[Type]Value[Type]

func (m Map[Type]) New(id Type) (Value[Type], bool) {
	v, ok := m[id]
	return v, ok
}

// Value represent the value.
type Value[Type ~string] interface {
	Type() Type
}

type oneof[Type ~string] struct {
	Type Type `json:"type"`
}

// Object is a type that can be one of the given types.
type Object[Type ~string, New Oneof[Type]] struct {
	Value[Type]
}

func (o Object[Type, New]) MarshalJSON() ([]byte, error) {
	return internal.MergeJSON(oneof[Type]{
		Type: o.Value.Type(),
	}, o.Value)
}

func (o *Object[Type, New]) UnmarshalJSON(p []byte) error {
	var typ oneof[Type]
	if err := json.Unmarshal(p, &typ); err != nil {
		return err
	}

	var allocator New
	var ok bool
	o.Value, ok = allocator.New(typ.Type)
	if !ok {
		return errors.New("invalid identifier for oneof.Object")
	}
	return json.Unmarshal(p, &o.Value)
}
