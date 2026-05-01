package rcl

import (
	"errors"
	"fmt"
)

type RclError int

var ErrCgoDisabled = RclError(CodeCgoDisabled)

func (e RclError) Error() string {
	code := e.Code()
	return fmt.Sprintf("[%d] %s", code, code.String())
}

func (e RclError) Code() ReturnCode {
	return ReturnCode(e)
}

type ReturnCode int

func Code(err error) ReturnCode {
	var e RclError
	if err == nil {
		return CodeOk
	}
	if errors.As(err, &e) {
		return e.Code()
	}
	return CodeUnspecified
}

const (
	CodeCgoDisabled = -2
	CodeUnspecified = -1

	// The operation ran as expected.
	//
	//	RCL_RET_OK, or RMW_RET_OK
	CodeOk = ReturnCode(0)
	// Generic error to indicate operation could not complete successfully.
	//
	//	RCL_RET_ERROR, or RMW_RET_ERROR
	CodeError = ReturnCode(1)
	// The operation was halted early because it exceeded its timeout criteria.
	//
	//	RCL_RET_TIMEOUT, or RMW_RET_TIMEOUT
	CodeTimeout = ReturnCode(2)
	// The operation or event handling is not supported.
	//
	//	RCL_RET_UNSUPPORTED, or RMW_RET_UNSUPPORTED
	CodeUnsupported = ReturnCode(3)
	// Failed to allocate memory.
	//
	//	RCL_RET_BAD_ALLOC, or RMW_RET_BAD_ALLOC
	CodeBadAlloc = ReturnCode(10)
	// Argument to function was invalid.
	//
	//	RCL_RET_INVALID_ARGUMENT, or RMW_RET_INVALID_ARGUMENT
	CodeInvalidArgument = ReturnCode(11)
	// Incorrect rmw implementation.
	//
	//	RMW_RET_INCORRECT_RMW_IMPLEMENTATION
	CodeIncorrectRmwImplementation = ReturnCode(12)

	// rcl_init() already called return code.
	//
	//	RCL_RET_ALREADY_INIT
	CodeAlreadyInit = ReturnCode(100)
	// rcl_init() not yet called return code.
	//
	//	RCL_RET_NOT_INIT
	CodeNotInit = ReturnCode(101)
	// Mismatched rmw identifier return code.
	//
	//	RCL_RET_MISMATCHED_RMW_ID
	CodeMismatchedRmwId = ReturnCode(102)
	// Topic name does not pass validation.
	//
	//	RCL_RET_TOPIC_NAME_INVALID
	CodeTopicNameInvalid = ReturnCode(103)
	// Service name (same as topic name) does not pass validation.
	//
	//	RCL_RET_SERVICE_NAME_INVALID
	CodeServiceNameInvalid = ReturnCode(104)
	// Topic name substitution is unknown.
	//
	//	RCL_RET_UNKNOWN_SUBSTITUTION
	CodeUnknownSubstitution = ReturnCode(105)
	// rcl_shutdown() already called return code.
	//
	//	RCL_RET_ALREADY_SHUTDOWN
	CodeAlreadyShutdown = ReturnCode(106)
	// Resource not found.
	//
	//	RCL_RET_NOT_FOUND
	CodeNotFound = ReturnCode(107)

	/* Node */

	// Invalid rcl_node_t given return code.
	//
	//	RCL_RET_NODE_INVALID
	CodeNodeInvalid = ReturnCode(200)
	// Invalid node name return code.
	//
	//	RCL_RET_NODE_INVALID_NAME
	CodeNodeInvalidName = ReturnCode(201)
	// Invalid node namespace return code.
	//
	//	RCL_RET_NODE_INVALID_NAMESPACE
	CodeNodeInvalidNamespace = ReturnCode(202)
	// Failed to find node name.
	//
	//	RCL_RET_NODE_NAME_NON_EXISTENT, or RMW_RET_NODE_NAME_NON_EXISTENT
	CodeNodeNameNonExistent = ReturnCode(203)

	/* Publisher */

	// Invalid rcl_publisher_t given return code.
	//
	//	RCL_RET_PUBLISHER_INVALID
	CodePublisherInvalid = ReturnCode(300)

	/* Subscription */

	// Invalid rcl_subscription_t given return code.
	//
	//	RCL_RET_SUBSCRIPTION_INVALID
	CodeSubscriptionInvalid = ReturnCode(400)
	// Failed to take a message from the subscription return code.
	//
	//	RCL_RET_SUBSCRIPTION_TAKE_FAILED
	CodeSubscriptionTakeFailed = ReturnCode(401)

	/* Service client */

	// Invalid rcl_client_t given return code.
	//
	//	RCL_RET_CLIENT_INVALID
	CodeClientInvalid = ReturnCode(500)
	// Failed to take a response from the client return code.
	//
	//	RCL_RET_CLIENT_TAKE_FAILED
	CodeClientTakeFailed = ReturnCode(501)

	/* Service server */

	// Invalid rcl_service_t given return code.
	//
	//	RCL_RET_SERVICE_INVALID
	CodeServiceInvalid = ReturnCode(600)
	// Failed to take a request from the service return code.
	//
	//	RCL_RET_SERVICE_TAKE_FAILED
	CodeServiceTakeFailed = ReturnCode(601)

	/* Timer */

	// Invalid rcl_timer_t given return code.
	//
	//	RCL_RET_TIMER_INVALID
	CodeTimerInvalid = ReturnCode(800)
	// Given timer was canceled return code.
	//
	//	RCL_RET_TIMER_CANCELED
	CodeTimerCanceled = ReturnCode(801)

	/* Wait Set */

	// Invalid rcl_wait_set_t given return code.
	//
	//	RCL_RET_WAIT_SET_INVALID
	CodeWaitSetInvalid = ReturnCode(900)
	// Given rcl_wait_set_t is empty return code.
	//
	//	RCL_RET_WAIT_SET_EMPTY
	CodeWaitSetEmpty = ReturnCode(901)
	// Given rcl_wait_set_t is full return code.
	//
	//	RCL_RET_WAIT_SET_FULL
	CodeWaitSetFull = ReturnCode(902)

	/* Argument parsing */

	// Argument is not a valid remap rule.
	//
	//	RCL_RET_INVALID_REMAP_RULE
	CodeInvalidRemapRule = ReturnCode(1001)
	// Expected one type of lexeme but got another.
	//
	//	RCL_RET_WRONG_LEXEME
	CodeWrongLexeme = ReturnCode(1002)
	// Found invalid ros argument while parsing.
	//
	//	RCL_RET_INVALID_ROS_ARGS
	CodeInvalidRosArgs = ReturnCode(1003)
	// Argument is not a valid parameter rule.
	//
	//	RCL_RET_INVALID_PARAM_RULE
	CodeInvalidParamRule = ReturnCode(1010)
	// Argument is not a valid log level rule.
	//
	//	RCL_RET_INVALID_LOG_LEVEL_RULE
	CodeInvalidLogLevelRule = ReturnCode(1020)

	/* Event */

	// Invalid rcl_event_t given return code.
	//
	//	RCL_RET_EVENT_INVALID
	CodeEventInvalid = ReturnCode(2000)
	// Failed to take an event from the event handle.
	//
	//	RCL_RET_EVENT_TAKE_FAILED
	CodeEventTakeFailed = ReturnCode(2001)

	/* Lifecycle State */

	// rcl_lifecycle state registered.
	//
	//	RCL_RET_LIFECYCLE_STATE_REGISTERED
	CodeLifecycleStateRegistered = ReturnCode(3000)
	// rcl_lifecycle state not registered.
	//
	//	RCL_RET_LIFECYCLE_STATE_NOT_REGISTERED
	CodeLifecycleStateNotRegistered = ReturnCode(3001)

	/* Action */

	// No terminal timestamp for the goal as it has not reached a terminal state.
	//
	//	RCL_ACTION_RET_NOT_TERMINATED_YET
	CodeActionNotTerminatedYet = ReturnCode(4001)
)

func (e ReturnCode) String() string {
	switch e {
	case CodeCgoDisabled:
		return "cgo is disabled"
	case CodeUnspecified:
		return "unspecified; maybe not an rcl error"

	case CodeOk:
		return "ok"
	case CodeError:
		return "RMW error"
	case CodeTimeout:
		return "timeout"
	case CodeUnsupported:
		return "unsupported"
	case CodeBadAlloc:
		return "bad allocation"
	case CodeInvalidArgument:
		return "invalid argument"
	case CodeIncorrectRmwImplementation:
		return "incorrect RMW implementation"
	case CodeAlreadyInit:
		return "already initialized"
	case CodeNotInit:
		return "not initialized"
	case CodeMismatchedRmwId:
		return "RMW ID mismatched"
	case CodeTopicNameInvalid:
		return "invalid topic name"
	case CodeServiceNameInvalid:
		return "invalid service name"
	case CodeUnknownSubstitution:
		return "unknown substitution"
	case CodeAlreadyShutdown:
		return "already shutdown"
	case CodeNotFound:
		return "not found"
	case CodeNodeInvalid:
		return "invalid node"
	case CodeNodeInvalidName:
		return "invalid node name"
	case CodeNodeInvalidNamespace:
		return "invalid node namespace"
	case CodeNodeNameNonExistent:
		return "node name does not exist"
	case CodePublisherInvalid:
		return "invalid publisher"
	case CodeSubscriptionInvalid:
		return "invalid subscription"
	case CodeSubscriptionTakeFailed:
		return "subscription take failed"
	case CodeClientInvalid:
		return "invalid client"
	case CodeClientTakeFailed:
		return "client take failed"
	case CodeServiceInvalid:
		return "invalid service"
	case CodeServiceTakeFailed:
		return "service take failed"
	case CodeTimerInvalid:
		return "invalid timer"
	case CodeTimerCanceled:
		return "timer canceled"
	case CodeWaitSetInvalid:
		return "invalid wait set"
	case CodeWaitSetEmpty:
		return "wait set empty"
	case CodeWaitSetFull:
		return "wait set full"
	case CodeInvalidRemapRule:
		return "invalid remap rule"
	case CodeWrongLexeme:
		return "wrong lexeme"
	case CodeInvalidRosArgs:
		return "invalid ROS arguments"
	case CodeInvalidParamRule:
		return "invalid parameter rule"
	case CodeInvalidLogLevelRule:
		return "invalid log level rule"
	case CodeEventInvalid:
		return "invalid event"
	case CodeEventTakeFailed:
		return "event take failed"
	case CodeLifecycleStateRegistered:
		return "lifecycle state registered"
	case CodeLifecycleStateNotRegistered:
		return "lifecycle state not registered"

	default:
		return fmt.Sprintf("unknown error code: %d", int(e))
	}
}
