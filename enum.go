// Package enum provides a way to define and work with enums in Go.
// The main benefit is to have IsValid function with zero-maintenance.
package enum

import (
	"reflect"
	"slices"
	"sync"
)

// In future, this package might relax constraint on enumType to also permit types that implement Equal(T) bool.
// https://github.com/golang/go/issues/49054
type enumType interface {
	// comparable alone is too broad constraint for enums (e.g channels are comparable),
	// thus we restrict further to ensure that only "suitable" type to be used as enum.
	comparable // currently redundant, but it is a gatekeeper to ensure enumType can always be used as map key.
	// AND any of
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~string
}

type typeID = reflect.Type

// gKey is a comparable composite key type for fast lookups across all enum values.
type gKey[T enumType] struct {
	typeID typeID
	val    T
}

var mu sync.RWMutex

// values are always slices of enumType, but can't be defined at compile time:
// https://github.com/golang/go/issues/51338
var groups = map[typeID]any{}

// keys are always gKey[enumType], but can't be defined at compile time:
// https://github.com/golang/go/issues/51338
var defs = map[any]struct{}{}

// keyOf returns unique key for each enumType T without allocating memory.
func keyOf[T enumType]() typeID {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// Def is an identity function that defines v as a valid value of enum T.
func Def[T enumType](v T) T {
	tKey := keyOf[T]()
	vKey := gKey[T]{val: v, typeID: tKey}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := defs[vKey]; ok {
		return v // already defined
	}
	defs[vKey] = struct{}{}
	vals, _ := groups[tKey].([]T)
	groups[tKey] = append(vals, v)
	return v
}

// IsValid checks if v is defined of enum T.
func IsValid[T enumType](v T) bool {
	mu.RLock()
	defer mu.RUnlock()
	k := gKey[T]{typeID: keyOf[T](), val: v}
	_, ok := defs[k]
	return ok
}

// ValuesOf returns defined values of enum T.
func ValuesOf[T enumType]() []T {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := groups[keyOf[T]()]; ok {
		return slices.Clone(v.([]T))
	}
	return nil
}
