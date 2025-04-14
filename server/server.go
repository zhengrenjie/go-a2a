package server

import (
	"context"
	"encoding/json"
	"errors"

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

func (s *A2AServer) HandleMessage(ctx context.Context, msg json.RawMessage) *protocol.JsonRpcResponse {
	var base struct {
		Version string             `json:"jsonrpc"`
		ID      uint64             `json:"id"`
		Method  protocol.A2AMethod `json:"method"`
		Params  json.RawMessage    `json:"params"`
	}

	// unmarshal raw message from request body to 'base'
	err := json.Unmarshal(msg, &base)
	if err != nil {
		return protocol.ErrJsonRpcParse.New().ToJsonRpc(0)
	}

	// check json rpc version
	if base.Version != protocol.JsonRpcVersion {
		return protocol.ErrInvalidVersion.New().
			Args(base.Version, protocol.JsonRpcVersion).
			ToJsonRpc(base.ID)
	}

	var params any
	switch base.Method {
	case protocol.MethodSendTask:
		params = new(protocol.TaskSendParams)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(base.ID)
		}

		ret, err := s.handler.SendTask(ctx, params.(*protocol.TaskSendParams))
		if err != nil {
			return s.handleError(base.ID, err)
		}

		return s.response(base.ID, ret)
	case protocol.MethodGetTask:
		params = new(protocol.TaskSendParams)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(base.ID)
		}

		ret, err := s.handler.GetTask(ctx, params.(*protocol.TaskSendParams))
		if err != nil {
			return s.handleError(base.ID, err)
		}

		return s.response(base.ID, ret)
	case protocol.MethodCancelTask:
		params = new(protocol.TaskSendParams)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(base.ID)
		}

		ret, err := s.handler.CancelTask(ctx, params.(*protocol.TaskSendParams))
		if err != nil {
			return s.handleError(base.ID, err)
		}

		return s.response(base.ID, ret)
	case protocol.MethodSetTaskPushNotifications:
		params = new(protocol.TaskPushNotificationConfig)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(base.ID)
		}

		ret, err := s.handler.SetTaskPushNotifications(ctx, params.(*protocol.TaskPushNotificationConfig))
		if err != nil {
			return s.handleError(base.ID, err)
		}

		return s.response(base.ID, ret)
	case protocol.MethodGetTaskPushNotifications:
		params = new(protocol.TaskPushNotificationConfig)
		err = json.Unmarshal(msg, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(base.ID)
		}

		ret, err := s.handler.GetTaskPushNotifications(ctx, params.(*protocol.TaskPushNotificationConfig))
		if err != nil {
			return s.handleError(base.ID, err)
		}

		return s.response(base.ID, ret)
	}

	return protocol.ErrMethodNotFound.New().
		Args(base.Method).
		ToJsonRpc(base.ID)
}

func (s *A2AServer) response(id uint64, ret any) *protocol.JsonRpcResponse {
	return &protocol.JsonRpcResponse{
		JsonRpc: protocol.JsonRpcVersion,
		ID:      id,
		Result:  ret,
	}
}

func (s *A2AServer) handleError(id uint64, err error) *protocol.JsonRpcResponse {
	var ret *protocol.Error
	if errors.As(err, &ret) {
		return ret.ToJsonRpc(id)
	}

	return protocol.ErrInternalError.New().
		Args(err.Error()).
		ToJsonRpc(id)
}
