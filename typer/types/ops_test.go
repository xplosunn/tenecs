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

func TestVariableTypeEqFunction(t *testing.T) {
	aToA := &Function{
		Generics: []string{"A"},
		Arguments: []FunctionArgument{
			FunctionArgument{
				Name:         "a",
				VariableType: &TypeArgument{Name: "A"},
			},
		},
		ReturnType: &TypeArgument{Name: "A"},
	}
	bToB := &Function{
		Generics: []string{"B"},
		Arguments: []FunctionArgument{
			FunctionArgument{
				Name:         "b",
				VariableType: &TypeArgument{Name: "B"},
			},
		},
		ReturnType: &TypeArgument{Name: "B"},
	}

	assert.True(t, VariableTypeEq(aToA, aToA))
	assert.True(t, VariableTypeEq(bToB, bToB))
	assert.True(t, VariableTypeEq(aToA, bToB))

	aToBoolean := &Function{
		Generics: []string{"A"},
		Arguments: []FunctionArgument{
			FunctionArgument{
				Name:         "a",
				VariableType: &TypeArgument{Name: "A"},
			},
		},
		ReturnType: Boolean(),
	}
	bToBoolean := &Function{
		Generics: []string{"B"},
		Arguments: []FunctionArgument{
			FunctionArgument{
				Name:         "b",
				VariableType: &TypeArgument{Name: "B"},
			},
		},
		ReturnType: Boolean(),
	}
	assert.True(t, VariableTypeEq(aToBoolean, aToBoolean))
	assert.True(t, VariableTypeEq(bToBoolean, bToBoolean))
	assert.True(t, VariableTypeEq(aToBoolean, bToBoolean))

	assert.False(t, VariableTypeEq(aToA, aToBoolean))
}
