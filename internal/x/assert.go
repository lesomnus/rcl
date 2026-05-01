package x

import (
	"reflect"
	"testing"
)

func NoError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	t.Fatalf("unexpected error: %v", err)
}

func Eq[T any](t *testing.T, expected, actual T) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		return
	}
	t.Fatalf("unexpected value: got %v, want %v", actual, expected)
}

func Len[T any](t *testing.T, value []T, size int) {
	t.Helper()
	if len(value) == size {
		return
	}
	t.Fatalf("unexpected length: got %d, want %d", len(value), size)
}
