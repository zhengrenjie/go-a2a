package protocol

import "encoding/json"

const JsonRpcVersion = "2.0"

type JsonRpcRequest struct {
	JsonRpcVersion string    `json:"jsonrpc"`
	ID             uint64    `json:"id"`
	Method         A2AMethod `json:"method"`
	Params         any       `json:"params"`
}

type JsonRpcResponse struct {
	JsonRpcVersion string `json:"jsonrpc"`
	ID             uint64 `json:"id"`
	Result         any    `json:"result,omitempty"`
	Error          any    `json:"error,omitempty"`
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
