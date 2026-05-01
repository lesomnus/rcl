#include "type.h"

#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include <dlfcn.h>

#include <rosidl_runtime_c/message_type_support_struct.h>
#include <rosidl_runtime_c/service_type_support_struct.h>
#include <rosidl_typesupport_introspection_c/message_introspection.h>
#include <rosidl_typesupport_introspection_c/service_introspection.h>

void*
rclgo_dlopen(const char* lib_name, char* err, size_t err_size) {
    dlerror();
    void* lib = dlopen(lib_name, RTLD_LAZY | RTLD_GLOBAL);

    const char* e = dlerror();
    if (e) {
        snprintf(err, err_size, "%s", e);
        return NULL;
    }
    return lib;
}

const rosidl_message_type_support_t*
rclgo_get_message_type_support(
	void* libts,
	const char* name,
	char* err, size_t err_size
) {
    dlerror();
    typedef const rosidl_message_type_support_t* (*f_t)(void);
    f_t f = (f_t)dlsym(libts, name);

    const char* e = dlerror();
    if (e) {
        snprintf(err, err_size, "%s", e);
        return NULL;
    }
    return f();
}

const rosidl_service_type_support_t*
rclgo_get_service_type_support(
	void* libts,
	const char* name,
	char* err, size_t err_size
) {
    dlerror();
    typedef const rosidl_service_type_support_t* (*f_t)(void);
    f_t f = (f_t)dlsym(libts, name);

    const char* e = dlerror();
    if (e) {
        snprintf(err, err_size, "%s", e);
        return NULL;
    }
    return f();
}

void*
rclgo_msg_create(
	void* libts,
	const char* name,
	char* err, size_t err_size
) {
    dlerror();
    typedef void* (*f_t)(void);
    f_t f = (f_t)dlsym(libts, name);

    const char* e = dlerror();
    if (e) {
        snprintf(err, err_size, "%s", e);
        return NULL;
    }
    return f();
}

void
rclgo_msg_destroy(
	void* libts,
	const char* name,
	void* msg,
	char* err, size_t err_size
) {
    dlerror();
    typedef void (*f_t)(void*);
    f_t f = (f_t)dlsym(libts, name);

    const char* e = dlerror();
    if (e) {
        snprintf(err, err_size, "%s", e);
        return;
    }
    f(msg);
}
