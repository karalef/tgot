package oneof_test

import (
	"encoding/json"
	"testing"

	oneof "github.com/karalef/tgot/api/internal/oneof"
)

type oneofTestType string

const one oneofTestType = "one"

var registry = oneof.NewMap[oneofTestType](oneT{})

func (oneofTestType) TypeFor(t oneofTestType) oneof.Type {
	return registry.TypeFor(t)
}

type oneT struct {
	A int `json:"a"`
}

func (oneT) Type() oneofTestType { return one }

type TestType = oneof.Object[oneofTestType, oneof.IDTypeType]

func TestOneof(t *testing.T) {
	v := TestType{
		Value: oneT{A: 1},
	}
	j, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	checkType(t, j, one)
	var d TestType
	if err = json.Unmarshal(j, &d); err != nil {
		t.Fatal(err)
	}
	if d.Type() != one {
		t.Fatalf("expected %s, got %s", one, d.Type())
	}
	if d.Value.(oneT).A != 1 {
		t.Fatalf("expected str, got %s", d.Value.(oneT).A)
	}
}

func checkType(t *testing.T, j []byte, expected oneofTestType) {
	var v map[string]any
	if err := json.Unmarshal(j, &v); err != nil {
		t.Fatal(err)
	}
	typ, ok := v["source"]
	if !ok {
		t.Fatal("no type field", v)
	}
	if typ.(string) != string(expected) {
		t.Fatalf("expected %s, got %s %v", expected, typ, v)
	}
}
