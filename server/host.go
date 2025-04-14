package server

import (
	"encoding/json"
	"io"
	"net/http"
)

type StandardA2AServerHost struct {
	addr string
}

// Host implements IA2AServerHost.
func (s *StandardA2AServerHost) Host(server *A2AServer) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: &standardHander{server: server},
	}

	return srv.ListenAndServe()
}

type standardHander struct {
	server *A2AServer
}

// ServeHTTP implements http.Handler.
func (s *standardHander) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		// TODO: handle error
	}

	ret := s.server.HandleMessage(req.Context(), json.RawMessage(body))
	respBody, err := json.Marshal(ret)
	if err != nil {
		// TODO: handle error
	}

	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(respBody)
}

var _ http.Handler = (*standardHander)(nil)
var _ IA2AServerHost = (*StandardA2AServerHost)(nil)
