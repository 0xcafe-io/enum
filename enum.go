// Package enum provides a way to define and work with enums in Go.
// The main benefit is to have IsValid function with zero-maintenance.
package enum

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
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

// typeID is a unique identifier for each enum type.
type typeID reflect.Type

// typeValue is used as composite key for fast lookups across all enum values.
type typeValue[T enumType] struct {
	typ typeID
	val T
}

var mu sync.RWMutex

// values are always slices of enumType, but can't be defined at compile time:
// https://github.com/golang/go/issues/51338
var groups = map[typeID]any{}

// keys are always typeValue[enumType], but can't be defined at compile time:
// https://github.com/golang/go/issues/51338
var defs = map[any]struct{}{}

// Def defines v as a valid value of enum T and returns it.
// Value is returned as-is, without any wrapping or conversion.
// Duplicate definitions are ignored.
// Usage:
//   type Status string
//   var (
//     StatusDraft  = enum.Def(Status("draft")) // type of StatusDraft is Status
//     StatusOpen   = enum.Def[Status]("open") // same thing, alternative syntax
//   )
func Def[T enumType](v T) T {
	typID := idOf[T]()
	vKey := typeValue[T]{val: v, typ: typID}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := defs[vKey]; ok {
		return v // already defined
	}
	defs[vKey] = struct{}{}
	vals, _ := groups[typID].([]T)
	groups[typID] = append(vals, v)
	return v
}

// IsValid reports whether v is defined for enum T.
func IsValid[T enumType](v T) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := defs[typeValue[T]{typ: idOf[T](), val: v}]
	return ok
}

// Validate checks whether v is defined for enum T.
// If not, returns an error, otherwise returns nil.
func Validate[T enumType](v T) error {
	// TODO cache error msg to avoid constructing it every time.
	typ := idOf[T]()
	mu.RLock()
	defer mu.RUnlock()
	_, valueExists := defs[typeValue[T]{typ: typ, val: v}]
	if !valueExists {
		vals, enumExists := groups[typ]
		if !enumExists {
			return fmt.Errorf("%s doesn't have any definition", typ.Name())
		}
		s, _ := vals.([]T)
		fmtVerb := "%v"
		if typ.Kind() == reflect.String {
			fmtVerb = "%q" // use quotes for strings to visually distinguish them from integers
		}
		return errors.New(errMsg(fmtVerb, v, s))
	}
	return nil
}

// ValuesOf returns defined values of enum T.
// Values are returned in the order they were mentioned (see https://go.dev/ref/spec#Package_initialization).
// It is safe to modify the returned slice.
func ValuesOf[T enumType]() []T {
	mu.RLock()
	defer mu.RUnlock()
	if vals, ok := groups[idOf[T]()]; ok {
		return slices.Clone(vals.([]T))
	}
	return nil
}

// idOf returns unique typeID for each T without instantiating.
func idOf[T enumType]() typeID {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func errMsg[T enumType](fmtVerb string, invalidVal T, vals []T) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(fmtVerb+" is not a valid choice, allowed values are: ", invalidVal))
	// vals are guaranteed to be non-empty for defined enums
	sb.WriteString(fmt.Sprintf(fmtVerb, vals[0]))
	tail := ", " + fmtVerb
	for _, v := range vals[1:] {
		sb.WriteString(fmt.Sprintf(tail, v))
	}
	return sb.String()
}

// Clear removes all definitions for enum T.
func Clear[T enumType]() {
	mu.Lock()
	defer mu.Unlock()
	typID := idOf[T]()
	delete(groups, typID)
	for k := range defs {
		if v, ok := k.(typeValue[T]); ok && v.typ == typID {
			delete(defs, k)
		}
	}
}
