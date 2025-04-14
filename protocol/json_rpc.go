package protocol

const JsonRpcVersion = "2.0"

type JsonRpcRequest struct {
	JsonRpc string    `json:"jsonrpc"`
	ID      uint64    `json:"id"`
	Method  A2AMethod `json:"method"`
	Params  any       `json:"params"`
}

type JsonRpcResponse struct {
	JsonRpc string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Result  any    `json:"result"`
	Error   any    `json:"error"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Error(id uint64, code int, message string, data any) *JsonRpcResponse {
	return &JsonRpcResponse{
		JsonRpc: JsonRpcVersion,
		ID:      id,
		Error: &JsonRpcError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}
