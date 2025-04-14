package server

import (
	"encoding/json"
	"io"
	"net/http"
)

func NewA2AHost(addr string) *StandardA2AServerHost {
	return &StandardA2AServerHost{addr: addr}
}

type StandardA2AServerHost struct {
	addr string
}

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
	server *A2AServer
}

// ServeHTTP implements http.Handler.
func (s *standardHander) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	// if body cannot read, return 400
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	// handle request
	respBody := s.server.HandleMessage(req.Context(), json.RawMessage(body)).ToByte()

	// write json-rpc response to http-response
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(respBody)
}

var _ http.Handler = (*standardHander)(nil)
var _ IA2AServerHost = (*StandardA2AServerHost)(nil)
