package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zhengrenjie/go-a2a/protocol"
)

type (
	StandardA2AServerHost struct {
		addr string
	}

	JsonRpcRaw struct {
		Version string             `json:"jsonrpc"`
		ID      uint64             `json:"id"`
		Method  protocol.A2AMethod `json:"method"`
		Params  json.RawMessage    `json:"params"`
	}
)

// Host implements IA2AServerHost.
func (s *StandardA2AServerHost) Host(server *A2AServer) error {
	http.HandleFunc("/.well-known/agent.json", func(resp http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		agentCard := server.AgentCard()
		agentCardJson, err := json.Marshal(agentCard)
		if err != nil {
			// TODO: handle error
		}

		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(agentCardJson)
	})

	http.Handle("/", &standardHander{server: server})
	return http.ListenAndServe(s.addr, nil)
}

type standardHander struct {
	server            *A2AServer
	keepAlive         bool
	keepAliveInterval time.Duration
}

// ServeHTTP implements http.Handler.
func (s *standardHander) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	// if request body cannot read, return 400
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// try unmarshal json-rpc request to 'raw'
	// get basic rpc info:
	//  - Version
	//  - ID
	//  - Method
	raw := new(JsonRpcRaw)
	err = json.Unmarshal(body, raw)
	if err != nil {
		response(w, protocol.ErrJsonRpcParse.New().ToJsonRpc(0))
		return
	}

	// check json rpc version
	if raw.Version != protocol.JsonRpcVersion {
		response(
			w,
			protocol.
				ErrInvalidVersion.
				New().
				Args(raw.Version, protocol.JsonRpcVersion).
				ToJsonRpc(raw.ID),
		)

		return
	}

	// route by 'method' in rpc
	if isStreaming(raw.Method) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported on server", http.StatusInternalServerError)
			return
		}

		respCh := make(chan *protocol.JsonRpcResponse, 10)
		go func() {
			s.server.HandleStreaming(req.Context(), raw, respCh)
		}()

		eventCh := make(chan string, 10)

		go func() {
			for resp := range respCh {
				data, err := json.Marshal(resp)
				if err == nil {
					msg := fmt.Sprintf("event: message\ndata: %s\n\n", data)
					eventCh <- msg
				}
			}

			// no more response, close the event channel
			close(eventCh)
		}()

		if s.keepAlive {
			go func() {
				ticker := time.NewTicker(s.keepAliveInterval)

				defer func() {
					// if eventCh is closed, ping msg will panic.
					// recover it but ignore.
					if r := recover(); r != nil {
						// ignore panic
					}

					ticker.Stop()
				}()

				for {
					select {
					case <-ticker.C:
						//: ping - 2025-03-27 07:44:38.682659+00:00
						eventCh <- fmt.Sprintf(":ping - %s\n\n", time.Now().Format(time.RFC3339))
					case <-req.Context().Done():
						return
					}
				}
			}()
		}

		for {
			select {
			case event, more := <-eventCh:
				fmt.Fprint(w, event)
				flusher.Flush()

				if !more {
					return
				}
			case <-req.Context().Done():
				return
			}
		}
	}

	resp := s.server.HandleMessage(req.Context(), raw)
	response(w, resp)
}

func response(w http.ResponseWriter, resp *protocol.JsonRpcResponse) {
	// write json-rpc response to http-response
	// http, as the "transport" layer for json-rpc, the status code is 200 even if error occurs in RPC.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp.ToByte())
}

func isStreaming(method protocol.A2AMethod) bool {
	return method == protocol.MethodSubscribeTask || method == protocol.MethodResubscribeTask
}

func NewA2AHost(addr string) *StandardA2AServerHost {
	return &StandardA2AServerHost{addr: addr}
}

var _ http.Handler = (*standardHander)(nil)
var _ IA2AServerHost = (*StandardA2AServerHost)(nil)
