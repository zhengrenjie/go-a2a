package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/zhengrenjie/go-a2a/protocol"
)

var (
	ErrUmarshal   = errors.New("unmarshal request error")
	ErrBadRequest = errors.New("build http request error")
)

type (
	A2AClient struct {
		endpoint  url.URL
		header    map[string]string
		requestId atomic.Int64
		client    *http.Client
	}

	JsonRpcRaw struct {
		Version string                 `json:"jsonrpc"`
		ID      uint64                 `json:"id"`
		Method  protocol.A2AMethod     `json:"method"`
		Result  json.RawMessage        `json:"result,omitempty"`
		Error   *protocol.JsonRpcError `json:"error,omitempty"`
	}
)

// AgentCard implements protocol.IA2AProtocol.
func (a *A2AClient) AgentCard() protocol.AgentCard {
	panic("unimplemented")
}

// CancelTask implements protocol.IA2AProtocol.
func (a *A2AClient) CancelTask(ctx context.Context, params *protocol.TaskSendParams) (*protocol.Task, error) {
	ret, err := a.sendRequest(ctx, protocol.MethodCancelTask, params, false)
	if err != nil {
		return nil, err
	}

	raw := <-ret
	if raw.Error != nil {
		return nil, raw.Error
	}

	task := new(protocol.Task)
	err = json.Unmarshal(raw.Result, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetTask implements protocol.IA2AProtocol.
func (a *A2AClient) GetTask(ctx context.Context, params *protocol.TaskSendParams) (*protocol.Task, error) {
	ret, err := a.sendRequest(ctx, protocol.MethodGetTask, params, false)
	if err != nil {
		return nil, err
	}

	raw := <-ret
	if raw.Error != nil {
		return nil, raw.Error
	}

	task := new(protocol.Task)
	err = json.Unmarshal(raw.Result, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetTaskPushNotifications implements protocol.IA2AProtocol.
func (a *A2AClient) GetTaskPushNotifications(ctx context.Context, params *protocol.TaskPushNotificationConfig) (*protocol.TaskPushNotificationConfig, error) {
	ret, err := a.sendRequest(ctx, protocol.MethodGetTaskPushNotifications, params, false)
	if err != nil {
		return nil, err
	}

	raw := <-ret
	if raw.Error != nil {
		return nil, raw.Error
	}

	config := new(protocol.TaskPushNotificationConfig)
	err = json.Unmarshal(raw.Result, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SendTask implements protocol.IA2AProtocol.
func (a *A2AClient) SendTask(ctx context.Context, params *protocol.TaskSendParams) (*protocol.Task, error) {
	ret, err := a.sendRequest(ctx, protocol.MethodSendTask, params, false)
	if err != nil {
		return nil, err
	}

	raw := <-ret
	if raw.Error != nil {
		return nil, raw.Error
	}

	task := new(protocol.Task)
	err = json.Unmarshal(raw.Result, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// SetTaskPushNotifications implements protocol.IA2AProtocol.
func (a *A2AClient) SetTaskPushNotifications(ctx context.Context, params *protocol.TaskPushNotificationConfig) (*protocol.TaskPushNotificationConfig, error) {
	ret, err := a.sendRequest(ctx, protocol.MethodSetTaskPushNotifications, params, false)
	if err != nil {
		return nil, err
	}

	raw := <-ret
	if raw.Error != nil {
		return nil, raw.Error
	}

	config := new(protocol.TaskPushNotificationConfig)
	err = json.Unmarshal(raw.Result, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SubscribeTask implements protocol.IA2AProtocol.
func (a *A2AClient) SubscribeTask(ctx context.Context, params *protocol.TaskSendParams) (chan any, error) {
	// TODO: Test if stream supported by server.
	ret, err := a.sendRequest(ctx, protocol.MethodSubscribeTask, params, true)
	if err != nil {
		return nil, err
	}

	ch := make(chan any, 10)
	go func() {
		for raw := range ret {
			msg := make(map[string]any)
			err := json.Unmarshal(raw.Result, &msg)
			if err != nil {
				// TODO: log error.
			}

			if msg["status"] != nil {
				statusUpdate := new(protocol.TaskStatusUpdateEvent)
				err := json.Unmarshal(raw.Result, statusUpdate)
				if err != nil {
					// TODO: log error.
				}

				ch <- statusUpdate
			} else {
				artifactUpdate := new(protocol.TaskArtifactUpdateEvent)
				err := json.Unmarshal(raw.Result, artifactUpdate)
				if err != nil {
					// TODO: log error.
				}

				ch <- artifactUpdate
			}
		}

		close(ch)
	}()

	return ch, nil
}

// ResubscribeTask implements protocol.IA2AProtocol.
func (a *A2AClient) ResubscribeTask(ctx context.Context, params *protocol.TaskSendParams) (chan any, error) {
	return a.SubscribeTask(ctx, params)
}

// sendRequest handles both
func (a *A2AClient) sendRequest(
	ctx context.Context,
	method protocol.A2AMethod,
	params any,
	stream bool) (<-chan *JsonRpcRaw, error) {

	id := a.requestId.Add(1)
	request := &protocol.JsonRpcRequest{
		JsonRpcVersion: protocol.JsonRpcVersion,
		ID:             uint64(id),
		Method:         method,
		Params:         params,
	}

	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal request error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint.String(), bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range a.header {
		req.Header.Set(k, v)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("launch request error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("server error, http-code: %s, body: %s", resp.Status, string(body))
	}

	if !stream {
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response error: %w", err)
		}

		ch := make(chan *JsonRpcRaw, 1)

		raw := new(JsonRpcRaw)
		err = json.Unmarshal(respBody, raw)
		if err != nil {
			return nil, fmt.Errorf("unmarshal response error: %w", err)
		}

		ch <- raw
		close(ch)

		return ch, nil
	}

	ch := make(chan *JsonRpcRaw, 10)
	go a.readSSE(resp.Body, ch)

	return ch, nil
}

func (c *A2AClient) readSSE(reader io.ReadCloser, ch chan<- *JsonRpcRaw) {
	defer reader.Close()

	br := bufio.NewReader(reader)
	var data string

	for {
		// NOTICE:
		// Here we assume that one line is one event.
		// For JSON-RPC response, the format should be:
		//   - data: {"jsonrpc": "2.0", "id": 1, "result": {"taskId": "123"}}\n\n
		line, err := br.ReadString('\n')
		if err != nil {
			// TODO:
		}

		line = strings.TrimRight(line, "\r\n")

		// empty line indicates the end of the event.
		if len(line) == 0 {
			if len(data) == 0 {
				continue
			}

			raw := new(JsonRpcRaw)
			err = json.Unmarshal([]byte(data), raw)
			if err != nil {
				// TODO:
			}

			ch <- raw
			data = ""
			continue
		}

		// if not empty line, must be begin with 'data:'
		if !strings.HasPrefix(line, "data:") {
			// TODO:
			continue
		}

		data = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	}
}

var _ protocol.IA2AProtocol = (*A2AClient)(nil)
