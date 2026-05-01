//go:build cgo && asan

package asan

/*
#include <sanitizer/lsan_interface.h>
#include <sanitizer/common_interface_defs.h>
#include <stdlib.h>
#include <unistd.h>

void do_leak_check() {
    __lsan_do_leak_check();
}

int do_recoverable_leak_check() {
    return __lsan_do_recoverable_leak_check();
}

// Redirects sanitizer output to a pipe, runs recoverable leak check,
// restores stderr, then returns the captured output. Caller must free().
// Returns NULL if no leaks were found.
char* recoverable_leak_check_message() {
    int pipefd[2];
    if (pipe(pipefd) == -1) return NULL;

    __sanitizer_set_report_fd((void*)(intptr_t)pipefd[1]);
    int leaked = __lsan_do_recoverable_leak_check();
    __sanitizer_set_report_fd((void*)(intptr_t)STDERR_FILENO);
    close(pipefd[1]);

    if (!leaked) {
        close(pipefd[0]);
        return NULL;
    }

    char*  buf   = NULL;
    size_t total = 0;
    char   tmp[4096];
    ssize_t n;
    while ((n = read(pipefd[0], tmp, sizeof(tmp))) > 0) {
        char* next = realloc(buf, total + (size_t)n + 1);
        if (!next) { free(buf); close(pipefd[0]); return NULL; }
        buf = next;
        __builtin_memcpy(buf + total, tmp, (size_t)n);
        total += (size_t)n;
    }
    close(pipefd[0]);

    if (buf) buf[total] = '\0';
    return buf;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// RecoverableLeakCheck reports leaks without aborting. Returns true if leaks were found.
func RecoverableLeakCheck() bool {
	return C.do_recoverable_leak_check() != 0
}

func LeakCheck() error {
	msg := C.recoverable_leak_check_message()
	if msg == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(msg))
	return fmt.Errorf(C.GoString(msg))
}
