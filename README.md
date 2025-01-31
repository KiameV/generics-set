# generics-set
A set-container that stores generics based on the stored-object's type.

The intent of this container set is to allow passing of structs based off their `type` 
by using generics.

This does support asynchronous operations but is not meant to be performant. The main
purpose for this package is to pass around injectable structs when setting servers up.

An example of usage:
```go
package test

import (
	"testing"
	
	"github.com/stretchr/testify/assert"

	set "github.com/kiamev/generics-set/set"
)

type (
	Foo1 interface {
		i()
	}
	Foo2 interface {
		Foo1
		j()
	}
	impl1 struct {
		Foo1
	}
	impl2 struct {
		Foo2
	}
)

func Test(t *testing.T) {
	var (
		// Create a new instance of the container set
		s         = set.New()
		
		i1        = &impl1{}
		i2        = impl2{}
		foo1 Foo1 = i1
		foo2 Foo2 = i2
	)
	// Adds foo2 with a lookup key that matches struct `impl1`
	set.Add(s, i1)
	// Adds foo2 with a lookup key that matches struct `impl2`
	set.Add(s, i2)
	// Adds foo2 with a lookup key that matches interface `Foo1`
	set.Add(s, foo1)
	// Adds foo2 with a lookup key that matches interface `Foo2`
	set.Add(s, foo2)
	// Adds i1 with a lookup key that matches interface `Foo1`
	// This will override the previous `set.Add(s, foo1)`
	set.AddAs[*impl1, Foo1](s, i1)
	// Adds i1 with a lookup key that matches interface `Foo1`
	// This will override `set.AddAs[*impl1, Foo1](s, i1)`
	set.AddAs[impl2, Foo1](s, i2)

	
	var gottenI1 *impl1 = set.Get[*impl1](s)
	assert.Equal(t, i1, gottenI1)

	var gottenI2 impl2 = set.Get[impl2](s)
	assert.Equal(t, i2, gottenI2)

	var gottenF1 Foo1 = set.Get[Foo1](s)
	// The underlying `Foo1` was overridden and is currently i2
	assert.Equal(t, foo2, gottenF1)
	assert.Equal(t, i2, gottenF1)

	var gottenF2 Foo2 = set.Get[Foo2](s)
	assert.Equal(t, foo2, gottenF2)
	assert.Equal(t, i2, gottenF2)
}
```
