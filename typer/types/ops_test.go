package types

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestVariableTypeContainedIn(t *testing.T) {
	a := &TypeArgument{Name: "A"}
	b := &TypeArgument{Name: "B"}
	aOrB := &OrVariableType{Elements: []VariableType{a, b}}
	assert.True(t, VariableTypeContainedIn(a, a))
	assert.True(t, VariableTypeContainedIn(a, aOrB))
	assert.True(t, VariableTypeContainedIn(b, aOrB))
	assert.False(t, VariableTypeContainedIn(aOrB, a))
	assert.False(t, VariableTypeContainedIn(aOrB, b))
}

func TestVariableTypeEq(t *testing.T) {
	a := &TypeArgument{Name: "A"}
	b := &TypeArgument{Name: "B"}
	aOrB := &OrVariableType{Elements: []VariableType{a, b}}
	assert.True(t, VariableTypeEq(a, a))
	assert.True(t, VariableTypeEq(b, b))

	assert.False(t, VariableTypeEq(a, b))
	assert.False(t, VariableTypeEq(b, a))

	assert.False(t, VariableTypeEq(a, aOrB))
	assert.False(t, VariableTypeEq(aOrB, a))

	assert.False(t, VariableTypeEq(b, aOrB))
	assert.False(t, VariableTypeEq(aOrB, b))
}
