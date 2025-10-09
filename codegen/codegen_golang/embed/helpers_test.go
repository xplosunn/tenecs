package embed

import (
	"testing"
)

func TestTypeNameCollapseRecursive(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name:     "single string",
			input:    "string",
			expected: "string",
		},
		{
			name:     "single element array",
			input:    []any{"int"},
			expected: "int",
		},
		{
			name:     "two element array",
			input:    []any{"string", "int"},
			expected: "int | string",
		},
		{
			name:     "three element array",
			input:    []any{"bool", "string", "int"},
			expected: "bool | int | string",
		},
		{
			name:     "three element array with repeated",
			input:    []any{"bool", "string", "bool"},
			expected: "bool | string",
		},
		{
			name:     "nested array with strings",
			input:    []any{"string", []any{"int", "bool"}},
			expected: "bool | int | string",
		},
		{
			name:     "deeply nested arrays",
			input:    []any{[]any{"b", "a"}, []any{"d", "c"}},
			expected: "a | b | c | d",
		},
		{
			name:     "complex nested structure",
			input:    []any{"string", []any{"int", []any{"bool", "float"}}},
			expected: "bool | float | int | string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := typeNameCollapseRecursive(tt.input)
			if result != tt.expected {
				t.Errorf("typeNameCollapseRecursive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTypeNameCollapseRecursivePanic(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			name:  "invalid input - int",
			input: 42,
		},
		{
			name:  "invalid input - bool",
			input: true,
		},
		{
			name:  "invalid input - nil",
			input: nil,
		},
		{
			name:  "invalid input - empty array",
			input: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("typeNameCollapseRecursive() should have panicked for input %v", tt.input)
				}
			}()
			typeNameCollapseRecursive(tt.input)
		})
	}
}
