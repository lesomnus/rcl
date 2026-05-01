//go:build cgo && asan

package asan

/*
#cgo CFLAGS: -fsanitize=address -g
#cgo LDFLAGS: -fsanitize=address
*/
import "C"
