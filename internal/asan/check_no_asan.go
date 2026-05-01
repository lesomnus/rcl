//go:build !asan

package asan

import (
	"fmt"
	"os"
)

func printWarning() {
	fmt.Fprintln(os.Stderr, "AddressSanitizer is not enabled; to enable, build with 'asan' option.")
}

func LeakCheck() error {
	printWarning()
	return nil
}

func RecoverableLeakCheck() bool {
	printWarning()
	return false
}
