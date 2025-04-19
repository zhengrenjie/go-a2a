package protocol

import (
	"encoding/json"
	"fmt"
)

const JsonRpcVersion = "2.0"

type JsonRpcRequest struct {
	JsonRpcVersion string    `json:"jsonrpc"`
	ID             uint64    `json:"id"`
	Method         A2AMethod `json:"method"`
	Params         any       `json:"params"`
}

type JsonRpcResponse struct {
	JsonRpcVersion string        `json:"jsonrpc"`
	ID             uint64        `json:"id"`
	Result         any           `json:"result,omitempty"`
	Error          *JsonRpcError `json:"error,omitempty"`
}

func (r *JsonRpcResponse) ToByte() []byte {
	ret, _ := json.Marshal(r)
	return ret
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *JsonRpcError) Error() string {
	if e.Data == nil {
		return fmt.Sprintf("JSON-RPC Error, Code: [%d], Message: %s", e.Code, e.Message)
	}

	data, _ := json.Marshal(e.Data)
	return fmt.Sprintf(
		"JSON-RPC Error with data, Code: [%d], Message: %s, Data: %s",
		e.Code,
		e.Message,
		string(data),
	)
}

// Type assertion to ensure JsonRpcError implements the error interface.
var _ error = (*JsonRpcError)(nil)
