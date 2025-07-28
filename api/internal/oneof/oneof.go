package oneof

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/karalef/tgot/api/internal"
)

// Type is an alias for reflect.Type.
type Type = reflect.Type

// NewMap makes a new types dictionary.
func NewMap[Type ~string](typed ...Value[Type]) Map[Type] {
	m := make(Map[Type], len(typed))
	for _, t := range typed {
		if typ := reflect.TypeOf(t); typ != nil {
			m[t.Type()] = typ
		}
	}
	return m
}

// Map represent the map of types.
type Map[Type ~string] map[Type]reflect.Type

func (m Map[Type]) TypeFor(id Type) reflect.Type { return m[id] }

// Oneof represents the type and types registry.
type Oneof[Type ~string] interface {
	~string

	// TypeFor returns the type for the given identifier.
	// Returns nil is the type is not registered.
	TypeFor(Type) reflect.Type
}

// IDType represents the model of the type identitfier.
type IDType interface {
	SetTypeID(string) IDType
	GetTypeID() string
}

// IDTypeType represents the model of the type identitfier with "type" field.
type IDTypeType struct {
	Type string `json:"type"`
}

func (i IDTypeType) SetTypeID(id string) IDType { i.Type = id; return i }
func (i IDTypeType) GetTypeID() string          { return i.Type }

// Value represent the value.
type Value[Type ~string] interface {
	Type() Type
}

// Object is a type that can be one of the given types.
type Object[Type Oneof[Type], ID IDType] struct {
	// Value contains the value (not a pointer).
	Value[Type]
}

func (o Object[Type, ID]) MarshalJSON() ([]byte, error) {
	var em ID
	id := em.SetTypeID(string(o.Value.Type()))
	return internal.MergeJSON(id, o.Value)
}

func (o *Object[Type, ID]) UnmarshalJSON(p []byte) error {
	var typeID ID
	err := json.Unmarshal(p, &typeID)
	if err != nil {
		return err
	}

	o.Value, err = New(Type(typeID.GetTypeID()), p)
	return err
}

// New allocates a new value of the specified type and unmarshals the data into
// it.
func New[Type Oneof[Type]](id Type, data []byte) (Value[Type], error) {
	rtyp := id.TypeFor(id)
	if rtyp == nil {
		return nil, errors.New("unknown identifier for oneof.Oneof: " + string(id))
	}
	ptr := reflect.New(rtyp)
	err := json.Unmarshal(data, ptr.Interface())
	if err != nil {
		return nil, err
	}

	if v, ok := ptr.Elem().Interface().(Value[Type]); ok {
		return v, nil
	}
	panic("invalid type for oneof.New")
}
