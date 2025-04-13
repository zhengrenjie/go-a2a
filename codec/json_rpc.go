package codec

type (
	O interface {
	}

	I interface {
	}
)

type JsonRpcRequest[O any] struct {
	JsonRpc string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  O      `json:"params"`
}

type JsonRpcResponse[I any] struct {
	JsonRpc string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Result  I      `json:"result"`
}
