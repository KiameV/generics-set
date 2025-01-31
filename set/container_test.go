package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	_tester     interface{ Foo() }
	_testerImpl struct{ _tester }
)

func TestNew(t *testing.T) {
	c := New()
	if assert.NotNil(t, c) {
		assert.NotNil(t, c.(*container).m)
	}
}

func TestAdd(t *testing.T) {
	c := &container{m: make(map[string]any)}
	expected := &_testerImpl{}
	Add[_tester](c, expected)
	assert.Equal(t, 1, len(c.m))
	r, ok := c.m[name[_tester]()]
	assert.True(t, ok)
	assert.Equal(t, expected, r)
}

func TestAddAs(t *testing.T) {
	c := &container{m: make(map[string]any)}
	expected := &_testerImpl{}
	AddAs[_tester, _testerImpl](c, expected)
	assert.Equal(t, 1, len(c.m))
	r, ok := c.m[name[_testerImpl]()]
	assert.True(t, ok)
	assert.Equal(t, expected, r)
}

func TestGet(t *testing.T) {
	c := &container{m: make(map[string]any)}
	assert.Panics(t, func() { Get[_tester](c) })

	Add[_tester](c, &_testerImpl{})
	assert.NotNil(t, Get[_tester](c))
}

func TestGetAs(t *testing.T) {
	c := &container{m: make(map[string]any)}
	assert.Panics(t, func() { GetAs[_tester, _testerImpl](c) })

	AddAs[_tester, _testerImpl](c, &_testerImpl{})
	r := GetAs[_tester, _testerImpl](c)
	if assert.NotNil(t, r) {
		assert.IsType(t, &_testerImpl{}, r)
	}
}

func TestTryGet(t *testing.T) {
	c := &container{m: make(map[string]any)}
	r, ok := TryGetAs[_tester, _testerImpl](c)
	assert.False(t, ok)
	assert.Nil(t, r)

	expected := _testerImpl{}
	Add[_tester](c, expected)
	r, ok = TryGet[_tester](c)
	assert.True(t, ok)
	assert.Equal(t, expected, r)
}

func TestTryGetAs(t *testing.T) {
	c := &container{m: make(map[string]any)}
	r, ok := TryGetAs[_tester, _testerImpl](c)
	assert.False(t, ok)
	assert.Nil(t, r)

	expected := _testerImpl{}
	AddAs[_tester, _testerImpl](c, expected)
	r, ok = TryGetAs[_tester, _testerImpl](c)
	assert.True(t, ok)
	assert.Equal(t, expected, r)
}

func TestLen(t *testing.T) {
	c := &container{m: make(map[string]any)}
	assert.Equal(t, 0, c.Len())

	Add[_tester](c, &_testerImpl{})
	assert.Equal(t, 1, c.Len())
}

func TestRemove(t *testing.T) {
	c := &container{m: make(map[string]any)}
	Remove[_tester](c)

	Add[_tester](c, &_testerImpl{})
	Remove[_tester](c)
	assert.Equal(t, 0, len(c.m))
}

func Test_name(t *testing.T) {
	assert.Equal(t, "github.com/kiamev/generics-set/generics_tester", name[_tester]())
	assert.Equal(t, "github.com/kiamev/generics-set/generics_testerImpl", name[_testerImpl]())
}
