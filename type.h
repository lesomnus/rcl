#include <rosidl_runtime_c/message_type_support_struct.h>
#include <rosidl_runtime_c/service_type_support_struct.h>

void*
rclgo_dlopen(const char* lib_name, char* err, size_t err_size);

const rosidl_message_type_support_t*
rclgo_get_message_type_support(
	void* libts,
	const char* name,
	char* err, size_t err_size
);

const rosidl_service_type_support_t*
rclgo_get_service_type_support(
	void* libts,
	const char* name,
	char* err, size_t err_size
);

void*
rclgo_msg_create(
	void* libts,
	const char* name,
	char* err, size_t err_size
);

void
rclgo_msg_destroy(
	void* libts,
	const char* name,
	void* msg,
	char* err, size_t err_size
);
