package server

import (
	"context"
	"encoding/json"

	"github.com/zhengrenjie/go-a2a/protocol"
)

type IA2AServerHost interface {
	Host(server *A2AServer) error
}

func NewA2AServer(p protocol.IA2AProtocol) *A2AServer {
	return &A2AServer{handler: p}
}

type A2AServer struct {
	handler protocol.IA2AProtocol
}

func (s *A2AServer) AgentCard() protocol.AgentCard {
	return s.handler.AgentCard()
}

func (s *A2AServer) HandleMessage(ctx context.Context, msg json.RawMessage) protocol.JsonRpcResponse {
	var base struct {
		JsonRpc string             `json:"jsonrpc"`
		ID      uint64             `json:"id"`
		Method  protocol.A2AMethod `json:"method"`
		Params  json.RawMessage    `json:"params"`
	}

	err := json.Unmarshal(msg, &base)
	if err != nil {
		return *protocol.Error(0, protocol.ErrJSONParse, "Invalid JSON", nil)
	}

	if base.JsonRpc != protocol.JsonRpcVersion {
		return *protocol.Error(base.ID, protocol.ErrInvalidRequest, "Invalid JSON-RPC version", nil)
	}

	var params any
	switch base.Method {
	case protocol.MethodSendTask:
		params = new(protocol.TaskSendParams)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return *protocol.Error(base.ID, protocol.ErrJSONParse, "Invalid JSON", nil)
		}

		ret, _ := s.handler.SendTask(ctx, params.(*protocol.TaskSendParams))

		// TODO: handle error
		return s.response(base.ID, ret)
	case protocol.MethodGetTask:
		params = new(protocol.TaskSendParams)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return *protocol.Error(base.ID, protocol.ErrJSONParse, "Invalid JSON", nil)
		}

		ret, _ := s.handler.GetTask(ctx, params.(*protocol.TaskSendParams))

		// TODO: handle error
		return s.response(base.ID, ret)
	case protocol.MethodCancelTask:
		params = new(protocol.TaskSendParams)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return *protocol.Error(base.ID, protocol.ErrJSONParse, "Invalid JSON", nil)
		}

		ret, _ := s.handler.CancelTask(ctx, params.(*protocol.TaskSendParams))
		// TODO: handle error
		return s.response(base.ID, ret)
	case protocol.MethodSetTaskPushNotifications:
		params = new(protocol.TaskPushNotificationConfig)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return *protocol.Error(base.ID, protocol.ErrJSONParse, "Invalid JSON", nil)
		}

		ret, _ := s.handler.SetTaskPushNotifications(ctx, params.(*protocol.TaskPushNotificationConfig))
		// TODO: handle error
		return s.response(base.ID, ret)
	case protocol.MethodGetTaskPushNotifications:
		params = new(protocol.TaskPushNotificationConfig)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return *protocol.Error(base.ID, protocol.ErrJSONParse, "Invalid JSON", nil)
		}

		ret, _ := s.handler.GetTaskPushNotifications(ctx, params.(*protocol.TaskPushNotificationConfig))
		// TODO: handle error
		return s.response(base.ID, ret)
	}

	return *protocol.Error(base.ID, protocol.ErrMethodNotFound, "Method not found", nil)
}

func (s *A2AServer) response(id uint64, ret any) protocol.JsonRpcResponse {
	return protocol.JsonRpcResponse{
		JsonRpc: protocol.JsonRpcVersion,
		ID:      id,
		Result:  ret,
	}
}
