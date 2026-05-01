package asan

import (
	"runtime"
	"testing"
)

func NoLeak(t *testing.T) {
	t.Cleanup(func() {
		runtime.GC()
		if err := LeakCheck(); err != nil {
			t.Error(err)
		}
	})
}
