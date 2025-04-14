package protocol

import "fmt"

const (
	CodeJSONParse                  = -32700
	CodeInvalidRequest             = -32600
	CodeMethodNotFound             = -32601
	CodeInvalidParams              = -32602
	CodeInternalError              = -32603
	CodeTaskNotFound               = -32001
	CodeTaskCannotCancel           = -32002
	CodePushNotificationNotSupport = -32003
	CodeUnsupportedOperation       = -32004
	CodeIncompatibleContentTypes   = -32005
)

// Etyp is the shortcut for [NewErrorType].
var Etyp = NewErrorType

var (
	// CodeJSONParse errors.
	ErrJsonRpcParse       = Etyp(CodeJSONParse, "Parse HTTP request to JSON-RPC error")
	ErrJsonRpcParamsParse = Etyp(CodeJSONParse, "Parse JSON-RPC params error")

	// InvalidRequest errors.

	// ErrInvalidVersion
	// Args: [client version], [supported version]
	ErrInvalidVersion = Etyp(CodeInvalidRequest, "Invalid JSON-RPC version: [%s], expected [%s]")

	// MethodNotFound errors.
	// Args: [request method]
	ErrMethodNotFound = Etyp(CodeMethodNotFound, "Method [%s] not found")

	// CodeInternalError errors.
	// Args: [error message]
	ErrInternalError = Etyp(CodeInternalError, "Internal error occurred: [%s]")
)

type (
	ErrorType struct {
		code   int
		format string
	}

	Error struct {
		typ ErrorType

		code    int
		message string
		args    []any
		data    any

		stack error
	}
)

func (e ErrorType) New() *Error {
	return &Error{
		typ:     e,
		code:    e.code,
		message: e.format,
	}
}

func (e *Error) Args(args ...any) *Error {
	e.args = args
	return e
}

func (e *Error) Stack(err error) *Error {
	e.stack = err
	return e
}

func (e *Error) Data(data any) *Error {
	e.data = data
	return e
}

func (e *Error) ToJsonRpc(id uint64) *JsonRpcResponse {
	return &JsonRpcResponse{
		JsonRpc: JsonRpcVersion,
		ID:      id,
		Error: &JsonRpcError{
			Code:    e.code,
			Message: e.message,
			Data:    e.data,
		},
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf(e.message, e.args...)
}

func NewErrorType(code int, format string) ErrorType {
	return ErrorType{
		code:   code,
		format: format,
	}
}

func Is(err error, tpy ErrorType) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}

	return e.typ == tpy
}

var _ error = (*Error)(nil)
