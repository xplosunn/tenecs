package binding_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/typer/binding"
	"testing"
)

func TestSetScopedIfAbsent(t *testing.T) {
	m := binding.NewTwoLevelMap[string, string, string]()

	v, ok := m.Get("a", "b")
	assert.False(t, ok)
	assert.Zero(t, v)

	m, ok = m.SetScopedIfAbsent("a", "b", "v")
	assert.True(t, ok)
	v, ok = m.Get("a", "b")
	assert.True(t, ok)
	assert.Equal(t, "v", v)

	m, ok = m.SetScopedIfAbsent("a", "b", "changed")
	assert.False(t, ok)
	v, ok = m.Get("a", "b")
	assert.True(t, ok)
	assert.Equal(t, "v", v)

	m, ok = m.SetGlobalIfAbsent("b", "changed")
	assert.False(t, ok)
	v, ok = m.Get("a", "b")
	assert.True(t, ok)
	assert.Equal(t, "v", v)
}

func TestSetGlobalIfAbsent(t *testing.T) {
	m := binding.NewTwoLevelMap[string, string, string]()

	v, ok := m.Get("a", "b")
	assert.False(t, ok)
	assert.Zero(t, v)

	m, ok = m.SetGlobalIfAbsent("b", "v")
	assert.True(t, ok)
	v, ok = m.Get("a", "b")
	assert.True(t, ok)
	assert.Equal(t, "v", v)
	v, ok = m.Get("b", "b")
	assert.True(t, ok)
	assert.Equal(t, "v", v)

	m, ok = m.SetScopedIfAbsent("a", "b", "changed")
	assert.False(t, ok)
	v, ok = m.Get("a", "b")
	assert.True(t, ok)
	assert.Equal(t, "v", v)

	m, ok = m.SetGlobalIfAbsent("b", "changed")
	assert.False(t, ok)
	v, ok = m.Get("a", "b")
	assert.True(t, ok)
	assert.Equal(t, "v", v)
}
