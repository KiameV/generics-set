package set

import (
	"reflect"
	"sync"
)

type (
	// Container defines the container to be used in subsequent function calls that are defined in this package.
	// Note: This is not meant to be highly performant and is mainly designed to more-easily coordinate injection at
	// startup.
	Container interface {
		Len() int
	}
	// container is the default implementation of the Container interface.
	container struct {
		mu sync.Mutex
		m  map[string]any
	}
)

// New creates a new generic set container.
// Note: This is not meant to be highly performant and is mainly designed to more-easily coordinate injection at
// startup.
//
// Example:
//
//	type (
//	  fooer interface{}
//	  fooImpl struct {}
//	)
//	var (
//	  c = set.New
//	  foo fooer = &fooImpl
//	  bar := &fooImpl
//	)
//	set.Add(c, foo)
//	set.Add(c, foo)
//	i := set.Get[fooer](c) // `i` will be `foo`
//	j := set.Get[*fooImpl](c) // `j` will be `bar`
func New() Container {
	return &container{
		m: make(map[string]any),
	}
}

// Add adds an item to the Container.
func Add[T any](c Container, e T) {
	c.(*container).lock()
	defer c.(*container).unlock()
	c.(*container).m[name[T]()] = e
}

// AddAs adds an item of T to the Container with the key as `O`.
//
// Example:
//
//	type foo interface{}
//	type fooImpl1 struct {}
//	type fooImpl2 struct {}
//
//	func __() {
//	    c := set.New()
//	    set.AddAs[foo, fooImpl1](c, &fooImpl1{})
//	    set.AddAs[foo, fooImpl2](c, &fooImpl2{})
//
//	    fooImpl1 foo := set.GetAs[foo, fooImpl1](c).Foo()
//	    fooImpl2 foo := set.GetAs[foo, fooImpl2](c).Foo()
//	}
func AddAs[T any, O any](c Container, e T) {
	c.(*container).lock()
	defer c.(*container).unlock()
	c.(*container).m[name[O]()] = e
}

func (c *container) Len() int {
	return len(c.m)
}

func (c *container) lock() {
	c.mu.Lock()
}

func (c *container) unlock() {
	c.mu.Unlock()
}

// Get gets an item T from the container as a type of O and will panic if no container of the type is found.
func Get[T any](c Container) T {
	c.(*container).lock()
	defer c.(*container).unlock()
	n := name[T]()
	t, ok := c.(*container).m[n].(T)
	if !ok {
		panic("item not found of type [" + n + "]")
	}
	return t
}

// GetAs gets an item from the container and will panic if no container of the type is found.
//
// Example:
//
//	type foo interface{}
//	type fooImpl1 struct {}
//	type fooImpl2 struct {}
//
//	func __() {
//	    c := set.New()
//	    set.AddAs[foo, fooImpl1](c, &fooImpl1{})
//	    set.AddAs[foo, fooImpl2](c, &fooImpl2{})
//
//	    set.GetAs[foo, fooImpl1](c).Foo()
//	    set.GetAs[foo, fooImpl2](c).Foo()
//	}
func GetAs[T any, O any](c Container) T {
	c.(*container).lock()
	defer c.(*container).unlock()
	n := name[O]()
	t, ok := c.(*container).m[n].(T)
	if !ok {
		panic("item not found of type [" + n + "]")
	}
	return t
}

// TryGet gets an item from the container.
func TryGet[T any](c Container) (T, bool) {
	c.(*container).lock()
	defer c.(*container).unlock()
	t, ok := c.(*container).m[name[T]()].(T)
	return t, ok
}

// TryGetAs gets an item T from the container using the type O.
//
// Example:
//
//	type foo interface{}
//	type fooImpl1 struct {}
//	type fooImpl2 struct {}
//
//	func __() {
//		c := set.New()
//		set.AddAs[foo, fooImpl1](c, &fooImpl1{})
//		set.AddAs[foo, fooImpl2](c, &fooImpl2{})
//
//		if c, ok := set.TryGetAs[foo, fooImpl1](c); ok {
//			c.Foo()
//		}
//		if c, ok := set.TryGetAs[foo, fooImpl2](c); ok {
//			c.Foo()
//		}
//	}
func TryGetAs[T any, O any](c Container) (T, bool) {
	c.(*container).lock()
	defer c.(*container).unlock()
	t, ok := c.(*container).m[name[O]()].(T)
	return t, ok
}

// Remove removes an item from the container.
func Remove[T any](c Container) {
	c.(*container).lock()
	defer c.(*container).unlock()
	delete(c.(*container).m, name[T]())
}

// name gets the name of a type.
func name[T any]() string {
	var (
		t T
		i = reflect.TypeOf(&t)
	)
	for i.Kind() == reflect.Ptr {
		i = i.Elem()
	}
	return i.PkgPath() + i.Name()
}
