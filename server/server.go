package server

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/zhengrenjie/go-a2a/protocol"
)

func NewA2AServer(p protocol.IA2AProtocol) *A2AServer {
	return &A2AServer{handler: p}
}

type A2AServer struct {
	handler protocol.IA2AProtocol
}

func (s *A2AServer) AgentCard() protocol.AgentCard {
	return s.handler.AgentCard()
}

// HandleMessage handles the following no-streaming methods:
//   - tasks/send
//   - tasks/get
//   - tasks/cancel
//   - tasks/pushNotification/set
//   - tasks/pushNotification/get
func (s *A2AServer) HandleMessage(ctx context.Context, raw *JsonRpcRaw) *protocol.JsonRpcResponse {

	var params any
	switch raw.Method {
	case protocol.MethodSendTask:
		params = new(protocol.TaskSendParams)
		err := json.Unmarshal(raw.Params, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(raw.ID)
		}

		ret, err := s.handler.SendTask(ctx, params.(*protocol.TaskSendParams))
		if err != nil {
			return s.handleError(raw.ID, err)
		}

		return s.response(raw.ID, ret)
	case protocol.MethodGetTask:
		params = new(protocol.TaskSendParams)
		err := json.Unmarshal(raw.Params, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(raw.ID)
		}

		ret, err := s.handler.GetTask(ctx, params.(*protocol.TaskSendParams))
		if err != nil {
			return s.handleError(raw.ID, err)
		}

		return s.response(raw.ID, ret)
	case protocol.MethodCancelTask:
		params = new(protocol.TaskSendParams)
		err := json.Unmarshal(raw.Params, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(raw.ID)
		}

		ret, err := s.handler.CancelTask(ctx, params.(*protocol.TaskSendParams))
		if err != nil {
			return s.handleError(raw.ID, err)
		}

		return s.response(raw.ID, ret)
	case protocol.MethodSetTaskPushNotifications:
		params = new(protocol.TaskPushNotificationConfig)
		err := json.Unmarshal(raw.Params, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(raw.ID)
		}

		ret, err := s.handler.SetTaskPushNotifications(ctx, params.(*protocol.TaskPushNotificationConfig))
		if err != nil {
			return s.handleError(raw.ID, err)
		}

		return s.response(raw.ID, ret)
	case protocol.MethodGetTaskPushNotifications:
		params = new(protocol.TaskPushNotificationConfig)
		err := json.Unmarshal(raw.Params, params)
		if err != nil {
			return protocol.ErrJsonRpcParamsParse.New().ToJsonRpc(raw.ID)
		}

		ret, err := s.handler.GetTaskPushNotifications(ctx, params.(*protocol.TaskPushNotificationConfig))
		if err != nil {
			return s.handleError(raw.ID, err)
		}

		return s.response(raw.ID, ret)
	}

	return protocol.ErrMethodNotFound.New().
		Args(raw.Method).
		ToJsonRpc(raw.ID)
}

func (s *A2AServer) HandleStreaming(ctx context.Context, raw *JsonRpcRaw, streaming <-chan *protocol.JsonRpcResponse) {

}

func (s *A2AServer) response(id uint64, ret any) *protocol.JsonRpcResponse {
	return &protocol.JsonRpcResponse{
		JsonRpcVersion: protocol.JsonRpcVersion,
		ID:             id,
		Result:         ret,
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
